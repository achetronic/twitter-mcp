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

	"github.com/mark3labs/mcp-go/mcp"
)

// HandleToolPostTweet handles the post_tweet tool
func (tm *ToolsManager) HandleToolPostTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, _ := request.Params.Arguments["text"].(string)
	replyToID, _ := request.Params.Arguments["reply_to_id"].(string)

	tweet, err := tm.dependencies.TwitterClient.PostTweet(text, replyToID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweet)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolDeleteTweet handles the delete_tweet tool
func (tm *ToolsManager) HandleToolDeleteTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tweetID, _ := request.Params.Arguments["tweet_id"].(string)

	err := tm.dependencies.TwitterClient.DeleteTweet(tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet deleted"}`), nil
}

// HandleToolGetTimeline handles the get_timeline tool
func (tm *ToolsManager) HandleToolGetTimeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	maxResults := 10
	if mr, ok := request.Params.Arguments["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	// First get the authenticated user's ID
	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	timeline, err := tm.dependencies.TwitterClient.GetTimeline(me.ID, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(timeline)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetMentions handles the get_mentions tool
func (tm *ToolsManager) HandleToolGetMentions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	maxResults := 10
	if mr, ok := request.Params.Arguments["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	// First get the authenticated user's ID
	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	mentions, err := tm.dependencies.TwitterClient.GetMentions(me.ID, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(mentions)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolSearchTweets handles the search_tweets tool
func (tm *ToolsManager) HandleToolSearchTweets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, _ := request.Params.Arguments["query"].(string)
	maxResults := 10
	if mr, ok := request.Params.Arguments["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	tweets, err := tm.dependencies.TwitterClient.SearchTweets(query, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweets)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetTrends handles the get_trends tool
func (tm *ToolsManager) HandleToolGetTrends(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	woeid := 1 // Worldwide by default
	if w, ok := request.Params.Arguments["woeid"].(float64); ok {
		woeid = int(w)
	}

	trends, err := tm.dependencies.TwitterClient.GetTrends(woeid)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(trends)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolSearchTopics handles the search_topics tool
func (tm *ToolsManager) HandleToolSearchTopics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	maxResults := 5
	if mr, ok := request.Params.Arguments["max_results"].(float64); ok {
		maxResults = int(mr)
		if maxResults > 20 {
			maxResults = 20
		}
	}

	// Extract topics from the request
	var topics []string
	if topicsRaw, ok := request.Params.Arguments["topics"].([]interface{}); ok {
		for _, t := range topicsRaw {
			if topic, ok := t.(string); ok {
				topics = append(topics, topic)
			}
		}
	}

	if len(topics) == 0 {
		return mcp.NewToolResultError("no topics provided"), nil
	}

	results, err := tm.dependencies.TwitterClient.GetTrendsByTopic(topics, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(results)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetTopicsHeat handles the get_topics_heat tool
func (tm *ToolsManager) HandleToolGetTopicsHeat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sampleSize := 20
	if ss, ok := request.Params.Arguments["sample_size"].(float64); ok {
		sampleSize = int(ss)
		if sampleSize > 100 {
			sampleSize = 100
		}
	}

	// Extract topics from the request
	var topics []string
	if topicsRaw, ok := request.Params.Arguments["topics"].([]interface{}); ok {
		for _, t := range topicsRaw {
			if topic, ok := t.(string); ok {
				topics = append(topics, topic)
			}
		}
	}

	if len(topics) == 0 {
		return mcp.NewToolResultError("no topics provided"), nil
	}

	heatResults, err := tm.dependencies.TwitterClient.GetTopicsHeat(topics, sampleSize)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(heatResults)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetMe handles the get_me tool
func (tm *ToolsManager) HandleToolGetMe(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(me)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolLikeTweet handles the like_tweet tool
func (tm *ToolsManager) HandleToolLikeTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tweetID, _ := request.Params.Arguments["tweet_id"].(string)

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.LikeTweet(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet liked"}`), nil
}

// HandleToolUnlikeTweet handles the unlike_tweet tool
func (tm *ToolsManager) HandleToolUnlikeTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tweetID, _ := request.Params.Arguments["tweet_id"].(string)

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.UnlikeTweet(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet unliked"}`), nil
}

// HandleToolRetweet handles the retweet tool
func (tm *ToolsManager) HandleToolRetweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tweetID, _ := request.Params.Arguments["tweet_id"].(string)

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.Retweet(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet retweeted"}`), nil
}

// HandleToolUndoRetweet handles the undo_retweet tool
func (tm *ToolsManager) HandleToolUndoRetweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tweetID, _ := request.Params.Arguments["tweet_id"].(string)

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.UndoRetweet(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Retweet removed"}`), nil
}
