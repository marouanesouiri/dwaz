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

	"github.com/stretchr/testify/assert"
)

func TestReusableBuilder_Reset(t *testing.T) {
	// slashCommandBuilder reuse
	slashBuilder := NewSlashCommandBuilder("first", "first description")
	cmd1 := slashBuilder.Build()
	assert.Equal(t, "first", cmd1.Name)

	slashBuilder.Reset()
	slashBuilder.SetName("second").SetDescription("second description")
	cmd2 := slashBuilder.Build()
	assert.Equal(t, "second", cmd2.Name)
	assert.Equal(t, "second description", cmd2.Description)
	assert.Empty(t, cmd2.Options)

	// userCommandBuilder reuse
	userBuilder := NewUserCommandBuilder("User One")
	uCmd1 := userBuilder.Build()
	assert.Equal(t, "User One", uCmd1.Name)

	userBuilder.Reset()
	userBuilder.SetName("User Two")
	uCmd2 := userBuilder.Build()
	assert.Equal(t, "User Two", uCmd2.Name)

	// messageCommandBuilder reuse
	msgBuilder := NewMessageCommandBuilder("Message One")
	mCmd1 := msgBuilder.Build()
	assert.Equal(t, "Message One", mCmd1.Name)

	msgBuilder.Reset()
	msgBuilder.SetName("Message Two")
	mCmd2 := msgBuilder.Build()
	assert.Equal(t, "Message Two", mCmd2.Name)

	// embedBuilder reuse
	embedBuilder := NewEmbedBuilder().SetTitle("Embed 1")
	e1 := embedBuilder.Build()
	assert.Equal(t, "Embed 1", e1.Title)

	embedBuilder.Reset()
	embedBuilder.SetTitle("Embed 2")
	e2 := embedBuilder.Build()
	assert.Equal(t, "Embed 2", e2.Title)
}
