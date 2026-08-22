package wumpkit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_FollowedChannelInterface(t *testing.T) {
	fc := &FollowedChannel{
		ChannelID: Snowflake(123456789012345678),
		WebhookID: Snowflake(987654321098765432),
	}

	// verify Channel interface compliance
	var ch Channel = fc
	assert.Equal(t, Snowflake(123456789012345678), ch.GetID())
	assert.Equal(t, ChannelTypeGuildAnnouncement, ch.GetType())
	assert.Equal(t, "<#123456789012345678>", ch.Mention())
	assert.Equal(t, "<#123456789012345678>", ch.String())

	data, err := ch.MarshalJSON()
	require.NoError(t, err)
	assert.True(t, json.Valid(data))
}
