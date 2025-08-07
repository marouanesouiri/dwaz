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

import "time"

// Emoji represents a custom emoji object used within a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	// ID is the unique Discord snowflake ID of the emoji.
	//
	// Optional:
	// - May be nil (zero value) for unicode emojis
	ID Snowflake `json:"id,omitempty"`

	// Name is the emoji's name.
	//
	// Optional:
	// - May be empty in deleted emojis.
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
		return e.Name
	}
	mention := "<"
	if e.Animated {
		mention += "a"
	}
	mention += e.Name + ":" + e.ID.String() + ":>"
	return mention
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
	var format EmojiFormat = EmojiFormatPNG
	if e.Animated {
		format = EmojiFormatGIF
	}
	return EmojiURL(e.ID, format, ImageSize32)
}

// URLWith returns the URL to the emoji's image.
// allowing explicit specification of image format and size.
func (e *Emoji) URLWith(format EmojiFormat, size ImageSize) string {
	return EmojiURL(e.ID, format, size)
}

