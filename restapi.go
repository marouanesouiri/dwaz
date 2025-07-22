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
//
// Parameters:
//   - callback — function to invoke with (*T, error) once the request completes.
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
//
// Returns: error — if the request failed.
func (c *callWithNoData) Wait() error {
	c.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.requester.do(c.method, c.endpoint, c.body, c.authWithToken)
	if err != nil {
		c.logger.Error(
			"Request failed for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.logger.Error("Request failed for endpoint " + c.method + c.endpoint + ": Invalid Token")
		return errors.New("Invalid Token !!")
	}

	c.logger.Debug("Successfully called endpoint: " + c.method + c.endpoint)
	return nil
}

// Submit runs the request asynchronously and calls the provided callback.
//
// Parameters:
//   - callback — function to invoke with (error) once the request completes.
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

// Shutdown gracefully shuts down the REST API client.
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
// Usage example:
//
//	gateway, err := .getGateway().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Gateway URL:", gateway.URL)
//
// Callers can use:
//   - Wait() to run synchronously.
//   - Submit(callback) to run asynchronously with a callback.
//
// Returns: *callWithData[gateway] — prepared request object to fetch gateway information.
func (r *restApi) getGateway() *callWithData[gateway] {
	return &callWithData[gateway]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/gateway",
		authWithToken: false,
		parse: func(b []byte) (*gateway, error) {
			obj := gateway{}
			err := sonic.Unmarshal(b, &obj)
			return &obj, err
		},
	}
}

// GetGatewayBot constructs a callWithData for the GET /gateway/bot endpoint.
//
// Usage example:
//
//	gatewayBot, err := .GetGatewayBot().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Recommended shards:", gatewayBot.Shards)
//
// Callers can use:
//   - Wait() to run synchronously.
//   - Submit(callback) to run asynchronously with a callback.
//
// Returns: *callWithData[GatewayBot] — prepared request object to fetch gateway bot information.
func (r *restApi) GetGatewayBot() *callWithData[GatewayBot] {
	return &callWithData[GatewayBot]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/gateway/bot",
		authWithToken: true,
		parse: func(b []byte) (*GatewayBot, error) {
			obj := GatewayBot{}
			err := sonic.Unmarshal(b, &obj)
			return &obj, err
		},
	}
}

/***********************
 *    User Endpoints   *
 ***********************/

// GetSelfUser retrieves the current user's data.
//
// Usage example:
//
//	selfUser, err := .GetSelfUser().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Bot username:", selfUser.Username)
//
// Callers can use:
//   - Wait() to run synchronously.
//   - Submit(callback) to run asynchronously with a callback.
//
// Returns: *callWithData[User] — prepared request object to fetch self user data.
func (r *restApi) GetSelfUser() *callWithData[User] {
	return &callWithData[User]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/users/@me",
		authWithToken: true,
		parse: func(b []byte) (*User, error) {
			obj := User{restApi: r}
			err := sonic.Unmarshal(b, &obj)
			return &obj, err
		},
	}
}

// ModifySelfUserParams defines fields to modify in the current user account.
type ModifySelfUserParams struct {
	// Username is the new username.
	//
	// Optional: leave empty to keep unchanged.
	Username string `json:"username,omitempty"`

	// Avatar is the new avatar image data.
	//
	// Optional: leave nil to keep unchanged.
	Avatar *Attachment `json:"avatar,omitempty"`

	// Banner is the new banner image data.
	//
	// Optional: leave nil to keep unchanged.
	Banner *Attachment `json:"banner,omitempty"`
}

// MarshalJSON is a method used by yada internaly.
func (p *ModifySelfUserParams) MarshalJSON() ([]byte, error) {
	type Alias ModifySelfUserParams

	aux := struct {
		*Alias
		Avatar *string `json:"avatar,omitempty"`
		Banner *string `json:"banner,omitempty"`
	}{
		Alias: (*Alias)(p),
	}

	if p.Avatar != nil {
		aux.Avatar = &p.Avatar.DataURI
	}
	if p.Banner != nil {
		aux.Banner = &p.Banner.DataURI
	}

	return sonic.Marshal(aux)
}

// ModifySelfUser updates the current (self) user account settings.
//
// Usage example: (update the username and avatar and leave the current banner)
//
//	update := &ModifySelfUserParams{
//	    Username: "new_username",
//	    Avatar:   yada.NewAttachment("path/to/avatar.png"),
//	}
//	user, err := .ModifySelfUser(update).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Updated username:", user.Username)
//
// Callers can use:
//   - Wait() to run synchronously.
//   - Submit(callback) to run asynchronously with a callback.
//
// Parameters:
//   - update — pointer to ModifySelfUserParams containing fields to update.
//
// Returns: *callWithNoData — prepared request object to modify self user data.
func (r *restApi) ModifySelfUser(update *ModifySelfUserParams) *callWithNoData {
	body, _ := sonic.Marshal(update)

	return &callWithNoData{
		requester:     r.requester,
		logger:        r.logger,
		method:        "PATCH",
		endpoint:      "/users/@me",
		authWithToken: true,
		body:          body,
	}
}

// GetUser retrieves a user by their Snowflake ID.
//
// Usage example:
//
//	userID := 123456789012345678
//	user, err := .GetUser(userID).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User username:", user.Username)
//
// Callers can use:
//   - Wait() to run synchronously.
//   - Submit(callback) to run asynchronously with a callback.
//
// Parameters:
//   - userID — Snowflake ID of the user.
//
// Returns: *callWithData[User] — prepared request object to fetch user data.
func (r *restApi) GetUser(userID Snowflake) *callWithData[User] {
	return &callWithData[User]{
		requester:     r.requester,
		logger:        r.logger,
		method:        "GET",
		endpoint:      "/users/" + userID.String(),
		authWithToken: true,
		parse: func(b []byte) (*User, error) {
			obj := User{restApi: r}
			err := sonic.Unmarshal(b, &obj)
			return &obj, err
		},
	}
}
