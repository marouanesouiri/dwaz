package wumpkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColor_Parsing(t *testing.T) {
	c, err := ParseColor("#1ABC9C")
	require.NoError(t, err)
	assert.Equal(t, "#1ABC9C", c.String())

	c2, err := ParseColor("FF0000")
	require.NoError(t, err)
	assert.Equal(t, "#FF0000", c2.String())

	assert.Equal(t, Color(0x5865F2), ColorBlurple)
}
