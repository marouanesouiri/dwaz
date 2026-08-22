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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessagePaginator(t *testing.T) {
	// generate mock messages 1..150
	var allMessages []Message
	for i := range 150 {
		allMessages = append(allMessages, Message{
			ID:      Snowflake(150 - i),
			Content: fmt.Sprintf("Message %d", 150-i),
		})
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		limitStr := r.URL.Query().Get("limit")
		beforeStr := r.URL.Query().Get("before")

		limit := 50
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}

		var before Snowflake
		if beforeStr != "" {
			b, _ := strconv.ParseUint(beforeStr, 10, 64)
			before = Snowflake(b)
		}

		var page []Message
		for _, m := range allMessages {
			if before == 0 || m.ID < before {
				page = append(page, m)
				if len(page) >= limit {
					break
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(page)
	}))
	defer server.Close()

	client := New(
		context.Background(),
		WithToken("test-token"),
		WithRESTConfig(RequesterConfig{
			BaseURL: server.URL,
		}),
	)

	ctx := context.Background()

	// collect 120 messages (should fetch 2 pages: 100 then 20)
	paginator := client.NewMessagePaginator(12345)
	messages, err := paginator.Collect(ctx, 120)
	require.NoError(t, err)
	assert.Len(t, messages, 120)
	assert.Equal(t, Snowflake(150), messages[0].ID)
	assert.Equal(t, Snowflake(31), messages[119].ID)
}
