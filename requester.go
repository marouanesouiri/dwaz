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
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/***********************
 *   Constants         *
 ***********************/

const (
	apiVersion       = "v10"
	baseApiUrl       = "https://discord.com/api/" + apiVersion
	maxRetries       = 5
	maxRequestAge    = 10 * time.Second
	headerRetryAfter = "Retry-After"
	headerGlobal     = "X-RateLimit-Global"
	headerRemaining  = "X-RateLimit-Remaining"
	headerResetAfter = "X-RateLimit-Reset-After"
	headerBucket     = "X-RateLimit-Bucket"
	headerScope      = "X-RateLimit-Scope"
)

var retryableStatusCodes = map[int]struct{}{
	429: {}, 500: {}, 502: {}, 503: {}, 504: {},
}

/***********************
 *   GlobalRateLimit   *
 ***********************/

// globalRateLimit stores the earliest time global requests can resume.
type globalRateLimit int64

// set updates the global reset time if the new time is later.
func (g *globalRateLimit) set(t time.Time) {
	newVal := t.UnixNano()
	for {
		oldVal := atomic.LoadInt64((*int64)(g))
		if newVal <= oldVal {
			return
		}
		if atomic.CompareAndSwapInt64((*int64)(g), oldVal, newVal) {
			return
		}
	}
}

// get returns the current global reset time.
func (g *globalRateLimit) get() time.Time {
	return time.Unix(0, atomic.LoadInt64((*int64)(g)))
}

/***********************
 *   ratelimitBucket   *
 ***********************/

// ratelimitBucket holds per-route rate limit info but no methods.
type ratelimitBucket struct {
	sync.Mutex
	remaining int
	resetAt   time.Time
}

/***********************
 *   Requester         *
 ***********************/

// requester handles HTTP requests with Discord rate limit compliance.
type requester struct {
	client    *http.Client
	token     string
	buckets   sync.Map // map[bucketRoute]*Bucket
	queues    sync.Map // map[bucketRoute:majorParam]*sync.Mutex
	global    globalRateLimit
	userAgent string
	logger    Logger
}

// newRequester creates a new Requester with the given bot token and logger.
func newRequester(client *http.Client, token string, logger Logger) *requester {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,

				MaxIdleConns:        500,
				MaxIdleConnsPerHost: 100,
				MaxConnsPerHost:     200,

				IdleConnTimeout:       120 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,

				DisableKeepAlives: false,
				ForceAttemptHTTP2: true,
			},
		}
	}

	return &requester{
		client:    client,
		token:     "Bot " + token,
		userAgent: "DiscordBot (yada)",
		logger:    logger,
	}
}

// updateBucket updates bucket state from headers.
func (r *requester) updateBucket(b *ratelimitBucket, h http.Header) {
	b.Lock()
	defer b.Unlock()

	if rem := h.Get(headerRemaining); rem != "" {
		if n, err := strconv.Atoi(rem); err == nil {
			b.remaining = n
		}
	}
	if resetAfter := h.Get(headerResetAfter); resetAfter != "" {
		if dur, err := strconv.ParseFloat(resetAfter, 64); err == nil {
			b.resetAt = time.Now().Add(time.Duration(dur * float64(time.Second)))
		}
	}
}

// do sends an HTTP request with automatic rate limit and retry handling.
func (r *requester) do(method, url string, body []byte, authenticateWithToken bool) (*http.Response, error) {
	routeData := r.generateRouteData(method, url)

	queueKey := routeData.bucketRoute + ":" + routeData.majorParam
	bucketKey := routeData.bucketRoute

	// Get or create per-resource queue mutex
	queue, _ := r.queues.LoadOrStore(queueKey, &sync.Mutex{})
	q := queue.(*sync.Mutex)

	// Get or create bucket for rate limit tracking
	bucket, _ := r.buckets.LoadOrStore(bucketKey, &ratelimitBucket{remaining: 1})
	b := bucket.(*ratelimitBucket)

	q.Lock()
	defer q.Unlock()

	for tries := range maxRetries {
		r.logger.Debug(fmt.Sprintf("Attempt #%d %s %s", tries+1, method, url))

		b.Lock()

		// Wait for bucket reset if exhausted
		if b.remaining == 0 && time.Now().Before(b.resetAt) {
			wait := b.resetAt.Sub(time.Now()) + 100*time.Millisecond
			r.logger.Debug(fmt.Sprintf("Bucket rate limited on route %s: waiting %v before retrying", bucketKey, wait))
			b.Unlock()
			time.Sleep(wait)
			b.Lock()
		}

		// Wait for global rate limit if active
		if now, globalReset := time.Now(), r.global.get(); globalReset.After(now) {
			wait := globalReset.Sub(now) + 100*time.Millisecond
			r.logger.Debug(fmt.Sprintf("Global rate limit active: waiting %v before retrying request %s %s", wait, method, url))
			b.Unlock()
			time.Sleep(wait)
			b.Lock()
		}

		b.Unlock()

		// Build HTTP request
		req, err := http.NewRequest(method, baseApiUrl+url, bytes.NewReader(body))
		if err != nil {
			r.logger.Error(fmt.Sprintf("Failed building request for %s %s: %v", method, url, err))
			return nil, err
		}

		if authenticateWithToken {
			req.Header.Set("Authorization", r.token)
		}
		req.Header.Set("User-Agent", r.userAgent)
		if method == "POST" || method == "PUT" || method == "PATCH" {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Accept", "application/json")

		// Execute request
		resp, err := r.client.Do(req)
		if err != nil {
			r.logger.Warn(fmt.Sprintf("HTTP request error for %s %s: %v", method, url, err))
			time.Sleep(time.Second)
			continue
		}

		// Handle rate limits and retryable errors
		if resp.StatusCode == 429 {
			retry := resp.Header.Get(headerRetryAfter)
			retryAfter := time.Second
			if retry != "" {
				if sec, err := strconv.ParseFloat(retry, 64); err == nil {
					whole, frac := math.Modf(sec)
					retryAfter = time.Duration(whole)*time.Second + time.Duration(frac*1000)*time.Millisecond
				}
			}

			r.logger.Debug(fmt.Sprintf("429 rate limit hit on route %s, retrying after %v", bucketKey, retryAfter))

			r.updateBucket(b, resp.Header)

			if resp.Header.Get(headerGlobal) == "true" || resp.Header.Get(headerScope) == "shared" {
				r.global.set(time.Now().Add(retryAfter))
			}

			resp.Body.Close()
			time.Sleep(retryAfter)
			continue
		}

		if _, retry := retryableStatusCodes[resp.StatusCode]; retry {
			r.logger.Warn(fmt.Sprintf("Retryable status %d for %s %s, retrying...", resp.StatusCode, method, url))
			resp.Body.Close()
			time.Sleep(time.Second)
			continue
		}

		// Update bucket state from response headers
		r.updateBucket(b, resp.Header)

		return resp, nil
	}

	r.logger.Error(fmt.Sprintf("Max retries reached for %s %s", method, url))
	return nil, errors.New("max retries reached")
}

// routeData stores normalized bucket route and major parameter.
type routeData struct {
	bucketRoute string
	majorParam  string
}

// generateRouteData normalizes the route and extracts major parameters for bucketing.
func (r *requester) generateRouteData(method, endpoint string) routeData {
	if strings.HasPrefix(endpoint, "/interactions/") && strings.HasSuffix(endpoint, "/callback") {
		return routeData{
			bucketRoute: method + ":/interactions/:id/:token/callback",
			majorParam:  "global",
		}
	}
	reSnowflake := regexp.MustCompile(`\d{17,19}`)
	majorMatch := reSnowflake.FindString(endpoint)
	baseRoute := reSnowflake.ReplaceAllString(endpoint, ":id")
	baseRoute = regexp.MustCompile(`/reactions/.*`).ReplaceAllString(baseRoute, "/reactions/:reaction")
	baseRoute = regexp.MustCompile(`/webhooks/:id/[^/?]+`).ReplaceAllString(baseRoute, "/webhooks/:id/:token")

	if method == "DELETE" && strings.HasPrefix(baseRoute, "/channels/:id/messages/:id") {
		if messageId, err := strconv.ParseInt(strings.Split(endpoint, "/")[len(strings.Split(endpoint, "/"))-1], 10, 64); err == nil {
			if time.Now().UnixMilli()-Snowflake(messageId).Timestamp().UnixMilli() > 14*24*60*60*1000 {
				baseRoute += "/DELETE_Old_MESSAGE"
			}
		}
	}
	return routeData{
		bucketRoute: method + ":" + baseRoute,
		majorParam:  majorMatch,
	}
}
