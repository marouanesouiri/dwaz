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
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AttachmentFlags represents bit flags for Discord attachment metadata.
//
// Reference: https://discord.com/developers/docs/resources/channel#attachment-object-attachment-flags
type AttachmentFlags int

const (
	// AttachmentFlag_IsRemix means this attachment has been edited using the remix feature on mobile.
	AttachmentFlag_IsRemix AttachmentFlags = 1 << 2
)

// Attachment represents a Discord attachment object.
//
// Reference: https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	// ID is the unique Discord snowflake ID of this attachment.
	//
	// Always present when received from Discord.
	ID Snowflake `json:"id,omitempty"`

	// Filename is the name of the attached file.
	//
	// Always present.
	Filename string `json:"filename"`

	// Title is the optional title of the attachment.
	//
	// Optional: may be empty or omitted.
	Title string `json:"title,omitempty"`

	// Description is the optional description of the attachment (max 1024 characters).
	//
	// Optional: may be empty or omitted.
	Description string `json:"description,omitempty"`

	// ContentType is the media type (MIME type) of the attachment.
	//
	// Optional: may be omitted or empty.
	ContentType string `json:"content_type,omitempty"`

	// Size is the size of the file in bytes.
	//
	// Always present when received from Discord or set by NewAttachment.
	Size int `json:"size,omitempty"`

	// URL is the source URL of the attachment file.
	//
	// Always present when received from Discord.
	URL string `json:"url,omitempty"`

	// ProxyURL is a proxied URL of the attachment file.
	//
	// Always present when received from Discord.
	ProxyURL string `json:"proxy_url,omitempty"`

	// Height is the height of the image file, if applicable.
	//
	// Always present: 0 if the attachment is not an image.
	Height int `json:"height,omitempty"`

	// Width is the width of the image file, if applicable.
	//
	// Always present: 0 if the attachment is not an image.
	Width int `json:"width,omitempty"`

	// Ephemeral indicates whether this attachment is ephemeral.
	//
	// Always present: false if not provided.
	Ephemeral bool `json:"ephemeral,omitempty"`

	// DurationSec is the duration of the audio file in seconds, if applicable.
	//
	// Optional: present only for audio or voice message attachments.
	DurationSec *float64 `json:"duration_secs,omitempty"`

	// Waveform is a base64 encoded byte array representing a sampled waveform.
	//
	// Optional: present only for voice messages.
	Waveform string `json:"waveform,omitempty"`

	// Flags is a bitfield combining attachment flags.
	//
	// Always present: 0 if no flags are set.
	Flags int `json:"flags,omitempty"`

	// DataURI holds a base64 encoded data URI of the file contents.
	//
	// Used internally for local files to upload.
	// Not sent or received from Discord.
	DataURI string `json:"-"`
}

// NewAttachment creates an Attachment from a local file path.
//
// Reads the file and encodes it as a base64 data URI to prepare for upload.
//
// Automatically sets Filename, Size, ContentType, DataURI, Flags (0), Ephemeral (false),
// and Height/Width if the file is an image.
//
// Example:
//
//	att, err := NewAttachment("path/to/file.png")
//	if err != nil {
//	    // handle error
//	}
//	// Use att.DataURI to send the file data.
func NewAttachment(path string) (*Attachment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	// Efficiently build base64 data URI
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	buf := make([]byte, 0, len("data:")+len(mimeType)+len(";base64,")+encodedLen)
	buf = append(buf, "data:"...)
	buf = append(buf, mimeType...)
	buf = append(buf, ";base64,"...)
	encoded := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encoded, data)
	buf = append(buf, encoded...)

	height, width := 0, 0
	if strings.HasPrefix(mimeType, "image/") {
		cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err == nil {
			height = cfg.Height
			width = cfg.Width
		}
	}

	_, filename := filepath.Split(path)
	return &Attachment{
		Filename:    filename,
		ContentType: mimeType,
		Size:        len(data),
		DataURI:     string(buf),
		Flags:       0,
		Ephemeral:   false,
		Height:      height,
		Width:       width,
	}, nil
}

// Save downloads the attachment from its URL and saves it to disk.
//
// It saves the file in the given directory with its Attachment.Filename.
// The extension is replaced based on ContentType if available.
//
// Returns an error if any operation fails.
//
// Example:
//
//	err := attachment.Save("./downloads")
//	if err != nil {
//	    // handle error
//	}
func (a *Attachment) Save(dir string) error {
	if a.URL == "" {
		return errors.New("attachment URL is empty")
	}

	exts, err := mime.ExtensionsByType(a.ContentType)
	if err != nil || len(exts) == 0 {
		exts = []string{filepath.Ext(a.Filename)}
	}
	ext := exts[0]

	baseName := strings.TrimSuffix(a.Filename, filepath.Ext(a.Filename))
	finalName := baseName + ext

	fullPath := filepath.Join(dir, finalName)

	resp, err := http.Get(a.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch attachment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch attachment: status %d", resp.StatusCode)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
