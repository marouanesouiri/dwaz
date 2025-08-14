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

// StickerType defines whether the sticker is a standard or guild sticker.
//
// Reference: https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
type StickerType int

const (
	// StickerTypeStandard represents an official sticker in a pack.
	StickerTypeStandard StickerType = iota + 1

	// StickerTypeGuild represents a sticker uploaded to a guild.
	StickerTypeGuild
)

// StickerFormatType defines the format of a sticker's image.
//
// Reference: https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
type StickerFormatType int

const (
	// StickerFormatTypePNG represents a PNG format sticker.
	StickerFormatTypePNG StickerFormatType = iota + 1

	// StickerFormatTypeAPNG represents an APNG format sticker.
	StickerFormatTypeAPNG

	// StickerFormatTypeLottie represents a Lottie format sticker.
	StickerFormatTypeLottie

	// StickerFormatTypeGIF represents a GIF format sticker.
	StickerFormatTypeGIF
)

// Sticker represents a sticker that can be sent in messages.
//
// Reference: https://discord.com/developers/docs/resources/sticker#sticker-object
type Sticker struct {
	// Unique ID of the sticker.
	ID Snowflake `json:"id"`

	// ID of the pack for standard stickers.
	PackID Snowflake `json:"pack_id,omitempty"`

	// Name of the sticker.
	Name string `json:"name"`

	// Description of the sticker (optional).
	Description string `json:"description,omitempty"`

	// Autocomplete/suggestion tags (max 200 characters).
	Tags string `json:"tags"`

	// Type of the sticker (standard or guild).
	Type StickerType `json:"type"`

	// Format type of the sticker.
	FormatType StickerFormatType `json:"format_type"`

	// Whether the guild sticker is available for use.
	Available *bool `json:"available,omitempty"`

	// ID of the guild that owns this sticker.
	GuildID Snowflake `json:"guild_id,omitempty"`

	// The user that uploaded the guild sticker.
	User *User `json:"user,omitempty"`

	// Sort order of the standard sticker in its pack.
	SortValue *int `json:"sort_value,omitempty"`
}

// URL returns the URL to the sticker's image.
func (s *Sticker) URL() string {
	var format ImageFormat
	switch s.FormatType {
	case StickerFormatTypeLottie:
		format = ImageFormatLottie
	case StickerFormatTypeGIF:
		format = ImageFormatGIF
	default:
		format = ImageFormatPNG
	}
	return StickerURL(s.ID, format)
}

// URLWith returns the URL to the sticker's image with custom format.
func (s *Sticker) URLWith(format ImageFormat) string {
	return StickerURL(s.ID, format)
}

// CreatedAt returns the time when this sticker was created.
func (s *Sticker) CreatedAt() time.Time {
	return s.ID.Timestamp()
}

// StickerPack represents a pack of standard stickers.
//
// Reference: https://discord.com/developers/docs/resources/sticker#sticker-pack-object
type StickerPack struct {
	// Unique ID of the sticker pack.
	ID Snowflake `json:"id"`

	// Array of stickers in the pack.
	Stickers []Sticker `json:"stickers"`

	// Name of the sticker pack.
	Name string `json:"name"`

	// SKU ID of the pack.
	SkuID Snowflake `json:"sku_id"`

	// ID of a sticker shown as icon.
	CoverStickerID Snowflake `json:"cover_sticker_id,omitempty"`

	// Description of the sticker pack.
	Description string `json:"description"`

	// Banner image ID.
	BannerAssetID Snowflake `json:"banner_asset_id,omitempty"`
}

// BannerURL returns the banner URL in PNG format, or empty string if none.
func (p *StickerPack) BannerURL() string {
	if p.BannerAssetID != 0 {
		return StickerPackBannerURL(p.BannerAssetID, ImageFormatPNG, ImageSize2048)
	}
	return ""
}

// BannerURLWith returns the banner URL with a custom format and size.
func (p *StickerPack) BannerURLWith(format ImageFormat, size ImageSize) string {
	if p.BannerAssetID != 0 {
		return StickerPackBannerURL(p.BannerAssetID, format, size)
	}
	return ""
}
