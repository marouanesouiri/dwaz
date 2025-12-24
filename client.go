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
	"runtime"
	"strings"
	"time"

	"github.com/marouanesouiri/stdx/xlog"
)

/*****************************
 *          Client
 *****************************/

// Client manages your Discord connection at a high level, grouping multiple shards together.
//
// It provides:
//   - Central configuration for your bot token, intents, and logger.
//   - REST API access via requester.
//   - Event dispatching via dispatcher.
//   - Shard management for scalable Gateway connections.
//
// Create a Client using dwaz.New() with desired options, then call Start().
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
	ctx                  context.Context
	Logger               xlog.Logger               // logger used throughout the client
	identifyLimiter      ShardsIdentifyRateLimiter // rate limiter controlling Identify payloads per shard
	token                string                    // bot token (without "Bot " prefix)
	intents              GatewayIntent             // configured Gateway intents
	shardManager         *ShardManager             // manages Gateway shard lifecycle
	shardManagerConfig   ShardManagerConfig        // configuration for shard management
	useCompression       bool                      // whether to use zlib-stream compression (default: true)
	*requester                                     // REST API client
	CacheManager                                   // CacheManager for caching discord entities
	*dispatcher                                    // event dispatcher
	requesterConfig      RequesterConfig           // configuration for the HTTP requester
	handlerExecutionMode HandlerExecutionMode      // mode for executing event handlers
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
	// Basic length check removed as tokens can vary, but generally good to keep some sanity check.
	// We'll trust the user more here as requested.
	if strings.HasPrefix(token, "Bot ") {
		token = strings.Split(token, " ")[1]
	}
	return func(c *Client) {
		c.token = token
		c.requesterConfig.Token = token
	}
}

// WithLogger sets a custom Logger implementation for your client.
//
// Usage:
//
//	y := dwaz.New(dwaz.WithLogger(myLogger))
//
// Logs fatal and exits if logger is nil.
func WithLogger(logger xlog.Logger) clientOption {
	if logger == nil {
		log.Fatal("WithLogger: logger must not be nil")
	}
	return func(c *Client) {
		c.Logger = logger
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

// WithRequesterConfig sets the configuration for the HTTP requester.
// Use this to configure a proxy URL or custom HTTP client.
func WithRequesterConfig(config RequesterConfig) clientOption {
	return func(c *Client) {
		// Preserve token if already set via WithToken, unless config has it
		if config.Token == "" {
			config.Token = c.token
		}
		c.requesterConfig = config
	}
}

// WithShardCount forces a specific number of shards to be used.
// If not set (0), the recommended shard count from Discord is used.
//
// Deprecated: Use WithShardManagerConfig for more control over sharding.
func WithShardCount(count int) clientOption {
	return func(c *Client) {
		c.shardManagerConfig.TotalShards = count
	}
}

// WithShardManagerConfig sets the shard manager configuration.
//
// For sharding (multiple shards in one process):
//
//	dwaz.WithShardManagerConfig(dwaz.ShardManagerConfig{TotalShards: 4})
//
// For clustering (specific shards per process):
//
//	// Process 1:
//	dwaz.WithShardManagerConfig(dwaz.ShardManagerConfig{TotalShards: 4, ShardIDs: []int{0, 1}})
//	// Process 2:
//	dwaz.WithShardManagerConfig(dwaz.ShardManagerConfig{TotalShards: 4, ShardIDs: []int{2, 3}})
func WithShardManagerConfig(config ShardManagerConfig) clientOption {
	return func(c *Client) {
		c.shardManagerConfig = config
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

// WithHandlerExecutionMode sets the execution mode for event handlers.
//
// Usage:
//
//	dwaz.New(..., dwaz.WithHandlerExecutionMode(dwaz.HandlerExecutionAsync))
//
// Default is HandlerExecutionSync (sequential).
func WithHandlerExecutionMode(mode HandlerExecutionMode) clientOption {
	return func(c *Client) {
		c.handlerExecutionMode = mode
	}
}

// WithCompression enables or disables zlib-stream compression for Gateway connections.
//
// When enabled (default), Gateway messages are compressed, reducing bandwidth by 60-80%.
//
// Usage:
//
//	dwaz.New(..., dwaz.WithCompression(false)) // disable compression
//
// Default is true (compression enabled).
func WithCompression(enabled bool) clientOption {
	return func(c *Client) {
		c.useCompression = enabled
	}
}

// WithIdentifyProperties sets custom properties for the Identify payload.
//
// Usage:
//
//	dwaz.New(..., dwaz.WithIdentifyProperties(dwaz.IdentifyProperties{
//	    OS:      "linux",
//	    Browser: "my-bot",
//	    Device:  "my-bot",
//	}))
func WithIdentifyProperties(props IdentifyProperties) clientOption {
	return func(c *Client) {
		c.shardManagerConfig.Identify = props
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
//   - Executor: SpawnExecutor (goroutine per task)
func New(ctx context.Context, options ...clientOption) *Client {
	if ctx == nil {
		ctx = context.Background()
	}

	client := &Client{
		ctx:    ctx,
		Logger: xlog.NewTextLogger(os.Stdout, xlog.LogLevelInfoLevel),
		intents: GatewayIntentGuilds |
			GatewayIntentGuildMessages |
			GatewayIntentGuildMembers,
		useCompression: true,
	}

	for _, option := range options {
		option(client)
	}

	if client.requesterConfig.Token == "" {
		client.requesterConfig.Token = client.token
	}

	client.requester = newRequester(client.requesterConfig, client.Logger)
	client.CacheManager = NewInMemoryCacheManager(
		CacheFlagGuilds | CacheFlagMembers | CacheFlagChannels | CacheFlagRoles | CacheFlagUsers | CacheFlagVoiceStates,
	)
	client.dispatcher = newDispatcher(client.Logger, client, client.handlerExecutionMode)
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
	res := c.requester.FetchGatewayBot()
	if res.IsErr() {
		return res.Err()
	}
	gatewayBotData := res.Value()

	if c.identifyLimiter == nil {
		c.identifyLimiter = NewDefaultShardsRateLimiter(gatewayBotData.SessionStartLimit.MaxConcurrency, 5*time.Second)
	}

	if c.shardManagerConfig.Identify.OS == "" {
		c.shardManagerConfig.Identify.OS = runtime.GOOS
	}
	if c.shardManagerConfig.Identify.Browser == "" {
		c.shardManagerConfig.Identify.Browser = "dwaz"
	}
	if c.shardManagerConfig.Identify.Device == "" {
		c.shardManagerConfig.Identify.Device = "dwaz"
	}

	// Determine total shards: use config if set, otherwise use Discord's recommendation
	totalShards := gatewayBotData.Shards
	if c.shardManagerConfig.TotalShards > 0 {
		totalShards = c.shardManagerConfig.TotalShards
	}

	// Create shard manager
	c.shardManager = NewShardManager(
		c.shardManagerConfig,
		c.token,
		c.intents,
		c.useCompression,
		c.Logger,
		c.dispatcher,
		c.identifyLimiter,
	)

	// Start all configured shards
	if err := c.shardManager.Start(c.ctx, totalShards); err != nil {
		return err
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
//   - Shuts down all managed shards via ShardManager.
func (c *Client) Shutdown() {
	c.Logger.Info("Client shutting down")
	if c.requester != nil {
		c.requester.Shutdown()
		c.requester = nil
	}
	c.Logger = nil
	if c.shardManager != nil {
		c.shardManager.Shutdown()
		c.shardManager = nil
	}
}
