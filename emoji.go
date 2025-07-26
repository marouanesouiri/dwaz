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

// Emoji represents a custom emoji object used within a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	// ID is the unique Discord snowflake ID of the emoji.
	//
	// Optional:
	// - May be nil (zero value) for unicode emojis
	ID Snowflake `json:"id"`

	// Name is the emoji's name.
	//
	// Optional:
	// - May be empty in reaction emoji objects.
	// - Always present in full emoji objects returned from guild endpoints.
	Name string `json:"name"`

	// Roles is a list of role IDs allowed to use this emoji.
	Roles []Snowflake `json:"roles"`

	// RequireColons indicates whether the emoji must be wrapped in colons to be used.
	RequireColons bool `json:"require_colons"`

	// Managed indicates whether the emoji is managed by an integration.
	Managed bool `json:"managed"`

	// Animated indicates whether the emoji is an animated emoji (.gif).
	Animated bool `json:"animated"`

	// Available indicates whether the emoji can currently be used.
	Available bool `json:"available"`
}

// Mention returns a Discord mention string for the emoji.
//
// Example output: "<:sliming:123456789012345678>"
func (e *Emoji) Mention() string {
	mention := "<"
	if e.Animated {
		mention += "a"
	}
	mention += e.Name + ":" + e.ID.String() + ":>"
	return mention
}
