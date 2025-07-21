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
	"errors"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
)

/***********************
 *	  callWithData	   *
 ***********************/

// callWithData represents a REST API request returning typed decoded data.
type callWithData[T any] struct {
	requester     *requester
	logger        Logger
	method        string
	endpoint      string
	body          []byte
	authWithToken bool
	parse         func([]byte) (*T, error)
}

// Wait executes the request synchronously and parses the response.
func (c *callWithData[T]) Wait() (*T, error) {
	c.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.requester.do(c.method, c.endpoint, c.body, c.authWithToken)
	if err != nil {
		c.logger.Error(
			"Request failed for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.logger.Error("Request failed for endpoint " + c.method + c.endpoint + ": Invalid Token")
		return nil, errors.New("Invalid Token !!")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error(
			"Failed reading response body for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return nil, err
	}

	data, err := c.parse(bodyBytes)
	if err != nil {
		c.logger.Error(
			"Failed parsing response for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return nil, err
	}

	c.logger.Debug("Successfully called endpoint: " + c.method + c.endpoint)
	return data, nil
}

// Submit runs the request asynchronously and calls the provided callback.
func (c *callWithData[T]) Submit(callback func(*T, error)) {
	// TODO: run callback using a worker pool
	go func() { callback(c.Wait()) }()
}

/***********************
 *	 callWithNoData	   *
 ***********************/

// callWithNoData represents a REST API request with no data parsing.
type callWithNoData struct {
	requester     *requester
	logger        Logger
	method        string
	endpoint      string
	body          []byte
	authWithToken bool
}

// Wait executes the request synchronously.
func (c *callWithNoData) Wait() error {
	c.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.requester.do(c.method, c.endpoint, c.body, c.authWithToken)
	if err != nil {
		c.logger.Error(
			"Request failed for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return err
	}
	res.Body.Close()

	c.logger.Debug("Successfully called endpoint: " + c.method + c.endpoint)
	return nil
}

// Submit runs the request asynchronously and calls the provided callback.
func (c *callWithNoData) Submit(callback func(error)) {
	// TODO: run callback using a worker pool
	go func() { callback(c.Wait()) }()
}

/***********************
 *       RestAPI       *
 ***********************/

// restApi provides methods for Discord REST API endpoints.
type restApi struct {
	requester *requester
	logger    Logger
}

// newRestApi creates a new restApi instance with optional custom requester and logger.
func newRestApi(requester *requester, token string, logger Logger) *restApi {
	if logger == nil {
		logger = NewDefaultLogger(nil, LogLevel_DebugLevel)
	}
	if requester == nil {
		requester = newRequester(nil, token, logger)
	}

	return &restApi{
		requester: requester,
		logger:    logger,
	}
}

// Shutdown call requester.Shutdown() which gracefully closes the underlying HTTP client's idle connections.
//
// It should be called before exiting your application to ensure
// that any idle connections in the HTTP transport are closed cleanly,
// preventing resource leaks and keeping a clean shutdown process.
func (r *restApi) Shutdown() {
	r.logger.Info("RestAPI shutting down")
	r.requester.Shutdown()
	r.logger = nil
	r.requester = nil
}

/***********************
 *   Gateway Endpoint  *
 ***********************/

// getGateway constructs a callWithData for the GET /gateway endpoint.
//
// This endpoint returns the WebSocket URL used to connect to the Discord Gateway.
// No authentication token is required.
//
// The returned callWithData value can be executed by calling either:
//   - Wait(): runs the request synchronously on the current goroutine (preferred).
//   - Submit(callback): runs the request asynchronously in a new goroutine, invoking the callback upon completion.
//
// Returns:
//
//	callWithData[gateway] — a prepared request object that can be executed to fetch gateway information.
func (r *restApi) getGateway() callWithData[gateway] {
	return callWithData[gateway]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/gateway",
		authWithToken: false,
		parse: func(b []byte) (*gateway, error) {
			obj := gateway{}
			err := sonic.Unmarshal(b, obj)
			return &obj, err
		},
	}
}

// GetGatewayBot constructs a callWithData for the GET /gateway/bot endpoint.
//
// This endpoint returns information about the current bot's gateway, including recommended shard count and session limits.
// Requires authentication via bot token.
//
// The returned callWithData value can be executed by calling either:
//   - Wait(): runs the request synchronously on the current goroutine (preferred).
//   - Submit(callback): runs the request asynchronously in a new goroutine, invoking the callback upon completion.
//
// Returns:
//
//	callWithData[GatewayBot] — a prepared request object for fetching the gateway bot information.
func (r *restApi) GetGatewayBot() callWithData[GatewayBot] {
	return callWithData[GatewayBot]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/gateway/bot",
		authWithToken: true,
		parse: func(b []byte) (*GatewayBot, error) {
			obj := GatewayBot{}
			err := sonic.Unmarshal(b, obj)
			return &obj, err
		},
	}
}

/***********************
 *    User Endpoints    *
 ***********************/

// GetSelfUser retrieves the current user's data.
//
// Usage example for beginners:
//
//	user, err := restApi.GetSelfUser().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Your username:", user.Username)
//
// Callers can use:
//   - Wait() to run synchronously (recommended for simplicity)
//   - Submit(callback) to run asynchronously with a callback
//
// Detailed info:
//
//	Endpoint: GET /users/@me
//	Requires OAuth2 identify scope; optionally includes email if email scope granted.
//
// Returns: callWithData[User] — prepared request object to fetch current user data.
func (r *restApi) GetSelfUser() callWithData[User] {
	return callWithData[User]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/users/@me",
		authWithToken: true,
		parse: func(b []byte) (*User, error) {
			obj := User{restApi: r}
			err := sonic.Unmarshal(b, obj)
			return &obj, err
		},
	}
}

// GetUser retrieves a user by their Snowflake ID.
//
// Usage example:
//
//	userID := 123456789012345678
//	user, err := restApi.GetUser(userID).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User username:", user.Username)
//
// Use Wait() for blocking call or Submit() for async callback.
//
// Detailed info:
//
//	Endpoint: GET /users/{userID}
//	Requires authentication token.
//
// Parameters:
//
//	userID — Snowflake ID of the user.
//
// Returns: callWithData[User] — prepared request object to fetch user data.
func (r *restApi) GetUser(userID Snowflake) callWithData[User] {
	return callWithData[User]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/users/" + userID.String(),
		authWithToken: true,
		parse: func(b []byte) (*User, error) {
			obj := User{restApi: r}
			err := sonic.Unmarshal(b, obj)
			return &obj, err
		},
	}
}
