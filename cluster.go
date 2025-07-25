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
	"log"
	"os"
	"strings"
	"time"
)

/*****************************
 *          Cluster
 *****************************/

// Cluster manages your Discord connection at a high level, grouping multiple shards together.
//
// It provides:
//   - Central configuration for your bot token, intents, and logger.
//   - REST API access via restApi.
//   - Event dispatching via dispatcher.
//   - Shard management for scalable Gateway connections.
//
// Create a Cluster using yada.New() with desired options, then call Start().
type Cluster struct {
	Logger      Logger        // logger used throughout the cluster
	token       string        // bot token (without "Bot " prefix)
	intents     GatewayIntent // configured Gateway intents
	shards      []*Shard      // managed Gateway shards
	*restApi                  // REST API client
	*dispatcher               // event dispatcher
}

// clusterOption defines a function used to configure Cluster during creation.
type clusterOption func(*Cluster)

/*****************************
 *       Options
 *****************************/

// WithToken sets the bot token for your cluster.
//
// Usage:
//
//	y := yada.New(yada.WithToken("your_bot_token"))
//
// Notes:
//   - Logs fatal and exits if token is empty or obviously invalid (< 50 chars).
//   - Removes "Bot " prefix automatically if provided.
//
// Warning: Never share your bot token publicly.
func WithToken(token string) clusterOption {
	if token == "" {
		log.Fatal("WithToken: token must not be empty")
	}
	if len(token) < 50 {
		log.Fatal("WithToken: token invalid")
	}
	if strings.HasPrefix(token, "Bot ") {
		token = strings.Split(token, " ")[1]
	}
	return func(c *Cluster) {
		c.token = token
	}
}

// WithLogger sets a custom Logger implementation for your cluster.
//
// Usage:
//
//	y := yada.New(yada.WithLogger(myLogger))
//
// Logs fatal and exits if logger is nil.
func WithLogger(logger Logger) clusterOption {
	if logger == nil {
		log.Fatal("WithLogger: logger must not be nil")
	}
	return func(c *Cluster) {
		c.Logger = logger
	}
}

// WithIntents sets Gateway intents for the cluster shards.
//
// Usage:
//
//	y := yada.New(yada.WithIntents(GatewayIntentGuilds, GatewayIntentMessageContent))
//
// Also supports bitwise OR usage:
//
//	y := yada.New(yada.WithIntents(GatewayIntentGuilds | GatewayIntentMessageContent))
func WithIntents(intents ...GatewayIntent) clusterOption {
	var totalIntents GatewayIntent
	for _, intent := range intents {
		totalIntents |= intent
	}
	return func(c *Cluster) {
		c.intents = totalIntents
	}
}

/*****************************
 *       Constructor
 *****************************/

// New creates a new Cluster instance with provided options.
//
// Example:
//
//	y := yada.New(
//	    yada.WithToken("my_bot_token"),
//	    yada.WithIntents(GatewayIntentGuilds, GatewayIntentMessageContent),
//	    yada.WithLogger(myLogger),
//	)
//
// Defaults:
//   - Logger: stdout logger at Info level.
//   - Intents: GatewayIntentGuilds | GatewayIntentGuildMessages
func New(options ...clusterOption) *Cluster {
	cluster := &Cluster{
		Logger:  NewDefaultLogger(os.Stdout, LogLevel_InfoLevel),
		intents: GatewayIntentGuilds | GatewayIntentGuildMessages,
	}

	for _, option := range options {
		option(cluster)
	}

	cluster.restApi = newRestApi(
		newRequester(nil, cluster.token, cluster.Logger),
		cluster.Logger,
	)
	cluster.dispatcher = newDispatcher(cluster.Logger)
	return cluster
}

/*****************************
 *       Start
 *****************************/

// Start initializes and connects all shards for the cluster.
//
// It performs the following steps:
//  1. Retrieves Gateway information from Discord.
//  2. Creates and connects shards with appropriate rate limiting.
//  3. Starts listening to Gateway events.
//
// The lifetime of the cluster is controlled by the provided context `ctx`:
//   - If `ctx` is `nil` or `context.Background()`, Start will block forever,
//     running the cluster until the program exits or Shutdown is called externally.
//   - If `ctx` is cancellable (e.g., created via context.WithCancel or context.WithTimeout),
//     the cluster will run until the context is cancelled or times out.
//     When the context is done, the cluster will shutdown gracefully and Start will return.
//
// This design gives you full control over the cluster's lifecycle.
// For typical usage where you want the bot to run continuously,
// simply pass `nil` as the context (recommended for beginners).
//
// Example usage:
//
//	// Run the cluster indefinitely (blocks forever)
//	err := cluster.Start(nil)
//
//	// Run the cluster with manual cancellation control
//	ctx, cancel := context.WithCancel(context.Background())
//	go func() {
//	    time.Sleep(time.Hour)
//	    cancel() // stops the cluster after 1 hour
//	}()
//	err := cluster.Start(ctx)
//
// Returns an error if Gateway information retrieval or shard connection fails.
func (c *Cluster) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	gatewayBotData, err := c.restApi.GetGatewayBot().Wait()
	if err != nil {
		return err
	}

	shardsLimiter := newDefaultShardsRateLimiter(
		gatewayBotData.SessionStartLimit.MaxConcurrency,
		5*time.Second,
	)

	for i := range gatewayBotData.Shards {
		shard := newShard(
			i, gatewayBotData.Shards, c.token, gatewayBotData.URL, c.intents,
			c.Logger, c.dispatcher, shardsLimiter,
		)
		if err := shard.connect(ctx); err != nil {
			return err
		}
		c.shards = append(c.shards, shard)
	}

	<-ctx.Done()
	c.Shutdown()
	return nil
}

/*****************************
 *       Shutdown
 *****************************/

// Shutdown cleanly shuts down the Cluster.
//
// It:
//   - Logs shutdown message.
//   - Shuts down the REST API client (closes idle connections).
//   - Shuts down all managed shards.
func (c *Cluster) Shutdown() {
	c.Logger.Info("Cluster shutting down")
	c.restApi.Shutdown()
	c.restApi = nil
	c.Logger = nil
	for _, shard := range c.shards {
		shard.Shutdown()
	}
	c.shards = nil
}
