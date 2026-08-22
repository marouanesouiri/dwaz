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
)

// MessagePaginator facilitates fetching large or multi-page message histories
// with automatic cursor tracking, inspired by JDA's PaginationAction.
type MessagePaginator struct {
	client    *Client
	channelID Snowflake
	before    Snowflake
	after     Snowflake
	around    Snowflake
}

// NewMessagePaginator creates a new paginator for retrieving messages in a channel.
func (c *Client) NewMessagePaginator(channelID Snowflake) *MessagePaginator {
	return &MessagePaginator{
		client:    c,
		channelID: channelID,
	}
}

// Before sets the initial message ID cursor to fetch messages before (older than) this ID.
func (p *MessagePaginator) Before(messageID Snowflake) *MessagePaginator {
	p.before = messageID
	return p
}

// After sets the initial message ID cursor to fetch messages after (newer than) this ID.
func (p *MessagePaginator) After(messageID Snowflake) *MessagePaginator {
	p.after = messageID
	return p
}

// Next fetches the next batch of up to limit messages (1 to 100) and advances the cursor automatically.
// Returns an empty slice when the channel history has no more messages.
func (p *MessagePaginator) Next(ctx context.Context, limit int) ([]Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	opts := FetchMessagesOptions{
		Limit:  limit,
		Before: p.before,
		After:  p.after,
		Around: p.around,
	}

	msgs, err := p.client.FetchChannelMessages(ctx, p.channelID, opts)
	if err != nil {
		return nil, err
	}

	if len(msgs) > 0 {
		if p.after != 0 {
			// moving forward (newest message is at index 0)
			p.after = msgs[0].ID
		} else {
			// moving backward (oldest message is at the end of the batch)
			p.before = msgs[len(msgs)-1].ID
		}
	}

	return msgs, nil
}

// Collect sequentially fetches multiple pages until maxCount messages are accumulated or history is exhausted.
func (p *MessagePaginator) Collect(ctx context.Context, maxCount int) ([]Message, error) {
	if maxCount <= 0 {
		return nil, nil
	}

	results := make([]Message, 0, min(maxCount, 100))
	for len(results) < maxCount {
		batchSize := min(100, maxCount-len(results))

		page, err := p.Next(ctx, batchSize)
		if err != nil {
			return results, err
		}
		if len(page) == 0 {
			break
		}
		results = append(results, page...)
	}

	return results, nil
}
