/***********************************************************************************
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

// EmbedType represents the type of an embed.
type EmbedType string

const (
	EmbedTypeRich       EmbedType = "rich"
	EmbedTypeImage      EmbedType = "image"
	EmbedTypeVideo      EmbedType = "video"
	EmbedTypeGifv       EmbedType = "gifv"
	EmbedTypeArticle    EmbedType = "article"
	EmbedTypeLink       EmbedType = "link"
	EmbedTypePollResult EmbedType = "poll_result"
)

// Embed represents a Discord embed object.
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object
//
// Limits:
//   - The combined sum of characters in all title, description, field.name, field.value,
//     footer.text, and author.name fields across all embeds in a message must not exceed 6000.
type Embed struct {
	// Title is the title of the embed.
	//
	// Optional, max 256 characters, empty string if not set.
	Title string `json:"title,omitempty"`

	// Type is the type of the embed.
	//
	// Optional, always "rich" for webhook embeds.
	Type EmbedType `json:"type,omitempty"`

	// Description is the description text of the embed.
	//
	// Optional, max 4096 characters, empty string if not set.
	Description string `json:"description,omitempty"`

	// URL is the URL of the embed.
	//
	// Optional, empty string if not set.
	URL string `json:"url,omitempty"`

	// Timestamp is the timestamp of the embed content in ISO8601 format.
	//
	// Optional, zero value if not set.
	Timestamp *time.Time `json:"timestamp"`

	// Color is the color code of the embed (decimal integer).
	//
	// Optional, 0 if not set.
	Color Color `json:"color,omitempty"`

	// Footer contains footer information.
	//
	// Optional, may be nil if not set.
	Footer *EmbedFooter `json:"footer,omitempty"`

	// Image contains image information.
	//
	// Optional, may be nil if not set.
	Image *EmbedImage `json:"image,omitempty"`

	// Thumbnail contains thumbnail information.
	//
	// Optional, may be nil if not set.
	Thumbnail *EmbedThumbnail `json:"thumbnail,omitempty"`

	// Video contains video information.
	//
	// Optional, may be nil if not set.
	Video *EmbedVideo `json:"video,omitempty"`

	// Provider contains provider information.
	//
	// Optional, may be nil if not set.
	Provider *EmbedProvider `json:"provider,omitempty"`

	// Author contains author information.
	//
	// Optional, may be nil if not set.
	// author.name max 256 characters
	Author *EmbedAuthor `json:"author,omitempty"`

	// Fields contains an array of embed fields.
	//
	// Optional, max 25 fields.
	// field.name max 256 characters, field.value max 1024 characters
	Fields []EmbedField `json:"fields,omitempty"`
}

// Builder returns a new EmbedBuilder initialized with a copy of the current embed.
func (e *Embed) Builder() EmbedBuilder {
	return EmbedBuilder{embed: *e}
}

// EmbedFooter represents the footer object of an embed.
//
// Limits:
// - text max 2048 characters
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	// Text is the footer text.
	//
	// Always present, max 2048 characters.
	Text string `json:"text"`

	// IconURL is the URL of the footer icon.
	//
	// Optional, may be empty string if not set.
	IconURL string `json:"icon_url,omitempty"`

	// ProxyIconURL is a proxied URL of the footer icon.
	//
	// Optional, may be empty string if not set.
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedImage represents the image object of an embed.
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	// URL is the source URL of the image.
	//
	// Always present. Supports only http(s) and attachments.
	URL string `json:"url"`

	// ProxyURL is a proxied URL of the image.
	//
	// Optional, may be empty string if not set.
	ProxyURL string `json:"proxy_url,omitempty"`

	// Height is the height of the image.
	//
	// Optional, 0 if not set.
	Height int `json:"height,omitempty"`

	// Width is the width of the image.
	//
	// Optional, 0 if not set.
	Width int `json:"width,omitempty"`
}

// EmbedThumbnail represents the thumbnail object of an embed.
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	// URL is the source URL of the thumbnail.
	//
	// Always present. Supports only http(s) and attachments.
	URL string `json:"url"`

	// ProxyURL is a proxied URL of the thumbnail.
	//
	// Optional, may be empty string if not set.
	ProxyURL string `json:"proxy_url,omitempty"`

	// Height is the height of the thumbnail.
	//
	// Optional, 0 if not set.
	Height int `json:"height,omitempty"`

	// Width is the width of the thumbnail.
	//
	// Optional, 0 if not set.
	Width int `json:"width,omitempty"`
}

// EmbedVideo represents the video object of an embed.
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	// URL is the source URL of the video.
	//
	// Optional, may be empty string if not set.
	URL string `json:"url,omitempty"`

	// ProxyURL is a proxied URL of the video.
	//
	// Optional, may be empty string if not set.
	ProxyURL string `json:"proxy_url,omitempty"`

	// Height is the height of the video.
	//
	// Optional, 0 if not set.
	Height int `json:"height,omitempty"`

	// Width is the width of the video.
	//
	// Optional, 0 if not set.
	Width int `json:"width,omitempty"`
}

// EmbedProvider represents the provider object of an embed.
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	// Name is the name of the provider.
	//
	// Optional, may be empty string if not set.
	Name string `json:"name,omitempty"`

	// URL is the URL of the provider.
	//
	// Optional, may be empty string if not set.
	URL string `json:"url,omitempty"`
}

// EmbedAuthor represents the author object of an embed.
//
// Limits:
// - name max 256 characters
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	// Name is the name of the author.
	//
	// Always present, max 256 characters.
	Name string `json:"name"`

	// URL is the URL of the author.
	//
	// Optional, may be empty string if not set.
	URL string `json:"url,omitempty"`

	// IconURL is the URL of the author icon.
	//
	// Optional, may be empty string if not set.
	IconURL string `json:"icon_url,omitempty"`

	// ProxyIconURL is a proxied URL of the author icon.
	//
	// Optional, may be empty string if not set.
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedField represents a field object in an embed.
//
// Limits:
// - name max 256 characters
// - value max 1024 characters
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	// Name is the name of the field.
	//
	// Always present, max 256 characters.
	Name string `json:"name"`

	// Value is the value of the field.
	//
	// Always present, max 1024 characters.
	Value string `json:"value"`

	// Inline indicates whether this field should display inline.
	//
	// Optional, false if not set.
	Inline bool `json:"inline,omitempty"`
}

// EmbedBuilder helps build an Embed with chainable methods.
type EmbedBuilder struct {
	embed Embed
}

// NewEmbedBuilder creates a new EmbedBuilder instance.
func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{}
}

// SetTitle sets the embed title (max 256 chars).
func (b *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	if len(title) > 256 {
		title = title[:256]
	}
	b.embed.Title = title
	return b
}

// SetDescription sets the embed description (max 4096 chars).
func (b *EmbedBuilder) SetDescription(desc string) *EmbedBuilder {
	if len(desc) > 4096 {
		desc = desc[:4096]
	}
	b.embed.Description = desc
	return b
}

// SetURL sets the embed SetURL.
func (b *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	b.embed.URL = url
	return b
}

// SetTimestamp sets the embed timestamp.
func (b *EmbedBuilder) SetTimestamp(t time.Time) *EmbedBuilder {
	b.embed.Timestamp = &t
	return b
}

// SetColor sets the embed color.
func (b *EmbedBuilder) SetColor(color Color) *EmbedBuilder {
	b.embed.Color = color
	return b
}

// SetFooter sets the embed footer text and optional icon URL.
func (b *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	if len(text) > 2048 {
		text = text[:2048]
	}
	b.embed.Footer = &EmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return b
}

// SetImage sets the embed image URL.
func (b *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	b.embed.Image = &EmbedImage{URL: url}
	return b
}

// SetThumbnail sets the embed thumbnail URL.
func (b *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &EmbedThumbnail{URL: url}
	return b
}

// SetAuthor sets the embed author name and optional URL/icon.
func (b *EmbedBuilder) SetAuthor(name, url, iconURL string) *EmbedBuilder {
	if len(name) > 256 {
		name = name[:256]
	}
	b.embed.Author = &EmbedAuthor{
		Name:    name,
		URL:     url,
		IconURL: iconURL,
	}
	return b
}

// AddField appends a field to the embed fields slice.
func (b *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	if len(b.embed.Fields) >= 25 {
		return b
	}
	if len(name) > 256 {
		name = name[:256]
	}
	if len(value) > 1024 {
		value = value[:1024]
	}
	b.embed.Fields = append(b.embed.Fields, EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return b
}

// SetFields sets all embed fields at once.
//
// Note: This method does not enforce field limits or length constraints.
// It's recommended to use EmbedBuilder.AddField for validation.
func (e *EmbedBuilder) SetFields(fields ...EmbedField) {
	e.embed.Fields = fields
}

// RemoveField removes a field from the EmbedBuilder
func (b *EmbedBuilder) RemoveField(i int) *EmbedBuilder {
	if len(b.embed.Fields) > i {
		b.embed.Fields = append(b.embed.Fields[:i], b.embed.Fields[i+1:]...)
	}
	return b
}

// Build returns the final Embed object ready to send.
func (b *EmbedBuilder) Build() Embed {
	return b.embed
}
