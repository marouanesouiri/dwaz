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
 *       RestAPI       *
 ***********************/

// RestAPI provides methods for Discord REST API endpoints.
type RestAPI struct {
	req    *requester
	logger Logger
}

// newRestAPI creates a new RestAPI instance with optional custom requester and logger.
func newRestAPI(req *requester, token string, logger Logger) *RestAPI {
	return &RestAPI{
		req:    req,
		logger: logger,
	}
}

// Shutdown gracefully shuts down the REST API client.
func (r *RestAPI) Shutdown() {
	r.logger.Info("RestAPI shutting down")
	r.req.Shutdown()
	r.logger = nil
	r.req = nil
}

/***********************
 *       Calls         *
 ***********************/

// call contains common HTTP request data and logic.
type call struct {
	api           *RestAPI
	method        string
	endpoint      string
	body          []byte
	authWithToken bool
	reason        string
}

// doRequest performs the HTTP request and returns the response body bytes.
func (c *call) doRequest() ([]byte, error) {
	c.api.logger.Debug("Calling endpoint: " + c.method + c.endpoint)

	res, err := c.api.req.do(c.method, c.endpoint, c.body, c.authWithToken, c.reason)
	if err != nil {
		c.api.logger.Error("Request failed for endpoint " + c.method + c.endpoint + ": " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		c.api.logger.Error("Request failed for endpoint " + c.method + c.endpoint + ": Invalid Token")
		return nil, errors.New("invalid token")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.api.logger.Error("Failed reading response body for endpoint " + c.method + c.endpoint + ": " + err.Error())
		return nil, err
	}

	c.api.logger.Debug("Successfully called endpoint: " + c.method + c.endpoint)
	return bodyBytes, nil
}

/***********************
 *      Call[T]        *
 ***********************/

// Call represents a REST API request returning typed decoded data.
type Call[T any] struct {
	call
}

// Wait executes the request synchronously and parses the response into type T.
func (c *Call[T]) Wait() (*T, error) {
	body, err := c.doRequest()
	if err != nil {
		return nil, err
	}

	var obj T
	if err := sonic.Unmarshal(body, &obj); err != nil {
		c.api.logger.Error("Failed parsing response for endpoint " + c.method + c.endpoint + ": " + err.Error())
		return nil, err
	}

	return &obj, nil
}

// Submit runs the request asynchronously and calls the provided callback.
func (c *Call[T]) Submit(callback func(*T, error)) {
	go func() { callback(c.Wait()) }()
}

/***********************
 *    CallNoData       *
 ***********************/

// CallNoData represents a REST API request with no data parsing.
type CallNoData struct {
	call
}

// Wait executes the request synchronously.
func (c *CallNoData) Wait() error {
	_, err := c.doRequest()
	return err
}

// Submit runs the request asynchronously and calls the provided callback.
func (c *CallNoData) Submit(callback func(error)) {
	go func() { callback(c.Wait()) }()
}

/***********************
 *   CallDecoded[T]    *
 ***********************/

// CallDecoded represents a REST API request returning decoded data via custom decode function.
type CallDecoded[T any] struct {
	call
	decode func([]byte) (T, error)
}

// Wait executes the request synchronously and parses the response.
func (c *CallDecoded[T]) Wait() (T, error) {
	var zero T

	body, err := c.doRequest()
	if err != nil {
		return zero, err
	}

	obj, err := c.decode(body)
	if err != nil {
		c.api.logger.Error("Failed decoding response for endpoint " + c.method + c.endpoint + ": " + err.Error())
		return zero, err
	}

	if withAPI, ok := any(&obj).(interface{ setRestApi(*RestAPI) }); ok {
		withAPI.setRestApi(c.api)
	}

	return obj, nil
}

// Submit runs the request asynchronously and calls the provided callback.
func (c *CallDecoded[T]) Submit(callback func(T, error)) {
	go func() { callback(c.Wait()) }()
}

// GetGatewayBot retrieves bot gateway information including recommended shard count and session limits.
//
// Usage example:
//
//	gateway, err := api.GetGatewayBot().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Recommended shards:", gateway.Shards)
//
// Callers can use:
//   - Wait() to execute synchronously.
//   - Submit(callback) to execute asynchronously with a callback.
//
// Returns: *Call[GatewayBot] — prepared request object to fetch gateway bot info.
func (r *RestAPI) GetGatewayBot() *Call[GatewayBot] {
	return &Call[GatewayBot]{
		call: call{
			api:           r,
			method:        "GET",
			endpoint:      "/gateway/bot",
			authWithToken: true,
		},
	}
}

// GetSelfUser retrieves the current bot user's data including username, ID, avatar, and flags.
//
// Usage example:
//
//	user, err := api.GetSelfUser().Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Bot username:", user.Username)
//
// Returns: *Call[User] — prepared request object to fetch self user data.
func (r *RestAPI) GetSelfUser() *Call[User] {
	return &Call[User]{
		call: call{
			api:           r,
			method:        "GET",
			endpoint:      "/users/@me",
			authWithToken: true,
		},
	}
}

// ModifySelfUserParams defines fields to modify in the current user account.
//
// Fields:
//   - Username: new username (optional).
//   - Avatar: new avatar image data as Attachment (optional).
//   - Banner: new banner image data as Attachment (optional).
type ModifySelfUserParams struct {
	Username string      `json:"username,omitempty"`
	Avatar   *Attachment `json:"avatar,omitempty"`
	Banner   *Attachment `json:"banner,omitempty"`
}

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

// ModifySelfUser updates the current bot user's username, avatar, or banner.
//
// Usage example:
//
//	update := &ModifySelfUserParams{
//	    Username: "new_username",
//	    Avatar:   yada.NewAttachment("path/to/avatar.png"),
//	}
//	err := api.ModifySelfUser(update).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User updated successfully")
//
// Returns: *CallNoData — prepared request object to modify self user data.
func (r *RestAPI) ModifySelfUser(update *ModifySelfUserParams) *CallNoData {
	body, _ := sonic.Marshal(update)

	return &CallNoData{
		call: call{
			api:           r,
			method:        "PATCH",
			endpoint:      "/users/@me",
			authWithToken: true,
			body:          body,
		},
	}
}

// GetUser retrieves a user by their Snowflake ID including username, avatar, and flags.
//
// Usage example:
//
//	user, err := api.GetUser(123456789012345678).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Username:", user.Username)
//
// Returns: *Call[User] — prepared request object to fetch user data.
func (r *RestAPI) GetUser(userID Snowflake) *Call[User] {
	return &Call[User]{
		call: call{
			api:           r,
			method:        "GET",
			endpoint:      "/users/" + userID.String(),
			authWithToken: true,
		},
	}
}

// GetChannel retrieves a channel by its Snowflake ID and decodes it into its concrete type
// (e.g. TextChannel, VoiceChannel, CategoryChannel).
//
// Usage example:
//
//	channel, err := api.GetChannel(123456789012345678).Wait()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Channel ID:", channel.GetID())
//
// Returns: *CallDecoded[Channel] — prepared request object to fetch channel data.
func (r *RestAPI) GetChannel(channelID Snowflake) *CallDecoded[Channel] {
	return &CallDecoded[Channel]{
		call: call{
			api:           r,
			method:        "GET",
			endpoint:      "/channels/" + channelID.String(),
			authWithToken: true,
		},
		decode: func(body []byte) (Channel, error) {
			var u struct{ Type ChannelType }
			if err := sonic.Unmarshal(body, &u); err != nil {
				return nil, err
			}

			var obj Channel
			var err error

			switch u.Type {
			case ChannelTypeGuildCategory:
				var c CategoryChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildText:
				var c TextChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildVoice:
				var c VoiceChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildAnnouncement:
				var c AnnouncementChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildStageVoice:
				var c StageVoiceChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildForum:
				var c ForumChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeGuildMedia:
				var c MediaChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypeAnnouncementThread:
				var c AnnouncementThreadChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypePrivateThread:
				var c PrivateThreadChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			case ChannelTypePublicThread:
				var c PublicThreadChannel
				err = sonic.Unmarshal(body, &c)
				obj = &c
			default:
				err = errors.New("unknown channel type")
			}
			return obj, err
		},
	}
}
