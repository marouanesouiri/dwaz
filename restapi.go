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

import "io"

/***********************
 *	  callWithData	   *
 ***********************/

// callWithData represents a REST API request returning typed decoded data.
type callWithData[T any] struct {
	requester       *requester
	logger          Logger
	method          string
	endpoint        string
	body            []byte
	authNotRequired bool
	parse           func([]byte) (*T, error)
}

// wait executes the request synchronously and parses the response.
func (c *callWithData[T]) wait() (*T, error) {
	c.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.requester.do(c.method, c.endpoint, c.body, c.authNotRequired)
	if err != nil {
		c.logger.Error(
			"Request failed for endpoint " + c.method + c.endpoint + ": " + err.Error(),
		)
		return nil, err
	}
	defer res.Body.Close()

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

// submit runs the request asynchronously and calls the provided callback.
func (c *callWithData[T]) submit(callback func(*T, error)) {
	go func() { callback(c.wait()) }()
}

/***********************
 *	 callWithNoData	   *
 ***********************/

// callWithNoData represents a REST API request with no data parsing.
type callWithNoData struct {
	requester       *requester
	logger          Logger
	method          string
	endpoint        string
	body            []byte
	authNotRequired bool
}

// wait executes the request synchronously.
func (c *callWithNoData) wait() error {
	c.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.requester.do(c.method, c.endpoint, c.body, c.authNotRequired)
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

// submit runs the request asynchronously and calls the provided callback.
func (c *callWithNoData) submit(callback func(error)) {
	go func() { callback(c.wait()) }()
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

/***********************
 *   Gateway Endpoint  *
 ***********************/

// getGateway returns a callWithData for the GET /gateway endpoint.
func (r *restApi) getGateway() *callWithData[gateway] {
	return &callWithData[gateway]{
		requester:       r.requester,
		logger:          r.logger,
		method:          "GET",
		endpoint:        "/gateway",
		authNotRequired: true,
		parse: func(b []byte) (*gateway, error) {
			obj := gateway{}
			return &obj, obj.fillFromJson(b)
		},
	}
}

// getGatewayBot returns a callWithData for the GET /gateway/bot endpoint.
func (r *restApi) getGatewayBot() *callWithData[gatewayBot] {
	return &callWithData[gatewayBot]{
		requester: r.requester,
		logger:    r.logger,
		method:    "GET",
		endpoint:  "/gateway/bot",
		parse: func(b []byte) (*gatewayBot, error) {
			obj := gatewayBot{}
			return &obj, obj.fillFromJson(b)
		},
	}
}
