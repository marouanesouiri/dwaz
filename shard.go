/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import (
	"context"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

/*******************************
 * Shards Identify Rate Limiter
 *******************************/

// ShardsIdentifyRateLimiter defines the interface for a rate limiter
// that controls the frequency of Identify payloads sent per shard.
//
// Implementations block the caller in Wait() until an Identify token is available.
type ShardsIdentifyRateLimiter interface {
	// Wait blocks until the shard is allowed to send an Identify payload.
	Wait()
}

// DefaultShardsRateLimiter implements a simple token bucket
// rate limiter using a buffered channel of tokens.
//
// The capacity and refill interval control the max burst and rate.
type DefaultShardsRateLimiter struct {
	tokens chan struct{}
}

var _ ShardsIdentifyRateLimiter = (*DefaultShardsRateLimiter)(nil)

// NewDefaultShardsRateLimiter creates a new token bucket rate limiter.
//
// r specifies the maximum burst tokens allowed.
// interval specifies how frequently tokens are refilled.
func NewDefaultShardsRateLimiter(r int, interval time.Duration) *DefaultShardsRateLimiter {
	rl := &DefaultShardsRateLimiter{tokens: make(chan struct{}, r)}
	// fill initial tokens
	for range r {
		rl.tokens <- struct{}{}
	}
	// refill tokens periodically in a goroutine
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()
	return rl
}

// Wait blocks until a token is available for sending Identify.
func (rl *DefaultShardsRateLimiter) Wait() {
	<-rl.tokens
}

/*************************************
 * Shard: a single Gateway connection
 *************************************/

// Shard manages a single WebSocket connection to Discord Gateway,
// including session state, event handling, heartbeats, and reconnects.
type Shard struct {
	shardID     int           // shard number (zero-based)
	totalShards int           // total number of shards in the bot
	token       string        // Discord bot token
	intents     GatewayIntent // Gateway intents bitmask
	url         string        // Gateway URL to connect to

	logger          Logger                    // logger interface for informational and error messages
	dispatcher      *dispatcher               // event dispatcher for received Gateway events
	identifyLimiter ShardsIdentifyRateLimiter // rate limiter controlling Identify payloads

	conn net.Conn // websocket connection

	seq       int64  // last received sequence number from Gateway
	sessionID string // current session id for resuming
	resumeURL string // Gateway URL to resume session on

	latency          int64       // heartbeat latency in milliseconds
	lastHeartbeatACK atomic.Bool // true if last heartbeat was acknowledged
}

// newShard constructs a new Shard instance with the specified parameters.
//
// shardID and totalShards configure the sharding info,
// token and url set authentication and gateway endpoint,
// intents specify Gateway events to receive,
// logger and dispatcher handle logging and event dispatching,
// limiter enforces Identify rate limits.
func newShard(
	shardID, totalShards int, token, url string, intents GatewayIntent,
	logger Logger, dispatcher *dispatcher, limiter ShardsIdentifyRateLimiter,
) *Shard {
	return &Shard{
		shardID:         shardID,
		totalShards:     totalShards,
		token:           token,
		url:             url,
		intents:         intents,
		logger:          logger,
		dispatcher:      dispatcher,
		identifyLimiter: limiter,
	}
}

// Connect establishes or resumes a WebSocket connection to Discord Gateway
//
// The shard attempts to connect to the resumeURL if set, otherwise
// to the default gateway url.
//
// It spawns a goroutine to read messages asynchronously.
func (s *Shard) connect(ctx context.Context) error {
	if s.conn != nil {
		s.conn.Close()
	}

	url := s.resumeURL
	if url == "" {
		url = s.url
	}

	dialer := ws.Dialer{}

	conn, _, _, err := dialer.Dial(ctx, url)
	if err != nil {
		return err
	}

	s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " connected")
	s.conn = conn
	s.lastHeartbeatACK.Store(true)

	go s.readLoop()
	return nil
}

// readLoop continuously reads messages from the Gateway WebSocket
//
// It handles Gateway opcodes, dispatches events, and triggers reconnects as needed.
func (s *Shard) readLoop() {
	for {
		msg, op, err := wsutil.ReadServerData(s.conn)
		if err != nil {
			s.logger.Error("Shard " + strconv.Itoa(s.shardID) + " read error: " + err.Error())
			s.reconnect()
			return
		}

		if op != ws.OpText {
			continue
		}

		var payload gatewayPayload
		if err := sonic.Unmarshal(msg, &payload); err != nil {
			s.logger.Error("Shard " + strconv.Itoa(s.shardID) + " unmarshal error: " + err.Error())
			continue
		}

		switch payload.Op {
		case gatewayOpcodeDispatch:
			atomic.StoreInt64(&s.seq, payload.S)
			s.dispatcher.dispatch(s.shardID, payload.T, payload.D)

			if payload.T == "READY" {
				var ready struct {
					SessionID string `json:"session_id"`
					ResumeURL string `json:"resume_gateway_url"`
				}
				sonic.Unmarshal(payload.D, &ready)
				s.sessionID = ready.SessionID
				s.resumeURL = ready.ResumeURL
				s.logger.Debug("Shard " + strconv.Itoa(s.shardID) + " session established")
			}

		case gatewayOpcodeReconnect:
			s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " RECONNECT received")
			s.reconnect()

		case gatewayOpcodeInvalidSession:
			var resumable bool
			sonic.Unmarshal(payload.D, &resumable)
			time.Sleep(time.Second)
			if resumable {
				s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " session invalid (resumable), resuming")
				s.sendResume()
			} else {
				s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " session invalid (non-resumable), identifying")
				s.sessionID = ""
				s.seq = 0
				s.sendIdentify()
			}

		case gatewayOpcodeHello:
			var hello struct {
				HeartbeatInterval float64 `json:"heartbeat_interval"`
			}
			sonic.Unmarshal(payload.D, &hello)
			interval := time.Duration(hello.HeartbeatInterval) * time.Millisecond
			s.logger.Debug("Shard " + strconv.Itoa(s.shardID) + " HELLO received, heartbeat " + interval.String())
			go s.startHeartbeat(interval)

			if s.sessionID != "" && atomic.LoadInt64(&s.seq) > 0 {
				s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " resuming session")
				s.sendResume()
			} else {
				s.logger.Debug("Shard " + strconv.Itoa(s.shardID) + " identifying new session")
				s.sendIdentify()
			}

		case gatewayOpcodeHeartbeatACK:
			s.lastHeartbeatACK.Store(true)
			atomic.StoreInt64(&s.latency, time.Now().UnixMilli())
			s.logger.Debug("Shard " + strconv.Itoa(s.shardID) + " heartbeatACK received")

		case gatewayOpcodeHeartbeat:
			s.sendHeartbeat()
		}
	}
}

// sendIdentify sends an Identify payload to Discord Gateway
//
// This authenticates the shard as a new session and requests events based on intents.
//
// Identify payloads are rate limited via identifyLimiter.
func (s *Shard) sendIdentify() error {
	payload, _ := sonic.Marshal(map[string]any{
		"op": gatewayOpcodeIdentify,
		"d": map[string]any{
			"token": s.token,
			"properties": map[string]string{
				"os":      "linux",
				"browser": LIB_NAME,
				"device":  LIB_NAME,
			},
			"shards":  [2]int{s.shardID, s.totalShards},
			"intents": s.intents,
		},
	})
	s.identifyLimiter.Wait()
	return wsutil.WriteClientMessage(s.conn, ws.OpText, payload)
}

// sendResume sends a Resume payload to Discord Gateway
//
// This attempts to resume a previous session using sessionID and sequence number.
func (s *Shard) sendResume() error {
	payload, _ := sonic.Marshal(map[string]any{
		"op": gatewayOpcodeResume,
		"d": map[string]any{
			"token":      s.token,
			"session_id": s.sessionID,
			"seq":        atomic.LoadInt64(&s.seq),
		},
	})
	return wsutil.WriteClientMessage(s.conn, ws.OpText, payload)
}

// sendHeartbeat sends a Heartbeat payload to Discord Gateway
//
// The payload data is the last sequence number received.
func (s *Shard) sendHeartbeat() error {
	payload, _ := sonic.Marshal(map[string]any{
		"op": gatewayOpcodeHeartbeat,
		"d":  atomic.LoadInt64(&s.seq),
	})
	return wsutil.WriteClientMessage(s.conn, ws.OpText, payload)
}

// startHeartbeat begins sending heartbeats at the given interval
//
// If a heartbeat ACK is not received before the next heartbeat,
// the shard reconnects automatically.
func (s *Shard) startHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if !s.lastHeartbeatACK.Load() {
			s.logger.Error("Shard " + strconv.Itoa(s.shardID) + " heartbeat not ACKed, reconnecting")
			s.reconnect()
			return
		}

		s.lastHeartbeatACK.Store(false)

		start := time.Now()
		if err := s.sendHeartbeat(); err != nil {
			s.logger.Error("Shard " + strconv.Itoa(s.shardID) + " heartbeat error: " + err.Error())
			s.reconnect()
			return
		}

		atomic.StoreInt64(&s.latency, time.Since(start).Milliseconds())
	}
}

// reconnect closes the current connection and attempts to reconnect
//
// Uses exponential backoff on reconnect failures, maxing out at 1 minute.
func (s *Shard) reconnect() {
	if s.conn != nil {
		s.conn.Close()
	}

	backoff := time.Second
	for {
		time.Sleep(backoff)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := s.connect(ctx)
		cancel()

		if err == nil {
			s.logger.Debug("Shard " + strconv.Itoa(s.shardID) + " reconnected")
			return
		}

		s.logger.Error("Shard " + strconv.Itoa(s.shardID) + " reconnect failed, retrying")
		if backoff < 10*time.Second {
			backoff += 2
		}
	}
}

// Latency returns the current heartbeat latency in milliseconds
func (s *Shard) Latency() int64 {
	return atomic.LoadInt64(&s.latency)
}

// Shutdown cleanly closes the shard's websocket connection.
//
// Call this when you want to stop the shard gracefully.
func (s *Shard) Shutdown() error {
	if s.conn != nil {
		s.logger.Info("Shard " + strconv.Itoa(s.shardID) + " shutting down")
		return s.conn.Close()
	}
	s.conn = nil
	s.logger = nil
	s.dispatcher = nil
	return nil
}
