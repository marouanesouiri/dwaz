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
	"context"
	"log"
	"os"
	"strings"
	"time"
)

/*****************************
 *          Client
 *****************************/

// Client manages your Discord connection at a high level, grouping multiple shards together.
//
// It provides:
//   - Central configuration for your bot token, intents, and logger.
//   - REST API access via restApi.
//   - Event dispatching via dispatcher.
//   - Shard management for scalable Gateway connections.
//
// Create a Client using dwaz.New() with desired options, then call Start().
type Client struct {
	ctx             context.Context
	Logger          Logger                    // logger used throughout the client
	workerPool      WorkerPool                // worker pool used to run tasks asynchronously
	identifyLimiter ShardsIdentifyRateLimiter // rate limiter controlling Identify payloads per shard
	token           string                    // bot token (without "Bot " prefix)
	intents         GatewayIntent             // configured Gateway intents
	shards          []*Shard                  // managed Gateway shards
	*restApi                                  // REST API client
	CacheManager                              // CacheManager for caching discord entities
	*dispatcher                               // event dispatcher
}

// clientOption defines a function used to configure Client during creation.
type clientOption func(*Client)

/*****************************
 *       Options
 *****************************/

// WithToken sets the bot token for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithToken("your_bot_token"))
//
// Notes:
//   - Logs fatal and exits if token is empty or obviously invalid (< 50 chars).
//   - Removes "Bot " prefix automatically if provided.
//
// Warning: Never share your bot token publicly.
func WithToken(token string) clientOption {
	if token == "" {
		log.Fatal("WithToken: token must not be empty")
	}
	if len(token) < 50 {
		log.Fatal("WithToken: token invalid")
	}
	if strings.HasPrefix(token, "Bot ") {
		token = strings.Split(token, " ")[1]
	}
	return func(c *Client) {
		c.token = token
	}
}

// WithLogger sets a custom Logger implementation for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithLogger(myLogger))
//
// Logs fatal and exits if logger is nil.
func WithLogger(logger Logger) clientOption {
	if logger == nil {
		log.Fatal("WithLogger: logger must not be nil")
	}
	return func(c *Client) {
		c.Logger = logger
	}
}

// WithWorkerPool sets a custom workerpool implementation for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithWorkerPool(myWorkerPool))
//
// Logs fatal and exits if workerpool is nil.
func WithWorkerPool(workerPool WorkerPool) clientOption {
	if workerPool == nil {
		log.Fatal("WithWorkerPool: workerPool must not be nil")
	}
	return func(c *Client) {
		c.workerPool = workerPool
	}
}

// WithCacheManager sets a custom CacheManager implementation for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithCacheManager(myCacheManager))
//
// Logs fatal and exits if cacheManager is nil.
func WithCacheManager(cacheManager CacheManager) clientOption {
	if cacheManager == nil {
		log.Fatal("WithCacheManager: cacheManager must not be nil")
	}
	return func(c *Client) {
		c.CacheManager = cacheManager
	}
}

// WithShardsIdentifyRateLimiter sets a custom ShardsIdentifyRateLimiter
// implementation for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithShardsIdentifyRateLimiter(myRateLimiter))
//
// Logs fatal and exits if the provided rateLimiter is nil.
func WithShardsIdentifyRateLimiter(rateLimiter ShardsIdentifyRateLimiter) clientOption {
	if rateLimiter == nil {
		log.Fatal("ShardsIdentifyRateLimiter: shardsIdentifyRateLimiter must not be nil")
	}
	return func(c *Client) {
		c.identifyLimiter = rateLimiter
	}
}

// WithIntents sets Gateway intents for the client shards.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithIntents(GatewayIntentGuilds, GatewayIntentMessageContent))
//
// Also supports bitwise OR usage:
//
//	y := dwaz.New(dwaz.WithIntents(GatewayIntentGuilds | GatewayIntentMessageContent))
func WithIntents(intents ...GatewayIntent) clientOption {
	var totalIntents GatewayIntent
	for _, intent := range intents {
		totalIntents |= intent
	}
	return func(c *Client) {
		c.intents = totalIntents
	}
}

/*****************************
 *       Constructor
 *****************************/

// New creates a new Client instance with provided options.
//
// Example:
//
//	y := dwaz.New(
//	    dwaz.WithToken("my_bot_token"),
//	    dwaz.WithIntents(GatewayIntentGuilds, GatewayIntentMessageContent),
//	    dwaz.WithLogger(myLogger),
//	)
//
// Defaults:
//   - Logger: stdout logger at Info level.
//   - Intents: GatewayIntentGuilds | GatewayIntentGuildMessages | GatewayIntentGuildMembers
func New(ctx context.Context, options ...clientOption) *Client {
	if ctx == nil {
		ctx = context.Background()
	}

	client := &Client{
		ctx:    ctx,
		Logger: NewDefaultLogger(os.Stdout, LogLevelInfoLevel),
		intents: GatewayIntentGuilds |
			GatewayIntentGuildMessages |
			GatewayIntentGuildMembers,
	}

	for _, option := range options {
		option(client)
	}

	if client.workerPool == nil {
		client.workerPool = NewDefaultWorkerPool(client.Logger)
	}

	client.restApi = newRestApi(
		newRequester(nil, client.token, client.Logger),
		client.Logger,
	)
	client.CacheManager = NewDefaultCache(
		CacheFlagGuilds | CacheFlagMembers | CacheFlagChannels | CacheFlagRoles | CacheFlagUsers,
	)
	client.dispatcher = newDispatcher(client.Logger, client.workerPool, client.CacheManager)
	return client
}

/*****************************
 *       Start
 *****************************/

// Start initializes and connects all shards for the client.
//
// It performs the following steps:
//  1. Retrieves Gateway information from Discord.
//  2. Creates and connects shards with appropriate rate limiting.
//  3. Starts listening to Gateway events.
//
// The lifetime of the client is controlled by the provided context `ctx`:
//   - If `ctx` is `nil` or `context.Background()`, Start will block forever,
//     running the client until the program exits or Shutdown is called externally.
//   - If `ctx` is cancellable (e.g., created via context.WithCancel or context.WithTimeout),
//     the client will run until the context is cancelled or times out.
//     When the context is done, the client will shutdown gracefully and Start will return.
//
// This design gives you full control over the client's lifecycle.
// For typical usage where you want the bot to run continuously,
// simply pass `nil` as the context (recommended for beginners).
//
// Example usage:
//
//	// Run the client indefinitely (blocks forever)
//	err := client.Start(nil)
//
//	// Run the client with manual cancellation control
//	ctx, cancel := context.WithCancel(context.Background())
//	go func() {
//	    time.Sleep(time.Hour)
//	    cancel() // stops the client after 1 hour
//	}()
//	err := client.Start(ctx)
//
// Returns an error if Gateway information retrieval or shard connection fails.
func (c *Client) Start() error {
	gatewayBotData, err := c.restApi.FetchGatewayBot()
	if err != nil {
		return err
	}

	if c.identifyLimiter == nil {
		c.identifyLimiter = NewDefaultShardsRateLimiter(gatewayBotData.SessionStartLimit.MaxConcurrency, 5*time.Second)
	}

	for i := range gatewayBotData.Shards {
		shard := newShard(
			i, gatewayBotData.Shards, c.token, c.intents,
			c.Logger, c.dispatcher, c.identifyLimiter,
		)
		if err := shard.connect(c.ctx); err != nil {
			return err
		}
		c.shards = append(c.shards, shard)
	}

	<-c.ctx.Done()
	if err := c.ctx.Err(); err != nil {
		c.Logger.WithField("err", err).Error("Client shutdown due to context error")
	}
	c.Shutdown()
	return nil
}

/*****************************
 *       Shutdown
 *****************************/

// Shutdown cleanly shuts down the Client.
//
// It:
//   - Logs shutdown message.
//   - Shuts down the REST API client (closes idle connections).
//   - Shuts down all managed shards.
func (c *Client) Shutdown() {
	c.Logger.Info("Client shutting down")
	c.restApi.Shutdown()
	c.restApi = nil
	c.Logger = nil
	c.workerPool = nil
	for _, shard := range c.shards {
		shard.Shutdown()
	}
	c.shards = nil
}
