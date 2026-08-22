package wumpkit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_EventHandling(t *testing.T) {
	client := New(context.Background(), WithToken("dummy_token"))
	d := client.dispatcher

	var wg sync.WaitGroup
	wg.Add(1)

	var receivedMsg string
	d.OnMessageCreate(func(ctx context.Context, e MessageCreateEvent) {
		receivedMsg = e.Message.Content
		wg.Done()
	})

	msgJSON := `{"id":"123","channel_id":"456","content":"hello wumpkit"}`
	d.Dispatch(0, "MESSAGE_CREATE", []byte(msgJSON))

	wg.Wait()

	assert.Equal(t, "hello wumpkit", receivedMsg)
}

func TestDispatcher_CachePopulatedWithoutUserHandlers(t *testing.T) {
	client := New(context.Background(), WithToken("dummy_token"))

	// dispatch GUILD_CREATE without registering any OnGuildCreate user handlers
	guildCreateJSON := `{
		"id": "999",
		"name": "Cached Guild",
		"channels": [
			{"id": "1001", "type": 0, "name": "general"}
		],
		"roles": [
			{"id": "2001", "name": "Member"}
		]
	}`

	client.dispatcher.Dispatch(0, "GUILD_CREATE", []byte(guildCreateJSON))

	// give the dispatch goroutine a moment to finish
	var g Guild
	var found bool
	for range 50 {
		g, found = client.GetGuild(context.Background(), 999)
		if found {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, found, "Guild should be cached even with no user handlers registered")
	assert.Equal(t, "Cached Guild", g.Name)

	ch, chFound := client.GetChannel(context.Background(), 1001)
	assert.True(t, chFound, "Channel should be cached even with no user handlers registered")
	assert.Equal(t, Snowflake(1001), ch.GetID())

	role, roleFound := client.GetRole(context.Background(), 2001)
	assert.True(t, roleFound, "Role should be cached even with no user handlers registered")
	assert.Equal(t, "Member", role.Name)
}
