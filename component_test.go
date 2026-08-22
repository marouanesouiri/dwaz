package wumpkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponent_Builders(t *testing.T) {
	btn := NewButtonBuilder().
		SetCustomID("test_btn").
		SetLabel("Click me").
		SetStyle(ButtonStylePrimary).
		Build()

	row := NewActionRowBuilder().
		SetComponent(btn).
		Build()

	require.Len(t, row.Components, 1)
	assert.Equal(t, ComponentTypeActionRow, row.Type)
}
