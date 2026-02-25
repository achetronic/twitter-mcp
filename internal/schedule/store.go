// Copyright 2024 Alby Hernández
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schedule

import (
	"fmt"
	"os"
	"sync"
	"time"
	"twitter-mcp/api"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Store manages persistence of scheduled tweets
type Store struct {
	mu       sync.Mutex
	filepath string
	data     api.ScheduleStore
}

// NewStore creates a new Store and loads existing data from disk
func NewStore(filepath string) (*Store, error) {
	s := &Store{filepath: filepath}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// load reads the YAML file from disk into memory
func (s *Store) load() error {
	s.data = api.ScheduleStore{}

	fileBytes, err := os.ReadFile(s.filepath)
	if os.IsNotExist(err) {
		// File doesn't exist yet, start with empty store
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read schedule file: %w", err)
	}

	if err := yaml.Unmarshal(fileBytes, &s.data); err != nil {
		return fmt.Errorf("failed to parse schedule file: %w", err)
	}

	return nil
}

// save writes the current data to disk
func (s *Store) save() error {
	fileBytes, err := yaml.Marshal(&s.data)
	if err != nil {
		return fmt.Errorf("failed to marshal schedule data: %w", err)
	}

	if err := os.WriteFile(s.filepath, fileBytes, 0644); err != nil {
		return fmt.Errorf("failed to write schedule file: %w", err)
	}

	return nil
}

// Add adds a new scheduled tweet to the store
func (s *Store) Add(tweetType api.ScheduledTweetType, content []string, scheduledAt time.Time) (*api.ScheduledTweet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tweet := api.ScheduledTweet{
		ID:          uuid.New().String(),
		Type:        tweetType,
		Content:     content,
		ScheduledAt: scheduledAt,
		Reviewed:    false,
		Status:      api.ScheduledTweetStatusPending,
		CreatedAt:   time.Now().UTC(),
	}

	s.data.ScheduledTweets = append(s.data.ScheduledTweets, tweet)

	if err := s.save(); err != nil {
		return nil, err
	}

	return &tweet, nil
}

// List returns all scheduled tweets, optionally filtered by status
func (s *Store) List(status api.ScheduledTweetStatus) []api.ScheduledTweet {
	s.mu.Lock()
	defer s.mu.Unlock()

	if status == "" {
		return s.data.ScheduledTweets
	}

	var result []api.ScheduledTweet
	for _, t := range s.data.ScheduledTweets {
		if t.Status == status {
			result = append(result, t)
		}
	}
	return result
}

// GetByID returns a scheduled tweet by ID
func (s *Store) GetByID(id string) (*api.ScheduledTweet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.data.ScheduledTweets {
		if t.ID == id {
			copy := t
			return &copy, nil
		}
	}

	return nil, fmt.Errorf("scheduled tweet with id '%s' not found", id)
}

// Update modifies an existing scheduled tweet
func (s *Store) Update(id string, fn func(*api.ScheduledTweet)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.data.ScheduledTweets {
		if t.ID == id {
			fn(&s.data.ScheduledTweets[i])
			return s.save()
		}
	}

	return fmt.Errorf("scheduled tweet with id '%s' not found", id)
}

// Delete removes a scheduled tweet by ID
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.data.ScheduledTweets {
		if t.ID == id {
			s.data.ScheduledTweets = append(s.data.ScheduledTweets[:i], s.data.ScheduledTweets[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("scheduled tweet with id '%s' not found", id)
}

// GetPublishable returns tweets that are reviewed, scheduled_at is past,
// and no other tweet was published within minHoursSinceLast hours
func (s *Store) GetPublishable(minHoursSinceLast int) []api.ScheduledTweet {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()

	// Find the last published tweet
	var lastPublishedAt time.Time
	for _, t := range s.data.ScheduledTweets {
		if t.Status == api.ScheduledTweetStatusPublished && t.PublishedAt != nil {
			if t.PublishedAt.After(lastPublishedAt) {
				lastPublishedAt = *t.PublishedAt
			}
		}
	}

	// Check if enough time has passed since last publish
	if minHoursSinceLast > 0 && !lastPublishedAt.IsZero() {
		minGap := time.Duration(minHoursSinceLast) * time.Hour
		if now.Sub(lastPublishedAt) < minGap {
			return nil
		}
	}

	// Return reviewed tweets whose scheduled time has passed
	var result []api.ScheduledTweet
	for _, t := range s.data.ScheduledTweets {
		if t.Reviewed && t.Status == api.ScheduledTweetStatusReviewed && t.ScheduledAt.Before(now) {
			result = append(result, t)
		}
	}

	return result
}
