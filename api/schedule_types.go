// SPDX-FileCopyrightText: 2026 Alby Hernández <hola@achetronic.com>
// SPDX-License-Identifier: Apache-2.0

package api

import "time"

// ScheduledTweetStatus represents the status of a scheduled tweet
type ScheduledTweetStatus string

const (
	ScheduledTweetStatusPending   ScheduledTweetStatus = "pending"
	ScheduledTweetStatusReviewed  ScheduledTweetStatus = "reviewed"
	ScheduledTweetStatusPublished ScheduledTweetStatus = "published"
	ScheduledTweetStatusFailed    ScheduledTweetStatus = "failed"
)

// ScheduledTweetType represents the type of a scheduled tweet
type ScheduledTweetType string

const (
	ScheduledTweetTypeTweet  ScheduledTweetType = "tweet"
	ScheduledTweetTypeThread ScheduledTweetType = "thread"
)

// ScheduledTweet represents a tweet or thread scheduled for publishing
type ScheduledTweet struct {
	ID          string               `yaml:"id"`
	Type        ScheduledTweetType   `yaml:"type"`
	Content     []string             `yaml:"content"`
	ScheduledAt time.Time            `yaml:"scheduled_at"`
	Reviewed    bool                 `yaml:"reviewed"`
	Status      ScheduledTweetStatus `yaml:"status"`
	CreatedAt   time.Time            `yaml:"created_at"`
	PublishedAt *time.Time           `yaml:"published_at,omitempty"`
	FailReason  string               `yaml:"fail_reason,omitempty"`
}

// ScheduleStore represents the full persistence file
type ScheduleStore struct {
	ScheduledTweets []ScheduledTweet `yaml:"scheduled_tweets"`
}
