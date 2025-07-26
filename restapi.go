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

// restApi provides methods for Discord REST API endpoints.
type restApi struct {
	req    *requester
	logger Logger
}

// newRestApi creates a new RestAPI instance with optional custom requester and logger.
func newRestApi(req *requester, logger Logger) *restApi {
	return &restApi{
		req:    req,
		logger: logger,
	}
}

// Shutdown gracefully shuts down the REST API client.
func (r *restApi) Shutdown() {
	r.logger.Info("RestAPI shutting down")
	r.req.Shutdown()
	r.logger = nil
	r.req = nil
}

/***********************
 *      Helpers        *
 ***********************/

func (r *restApi) doRequest(method, endpoint string, body []byte, authWithToken bool, reason string) ([]byte, error) {
	r.logger.Debug("Calling endpoint: " + method + endpoint)

	res, err := r.req.do(method, endpoint, body, authWithToken, reason)
	if err != nil {
		r.logger.Error("Request failed for endpoint " + method + endpoint + ": " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		r.logger.Error("Request failed for endpoint " + method + endpoint + ": Invalid Token")
		return nil, errors.New("invalid token")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		r.logger.Error("Failed reading response body for endpoint " + method + endpoint + ": " + err.Error())
		return nil, err
	}

	r.logger.Debug("Successfully called endpoint: " + method + endpoint)
	return bodyBytes, nil
}

/***********************
 *      Endpoints      *
 ***********************/

// FetchGatewayBot retrieves bot gateway information including recommended shard count and session limits.
//
// Usage example:
//
//	gateway, err := api.FetchGatewayBot()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Recommended shards:", gateway.Shards)
//
// Returns:
//   - GatewayBot: the bot gateway information.
//   - error: if the request failed or decoding failed.
func (r *restApi) FetchGatewayBot() (GatewayBot, error) {
	body, err := r.doRequest("GET", "/gateway/bot", nil, true, "")
	if err != nil {
		return GatewayBot{}, err
	}

	var obj GatewayBot
	if err := sonic.Unmarshal(body, &obj); err != nil {
		r.logger.Error("Failed parsing response for /gateway/bot: " + err.Error())
		return GatewayBot{}, err
	}
	return obj, nil
}

// FetchSelfUser retrieves the current bot user's data including username, ID, avatar, and flags.
//
// Usage example:
//
//	user, err := api.FetchSelfUser()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Bot username:", user.Username)
//
// Returns:
//   - User: the current user data.
//   - error: if the request failed or decoding failed.
func (r *restApi) FetchSelfUser() (User, error) {
	body, err := r.doRequest("GET", "/users/@me", nil, true, "")
	if err != nil {
		return User{}, err
	}

	var obj User
	if err := sonic.Unmarshal(body, &obj); err != nil {
		r.logger.Error("Failed parsing response for /users/@me: " + err.Error())
		return User{}, err
	}
	return obj, nil
}

// ModifySelfUserParams defines the parameters to update the current user account.
//
// All fields are optional:
//   - If a field is not set (left nil or empty), it will remain unchanged.
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
//	err := api.ModifySelfUser(update)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User updated successfully")
//
// Returns:
//   - error: if the request failed.
func (r *restApi) ModifySelfUser(update *ModifySelfUserParams) error {
	body, _ := sonic.Marshal(update)
	_, err := r.doRequest("PATCH", "/users/@me", body, true, "")
	return err
}

// FetchUser retrieves a user by their Snowflake ID including username, avatar, and flags.
//
// Usage example:
//
//	user, err := api.FetchUser(123456789012345678)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Username:", user.Username)
//
// Returns:
//   - User: the user data.
//   - error: if the request failed or decoding failed.
func (r *restApi) FetchUser(userID Snowflake) (User, error) {
	body, err := r.doRequest("GET", "/users/"+userID.String(), nil, true, "")
	if err != nil {
		return User{}, err
	}

	var obj User
	if err := sonic.Unmarshal(body, &obj); err != nil {
		r.logger.Error("Failed parsing response for /users/{id}: " + err.Error())
		return User{}, err
	}
	return obj, nil
}

// FetchChannel retrieves a channel by its Snowflake ID and decodes it into its concrete type
// (e.g. TextChannel, VoiceChannel, CategoryChannel).
//
// Usage example:
//
//	channel, err := api.FetchChannel(123456789012345678)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Channel ID:", channel.GetID())
//
// Returns:
//   - Channel: the decoded channel object.
//   - error: if the request failed or the type is unknown or decoding failed.
func (r *restApi) FetchChannel(channelID Snowflake) (Channel, error) {
	body, err := r.doRequest("GET", "/channels/"+channelID.String(), nil, true, "")
	if err != nil {
		return nil, err
	}

	var meta struct{ Type ChannelType }
	if err := sonic.Unmarshal(body, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ChannelTypeGuildCategory:
		var c CategoryChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildText:
		var c TextChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildVoice:
		var c VoiceChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildAnnouncement:
		var c AnnouncementChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildStageVoice:
		var c StageVoiceChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildForum:
		var c ForumChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeGuildMedia:
		var c MediaChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypeAnnouncementThread:
		var c AnnouncementThreadChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypePrivateThread:
		var c PrivateThreadChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	case ChannelTypePublicThread:
		var c PublicThreadChannel
		err = sonic.Unmarshal(body, &c)
		return &c, err
	default:
		return nil, errors.New("unknown channel type")
	}
}
