/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD-3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

import (
	"fmt"
	"strings"
	"time"
)

// TimestampStyle represents the formatting style for Discord dynamic timestamps.
//
// Reference: https://discord.com/developers/docs/reference#message-formatting-timestamp-styles
type TimestampStyle string

const (
	// TimestampStyleShortTime formats as "16:20" (e.g. <t:1620000000:t>).
	TimestampStyleShortTime TimestampStyle = "t"

	// TimestampStyleLongTime formats as "16:20:30" (e.g. <t:1620000000:T>).
	TimestampStyleLongTime TimestampStyle = "T"

	// TimestampStyleShortDate formats as "20/04/2021" (e.g. <t:1620000000:d>).
	TimestampStyleShortDate TimestampStyle = "d"

	// TimestampStyleLongDate formats as "20 April 2021" (e.g. <t:1620000000:D>).
	TimestampStyleLongDate TimestampStyle = "D"

	// TimestampStyleShortDateTime formats as "20 April 2021 16:20" (e.g. <t:1620000000:f>). Default style.
	TimestampStyleShortDateTime TimestampStyle = "f"

	// TimestampStyleLongDateTime formats as "Tuesday, 20 April 2021 16:20" (e.g. <t:1620000000:F>).
	TimestampStyleLongDateTime TimestampStyle = "F"

	// TimestampStyleRelative formats as "2 months ago" or "in 5 minutes" (e.g. <t:1620000000:R>).
	TimestampStyleRelative TimestampStyle = "R"
)

// FormatTimestamp formats a time.Time into Discord's dynamic timestamp format (<t:unix:style>).
// If style is omitted, Discord defaults to ShortDateTime.
func FormatTimestamp(t time.Time, style ...TimestampStyle) string {
	unix := t.Unix()
	if len(style) > 0 && style[0] != "" {
		return fmt.Sprintf("<t:%d:%s>", unix, style[0])
	}
	return fmt.Sprintf("<t:%d>", unix)
}

// Bold wraps text in Discord bold markdown (**text**).
func Bold(text string) string {
	return "**" + text + "**"
}

// Italic wraps text in Discord italic markdown (*text*).
func Italic(text string) string {
	return "*" + text + "*"
}

// Underline wraps text in Discord underline markdown (__text__).
func Underline(text string) string {
	return "__" + text + "__"
}

// Strikethrough wraps text in Discord strikethrough markdown (~~text~~).
func Strikethrough(text string) string {
	return "~~" + text + "~~"
}

// Spoiler wraps text in Discord spoiler tags (||text||).
func Spoiler(text string) string {
	return "||" + text + "||"
}

// InlineCode wraps text in inline code backticks (`code`).
func InlineCode(code string) string {
	return "`" + code + "`"
}

// CodeBlock wraps text in a code block with an optional language identifier.
func CodeBlock(language, code string) string {
	if language != "" {
		return "```" + language + "\n" + code + "\n```"
	}
	return "```\n" + code + "\n```"
}

// BlockQuote prefixes text with a blockquote indicator (> text).
func BlockQuote(text string) string {
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = "> " + l
	}
	return strings.Join(lines, "\n")
}

// Subtext wraps text in small subtext format (-# text).
func Subtext(text string) string {
	return "-# " + text
}

// Heading formats text as a markdown heading with levels 1 (#), 2 (##), or 3 (###).
func Heading(level int, text string) string {
	level = min(3, max(1, level))
	return strings.Repeat("#", level) + " " + text
}

// FormatUserMention formats a user snowflake into a mention string (<@123456789>).
func FormatUserMention(id Snowflake) string {
	return fmt.Sprintf("<@%s>", id)
}

// UserMention formats a user snowflake into a mention string (<@123456789>).
func UserMention(id Snowflake) string {
	return FormatUserMention(id)
}

// FormatRoleMention formats a role snowflake into a mention string (<@&123456789>).
func FormatRoleMention(id Snowflake) string {
	return fmt.Sprintf("<@&%s>", id)
}

// RoleMention formats a role snowflake into a mention string (<@&123456789>).
func RoleMention(id Snowflake) string {
	return FormatRoleMention(id)
}

// FormatChannelMention formats a channel snowflake into a mention string (<#123456789>).
func FormatChannelMention(id Snowflake) string {
	return fmt.Sprintf("<#%s>", id)
}

// ChannelMention formats a channel snowflake into a mention string (<#123456789>).
func ChannelMention(id Snowflake) string {
	return FormatChannelMention(id)
}

// FormatCustomEmoji formats an emoji name and snowflake into a custom emoji string (<:name:123456789> or <a:name:123456789> if animated).
func FormatCustomEmoji(name string, id Snowflake, animated bool) string {
	if animated {
		return fmt.Sprintf("<a:%s:%s>", name, id)
	}
	return fmt.Sprintf("<:%s:%s>", name, id)
}

// CustomEmojiMention formats an emoji name and snowflake into a custom emoji string.
func CustomEmojiMention(name string, id Snowflake, animated bool) string {
	return FormatCustomEmoji(name, id, animated)
}
