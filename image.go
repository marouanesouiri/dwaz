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
	"strconv"
)

// NOTE:
// Image format enums are duplicated per endpoint section because
// not all Discord image endpoints support the same set of image formats.
// This duplication ensures type safety and clear intent per endpoint.
// Also, the CDN base URLs are defined as constants for maintainability,
// since some media like Sticker GIFs do NOT use the CDN base URL and
// instead use a different base URL for serving content.

const (
	ImageBaseURL = "https://cdn.discordapp.com/"
	MediaBaseURL = "https://media.discordapp.net/"
)

type ImageSize int

const (
	ImageSize16   ImageSize = 16
	ImageSize32   ImageSize = 32
	ImageSize64   ImageSize = 64
	ImageSize128  ImageSize = 128
	ImageSize256  ImageSize = 256
	ImageSize512  ImageSize = 512
	ImageSize1024 ImageSize = 1024
	ImageSize2048 ImageSize = 2048
	ImageSize4096 ImageSize = 4096
)

/***********************
 *        Emoji        *
 ***********************/

type EmojiFormat string

const (
	EmojiFormatPNG  EmojiFormat = ".png"
	EmojiFormatJPEG EmojiFormat = ".jpeg"
	EmojiFormatWebP EmojiFormat = ".webp"
	EmojiFormatGIF  EmojiFormat = ".gif"
	EmojiFormatAVIF EmojiFormat = ".avif"
)

func EmojiURL(emojiID Snowflake, format EmojiFormat, size ImageSize) string {
	return ImageBaseURL + "emojis/" + emojiID.String() + string(format) + "?size=" + strconv.Itoa(int(size))
}

/***********************
 *   	  Guild        *
 ***********************/

type GuildIconFormat string

const (
	GuildIconFormatPNG  GuildIconFormat = ".png"
	GuildIconFormatJPEG GuildIconFormat = ".jpeg"
	GuildIconFormatWebP GuildIconFormat = ".webp"
	GuildIconFormatGIF  GuildIconFormat = ".gif"
)

func GuildIconURL(guildID Snowflake, iconHash string, format GuildIconFormat, size ImageSize) string {
	if format == GuildIconFormatGIF && (len(iconHash) < 2 || iconHash[:2] != "a_") {
		format = GuildIconFormatPNG
	}

	url := ImageBaseURL + "icons/" + guildID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildIconFormatWebP && len(iconHash) >= 2 && iconHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type GuildSplashFormat string

const (
	GuildSplashFormatPNG  GuildSplashFormat = ".png"
	GuildSplashFormatJPEG GuildSplashFormat = ".jpeg"
	GuildSplashFormatWebP GuildSplashFormat = ".webp"
)

func GuildSplashURL(guildID Snowflake, splashHash string, format GuildSplashFormat, size ImageSize) string {
	return ImageBaseURL + "splashes/" + guildID.String() + "/" + splashHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

type GuildBannerFormat string

const (
	GuildBannerFormatPNG  GuildBannerFormat = ".png"
	GuildBannerFormatJPEG GuildBannerFormat = ".jpeg"
	GuildBannerFormatWebP GuildBannerFormat = ".webp"
	GuildBannerFormatGIF  GuildBannerFormat = ".gif"
)

func GuildBannerURL(guildID Snowflake, bannerHash string, format GuildBannerFormat, size ImageSize) string {
	if format == GuildBannerFormatGIF && (len(bannerHash) < 2 || bannerHash[:2] != "a_") {
		format = GuildBannerFormatPNG
	}

	url := ImageBaseURL + "banners/" + guildID.String() + "/" + bannerHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildBannerFormatWebP && len(bannerHash) >= 2 && bannerHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

/***********************
 *         User        *
 ***********************/

// DefaultUserAvatarURL returns the default user avatar URL.
// Size param ignored, fixed size only.
func DefaultUserAvatarURL(index int) string {
	return ImageBaseURL + "embed/avatars/" + strconv.Itoa(index) + ".png"
}

type UserAvatarFormat string

const (
	UserAvatarFormatPNG  UserAvatarFormat = ".png"
	UserAvatarFormatJPEG UserAvatarFormat = ".jpeg"
	UserAvatarFormatWebP UserAvatarFormat = ".webp"
	UserAvatarFormatGIF  UserAvatarFormat = ".gif"
)

func UserAvatarURL(userID Snowflake, avatarHash string, format UserAvatarFormat, size ImageSize) string {
	if format == UserAvatarFormatGIF && (len(avatarHash) < 2 || avatarHash[:2] != "a_") {
		format = UserAvatarFormatPNG
	}

	url := ImageBaseURL + "avatars/" + userID.String() + "/" + avatarHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == UserAvatarFormatWebP && len(avatarHash) >= 2 && avatarHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type UserBannerFormat string

const (
	UserBannerFormatPNG  UserBannerFormat = ".png"
	UserBannerFormatJPEG UserBannerFormat = ".jpeg"
	UserBannerFormatWebP UserBannerFormat = ".webp"
	UserBannerFormatGIF  UserBannerFormat = ".gif"
)

func UserBannerURL(userID Snowflake, bannerHash string, format UserBannerFormat, size ImageSize) string {
	if format == UserBannerFormatGIF && (len(bannerHash) < 2 || bannerHash[:2] != "a_") {
		format = UserBannerFormatPNG
	}

	url := ImageBaseURL + "banners/" + userID.String() + "/" + bannerHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == UserBannerFormatWebP && len(bannerHash) >= 2 && bannerHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

/***********************
 * 	   Application     *
 ***********************/

type ApplicationIconFormat string

const (
	ApplicationIconFormatPNG  ApplicationIconFormat = ".png"
	ApplicationIconFormatJPEG ApplicationIconFormat = ".jpeg"
	ApplicationIconFormatWebP ApplicationIconFormat = ".webp"
)

func ApplicationIconURL(appID Snowflake, iconHash string, format ApplicationIconFormat, size ImageSize) string {
	return ImageBaseURL + "app-icons/" + appID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

type ApplicationCoverFormat string

const (
	ApplicationCoverFormatPNG  ApplicationCoverFormat = ".png"
	ApplicationCoverFormatJPEG ApplicationCoverFormat = ".jpeg"
	ApplicationCoverFormatWebP ApplicationCoverFormat = ".webp"
)

func ApplicationCoverURL(appID Snowflake, coverHash string, format ApplicationCoverFormat, size ImageSize) string {
	return ImageBaseURL + "app-icons/" + appID.String() + "/" + coverHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

/***********************
 *   	Sticker  	   *
 ***********************/

type StickerFormat string

const (
	StickerFormatPNG    StickerFormat = ".png"
	StickerFormatGIF    StickerFormat = ".gif"
	StickerFormatLottie StickerFormat = ".json"
)

// Stickers with GIF format are served from MediaBaseURL, not CDN base.
func StickerURL(stickerID Snowflake, format StickerFormat) string {
	base := ImageBaseURL + "stickers/" + stickerID.String()
	if format == StickerFormatGIF {
		base = MediaBaseURL + "stickers/" + stickerID.String()
	}
	return base + string(format)
}

type StickerPackBannerFormat string

const (
	StickerPackBannerFormatPNG    StickerPackBannerFormat = ".png"
	StickerPackBannerFormatGIF    StickerPackBannerFormat = ".gif"
	StickerPackBannerFormatLottie StickerPackBannerFormat = ".json"
)

func StickerPackBannerURL(stickerPackBannerAssetID Snowflake, format StickerPackBannerFormat, size ImageSize) string {
	return ImageBaseURL + "app-assets/710982414301790216/store/" + stickerPackBannerAssetID.String() + "/" + string(format) + "?size=" + strconv.Itoa(int(size))
}

/***********************
 *	  Guild Member     *
 ***********************/

type GuildMemberAvatarFormat string

const (
	GuildMemberAvatarFormatPNG  GuildMemberAvatarFormat = ".png"
	GuildMemberAvatarFormatJPEG GuildMemberAvatarFormat = ".jpeg"
	GuildMemberAvatarFormatWebP GuildMemberAvatarFormat = ".webp"
	GuildMemberAvatarFormatGIF  GuildMemberAvatarFormat = ".gif"
)

func GuildMemberAvatarURL(guildID, userID Snowflake, avatarHash string, format GuildMemberAvatarFormat, size ImageSize) string {
	if format == GuildMemberAvatarFormatGIF && (len(avatarHash) < 2 || avatarHash[:2] != "a_") {
		format = GuildMemberAvatarFormatPNG
	}

	url := ImageBaseURL + "guilds/" + guildID.String() + "/users/" + userID.String() + "/avatars/" + avatarHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildMemberAvatarFormatWebP && len(avatarHash) >= 2 && avatarHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type GuildMemberBannerFormat string

const (
	GuildMemberBannerFormatPNG  GuildMemberBannerFormat = ".png"
	GuildMemberBannerFormatJPEG GuildMemberBannerFormat = ".jpeg"
	GuildMemberBannerFormatWebP GuildMemberBannerFormat = ".webp"
	GuildMemberBannerFormatGIF  GuildMemberBannerFormat = ".gif"
)

func GuildMemberBannerURL(guildID, userID Snowflake, bannerHash string, format GuildMemberBannerFormat, size ImageSize) string {
	if format == GuildMemberBannerFormatGIF && (len(bannerHash) < 2 || bannerHash[:2] != "a_") {
		format = GuildMemberBannerFormatPNG
	}

	url := ImageBaseURL + "guilds/" + guildID.String() + "/users/" + userID.String() + "/banners/" + bannerHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildMemberBannerFormatWebP && len(bannerHash) >= 2 && bannerHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

/***********************
 *	    Guild Role     *
 ***********************/

type RoleIconFormat string

const (
	RoleIconFormatPNG  RoleIconFormat = ".png"
	RoleIconFormatJPEG RoleIconFormat = ".jpeg"
	RoleIconFormatWebP RoleIconFormat = ".webp"
)

func RoleIconURL(roleID Snowflake, iconHash string, format RoleIconFormat, size ImageSize) string {
	return ImageBaseURL + "role-icons/" + roleID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))
}
