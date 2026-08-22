/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

import "context"

// CDNSoundboardSoundsURL is the CDN base URL for soundboard audio files.
const CDNSoundboardSoundsURL = "https://cdn.discordapp.com/soundboard-sounds/"

// SoundboardSound is an alias for SoundBoardSound.
type SoundboardSound = SoundBoardSound

// SoundBoardSound represents a soundboard sound in Discord.
//
// Reference: https://discord.com/developers/docs/resources/soundboard#soundboard-sound-object
type SoundBoardSound struct {
	// SoundID is the sound's unique Discord snowflake ID.
	SoundID Snowflake `json:"sound_id"`

	// Name is the sound's name.
	Name string `json:"name"`

	// Volume is the volume of this sound, from 0 to 1.
	Volume float64 `json:"volume"`

	// EmojiID is the id of this sound's custom emoji.
	//
	// Optional:
	//   - Will be 0 if the emoji is Unicode (standard emoji).
	EmojiID Snowflake `json:"emoji_id,omitempty"`

	// EmojiName is the name of this sound's standard emoji.
	//
	// Optional:
	//   - Will be empty string if the emoji is custom (not standard emoji).
	EmojiName string `json:"emoji_name,omitempty"`

	// Available is whether this sound can be used, may be false due to loss of Server Boosts.
	Available bool `json:"available"`

	// GuildID is the ID of the guild this sound belongs to.
	//
	// Optional:
	//   - Will be absent if the sound is global (not guild-specific).
	GuildID Snowflake `json:"guild_id,omitempty"`

	// User is the user who created this sound.
	//
	// Optional:
	//   - Will be absent if the sound is global (not guild-specific).
	User *User `json:"user,omitempty"`
}

// URL returns the CDN URL of the soundboard sound.
func (s *SoundBoardSound) URL() string {
	return CDNSoundboardSoundsURL + s.SoundID.String()
}

// Save downloads the soundboard's sound from its URL and saves it to disk.
//
// If fileName is not provided (empty string), it saves the file in the given
// directory using Attachment.Filename
//
// Info:
//   - The extension is replaced based on the Content-Type of the file.
//
// Example:
//
//	err := sound.Save("mysound", "./downloads")
//	if err != nil {
//	    // handle error
//	}
//
// Returns:
//   - string: full path to the downloaded file.
//   - error: non-nil if any operation fails.
func (s *SoundBoardSound) Save(ctx context.Context, fileName, dir string) (string, error) {
	if fileName == "" {
		return DownloadFile(ctx, s.URL(), s.Name, dir)
	}
	return DownloadFile(ctx, s.URL(), fileName, dir)
}

// String implements the fmt.Stringer interface.
func (s *SoundBoardSound) String() string {
	return s.Name
}
