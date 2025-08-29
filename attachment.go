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

import "time"

// AttachmentFlags represents bit flags for Discord attachment metadata.
//
// Reference: https://discord.com/developers/docs/resources/channel#attachment-object-attachment-flags
type AttachmentFlags int

const (
	// AttachmentFlagIsRemix means this attachment has been edited using the remix feature on mobile.
	AttachmentFlagIsRemix AttachmentFlags = 1 << 2
)

// Attachment represents a Discord attachment object.
//
// Reference: https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	// ID is the unique Discord snowflake ID of this attachment.
	ID Snowflake `json:"id,omitempty"`

	// Filename is the name of the attached file.
	Filename string `json:"filename"`

	// Title is the optional title of the attachment.
	//
	// Optional
	//  - May be empty string if unset.
	Title string `json:"title,omitempty"`

	// Description is the optional description of the attachment (max 1024 characters).
	//
	// Optional:
	//  - May be empty string if unset.
	Description string `json:"description,omitempty"`

	// ContentType is the media type (MIME type) of the attachment.
	//
	// Optional:
	//  - May be empty string.
	ContentType string `json:"content_type,omitempty"`

	// Size is the size of the file in bytes.
	Size int `json:"size,omitempty"`

	// URL is the source URL of the attachment file.
	URL string `json:"url,omitempty"`

	// ProxyURL is a proxied URL of the attachment file.
	ProxyURL string `json:"proxy_url,omitempty"`

	// Height is the height of the image file, if applicable.
	//  - 0 if the attachment is not an image.
	Height int `json:"height,omitempty"`

	// Width is the width of the image file, if applicable.
	//  - 0 if the attachment is not an image.
	Width int `json:"width,omitempty"`

	// Ephemeral indicates whether this attachment is ephemeral.
	Ephemeral bool `json:"ephemeral,omitempty"`

	// Flags is a bitfield combining attachment flags.
	Flags AttachmentFlags `json:"flags,omitempty"`

	// DurationSec is the duration of the audio file in seconds, if applicable.
	//
	// Optional:
	// 	- Present only for audio or voice message attachments.
	DurationSec *float64 `json:"duration_secs,omitempty"`

	// Waveform is a base64 encoded byte array representing a sampled waveform.
	//
	// Optional:
	//  - present only for voice messages.
	Waveform *string `json:"waveform,omitempty"`
}

// CreatedAt returns the time when this attachment is created.
func (a *Attachment) CreatedAt() time.Time {
	return a.ID.Timestamp()
}

// Save downloads the attachment from its URL and saves it to disk.
//
// If fileName is not provided (empty string), it saves the file in the given
// directory using Attachment.Filename
//
// Info:
//   - The extension is replaced based on the Content-Type of the file.
//
// Example:
//
//	err := attachment.Save("myfile", "./downloads")
//	if err != nil {
//	    // handle error
//	}
//
// Returns:
//   - string: full path to the downloaded file.
//   - error: non-nil if any operation fails.
func (a *Attachment) Save(fileName, dir string) (string, error) {
	if fileName == "" {
		fileName = a.Filename
	}
	return DownloadFile(a.URL, fileName, dir)
}
