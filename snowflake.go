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
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

/***********************
 * Constants & Types   *
 ***********************/

// Discord epoch start time: January 1, 2015 UTC
const discordEpoch = 1420070400000

/***********************
 *     Snowflake       *
 ***********************/

// Snowflake is a Discord unique identifier.
type Snowflake uint64

var (
	_ json.Marshaler   = (*Snowflake)(nil)
	_ json.Unmarshaler = (*Snowflake)(nil)
)

func (s *Snowflake) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	str, err := strconv.Unquote(string(buf))
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	*s = Snowflake(id)
	return nil
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(s), 10) + `"`), nil
}

// UnSet returns true if the Snowflake is zero (unset).
func (s Snowflake) UnSet() bool {
	return s == 0
}

// String returns the Snowflake as string.
func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

// Timestamp returns the creation time of the snowflake as time.Time.
func (s Snowflake) Timestamp() time.Time {
	ms := (uint64(s) >> 22) + discordEpoch
	return time.UnixMilli(int64(ms))
}

// WorkerID extracts the internal Discord worker ID from the snowflake.
func (s Snowflake) WorkerID() uint64 {
	return (uint64(s) & 0x3E0000) >> 17
}

// ProcessID extracts the internal Discord process ID from the snowflake.
func (s Snowflake) ProcessID() uint64 {
	return (uint64(s) & 0x1F000) >> 12
}

// Sequence extracts the sequence number (increment part) of the snowflake.
func (s Snowflake) Sequence() uint64 {
	return uint64(s) & 0xFFF
}

// GuildShardID calculates the shard ID for a guild snowflake given the total number of shards.
//
// Formula: (guild_id >> 22) % num_shards
// Reference: https://discord.com/developers/docs/events/gateway#sharding-sharding-formula
func GuildShardID(guildID Snowflake, totalShards int) int {
	if totalShards <= 1 {
		return 0
	}
	return int((uint64(guildID) >> 22) % uint64(totalShards))
}

// GuildShardID returns the shard ID for this guild Snowflake given the total number of shards.
func (s Snowflake) GuildShardID(totalShards int) int {
	return GuildShardID(s, totalShards)
}

/***********************
 * Utilities           *
 ***********************/

// ParseSnowflake parses a string into a Snowflake.
func ParseSnowflake(id string) (Snowflake, error) {
	v, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid snowflake: %w", err)
	}
	return Snowflake(v), nil
}
