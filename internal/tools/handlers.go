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
	args := getArgs(request)
	text := getString(args, "text", "")
	replyToID := getString(args, "reply_to_id", "")

	tweet, err := tm.dependencies.TwitterClient.PostTweet(text, replyToID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweet)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolDeleteTweet handles the delete_tweet tool
func (tm *ToolsManager) HandleToolDeleteTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

	err := tm.dependencies.TwitterClient.DeleteTweet(tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet deleted"}`), nil
}

// HandleToolGetTimeline handles the get_timeline tool
func (tm *ToolsManager) HandleToolGetTimeline(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	maxResults := getInt(args, "max_results", 10)

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
	args := getArgs(request)
	maxResults := getInt(args, "max_results", 10)

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
	args := getArgs(request)
	query := getString(args, "query", "")
	maxResults := getInt(args, "max_results", 10)

	tweets, err := tm.dependencies.TwitterClient.SearchTweets(query, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweets)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetTrends handles the get_trends tool
func (tm *ToolsManager) HandleToolGetTrends(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	woeid := getInt(args, "woeid", 1)

	trends, err := tm.dependencies.TwitterClient.GetTrends(woeid)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(trends)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolSearchTopics handles the search_topics tool
func (tm *ToolsManager) HandleToolSearchTopics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	maxResults := getInt(args, "max_results", 10)
	if maxResults < 10 {
		maxResults = 10
	}
	if maxResults > 20 {
		maxResults = 20
	}

	// Extract topics from the request
	topics := getStringSlice(args, "topics")

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
	args := getArgs(request)
	sampleSize := getInt(args, "sample_size", 20)
	if sampleSize > 100 {
		sampleSize = 100
	}

	// Extract topics from the request
	topics := getStringSlice(args, "topics")

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
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

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
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

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
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

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
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

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

// HandleToolFollowUser handles the follow_user tool
func (tm *ToolsManager) HandleToolFollowUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	username := getString(args, "username", "")

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	targetUser, err := tm.dependencies.TwitterClient.GetUserByUsername(username)
	if err != nil {
		return mcp.NewToolResultError("failed to get target user: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.FollowUser(me.ID, targetUser.ID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "User followed"}`), nil
}

// HandleToolUnfollowUser handles the unfollow_user tool
func (tm *ToolsManager) HandleToolUnfollowUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	username := getString(args, "username", "")

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	targetUser, err := tm.dependencies.TwitterClient.GetUserByUsername(username)
	if err != nil {
		return mcp.NewToolResultError("failed to get target user: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.UnfollowUser(me.ID, targetUser.ID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "User unfollowed"}`), nil
}

// HandleToolGetUserProfile handles the get_user_profile tool
func (tm *ToolsManager) HandleToolGetUserProfile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	username := getString(args, "username", "")

	profile, err := tm.dependencies.TwitterClient.GetUserProfile(username)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(profile)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolGetUserTweets handles the get_user_tweets tool
func (tm *ToolsManager) HandleToolGetUserTweets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	username := getString(args, "username", "")
	maxResults := getInt(args, "max_results", 10)

	user, err := tm.dependencies.TwitterClient.GetUserByUsername(username)
	if err != nil {
		return mcp.NewToolResultError("failed to get user: " + err.Error()), nil
	}

	tweets, err := tm.dependencies.TwitterClient.GetUserTweets(user.ID, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(tweets)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolBookmarkTweet handles the bookmark_tweet tool
func (tm *ToolsManager) HandleToolBookmarkTweet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.BookmarkTweet(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Tweet bookmarked"}`), nil
}

// HandleToolRemoveBookmark handles the remove_bookmark tool
func (tm *ToolsManager) HandleToolRemoveBookmark(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	tweetID := getString(args, "tweet_id", "")

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	err = tm.dependencies.TwitterClient.RemoveBookmark(me.ID, tweetID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(`{"success": true, "message": "Bookmark removed"}`), nil
}

// HandleToolGetBookmarks handles the get_bookmarks tool
func (tm *ToolsManager) HandleToolGetBookmarks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	maxResults := getInt(args, "max_results", 10)

	me, err := tm.dependencies.TwitterClient.GetMe()
	if err != nil {
		return mcp.NewToolResultError("failed to get user info: " + err.Error()), nil
	}

	bookmarks, err := tm.dependencies.TwitterClient.GetBookmarks(me.ID, maxResults)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(bookmarks)
	return mcp.NewToolResultText(string(result)), nil
}

// HandleToolPostThread handles the post_thread tool
func (tm *ToolsManager) HandleToolPostThread(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request)
	tweets := getStringSlice(args, "tweets")

	if len(tweets) == 0 {
		return mcp.NewToolResultError("no tweets provided for thread"), nil
	}

	postedTweets, err := tm.dependencies.TwitterClient.PostThread(tweets)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, _ := json.Marshal(postedTweets)
	return mcp.NewToolResultText(string(result)), nil
}


