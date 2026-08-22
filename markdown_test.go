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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarkdown_Formatting(t *testing.T) {
	assert.Equal(t, "**bold text**", Bold("bold text"))
	assert.Equal(t, "*italic text*", Italic("italic text"))
	assert.Equal(t, "__underline text__", Underline("underline text"))
	assert.Equal(t, "~~strikethrough text~~", Strikethrough("strikethrough text"))
	assert.Equal(t, "||secret||", Spoiler("secret"))
	assert.Equal(t, "`foo()`", InlineCode("foo()"))
	assert.Equal(t, "```go\nfmt.Println()\n```", CodeBlock("go", "fmt.Println()"))
	assert.Equal(t, "```\nhello\n```", CodeBlock("", "hello"))
	assert.Equal(t, "> quote line", BlockQuote("quote line"))
	assert.Equal(t, "> line 1\n> line 2", BlockQuote("line 1\nline 2"))
	assert.Equal(t, "-# small note", Subtext("small note"))
	assert.Equal(t, "# Title", Heading(1, "Title"))
	assert.Equal(t, "## Subtitle", Heading(2, "Subtitle"))
	assert.Equal(t, "### Section", Heading(3, "Section"))
}

func TestMarkdown_Mentions(t *testing.T) {
	id := Snowflake(123456789012345678)
	assert.Equal(t, "<@123456789012345678>", UserMention(id))
	assert.Equal(t, "<@&123456789012345678>", RoleMention(id))
	assert.Equal(t, "<#123456789012345678>", ChannelMention(id))
	assert.Equal(t, "<:pepe:123456789012345678>", CustomEmojiMention("pepe", id, false))
	assert.Equal(t, "<a:pepedance:123456789012345678>", CustomEmojiMention("pepedance", id, true))
}

func TestMarkdown_Timestamps(t *testing.T) {
	tm := time.Unix(1620000000, 0).UTC()
	assert.Equal(t, "<t:1620000000>", FormatTimestamp(tm))
	assert.Equal(t, "<t:1620000000:t>", FormatTimestamp(tm, TimestampStyleShortTime))
	assert.Equal(t, "<t:1620000000:T>", FormatTimestamp(tm, TimestampStyleLongTime))
	assert.Equal(t, "<t:1620000000:d>", FormatTimestamp(tm, TimestampStyleShortDate))
	assert.Equal(t, "<t:1620000000:D>", FormatTimestamp(tm, TimestampStyleLongDate))
	assert.Equal(t, "<t:1620000000:f>", FormatTimestamp(tm, TimestampStyleShortDateTime))
	assert.Equal(t, "<t:1620000000:F>", FormatTimestamp(tm, TimestampStyleLongDateTime))
	assert.Equal(t, "<t:1620000000:R>", FormatTimestamp(tm, TimestampStyleRelative))
}
