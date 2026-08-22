/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD-3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/wumpkit-go/wumpkit"
	goredis "github.com/redis/go-redis/v9"
)

// RedisRateLimitStore implements wumpkit.RateLimitStore using Redis for multi-cluster synchronization.
type RedisRateLimitStore struct {
	client goredis.UniversalClient
	prefix string
}

// NewRedisRateLimitStore creates a new RedisRateLimitStore with an optional key prefix (defaults to "wumpkit:ratelimit:").
func NewRedisRateLimitStore(client goredis.UniversalClient, prefix ...string) *RedisRateLimitStore {
	p := "wumpkit:ratelimit:"
	if len(prefix) > 0 && prefix[0] != "" {
		p = prefix[0]
		if !strings.HasSuffix(p, ":") {
			p += ":"
		}
	}
	return &RedisRateLimitStore{
		client: client,
		prefix: p,
	}
}

func (s *RedisRateLimitStore) bucketKey(key string) string {
	return s.prefix + "bucket:" + key
}

func (s *RedisRateLimitStore) hashKey(key string) string {
	return s.prefix + "hash:" + key
}

func (s *RedisRateLimitStore) globalKey() string {
	return s.prefix + "global_lock"
}

// Reserve atomically checks remaining quota and consumes 1 token in Redis.
func (s *RedisRateLimitStore) Reserve(bucketKey string) time.Duration {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := s.client.Get(ctx, s.bucketKey(bucketKey)).Result()
	if err != nil {
		return 0
	}

	var state wumpkit.RateLimitBucketState
	if err := json.Unmarshal([]byte(val), &state); err != nil {
		return 0
	}

	now := time.Now()
	if now.After(state.ResetAt) {
		state.Remaining = state.Limit
	} else if state.Remaining <= 0 {
		return max(0, state.ResetAt.Sub(now))
	}

	if state.Remaining > 0 {
		state.Remaining--
		if data, err := json.Marshal(state); err == nil {
			ttl := max(time.Minute, time.Until(state.ResetAt)+time.Minute)
			_ = s.client.Set(ctx, s.bucketKey(bucketKey), data, ttl).Err()
		}
	}

	return 0
}

func (s *RedisRateLimitStore) GetBucket(bucketKey string) (wumpkit.RateLimitBucketState, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var zero wumpkit.RateLimitBucketState
	val, err := s.client.Get(ctx, s.bucketKey(bucketKey)).Result()
	if err != nil {
		return zero, false
	}

	var state wumpkit.RateLimitBucketState
	if err := json.Unmarshal([]byte(val), &state); err != nil {
		return zero, false
	}
	return state, true
}

func (s *RedisRateLimitStore) SetBucket(bucketKey string, state wumpkit.RateLimitBucketState) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := json.Marshal(state)
	if err != nil {
		return
	}

	ttl := time.Until(state.ResetAt) + 1*time.Minute
	if ttl <= 0 {
		ttl = 1 * time.Minute
	}

	_ = s.client.Set(ctx, s.bucketKey(bucketKey), data, ttl).Err()
}

func (s *RedisRateLimitStore) GetGlobalLock() time.Time {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := s.client.Get(ctx, s.globalKey()).Result()
	if err != nil {
		return time.Time{}
	}

	unixNano, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(0, unixNano)
}

func (s *RedisRateLimitStore) SetGlobalLock(until time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ttl := time.Until(until)
	if ttl <= 0 {
		return
	}

	_ = s.client.Set(ctx, s.globalKey(), strconv.FormatInt(until.UnixNano(), 10), ttl).Err()
}

func (s *RedisRateLimitStore) GetBucketHash(optimisticKey string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := s.client.Get(ctx, s.hashKey(optimisticKey)).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (s *RedisRateLimitStore) SetBucketHash(optimisticKey, bucketHash string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = s.client.Set(ctx, s.hashKey(optimisticKey), bucketHash, 24*time.Hour).Err()
}
