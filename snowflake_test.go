package wumpkit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnowflake_GuildShardID(t *testing.T) {
	guildID := Snowflake(41771983423143937)

	assert.Equal(t, 0, GuildShardID(guildID, 1))

	expectedShard1 := int((uint64(guildID) >> 22) % 2)
	assert.Equal(t, expectedShard1, GuildShardID(guildID, 2))
}

func TestSnowflake_Utilities(t *testing.T) {
	// parse Snowflake
	sf, err := ParseSnowflake("175928847299117063")
	require.NoError(t, err)
	assert.Equal(t, Snowflake(175928847299117063), sf)
	assert.Equal(t, "175928847299117063", sf.String())

	// invalid Snowflake
	_, err = ParseSnowflake("invalid")
	assert.Error(t, err)

	// timestamp
	ts := sf.Timestamp()
	assert.True(t, ts.After(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)))
}
