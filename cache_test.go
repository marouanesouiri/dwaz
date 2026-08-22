package wumpkit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheManager_CRUD(t *testing.T) {
	ctx := context.Background()
	cm := NewInMemoryCacheManager(CacheFlagsAll, DefaultCacheConfig())

	guildID := Snowflake(1001)
	userID := Snowflake(2001)
	channelID := Snowflake(3001)
	roleID := Snowflake(4001)

	// test Guild
	guild := Guild{ID: guildID, Name: "Test Guild"}
	cm.AddGuild(ctx, guild)
	g, ok := cm.GetGuild(ctx, guildID)
	require.True(t, ok)
	assert.Equal(t, "Test Guild", g.Name)
	assert.True(t, cm.HasGuild(ctx, guildID))

	// test Channel
	ch := &TextChannel{
		GuildChannelFields: GuildChannelFields{
			ChannelFields: ChannelFields{
				ID:   channelID,
				Type: ChannelTypeGuildText,
			},
			GuildID: guildID,
			Name:    "general",
		},
	}
	cm.AddChannel(ctx, ch)
	c, ok := cm.GetChannel(ctx, channelID)
	require.True(t, ok)
	nc, ok := c.(NamedChannel)
	require.True(t, ok)
	assert.Equal(t, "general", nc.GetName())

	channels, ok := cm.GetGuildChannels(ctx, guildID)
	require.True(t, ok)
	assert.Len(t, channels, 1)

	// test Member & User
	user := User{ID: userID, Username: "testuser"}
	member := Member{ID: userID, GuildID: guildID, Nickname: "Tester"}
	cm.AddUser(ctx, user)
	cm.AddMember(ctx, member)

	u, ok := cm.GetUser(ctx, userID)
	require.True(t, ok)
	assert.Equal(t, "testuser", u.Username)

	m, ok := cm.GetMember(ctx, guildID, userID)
	require.True(t, ok)
	assert.Equal(t, "Tester", m.Nickname)

	members, ok := cm.GetGuildMembers(ctx, guildID)
	require.True(t, ok)
	assert.Len(t, members, 1)

	// test Role
	role := Role{ID: roleID, GuildID: guildID, Name: "Admin"}
	cm.AddRole(ctx, role)
	r, ok := cm.GetRole(ctx, roleID)
	require.True(t, ok)
	assert.Equal(t, "Admin", r.Name)

	roles, ok := cm.GetGuildRoles(ctx, guildID)
	require.True(t, ok)
	assert.Len(t, roles, 1)

	// test Deletion
	assert.True(t, cm.DelRole(ctx, guildID, roleID))
	assert.False(t, cm.HasRole(ctx, roleID))

	assert.True(t, cm.DelMember(ctx, guildID, userID))
	assert.False(t, cm.HasMember(ctx, guildID, userID))

	assert.True(t, cm.DelChannel(ctx, channelID))
	assert.False(t, cm.HasChannel(ctx, channelID))

	assert.True(t, cm.DelGuild(ctx, guildID))
	assert.False(t, cm.HasGuild(ctx, guildID))
}
