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
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type mockRoundTripper struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

func newMockResponse(status int, body string, headers map[string]string) *http.Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     h,
	}
}

func newTestRequester(mockFn func(*http.Request) (*http.Response, error)) *requester {
	mockClient := &http.Client{
		Transport: &mockRoundTripper{fn: mockFn},
		Timeout:   5 * time.Second,
	}
	logger := NewDefaultLogger(nil, LogLevel_DebugLevel)
	return newRequester(mockClient, "testtoken", logger)
}

func TestRequester_Do_Success(t *testing.T) {
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		return newMockResponse(200, `{"ok":true}`, map[string]string{
			"X-RateLimit-Remaining":   "10",
			"X-RateLimit-Reset-After": "1",
		}), nil
	})

	resp, err := r.do("GET", "/channels/123/messages", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func TestRequester_Do_RateLimitRetry(t *testing.T) {
	attempts := int32(0)
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		if atomic.AddInt32(&attempts, 1) <= 2 {
			return newMockResponse(429, `{"message":"rate limited"}`, map[string]string{
				"Retry-After":             "0.1",
				"X-RateLimit-Remaining":   "0",
				"X-RateLimit-Reset-After": "0.1",
			}), nil
		}
		return newMockResponse(200, `{"ok":true}`, map[string]string{
			"X-RateLimit-Remaining":   "5",
			"X-RateLimit-Reset-After": "1",
		}), nil
	})

	resp, err := r.do("GET", "/channels/123/messages", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
	if attempts < 3 {
		t.Fatalf("expected at least 3 attempts, got %d", attempts)
	}
}

func TestRequester_Do_GlobalRateLimit(t *testing.T) {
	attempts := int32(0)
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		if atomic.AddInt32(&attempts, 1) == 1 {
			return newMockResponse(429, `{"message":"global rate limit"}`, map[string]string{
				"Retry-After":             "0.1",
				"X-RateLimit-Global":      "true",
				"X-RateLimit-Remaining":   "0",
				"X-RateLimit-Reset-After": "0.1",
			}), nil
		}
		return newMockResponse(200, `{"ok":true}`, nil), nil
	})

	resp, err := r.do("GET", "/channels/123/messages", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}

func TestRequester_Do_RetryableStatusCodes(t *testing.T) {
	attempts := int32(0)
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 3 {
			return newMockResponse(503, "Service Unavailable", nil), nil
		}
		return newMockResponse(200, `{"ok":true}`, nil), nil
	})

	resp, err := r.do("GET", "/channels/123/messages", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
	if attempts != 4 {
		t.Fatalf("expected 4 attempts, got %d", attempts)
	}
}

func TestRequester_Do_MaxRetriesExceeded(t *testing.T) {
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		return newMockResponse(503, "Service Unavailable", nil), nil
	})

	_, err := r.do("GET", "/channels/123/messages", nil, true)
	if err == nil || !strings.Contains(err.Error(), "max retries") {
		t.Fatalf("expected max retries error, got %v", err)
	}
}

func TestRequester_ConcurrencyStress(t *testing.T) {
	var total int64
	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		return newMockResponse(200, `{"ok":true}`, map[string]string{
			"X-RateLimit-Remaining":   "10",
			"X-RateLimit-Reset-After": "1",
		}), nil
	})

	const concurrency = 50
	const requestsPerGoroutine = 10
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for range concurrency {
		func() {
			defer wg.Done()
			for range requestsPerGoroutine {
				resp, err := r.do("GET", "/channels/123/messages", nil, true)
				if err != nil {
					t.Errorf("request error: %v", err)
					return
				}
				resp.Body.Close()
				atomic.AddInt64(&total, 1)
			}
		}()
	}
	wg.Wait()

	if total != concurrency*requestsPerGoroutine {
		t.Fatalf("expected %d successful requests, got %d", concurrency*requestsPerGoroutine, total)
	}
}

func TestRequester_ConcurrentRateLimitEnforcement(t *testing.T) {
	var attempts int32
	var mu sync.Mutex
	rateLimitedUntil := time.Time{}

	r := newTestRequester(func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		// Simulate a global rate limit window until rateLimitedUntil
		if now.Before(rateLimitedUntil) {
			return newMockResponse(429, "", map[string]string{
				"Retry-After":             fmt.Sprintf("%.1f", rateLimitedUntil.Sub(now).Seconds()),
				"X-RateLimit-Global":      "true",
				"X-RateLimit-Remaining":   "0",
				"X-RateLimit-Reset-After": fmt.Sprintf("%.1f", rateLimitedUntil.Sub(now).Seconds()),
			}), nil
		}

		n := atomic.AddInt32(&attempts, 1)
		// Trigger a global rate limit every 20 requests lasting 300ms
		if n%20 == 0 {
			rateLimitedUntil = now.Add(300 * time.Millisecond)
			return newMockResponse(429, "", map[string]string{
				"Retry-After":             "0.3",
				"X-RateLimit-Global":      "true",
				"X-RateLimit-Remaining":   "0",
				"X-RateLimit-Reset-After": "0.3",
			}), nil
		}

		return newMockResponse(200, `{"ok":true}`, map[string]string{
			"X-RateLimit-Remaining":   "10",
			"X-RateLimit-Reset-After": "1",
		}), nil
	})

	const concurrency = 10
	const requestsPerGoroutine = 30
	totalRequests := concurrency * requestsPerGoroutine

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for range concurrency {
		go func() {
			defer wg.Done()
			for range requestsPerGoroutine {
				resp, err := r.do("GET", "/channels/123/messages", nil, true)
				if err != nil {
					t.Errorf("request error: %v", err)
					return
				}
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	minExpected := time.Duration(totalRequests/20) * 300 * time.Millisecond
	if elapsed < minExpected {
		t.Errorf("expected total duration at least %v due to rate limits, got %v", minExpected, elapsed)
	}
}

func TestGenerateBucketKey(t *testing.T) {
	r := &requester{}

	// Old message snowflake (more than 14 days)
	oldMessageID := "1363358614089371648"
	// New message snowflake
	newMessageID := "1396987230249029793"

	cases := []struct {
		method   string
		endpoint string
	}{
		// Old message delete
		{"DELETE", "/channels/123456789012345678/messages/" + oldMessageID},

		// New message delete
		{"DELETE", "/channels/123456789012345678/messages/" + newMessageID},

		// Interaction callback
		{"POST", "/interactions/987654321098765432/abcdef/callback"},

		// Webhook with token
		{"POST", "/webhooks/123456789012345678/abcdef1234567890"},

		// Reaction add
		{"PUT", "/channels/123456789012345678/messages/234567890123456789/reactions/XXXXXXX/@me"},

		// Normal GET channel message
		{"GET", "/channels/123456789012345678/messages/234567890123456789"},

		// Other route with IDs
		{"PATCH", "/guilds/987654321098765432/members/123456789012345678"},

		// Route without IDs
		{"GET", "/gateway/bot"},
		{"GET", "/users/@me"},
	}

	for _, c := range cases {
		key := r.generateBucketKey(c.method, c.endpoint)
		fmt.Printf("Method: %s, Endpoint: %s\n => BucketKey: %s\n\n", c.method, c.endpoint, key)
	}
}
