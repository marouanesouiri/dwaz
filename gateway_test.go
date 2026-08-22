package wumpkit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGateway_PresenceSerialization(t *testing.T) {
	update := GatewayPresenceUpdate{
		Status: UserStatusOnline,
		AFK:    false,
		Activities: []Activity{
			{
				Name: "wumpkit bot",
				Type: ActivityTypeGame,
			},
		},
	}

	data, err := json.Marshal(update)
	require.NoError(t, err)

	var parsed struct {
		Status     string     `json:"status"`
		AFK        bool       `json:"afk"`
		Activities []Activity `json:"activities"`
	}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "online", parsed.Status)
	assert.False(t, parsed.AFK)
	require.Len(t, parsed.Activities, 1)
	assert.Equal(t, "wumpkit bot", parsed.Activities[0].Name)
}

func TestGateway_CloseCodeClassification(t *testing.T) {
	assert.True(t, isFatalCloseCode(4004), "4004 (AuthenticationFailed) should be fatal")
	assert.True(t, isFatalCloseCode(4014), "4014 (DisallowedIntents) should be fatal")
	assert.False(t, isFatalCloseCode(4000), "4000 (UnknownError) should not be fatal")

	assert.True(t, isReidentifyCloseCode(4007), "4007 (InvalidSeq) should be re-identify")
	assert.True(t, isReidentifyCloseCode(4009), "4009 (SessionTimedOut) should be re-identify")
	assert.False(t, isReidentifyCloseCode(4000), "4000 should not be re-identify")
}
