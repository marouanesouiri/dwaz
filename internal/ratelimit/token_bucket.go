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

package ratelimit

import (
	"sync/atomic"
	"time"
)

// ErrBucketExhausted is returned by Bucket.Add when no tokens are available
// and the caller must wait until the bucket resets.
var ErrBucketExhausted = errBucketExhausted{}

type errBucketExhausted struct{}

func (errBucketExhausted) Error() string { return "rate limit bucket exhausted" }

// Bucket is a generic rate limiting token bucket interface used for the global rate limit gate.
type Bucket interface {
	// Capacity returns the maximum number of tokens in the bucket.
	Capacity() uint
	// Remaining returns the current number of available tokens.
	Remaining() uint
	// Reset returns the time at which the bucket will next replenish.
	Reset() time.Time
	// Add attempts to consume the specified number of tokens.
	// Returns ErrBucketExhausted if insufficient tokens are available.
	Add(amount uint) error
}

// tokenBucket is a lock-free token bucket implementation using atomic operations.
// It replenishes all tokens at once after the rate period elapses, implementing
// a fixed-window counter suitable for Discord's global 50 req/s limit.
type tokenBucket struct {
	capacity    uint64
	rateNano    int64
	remaining   atomic.Int64
	resetAtNano atomic.Int64
}

// NewTokenBucket creates a new token bucket with the given capacity and replenishment rate.
// The bucket starts fully replenished.
func NewTokenBucket(capacity uint, rate time.Duration) Bucket {
	b := &tokenBucket{
		capacity: uint64(capacity),
		rateNano: rate.Nanoseconds(),
	}
	b.remaining.Store(int64(capacity))
	b.resetAtNano.Store(time.Now().Add(rate).UnixNano())
	return b
}

// Capacity returns the maximum number of tokens in the bucket.
func (b *tokenBucket) Capacity() uint {
	return uint(b.capacity)
}

// Remaining returns the current number of available tokens after replenishing if needed.
func (b *tokenBucket) Remaining() uint {
	b.replenish(time.Now().UnixNano())
	rem := b.remaining.Load()
	if rem < 0 {
		return 0
	}
	return uint(rem)
}

// Reset returns the time at which the bucket will next replenish.
func (b *tokenBucket) Reset() time.Time {
	b.replenish(time.Now().UnixNano())
	return time.Unix(0, b.resetAtNano.Load())
}

// replenish refills the bucket if the current window has elapsed.
// Uses CAS to ensure exactly one goroutine performs the replenishment.
func (b *tokenBucket) replenish(now int64) {
	resetAt := b.resetAtNano.Load()
	if now >= resetAt {
		newReset := now + b.rateNano
		if b.resetAtNano.CompareAndSwap(resetAt, newReset) {
			b.remaining.Store(int64(b.capacity))
		}
	}
}

// Add attempts to consume the specified number of tokens from the bucket.
// Returns nil on success, or ErrBucketExhausted if insufficient tokens are available.
func (b *tokenBucket) Add(amount uint) error {
	now := time.Now().UnixNano()
	b.replenish(now)

	req := int64(amount)
	for {
		cur := b.remaining.Load()
		if cur < req {
			return ErrBucketExhausted
		}
		if b.remaining.CompareAndSwap(cur, cur-req) {
			return nil
		}
	}
}
