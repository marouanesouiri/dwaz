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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInteractionResponse_Helpers(t *testing.T) {
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/interactions/999/test-token/callback")

		_ = json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(
		context.Background(),
		WithToken("test-token"),
		WithRESTConfig(RequesterConfig{
			BaseURL: server.URL,
		}),
	)

	interaction := &ChatInputCommandInteraction{
		ApplicationCommandInteractionFields: ApplicationCommandInteractionFields{
			InteractionFields: InteractionFields{
				ID:            999,
				ApplicationID: 111,
				Token:         "test-token",
				Type:          InteractionTypeApplicationCommand,
			},
		},
	}

	ctx := context.Background()

	// 1. RespondMessage
	err := client.RespondMessage(ctx, interaction, "Hello World", RespondMessageOptions{
		Ephemeral: true,
	})
	require.NoError(t, err)
	assert.Equal(t, float64(4), receivedBody["type"])

	data, ok := receivedBody["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Hello World", data["content"])
	assert.Equal(t, float64(MessageFlagEphemeral), data["flags"])

	// 2. DeferResponse
	err = client.DeferResponse(ctx, interaction, false)
	require.NoError(t, err)
	assert.Equal(t, float64(5), receivedBody["type"])

	// 3. DeferComponentUpdate
	err = client.DeferComponentUpdate(ctx, interaction)
	require.NoError(t, err)
	assert.Equal(t, float64(6), receivedBody["type"])
}

func TestInteraction_Author(t *testing.T) {
	// 1. Guild interaction (Member populated)
	guildInteraction := &ChatInputCommandInteraction{
		ApplicationCommandInteractionFields: ApplicationCommandInteractionFields{
			Member: &ResolvedMember{
				FullMember: FullMember{
					Member: Member{
						ID:       123456789,
						GuildID:  987654321,
						Nickname: "GuildNickname",
					},
					User: User{
						ID:            123456789,
						Username:      "GuildUser",
						Discriminator: "0001",
						GlobalName:    "GlobalName",
						Avatar:        "avatarhash",
						Bot:           true,
					},
				},
			},
		},
	}

	author := guildInteraction.Author()
	require.NotNil(t, author)
	assert.Equal(t, Snowflake(123456789), author.ID)
	assert.Equal(t, "GuildUser", author.Username)
	assert.Equal(t, "0001", author.Discriminator)
	assert.Equal(t, "GlobalName", author.DisplayName())
	assert.True(t, author.Bot)
	assert.Equal(t, "<@123456789>", author.Mention())
	assert.Contains(t, author.AvatarURL(), "avatarhash")

	// 2. DM interaction (User populated)
	dmInteraction := &ChatInputCommandInteraction{
		ApplicationCommandInteractionFields: ApplicationCommandInteractionFields{
			User: &User{
				ID:            555666777,
				Username:      "DMUser",
				Discriminator: "0002",
				GlobalName:    "DMGlobal",
				Bot:           false,
			},
		},
	}

	dmAuthor := dmInteraction.Author()
	require.NotNil(t, dmAuthor)
	assert.Equal(t, Snowflake(555666777), dmAuthor.ID)
	assert.Equal(t, "DMUser", dmAuthor.Username)
	assert.Equal(t, "0002", dmAuthor.Discriminator)
	assert.Equal(t, "DMGlobal", dmAuthor.DisplayName())
	assert.False(t, dmAuthor.Bot)
	assert.Equal(t, "<@555666777>", dmAuthor.Mention())

	// 3. Empty interaction (Neither populated)
	emptyInteraction := &ChatInputCommandInteraction{}
	assert.Nil(t, emptyInteraction.Author())
}
