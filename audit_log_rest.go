/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// FetchAuditLogOptions contains query parameters for fetching guild audit logs.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
type FetchAuditLogOptions struct {
	// UserID filters the log for actions made by a user.
	UserID Snowflake `json:"user_id,omitempty"`

	// ActionType is the type of audit log event.
	ActionType AuditLogEvent `json:"action_type,omitempty"`

	// Before filters the log before a certain entry ID.
	Before Snowflake `json:"before,omitempty"`

	// After filters the log after a certain entry ID.
	After Snowflake `json:"after,omitempty"`

	// Limit is the max number of entries to return (1-100, default 50).
	Limit int `json:"limit,omitempty"`
}

// FetchGuildAuditLog returns an audit log object for the guild. Requires 'VIEW_AUDIT_LOG' permission.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
func (r *requester) FetchGuildAuditLog(ctx context.Context, guildID Snowflake, opts FetchAuditLogOptions) (AuditLog, error) {
	query := url.Values{}
	if opts.UserID != 0 {
		query.Set("user_id", opts.UserID.String())
	}
	if opts.ActionType != 0 {
		query.Set("action_type", strconv.Itoa(int(opts.ActionType)))
	}
	if opts.Before != 0 {
		query.Set("before", opts.Before.String())
	}
	if opts.After != 0 {
		query.Set("after", opts.After.String())
	}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}

	endpoint := fmt.Sprintf("/guilds/%s/audit-logs", guildID)
	if q := query.Encode(); q != "" {
		endpoint += "?" + q
	}

	return requestJSON[AuditLog](ctx, r, Request{
		Method: http.MethodGet,
		URL:    endpoint,
	})
}
