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

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"twitter-mcp/api"

	"github.com/mark3labs/mcp-go/mcp"
)

// HandleToolScheduleTweet handles the schedule_tweet tool
func (tm *ToolsManager) HandleToolScheduleTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)

	tweetType := api.ScheduledTweetType(getString(args, "type", "tweet"))
	scheduledAtStr := getString(args, "scheduled_at", "")
	content := getStringSlice(args, "content")

	if len(content) == 0 {
		return mcp.NewToolResultError("content is required"), nil
	}

	if scheduledAtStr == "" {
		return mcp.NewToolResultError("scheduled_at is required"), nil
	}

	scheduledAt, err := time.Parse(time.RFC3339, scheduledAtStr)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid scheduled_at format, use RFC3339 (e.g. 2026-02-25T10:00:00Z): %s", err.Error())), nil
	}

	tweet, err := tm.dependencies.ScheduleStore.Add(tweetType, content, scheduledAt)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweet)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolScheduleUpdate handles the schedule_update tool
func (tm *ToolsManager) HandleToolScheduleUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	id := getString(args, "id", "")

	if id == "" {
		return mcp.NewToolResultError("id is required"), nil
	}

	err := tm.dependencies.ScheduleStore.Update(id, func(t *api.ScheduledTweet) {
		if v := getString(args, "type", ""); v != "" {
			t.Type = api.ScheduledTweetType(v)
		}
		if v := getStringSlice(args, "content"); len(v) > 0 {
			t.Content = v
		}
		if v := getString(args, "scheduled_at", ""); v != "" {
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				t.ScheduledAt = parsed
			}
		}
		if v, ok := args["reviewed"].(bool); ok {
			t.Reviewed = v
			if v {
				t.Status = api.ScheduledTweetStatusReviewed
			} else {
				t.Status = api.ScheduledTweetStatusPending
			}
		}
	})

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tweet, _ := tm.dependencies.ScheduleStore.GetByID(id)
	result, _ := json.Marshal(tweet)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolScheduleDelete handles the schedule_delete tool
func (tm *ToolsManager) HandleToolScheduleDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	id := getString(args, "id", "")

	if id == "" {
		return mcp.NewToolResultError("id is required"), nil
	}

	if err := tm.dependencies.ScheduleStore.Delete(id); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Scheduled tweet deleted"}`), nil
}

// HandleToolScheduleList handles the schedule_list tool
func (tm *ToolsManager) HandleToolScheduleList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	status := api.ScheduledTweetStatus(getString(args, "status", ""))

	tweets := tm.dependencies.ScheduleStore.List(status)

	result, _ := json.Marshal(tweets)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolScheduleGetPublishable handles the schedule_get_publishable tool
func (tm *ToolsManager) HandleToolScheduleGetPublishable(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	minHours := getInt(args, "min_hours_since_last", 1)

	tweets := tm.dependencies.ScheduleStore.GetPublishable(minHours)

	result, _ := json.Marshal(tweets)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolSchedulePublish handles the schedule_publish tool
func (tm *ToolsManager) HandleToolSchedulePublish(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	id := getString(args, "id", "")

	if id == "" {
		return mcp.NewToolResultError("id is required"), nil
	}

	tweet, err := tm.dependencies.ScheduleStore.GetByID(id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Publish all content items (tweet or thread)
	var lastTweetID string
	for _, text := range tweet.Content {
		posted, err := tm.dependencies.TwitterClient.PostTweet(text, lastTweetID)
		if err != nil {
			// Mark as failed
			tm.dependencies.ScheduleStore.Update(id, func(t *api.ScheduledTweet) {
				t.Status = api.ScheduledTweetStatusFailed
				t.FailReason = err.Error()
			})
			return mcp.NewToolResultError(fmt.Sprintf("failed to publish tweet: %s", err.Error())), nil
		}
		lastTweetID = posted.ID
	}

	// Mark as published
	now := time.Now().UTC()
	tm.dependencies.ScheduleStore.Update(id, func(t *api.ScheduledTweet) {
		t.Status = api.ScheduledTweetStatusPublished
		t.PublishedAt = &now
	})

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet published successfully"}`), nil
}
