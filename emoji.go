/************************************************************************************
 *
 * dwaz (Discord Wrapper API for Zwafriya), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package dwaz

import (
	"errors"
	"strings"
	"time"

	"encoding/json"

	"github.com/marouanesouiri/stdx/result"
)

// Emoji represents a custom emoji object used within a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	// ID is the unique Discord snowflake ID of the emoji.
	//
	// Optional:
	//   - May zero for unicode emojis
	ID Snowflake `json:"id,omitempty"`

	// Name is the emoji's name.
	//
	// Optional:
	//   - May be empty in deleted emojis.
	Name string `json:"name,omitempty"`

	// Roles is a list of role IDs allowed to use this emoji.
	Roles []Snowflake `json:"roles,omitempty"`

	// RequireColons indicates whether the emoji must be wrapped in colons to be used.
	RequireColons bool `json:"require_colons,omitempty"`

	// Managed indicates whether the emoji is managed by an integration.
	Managed bool `json:"managed,omitempty"`

	// Animated indicates whether the emoji is an animated emoji (.gif).
	Animated bool `json:"animated,omitempty"`

	// Available indicates whether the emoji can currently be used.
	Available bool `json:"available,omitempty"`
}

// Mention returns a Discord mention string for the emoji.
//
// Example output: "<:sliming:123456789012345678>"
func (e *Emoji) Mention() string {
	if e.ID == 0 { // (no id == unicode emoji)
		if e.RequireColons {
			return ":" + e.Name + ":"
		}
		return e.Name
	}
	mention := "<"
	if e.Animated {
		mention += "a"
	}
	mention += ":" + e.Name + ":" + e.ID.String() + ">"
	return mention
}

// String implements the fmt.Stringer interface.
func (e *Emoji) String() string {
	return e.Mention()
}

// ParseEmoji parses a Discord emoji mention string into an Emoji object.
//
// Supports:
//   - Custom emojis: <:name:id> or <a:name:id>
//   - Unicode emojis: just the string
func ParseEmoji(mention string) (Emoji, error) {
	if !strings.HasPrefix(mention, "<") || !strings.HasSuffix(mention, ">") {
		return Emoji{Name: mention}, nil
	}

	content := mention[1 : len(mention)-1]
	parts := strings.Split(content, ":")

	if len(parts) != 3 {
		return Emoji{}, errors.New("invalid emoji mention format")
	}

	if parts[0] != "" && parts[0] != "a" {
		return Emoji{}, errors.New("invalid emoji prefix")
	}

	animated := parts[0] == "a"
	name := parts[1]
	id, err := ParseSnowflake(parts[2])
	if err != nil {
		return Emoji{}, err
	}

	return Emoji{
		ID:            id,
		Name:          name,
		Animated:      animated,
		RequireColons: true,
	}, nil
}

// CreatedAt returns the time when this emojis is created at.
func (e *Emoji) CreatedAt() time.Time {
	if e.ID == 0 {
		return time.Time{}
	}
	return e.ID.Timestamp()
}

// URL returns the URL to the emoji's image.
func (e *Emoji) URL() string {
	format := ImageFormatPNG
	if e.Animated {
		format = ImageFormatGIF
	}
	return EmojiURL(e.ID, format, ImageSizeDefault)
}

// URLWith returns the URL to the emoji's image.
// allowing explicit specification of image format and size.
func (e *Emoji) URLWith(format ImageFormat, size ImageSize) string {
	return EmojiURL(e.ID, format, size)
}

// PartialEmoji represents a partial emoji object used in a Discord poll, typically within a PollMedia object for poll answers,
// or when sending a message with a poll request.
//
// When creating a poll answer, provide only the ID for a custom emoji or only the Name for a Unicode emoji.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-media-object-poll-media-object-structure
type PartialEmoji struct {
	// ID is the unique identifier for a custom emoji.
	// When sending a poll request with a custom emoji, provide only the ID and leave Name empty.
	//
	// Optional:
	//  - Will be 0 if no ID is set (e.g., for Unicode emojis or when not provided in a response).
	ID Snowflake `json:"id,omitempty"`

	// Name is the name of the emoji, used for Unicode emojis (e.g., "😊").
	// When sending a poll request with a Unicode emoji, provide only the Name and leave ID as 0.
	//
	// Optional:
	//  - Will be empty if no name is set (e.g., for custom emojis or when not provided in a response).
	Name string `json:"name,omitempty"`

	// Animated indicates whether the emoji is animated.
	Animated bool `json:"animated"`
}

// ListGuildEmojis retrieves a list of emoji objects for the given guild.
//
// Includes user fields if the bot has the PermissionCreateGuildExpressions or PermissionManageGuildExpressions permission.
//
// Reference: https://discord.com/developers/docs/resources/emoji#list-guild-emojis
func (r *requester) ListGuildEmojis(guildID Snowflake) result.Result[[]Emoji] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/guilds/" + guildID.String() + "/emojis",
	})
	if res.IsErr() {
		return result.Err[[]Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emojis []Emoji
	if err := json.NewDecoder(body).Decode(&emojis); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/emojis",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Emoji](err)
	}
	return result.Ok(emojis)
}

// FetchGuildEmoji retrieves an emoji object for the given guild and emoji IDs.
//
// Includes the user field if the bot has the PermissionManageGuildExpressions permission,
// or if the bot created the emoji and has the PermissionCreateGuildExpressions permission.
//
// Reference: https://discord.com/developers/docs/resources/emoji#get-guild-emoji
func (r *requester) FetchGuildEmoji(guildID, emojiID Snowflake) result.Result[Emoji] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/guilds/" + guildID.String() + "/emojis/" + emojiID.String(),
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/emojis/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// CreateGuildEmojiOptions contains parameters for creating a new guild emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#create-guild-emoji-json-params
type CreateGuildEmojiOptions struct {
	// Name is the name of the emoji.
	Name string `json:"name"`

	// Image is the emoji image data.
	// Emojis and animated emojis have a maximum file size of 256 KiB.
	Image Base64Image `json:"image"`

	// Roles are the roles allowed to use this emoji.
	Roles []Snowflake `json:"roles,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// CreateGuildEmoji creates a new emoji for the guild.
//
// Requires the PermissionCreateGuildExpressions permission.
//
// Emojis and animated emojis have a maximum file size of 256 KiB.
// Attempting to upload an emoji larger than this limit will fail with an error.
//
// Reference: https://discord.com/developers/docs/resources/emoji#create-guild-emoji
func (r *requester) CreateGuildEmoji(guildID Snowflake, opts CreateGuildEmojiOptions) result.Result[Emoji] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/guilds/" + guildID.String() + "/emojis",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/guilds/{id}/emojis",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// ModifyGuildEmojiOptions contains parameters for modifying an existing guild emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#modify-guild-emoji-json-params
type ModifyGuildEmojiOptions struct {
	// Name is the name of the emoji.
	Name string `json:"name,omitempty"`

	// Roles are the roles allowed to use this emoji.
	Roles []Snowflake `json:"roles,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyGuildEmoji updates the given emoji.
//
// For emojis created by the current user, requires either the PermissionCreateGuildExpressions or PermissionManageGuildExpressions permission.
// For other emojis, requires the PermissionManageGuildExpressions permission.
//
// Reference: https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
func (r *requester) ModifyGuildEmoji(guildID, emojiID Snowflake, opts ModifyGuildEmojiOptions) result.Result[Emoji] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/guilds/" + guildID.String() + "/emojis/" + emojiID.String(),
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/emojis/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// DeleteGuildEmoji deletes the given emoji.
//
// For emojis created by the current user, requires either the PermissionCreateGuildExpressions or PermissionManageGuildExpressions permission.
// For other emojis, requires the PermissionManageGuildExpressions permission.
//
// Reference: https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
func (r *requester) DeleteGuildEmoji(guildID, emojiID Snowflake, reason string) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/guilds/" + guildID.String() + "/emojis/" + emojiID.String(),
		Reason: reason,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

type applicationEmojisResponse struct {
	Items []Emoji `json:"items"`
}

// ListApplicationEmojis retrieves all emojis for an application.
//
// Includes a user object for the team member that uploaded the emoji from the app's settings,
// or for the bot user if uploaded using the API.
//
// Reference: https://discord.com/developers/docs/resources/emoji#list-application-emojis
func (r *requester) ListApplicationEmojis(applicationID Snowflake) result.Result[[]Emoji] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/applications/" + applicationID.String() + "/emojis",
	})
	if res.IsErr() {
		return result.Err[[]Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var response applicationEmojisResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/applications/{id}/emojis",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Emoji](err)
	}
	return result.Ok(response.Items)
}

// FetchApplicationEmoji retrieves a specific emoji from an application.
//
// Includes the user field.
//
// Reference: https://discord.com/developers/docs/resources/emoji#get-application-emoji
func (r *requester) FetchApplicationEmoji(applicationID, emojiID Snowflake) result.Result[Emoji] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/applications/" + applicationID.String() + "/emojis/" + emojiID.String(),
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/applications/{id}/emojis/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// CreateApplicationEmojiOptions contains parameters for creating a new application emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#create-application-emoji-json-params
type CreateApplicationEmojiOptions struct {
	// Name is the name of the emoji.
	Name string `json:"name"`

	// Image is the emoji image data.
	// Emojis and animated emojis have a maximum file size of 256 KiB.
	Image Base64Image `json:"image"`
}

// CreateApplicationEmoji creates a new emoji for an application.
//
// Emojis and animated emojis have a maximum file size of 256 KiB.
// Attempting to upload an emoji larger than this limit will fail with an error.
//
// Reference: https://discord.com/developers/docs/resources/emoji#create-application-emoji
func (r *requester) CreateApplicationEmoji(applicationID Snowflake, opts CreateApplicationEmojiOptions) result.Result[Emoji] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/applications/" + applicationID.String() + "/emojis",
		Body:   reqBody,
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/applications/{id}/emojis",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// ModifyApplicationEmojiOptions contains parameters for updating an existing application emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#modify-application-emoji-json-params
type ModifyApplicationEmojiOptions struct {
	// Name is the name of the emoji.
	Name string `json:"name,omitempty"`
}

// ModifyApplicationEmoji updates an existing application emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#modify-application-emoji
func (r *requester) ModifyApplicationEmoji(applicationID, emojiID Snowflake, opts ModifyApplicationEmojiOptions) result.Result[Emoji] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/applications/" + applicationID.String() + "/emojis/" + emojiID.String(),
		Body:   reqBody,
	})
	if res.IsErr() {
		return result.Err[Emoji](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var emoji Emoji
	if err := json.NewDecoder(body).Decode(&emoji); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/applications/{id}/emojis/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Emoji](err)
	}
	return result.Ok(emoji)
}

// DeleteApplicationEmoji deletes an application emoji.
//
// Reference: https://discord.com/developers/docs/resources/emoji#delete-application-emoji
func (r *requester) DeleteApplicationEmoji(applicationID, emojiID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/applications/" + applicationID.String() + "/emojis/" + emojiID.String(),
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}
