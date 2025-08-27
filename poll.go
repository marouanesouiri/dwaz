/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import "time"

// PollLayoutType represents the layout type of a Discord poll.
// It defines the visual or structural arrangement of the poll.
// Currently, only the default layout is supported, but additional layouts may be introduced in the future.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-object-poll-object-structure
type PollLayoutType int

const (
	// PollLayoutTypeDefault represents the default layout type for a poll, with an ID of 1.
	// This is currently the only supported layout type.
	PollLayoutTypeDefault PollLayoutType = iota + 1
)

// Is returns true if the pool layout's type matches the provided layout type.
func (t PollLayoutType) Is(pollLayoutType PollLayoutType) bool {
	return t == pollLayoutType
}

// PollMedia represents the media content of a poll question or answer in Discord.
// It encapsulates the text and optional emoji associated with a poll's question or answer.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-media-object-poll-media-object-structure
type PollMedia struct {
	// Text is the text content of the poll question or answer.
	// Currently, it is always non-empty, with a max length of 300
	// characters for questions and 55 for answers. Future Discord
	// updates may allow empty text to indicate no text content.
	// Use an empty string ("") to represent no text.
	//
	// Optional:
	//  - Will be empty if no text is set.
	Text string `json:"text"`

	// Emoji is an optional partial emoji for the poll question or answer.
	// When creating a poll answer with an emoji, only the emoji's ID
	// (for custom emojis) or name (for default emojis) needs to be provided.
	//
	// Optional:
	//  - Will be nil if no emoji is set.
	Emoji *PartialEmoji `json:"emoji,omitempty"`
}

// PollAnswer represents an answer option in a Discord poll.
// It contains the answer's ID and its media content.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-answer-object-poll-answer-object-structure
type PollAnswer struct {
	// AnswerID is the ID of the answer, a number that labels each answer.
	// Currently, it is always set for poll answers, but future updates may allow it to be
	// unset. Will be nil if not provided.
	// As an implementation detail, it currently starts at 1 for the first answer and
	// increments sequentially. It is recommended not to depend on this sequence.
	//
	// Optional:
	//  - Will be nil if no ID is set.
	AnswerID *int `json:"answer_id,omitempty"`

	// PollMedia is the data of the answer.
	PollMedia PollMedia `json:"poll_media"`
}

// PollAnswerCount represents the vote count and user voting status for a specific answer in a Discord poll.
// It is part of the Poll Results Object, which contains the number of votes for each answer.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-results-object-poll-answer-count-object-structure
type PollAnswerCount struct {
	// ID is the answer_id of the poll answer, corresponding to the unique identifier of the answer option.
	ID int `json:"id"`

	// Count is the number of votes cast for this answer.
	// Note:
	//   - While a poll is in progress, this count may not be perfectly accurate due to the complexities
	//     of counting at scale. Once the poll is finalized (as indicated by PollResults.IsFinalized),
	//     the count reflects the accurate tally.
	Count int `json:"count"`

	// MeVoted indicates whether the current user has voted for this answer.
	MeVoted bool `json:"me_voted"`
}

// PollResults represents the results of a Discord poll, including whether the votes
// have been finalized and the counts for each answer option.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-results-object-poll-results-object-structure
type PollResults struct {
	// IsFinalized indicates whether the votes for the poll have been precisely counted.
	// If true, the vote counts are final; if false, the counts may still be updating.
	IsFinalized bool `json:"is_finalized"`

	// AnswerCounts is a list of PollAnswerCount objects, each containing the count
	// of votes for a specific answer option in the poll.
	AnswerCounts []PollAnswerCount `json:"answer_counts"`
}

// Poll represents a message poll sent in a channel within Discord.
//
// Reference: https://discord.com/developers/docs/resources/poll#poll-object
type Poll struct {
	// Question is the question of the poll. Only text is supported.
	Question PollMedia `json:"question"`

	// Answers is a list of each answer available in the poll.
	//
	// Note:
	//   - Currently, there is a maximum of 10 answers per poll.
	Answers []PollAnswer `json:"answers"`

	// Expiry is the time when the poll ends. Nullable to support potential
	// future non-expiring polls. Will be nil if the poll has no expiry, but currently all polls expire.
	// This is designed for future Discord updates to support never-expiring polls.
	//
	// Optional:
	//  - Will be nil if the poll has no expiry.
	Expiry *time.Time `json:"expiry"`

	// AllowMultiselect indicates whether a user can select multiple answers.
	AllowMultiselect bool `json:"allow_multiselect"`

	// LayoutType is an integer defining the visual layout of the poll.
	LayoutType PollLayoutType `json:"layout_type"`

	// Results contains the results of the poll, if available. Optional and nullable.
	//
	// Optional:
	//  - Will be nil if the poll has no results.
	Results *PollResults `json:"results"`
}
