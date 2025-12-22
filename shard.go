/************************************************************************************
 *
 * dwaz (Discord Wrapper API for Zwafriya), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package dwaz

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/marouanesouiri/stdx/xlog"
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
	for range r {
		rl.tokens <- struct{}{}
	}
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
 * ShardManager: manages multiple shards
 *************************************/

// ShardManagerConfig configures how shards are managed.
//
// For sharding (multiple shards in one process):
//
//	config := ShardManagerConfig{TotalShards: 4}  // manages shards 0-3
//
// For clustering (specific shards per process):
//
//	// Process 1:
//	config := ShardManagerConfig{TotalShards: 4, ShardIDs: []int{0, 1}}
//	// Process 2:
//	config := ShardManagerConfig{TotalShards: 4, ShardIDs: []int{2, 3}}
//
// IdentifyProperties configures the "properties" field in the Identify payload.
type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type ShardManagerConfig struct {
	TotalShards int
	ShardIDs    []int
	Identify    IdentifyProperties
}

// ShardManager manages the lifecycle of multiple Gateway shards.
//
// It handles shard creation, connection, and shutdown with support for
// both sharding (multiple shards in one process) and clustering
// (distributing specific shards across multiple processes).
type ShardManager struct {
	config          ShardManagerConfig
	shards          []*Shard
	token           string
	intents         GatewayIntent
	useCompression  bool
	logger          xlog.Logger
	dispatcher      *dispatcher
	identifyLimiter ShardsIdentifyRateLimiter
}

// NewShardManager creates a new ShardManager with the given configuration.
func NewShardManager(
	config ShardManagerConfig,
	token string,
	intents GatewayIntent,
	useCompression bool,
	logger xlog.Logger,
	dispatcher *dispatcher,
	identifyLimiter ShardsIdentifyRateLimiter,
) *ShardManager {
	return &ShardManager{
		config:          config,
		token:           token,
		intents:         intents,
		useCompression:  useCompression,
		logger:          logger,
		dispatcher:      dispatcher,
		identifyLimiter: identifyLimiter,
	}
}

// Start connects all configured shards to Discord Gateway.
//
// If ShardIDs are specified in config, only those shards are started.
// Otherwise, all shards [0..TotalShards-1] are started.
//
// The totalShards parameter is the total shard count (from Discord or override).
func (sm *ShardManager) Start(ctx context.Context, totalShards int) error {
	var shardIDs []int
	if len(sm.config.ShardIDs) > 0 {
		shardIDs = sm.config.ShardIDs
	} else {
		shardIDs = make([]int, totalShards)
		for i := range totalShards {
			shardIDs[i] = i
		}
	}

	sm.logger.WithFields(map[string]any{
		"total_shards":   totalShards,
		"managed_shards": shardIDs,
	}).Info("starting shard manager")

	for _, shardID := range shardIDs {
		shard := newShard(
			shardID, totalShards, sm.token, sm.intents,
			sm.logger, sm.dispatcher, sm.identifyLimiter,
			sm.useCompression, sm.config.Identify,
		)
		if err := shard.connect(ctx); err != nil {
			return err
		}
		sm.shards = append(sm.shards, shard)
	}

	return nil
}

// Shutdown gracefully closes all managed shards.
func (sm *ShardManager) Shutdown() {
	sm.logger.Info("shard manager shutting down")
	for _, shard := range sm.shards {
		shard.Shutdown()
	}
	sm.shards = nil
}

// Shards returns the list of managed shards.
func (sm *ShardManager) Shards() []*Shard {
	return sm.shards
}

// ShardCount returns the number of shards currently managed.
func (sm *ShardManager) ShardCount() int {
	return len(sm.shards)
}

/*************************************
 * Shard: a single Gateway connection
 *************************************/

const (
	gatewayVersion = "10"
	gatewayURL     = "wss://gateway.discord.gg/?v=" + gatewayVersion + "&encoding=json"
	gatewayURLZlib = "wss://gateway.discord.gg/?v=" + gatewayVersion + "&encoding=json&compress=zlib-stream"
)

// Shard manages a single WebSocket connection to Discord Gateway,
// including session state, event handling, heartbeats, and reconnects.
type Shard struct {
	shardID     int           // shard number (zero-based)
	totalShards int           // total number of shards in the bot
	token       string        // Discord bot token
	intents     GatewayIntent // Gateway intents bitmask

	logger          xlog.Logger               // logger interface for informational and error messages
	dispatcher      *dispatcher               // event dispatcher for received Gateway events
	identifyLimiter ShardsIdentifyRateLimiter // rate limiter controlling Identify payloads

	conn net.Conn // websocket connection

	seq       int64  // last received sequence number from Gateway
	sessionID string // current session id for resuming
	resumeURL string // Gateway URL to resume session on

	latency           int64         // heartbeat latency in milliseconds
	lastHeartbeatSent int64         // timestamp (unix nano) of last heartbeat sent
	lastHeartbeatACK  atomic.Bool   // true if last heartbeat was acknowledged
	heartbeatStop     chan struct{} // signal to stop heartbeat goroutine

	// Compression support
	useCompression bool
	properties     IdentifyProperties
}

// newShard constructs a new Shard instance with the specified parameters.
//
// shardID and totalShards configure the sharding info,
// token and url set authentication and gateway endpoint,
// intents specify Gateway events to receive,
// logger and dispatcher handle logging and event dispatching,
// limiter enforces Identify rate limits,
// useCompression enables zlib-stream compression,
// properties configures Identify payload.
func newShard(
	shardID, totalShards int, token string, intents GatewayIntent,
	logger xlog.Logger, dispatcher *dispatcher, limiter ShardsIdentifyRateLimiter,
	useCompression bool, properties IdentifyProperties,
) *Shard {
	return &Shard{
		shardID:         shardID,
		totalShards:     totalShards,
		token:           token,
		intents:         intents,
		logger:          logger.WithField("shard_id", shardID),
		dispatcher:      dispatcher,
		identifyLimiter: limiter,
		useCompression:  useCompression,
		properties:      properties,
	}
}

// Connect establishes or resumes a WebSocket connection to Discord Gateway
//
// The shard attempts to connect to the resumeURL if set, otherwise
// to the default gateway url (with or without compression).
//
// It spawns a goroutine to read messages asynchronously.
func (s *Shard) connect(ctx context.Context) error {
	// Stop any existing heartbeat goroutine
	if s.heartbeatStop != nil {
		close(s.heartbeatStop)
	}
	s.heartbeatStop = make(chan struct{})

	if s.conn != nil {
		s.conn.Close()
	}

	connURL := s.resumeURL
	if connURL == "" {
		if s.useCompression {
			connURL = gatewayURLZlib
		} else {
			connURL = gatewayURL
		}
	} else {
		connURL = s.buildResumeURL(connURL)
	}

	dialer := ws.Dialer{}

	conn, _, _, err := dialer.Dial(ctx, connURL)
	if err != nil {
		return err
	}

	s.logger.Info("connected")
	s.conn = conn
	s.lastHeartbeatACK.Store(true)

	atomic.StoreInt64(&s.latency, 0)

	go s.readLoop()
	return nil
}

// buildResumeURL appends the required query params to the resumeURL.
func (s *Shard) buildResumeURL(resumeURL string) string {
	parsed, err := url.Parse(resumeURL)
	if err != nil {
		return resumeURL
	}

	q := parsed.Query()
	if q.Get("v") == "" {
		q.Set("v", gatewayVersion)
	}
	if q.Get("encoding") == "" {
		q.Set("encoding", "json")
	}
	if s.useCompression && q.Get("compress") == "" {
		q.Set("compress", "zlib-stream")
	}
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

// readLoop continuously reads messages from the Gateway WebSocket
//
// It handles Gateway opcodes, dispatches events, and triggers reconnects as needed.
func (s *Shard) readLoop() {
	var (
		decoder *json.Decoder
		z       io.ReadCloser
		err     error
	)

	if s.useCompression {
		gr := &gatewayReader{conn: s.conn, shard: s}
		z, err = zlib.NewReader(gr)
		if err != nil {
			s.logger.WithField("error", err).Error("zlib handshake failed")
			s.reconnect()
			return
		}
		defer z.Close()
		decoder = json.NewDecoder(z)
	}

	defer s.conn.Close()

	for {
		var payload gatewayPayload

		if s.useCompression {
			if err := decoder.Decode(&payload); err != nil {
				s.logger.WithField("error", err).Error("decode/read error")
				s.reconnect()
				return
			}
		} else {
			msg, op, err := wsutil.ReadServerData(s.conn)
			if err != nil {
				s.logger.WithField("error", err).Error("read error")
				s.reconnect()
				return
			}
			if op == ws.OpText {
				if err := json.Unmarshal(msg, &payload); err != nil {
					s.logger.WithField("error", err).Error("unmarshal error")
					continue
				}
			} else if op == ws.OpClose {
				s.reconnect()
				return
			} else {
				continue
			}
		}

		s.handleGatewayPayload(payload)
	}
}

// gatewayReader implements io.Reader to bridge WebSocket frames to a stream.
// It handles buffering binary frames and processing control frames internally.
type gatewayReader struct {
	conn  net.Conn
	shard *Shard
	buf   bytes.Buffer
}

func (gr *gatewayReader) Read(p []byte) (n int, err error) {
	if gr.buf.Len() > 0 {
		return gr.buf.Read(p)
	}

	for {
		msg, op, err := wsutil.ReadServerData(gr.conn)
		if err != nil {
			return 0, err
		}

		switch op {
		case ws.OpBinary:
			gr.buf.Write(msg)
			return gr.buf.Read(p)

		case ws.OpClose:
			return 0, io.EOF

		case ws.OpPing:
			wsutil.WriteClientMessage(gr.conn, ws.OpPong, msg)
			continue

		case ws.OpPong:
			continue

		case ws.OpText:
			continue
		}
	}
}

func (s *Shard) handleGatewayPayload(payload gatewayPayload) {
	if payload.S > 0 {
		atomic.StoreInt64(&s.seq, payload.S)
	}

	s.dispatcher.dispatch(s.shardID, payload.T, payload.D)

	switch payload.Op {
	case gatewayOpcodeDispatch:
		if payload.T == "READY" {
			var ready struct {
				SessionID        string `json:"session_id"`
				ResumeGatewayURL string `json:"resume_gateway_url"`
			}
			json.Unmarshal(payload.D, &ready)
			s.sessionID = ready.SessionID
			s.resumeURL = ready.ResumeGatewayURL
			s.logger.Info("READY received")
		} else if payload.T == "RESUMED" {
			s.logger.Info("RESUMED received")
		}

	case gatewayOpcodeReconnect:
		s.logger.Info("RECONNECT received")
		s.conn.Close()

	case gatewayOpcodeInvalidSession:
		var resumable bool
		json.Unmarshal(payload.D, &resumable)
		time.Sleep(time.Duration(100+s.shardID%500) * time.Millisecond)

		if resumable {
			s.logger.Info("session invalid (resumable), resuming")
			s.sendResume()
		} else {
			s.logger.Info("session invalid (non-resumable), identifying")
			s.sessionID = ""
			s.seq = 0
			s.sendIdentify()
		}

	case gatewayOpcodeHello:
		var hello struct {
			HeartbeatInterval float64 `json:"heartbeat_interval"`
		}
		json.Unmarshal(payload.D, &hello)
		interval := time.Duration(hello.HeartbeatInterval) * time.Millisecond
		s.logger.WithField("heartbeat_interval", interval.String()).Debug("HELLO received")
		go s.startHeartbeat(interval)

		if s.sessionID != "" && atomic.LoadInt64(&s.seq) > 0 {
			s.logger.Info("resuming session")
			s.sendResume()
		} else {
			s.logger.Debug("identifying new session")
			s.sendIdentify()
		}

	case gatewayOpcodeHeartbeatACK:
		s.lastHeartbeatACK.Store(true)
		sent := atomic.LoadInt64(&s.lastHeartbeatSent)
		if sent > 0 {
			rtt := time.Since(time.Unix(0, sent)).Milliseconds()
			atomic.StoreInt64(&s.latency, rtt)
			s.logger.WithField("rtt_ms", rtt).Debug("heartbeatACK")
		}

	case gatewayOpcodeHeartbeat:
		s.sendHeartbeat()
	}
}

// sendIdentify sends an Identify payload to Discord Gateway
//
// This authenticates the shard as a new session and requests events based on intents.
//
// Identify payloads are rate limited via identifyLimiter.
func (s *Shard) sendIdentify() error {
	payload, _ := json.Marshal(map[string]any{
		"op": gatewayOpcodeIdentify,
		"d": map[string]any{
			"token": s.token,
			"properties": map[string]string{
				"os":      s.properties.OS,
				"browser": s.properties.Browser,
				"device":  s.properties.Device,
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
	payload, _ := json.Marshal(map[string]any{
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
	payload, _ := json.Marshal(map[string]any{
		"op": gatewayOpcodeHeartbeat,
		"d":  atomic.LoadInt64(&s.seq),
	})
	return wsutil.WriteClientMessage(s.conn, ws.OpText, payload)
}

// startHeartbeat begins sending heartbeats at the given interval.
//
// Per Discord spec, the first heartbeat has a jitter delay (interval * random 0-1).
// If a heartbeat ACK is not received before the next heartbeat, the shard reconnects.
func (s *Shard) startHeartbeat(interval time.Duration) {
	jitter := time.Duration(rand.Float64() * float64(interval))
	select {
	case <-time.After(jitter):
	case <-s.heartbeatStop:
		return
	}

	if err := s.sendHeartbeat(); err != nil {
		s.logger.WithField("error", err).Error("first heartbeat error")
		return
	}
	s.lastHeartbeatACK.Store(false)
	atomic.StoreInt64(&s.lastHeartbeatSent, time.Now().UnixNano())

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.heartbeatStop:
			return
		case <-ticker.C:
			if s.conn == nil {
				return
			}

			if !s.lastHeartbeatACK.Load() {
				s.logger.Error("heartbeat not ACKed, reconnecting")
				s.conn.Close()
				return
			}

			s.lastHeartbeatACK.Store(false)
			atomic.StoreInt64(&s.lastHeartbeatSent, time.Now().UnixNano())

			if err := s.sendHeartbeat(); err != nil {
				s.logger.WithField("error", err).Error("heartbeat error")
				s.conn.Close()
				return
			}
		}
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
	maxBackoff := 60 * time.Second

	for {
		s.logger.WithField("backoff", backoff.String()).Info("attempting reconnect")
		time.Sleep(backoff)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		err := s.connect(ctx)
		cancel()

		if err == nil {
			s.logger.Debug("reconnected successfully")
			return
		}

		s.logger.WithField("error", err).Error("reconnect failed")
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
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
		s.logger.Info("shutting down")
		return s.conn.Close()
	}
	s.conn = nil
	s.logger = nil
	s.dispatcher = nil
	return nil
}
