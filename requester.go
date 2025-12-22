package dwaz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/marouanesouiri/stdx/result"
	"github.com/marouanesouiri/stdx/xlog"
)

/***********************
 *   Constants         *
 ***********************/

const (
	apiVersion       = "v10"
	defaultBaseURL   = "https://discord.com/api/" + apiVersion
	defaultUserAgent = "DiscordBot (dwaz, v0.0.1)"
	headerReason     = "X-Audit-Log-Reason"
)

/***********************
 *   Configuration     *
 ***********************/

// RequesterConfig holds configuration for the HTTP requester.
type RequesterConfig struct {
	// BaseURL is the base URL for API requests. Defaults to Discord API.
	// Change this to use a reverse proxy like nirn-proxy.
	BaseURL string

	// APIVersion is the Discord API version. Defaults to "v10".
	APIVersion string

	// UserAgent is the User-Agent header value.
	UserAgent string

	// MaxRetries is the maximum number of retries for failed requests.
	MaxRetries int

	// Token is the Bot token.
	Token string

	// HTTPClient is a custom HTTP client.
	HTTPClient *http.Client
}

func DefaultRequesterConfig() RequesterConfig {
	return RequesterConfig{
		BaseURL:    defaultBaseURL,
		APIVersion: apiVersion,
		UserAgent:  defaultUserAgent,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				MaxIdleConns:          500,
				MaxIdleConnsPerHost:   100,
				MaxConnsPerHost:       200,
				IdleConnTimeout:       120 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ForceAttemptHTTP2:     true,
			},
		},
	}
}

/***********************
 *   Requester         *
 ***********************/

// requester handles HTTP requests with basic retry logic, intended to be used with a proxy.
type requester struct {
	config RequesterConfig
	logger xlog.Logger
}

// newRequester creates a new Requester with the given config.
func newRequester(config RequesterConfig, logger xlog.Logger) *requester {
	return &requester{
		config: config,
		logger: logger,
	}
}

// Shutdown gracefully closes the underlying HTTP client's idle connections.
func (r *requester) Shutdown() {
	if r.config.HTTPClient != nil {
		if tr, ok := r.config.HTTPClient.Transport.(interface{ CloseIdleConnections() }); ok {
			tr.CloseIdleConnections()
		}
	}
}

func isRetryableStatus(code int) bool {
	switch code {
	case 429, 502, 503, 504:
		return true
	default:
		return false
	}
}

// Request represents an HTTP request to be executed by the requester.
type Request struct {
	// Body is the raw JSON request body or other data sent to the API.
	Body []byte
	// Method is the HTTP verb (GET, POST, etc.) to use for the request.
	Method string
	// URL is the endpoint path (e.g., "/guilds/123") relative to the base URL.
	URL string
	// Reason is the value for the "X-Audit-Log-Reason" header, useful for audit logs.
	Reason string
	// NoAuth indicates whether to skip token-based authentication for this request.
	// If false (default), the "Authorization" header will be automatically set using the bot token.
	NoAuth bool
}

// DoRequest executes a request with a raw body and returns the response body as a stream.
// The caller is responsible for closing the body.
func (r *requester) DoRequest(req Request) result.Result[io.ReadCloser] {
	fullURL := r.config.BaseURL + req.URL

	for i := 0; i <= r.config.MaxRetries; i++ {
		if i > 0 {
			r.logger.WithFields(map[string]any{
				"method":  req.Method,
				"url":     fullURL,
				"attempt": i + 1,
			}).Debug("retrying request")
			time.Sleep(100 * time.Millisecond)
		}

		httpRequest, err := http.NewRequest(req.Method, fullURL, bytes.NewReader(req.Body))
		if err != nil {
			return result.Err[io.ReadCloser](fmt.Errorf("failed creating request: %v", err))
		}

		if !req.NoAuth {
			httpRequest.Header.Set("Authorization", r.config.Token)
		}
		httpRequest.Header.Set("User-Agent", r.config.UserAgent)
		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Accept", "application/json")

		if req.Reason != "" {
			httpRequest.Header.Set(headerReason, req.Reason)
		}

		resp, err := r.config.HTTPClient.Do(httpRequest)
		if err != nil {
			r.logger.WithFields(map[string]any{
				"method":   req.Method,
				"endpoint": req.URL,
				"error":    err,
			}).Warn("request network error")
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return result.Ok(resp.Body)
		}

		if isRetryableStatus(resp.StatusCode) {
			resp.Body.Close()
			r.logger.WithFields(map[string]any{
				"method":   req.Method,
				"endpoint": req.URL,
				"status":   resp.StatusCode,
			}).Warn("request failed with retryable status")
			continue
		}

		resp.Body.Close()
		return result.Err[io.ReadCloser](fmt.Errorf("request failed with status %d", resp.StatusCode))
	}

	return result.Err[io.ReadCloser](fmt.Errorf("max retries reached, endpoint %s", req.URL))
}
