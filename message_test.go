package wumpkit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_EmbedBuilder(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	builder := NewEmbedBuilder().
		SetTitle("Test Embed").
		SetDescription("This is a test description").
		SetColor(ColorBlurple).
		SetTimestamp(now)

	embed := builder.Build()
	assert.Equal(t, "Test Embed", embed.Title)
	assert.Equal(t, "This is a test description", embed.Description)
	assert.Equal(t, ColorBlurple, embed.Color)
	assert.Equal(t, &now, embed.Timestamp)
}
