package wumpkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum_IsMethods(t *testing.T) {
	cmdType := ApplicationCommandTypeChatInput
	assert.True(t, cmdType.Is(ApplicationCommandTypeChatInput))
	assert.False(t, cmdType.Is(ApplicationCommandTypeUser))

	locale := LocaleEnglishUS
	assert.True(t, locale.Is(LocaleEnglishUS))
	assert.False(t, locale.Is(LocaleFrench))

	status := UserStatusOnline
	assert.True(t, status.Is(UserStatusOnline))
	assert.False(t, status.Is(UserStatusDoNotDisturb))

	event := GatewayEventMessageCreate
	assert.True(t, event.Is(GatewayEventMessageCreate))
	assert.False(t, event.Is(GatewayEventMessageDelete))

	comp := GatewayCompressionZlibStream
	assert.True(t, comp.Is(GatewayCompressionZlibStream))
	assert.False(t, comp.Is(GatewayCompressionNone))

	field := ChannelInfoFieldStatus
	assert.True(t, field.Is(ChannelInfoFieldStatus))
	assert.False(t, field.Is(ChannelInfoFieldVoiceStartTime))
}
