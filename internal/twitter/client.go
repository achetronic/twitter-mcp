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

package twitter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

const (
	baseURLv1 = "https://api.twitter.com/1.1"
	baseURLv2 = "https://api.twitter.com/2"
)

// Client represents a Twitter/X API client
type Client struct {
	// OAuth 1.0a client for v1.1 API (write operations)
	oauth1Client *http.Client
	// Bearer token for v2 API (read operations)
	bearerToken string
	httpClient  *http.Client
}

// NewClient creates a new Twitter client
func NewClient(apiKey, apiKeySecret, accessToken, accessTokenSecret, bearerToken string) *Client {
	// Setup OAuth 1.0a for v1.1 API
	config := oauth1.NewConfig(apiKey, apiKeySecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	oauth1Client := config.Client(oauth1.NoContext, token)

	return &Client{
		oauth1Client: oauth1Client,
		bearerToken:  bearerToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequestV2 performs an HTTP request to the Twitter v2 API using Bearer token
func (c *Client) doRequestV2(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURLv2+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doRequestV1 performs an HTTP request to the Twitter v1.1 API using OAuth 1.0a
func (c *Client) doRequestV1(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURLv1+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.oauth1Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doRequestV1Form performs a form-encoded POST request to the Twitter v1.1 API
func (c *Client) doRequestV1Form(endpoint string, params url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", baseURLv1+endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.oauth1Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// PublicMetrics represents engagement metrics for a tweet
type PublicMetrics struct {
	LikeCount    int `json:"like_count"`
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
	QuoteCount   int `json:"quote_count"`
}

// Tweet represents a tweet object
type Tweet struct {
	ID            string         `json:"id"`
	Text          string         `json:"text"`
	AuthorID      string         `json:"author_id,omitempty"`
	CreatedAt     string         `json:"created_at,omitempty"`
	PublicMetrics *PublicMetrics `json:"public_metrics,omitempty"`
}

// User represents a Twitter user
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// TweetResponse represents the response from tweet-related endpoints
type TweetResponse struct {
	Data     *Tweet `json:"data,omitempty"`
	Includes struct {
		Users []User `json:"users,omitempty"`
	} `json:"includes,omitempty"`
}

// TweetsResponse represents multiple tweets
type TweetsResponse struct {
	Data     []Tweet `json:"data,omitempty"`
	Includes struct {
		Users []User `json:"users,omitempty"`
	} `json:"includes,omitempty"`
	Meta struct {
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token,omitempty"`
	} `json:"meta,omitempty"`
}

// Trend represents a trending topic
type Trend struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	TweetVolume int    `json:"tweet_volume,omitempty"`
	Query       string `json:"query,omitempty"`
}

// TrendsResponse represents the trends response (v1.1 API)
type TrendsResponse []struct {
	Trends    []Trend `json:"trends"`
	AsOf      string  `json:"as_of"`
	CreatedAt string  `json:"created_at"`
	Locations []struct {
		Name  string `json:"name"`
		WOEID int    `json:"woeid"`
	} `json:"locations"`
}

// PostTweet posts a new tweet (v2 API)
func (c *Client) PostTweet(text string, replyToID string) (*Tweet, error) {
	payload := map[string]interface{}{
		"text": text,
	}

	if replyToID != "" {
		payload["reply"] = map[string]string{
			"in_reply_to_tweet_id": replyToID,
		}
	}

	body, err := c.doRequestV2("POST", "/tweets", payload)
	if err != nil {
		return nil, err
	}

	var response TweetResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse tweet response: %w", err)
	}

	return response.Data, nil
}

// DeleteTweet deletes a tweet (v2 API)
func (c *Client) DeleteTweet(tweetID string) error {
	_, err := c.doRequestV2("DELETE", "/tweets/"+tweetID, nil)
	return err
}

// GetTimeline gets the authenticated user's home timeline (v2 API)
func (c *Client) GetTimeline(userID string, maxResults int) (*TweetsResponse, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("/users/%s/timelines/reverse_chronological?max_results=%d&tweet.fields=created_at,author_id&expansions=author_id", userID, maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TweetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse timeline response: %w", err)
	}

	return &response, nil
}

// GetMentions gets mentions of the authenticated user (v2 API)
func (c *Client) GetMentions(userID string, maxResults int) (*TweetsResponse, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("/users/%s/mentions?max_results=%d&tweet.fields=created_at,author_id&expansions=author_id", userID, maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TweetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse mentions response: %w", err)
	}

	return &response, nil
}

// SearchTweets searches for tweets (v2 API)
func (c *Client) SearchTweets(query string, maxResults int) (*TweetsResponse, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	encodedQuery := url.QueryEscape(query)
	endpoint := fmt.Sprintf("/tweets/search/recent?query=%s&max_results=%d&tweet.fields=created_at,author_id,public_metrics&expansions=author_id", encodedQuery, maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TweetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	return &response, nil
}

// GetTrends gets trending topics for a location (v1.1 API)
// WOEID: 1 = Worldwide, 23424950 = Spain, 766273 = Madrid
func (c *Client) GetTrends(woeid int) ([]Trend, error) {
	if woeid <= 0 {
		woeid = 1 // Worldwide
	}

	endpoint := fmt.Sprintf("/trends/place.json?id=%d", woeid)

	body, err := c.doRequestV1("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TrendsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse trends response: %w", err)
	}

	if len(response) > 0 {
		return response[0].Trends, nil
	}

	return []Trend{}, nil
}

// GetTrendsByTopic searches tweets and returns them filtered by topics
// This is a workaround since Twitter API doesn't have topic-based trends directly
func (c *Client) GetTrendsByTopic(topics []string, maxResults int) (map[string]*TweetsResponse, error) {
	results := make(map[string]*TweetsResponse)

	for _, topic := range topics {
		tweets, err := c.SearchTweets(topic, maxResults)
		if err != nil {
			// Continue with other topics even if one fails
			continue
		}
		results[topic] = tweets
	}

	return results, nil
}

// TopicHeat represents the "heat" or popularity of a topic
type TopicHeat struct {
	Topic         string  `json:"topic"`
	TweetCount    int     `json:"tweet_count"`
	TotalLikes    int     `json:"total_likes"`
	TotalRetweets int     `json:"total_retweets"`
	TotalReplies  int     `json:"total_replies"`
	TotalQuotes   int     `json:"total_quotes"`
	AvgEngagement float64 `json:"avg_engagement"`
	HeatScore     float64 `json:"heat_score"` // 0-100 calculated score
}

// GetTopicsHeat searches topics and calculates a heat score for each
func (c *Client) GetTopicsHeat(topics []string, maxResults int) ([]TopicHeat, error) {
	var results []TopicHeat

	for _, topic := range topics {
		tweets, err := c.SearchTweets(topic, maxResults)
		if err != nil {
			// Add topic with zero heat if search fails
			results = append(results, TopicHeat{
				Topic:     topic,
				HeatScore: 0,
			})
			continue
		}

		heat := TopicHeat{
			Topic:      topic,
			TweetCount: len(tweets.Data),
		}

		// Sum up all metrics
		for _, tweet := range tweets.Data {
			if tweet.PublicMetrics != nil {
				heat.TotalLikes += tweet.PublicMetrics.LikeCount
				heat.TotalRetweets += tweet.PublicMetrics.RetweetCount
				heat.TotalReplies += tweet.PublicMetrics.ReplyCount
				heat.TotalQuotes += tweet.PublicMetrics.QuoteCount
			}
		}

		// Calculate average engagement per tweet
		if heat.TweetCount > 0 {
			totalEngagement := heat.TotalLikes + heat.TotalRetweets + heat.TotalReplies + heat.TotalQuotes
			heat.AvgEngagement = float64(totalEngagement) / float64(heat.TweetCount)
		}

		// Calculate heat score (0-100)
		// Formula: combines tweet count and engagement
		// - Tweet count contributes up to 40 points (maxed at 100 tweets)
		// - Avg engagement contributes up to 60 points (logarithmic scale)
		tweetScore := float64(heat.TweetCount) / float64(maxResults) * 40
		if tweetScore > 40 {
			tweetScore = 40
		}

		// Logarithmic scale for engagement (1 engagement = ~10 points, 100 = ~40 points, 1000 = ~60 points)
		engagementScore := 0.0
		if heat.AvgEngagement > 0 {
			import_math := heat.AvgEngagement + 1 // avoid log(0)
			engagementScore = 20 * (1 + logBase10(import_math))
			if engagementScore > 60 {
				engagementScore = 60
			}
		}

		heat.HeatScore = tweetScore + engagementScore

		results = append(results, heat)
	}

	// Sort by heat score descending
	sortTopicsByHeat(results)

	return results, nil
}

// logBase10 calculates log base 10
func logBase10(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// log10(x) = ln(x) / ln(10)
	// Using a simple approximation or math package
	result := 0.0
	for x >= 10 {
		x /= 10
		result++
	}
	// Linear interpolation for the fractional part
	if x > 1 {
		result += (x - 1) / 9
	}
	return result
}

// sortTopicsByHeat sorts topics by heat score in descending order
func sortTopicsByHeat(topics []TopicHeat) {
	for i := 0; i < len(topics)-1; i++ {
		for j := i + 1; j < len(topics); j++ {
			if topics[j].HeatScore > topics[i].HeatScore {
				topics[i], topics[j] = topics[j], topics[i]
			}
		}
	}
}

// GetMe gets the authenticated user's info (v2 API)
func (c *Client) GetMe() (*User, error) {
	body, err := c.doRequestV2("GET", "/users/me", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data User `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &response.Data, nil
}

// LikeTweet likes a tweet (v2 API)
func (c *Client) LikeTweet(userID, tweetID string) error {
	payload := map[string]string{
		"tweet_id": tweetID,
	}

	_, err := c.doRequestV2("POST", "/users/"+userID+"/likes", payload)
	return err
}

// UnlikeTweet removes a like from a tweet (v2 API)
func (c *Client) UnlikeTweet(userID, tweetID string) error {
	_, err := c.doRequestV2("DELETE", "/users/"+userID+"/likes/"+tweetID, nil)
	return err
}

// Retweet retweets a tweet (v2 API)
func (c *Client) Retweet(userID, tweetID string) error {
	payload := map[string]string{
		"tweet_id": tweetID,
	}

	_, err := c.doRequestV2("POST", "/users/"+userID+"/retweets", payload)
	return err
}

// UndoRetweet removes a retweet (v2 API)
func (c *Client) UndoRetweet(userID, tweetID string) error {
	_, err := c.doRequestV2("DELETE", "/users/"+userID+"/retweets/"+tweetID, nil)
	return err
}

// FollowUser follows a user (v2 API)
func (c *Client) FollowUser(sourceUserID, targetUserID string) error {
	payload := map[string]string{
		"target_user_id": targetUserID,
	}

	_, err := c.doRequestV2("POST", "/users/"+sourceUserID+"/following", payload)
	return err
}

// UnfollowUser unfollows a user (v2 API)
func (c *Client) UnfollowUser(sourceUserID, targetUserID string) error {
	_, err := c.doRequestV2("DELETE", "/users/"+sourceUserID+"/following/"+targetUserID, nil)
	return err
}

// GetUserByUsername gets a user's profile by username (v2 API)
func (c *Client) GetUserByUsername(username string) (*User, error) {
	endpoint := fmt.Sprintf("/users/by/username/%s?user.fields=description,public_metrics,created_at,profile_image_url", username)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data UserProfile `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &User{
		ID:       response.Data.ID,
		Name:     response.Data.Name,
		Username: response.Data.Username,
	}, nil
}

// UserProfile represents a detailed user profile
type UserProfile struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Username        string        `json:"username"`
	Description     string        `json:"description,omitempty"`
	ProfileImageURL string        `json:"profile_image_url,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	PublicMetrics   *UserMetrics  `json:"public_metrics,omitempty"`
}

// UserMetrics represents user engagement metrics
type UserMetrics struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	TweetCount     int `json:"tweet_count"`
	ListedCount    int `json:"listed_count"`
}

// GetUserProfile gets a user's full profile by username (v2 API)
func (c *Client) GetUserProfile(username string) (*UserProfile, error) {
	endpoint := fmt.Sprintf("/users/by/username/%s?user.fields=description,public_metrics,created_at,profile_image_url", username)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data UserProfile `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user profile: %w", err)
	}

	return &response.Data, nil
}

// GetUserTweets gets recent tweets from a specific user (v2 API)
func (c *Client) GetUserTweets(userID string, maxResults int) (*TweetsResponse, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("/users/%s/tweets?max_results=%d&tweet.fields=created_at,author_id,public_metrics&expansions=author_id", userID, maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TweetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse user tweets: %w", err)
	}

	return &response, nil
}

// BookmarkTweet bookmarks a tweet (v2 API)
func (c *Client) BookmarkTweet(userID, tweetID string) error {
	payload := map[string]string{
		"tweet_id": tweetID,
	}

	_, err := c.doRequestV2("POST", "/users/"+userID+"/bookmarks", payload)
	return err
}

// RemoveBookmark removes a bookmark from a tweet (v2 API)
func (c *Client) RemoveBookmark(userID, tweetID string) error {
	_, err := c.doRequestV2("DELETE", "/users/"+userID+"/bookmarks/"+tweetID, nil)
	return err
}

// GetBookmarks gets the authenticated user's bookmarks (v2 API)
func (c *Client) GetBookmarks(userID string, maxResults int) (*TweetsResponse, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("/users/%s/bookmarks?max_results=%d&tweet.fields=created_at,author_id,public_metrics&expansions=author_id", userID, maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response TweetsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse bookmarks: %w", err)
	}

	return &response, nil
}

// PostThread posts a thread of tweets (v2 API)
func (c *Client) PostThread(tweets []string) ([]*Tweet, error) {
	var postedTweets []*Tweet
	var replyToID string

	for _, text := range tweets {
		tweet, err := c.PostTweet(text, replyToID)
		if err != nil {
			return postedTweets, fmt.Errorf("failed to post tweet in thread: %w", err)
		}
		postedTweets = append(postedTweets, tweet)
		replyToID = tweet.ID
	}

	return postedTweets, nil
}

// SendDM sends a direct message to a user (v2 API)
func (c *Client) SendDM(participantID, text string) error {
	payload := map[string]interface{}{
		"text": text,
	}

	_, err := c.doRequestV2("POST", "/dm_conversations/with/"+participantID+"/messages", payload)
	return err
}

// DMConversation represents a DM conversation
type DMConversation struct {
	ID               string `json:"id"`
	Text             string `json:"text"`
	SenderID         string `json:"sender_id"`
	CreatedAt        string `json:"created_at,omitempty"`
}

// GetDMEvents gets recent DM events (v2 API)
func (c *Client) GetDMEvents(maxResults int) ([]DMConversation, error) {
	if maxResults <= 0 {
		maxResults = 10
	}
	if maxResults > 100 {
		maxResults = 100
	}

	endpoint := fmt.Sprintf("/dm_events?max_results=%d&dm_event.fields=text,sender_id,created_at", maxResults)

	body, err := c.doRequestV2("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []DMConversation `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse DM events: %w", err)
	}

	return response.Data, nil
}

// MediaUploadResponse represents the response from media upload
type MediaUploadResponse struct {
	MediaID       int64  `json:"media_id"`
	MediaIDString string `json:"media_id_string"`
}

// UploadMedia uploads media (image) to Twitter (v1.1 API)
func (c *Client) UploadMedia(imageData []byte) (*MediaUploadResponse, error) {
	// Base64 encode the image
	encoded := base64.StdEncoding.EncodeToString(imageData)

	params := url.Values{}
	params.Set("media_data", encoded)

	body, err := c.doRequestV1Form("/media/upload.json", params)
	if err != nil {
		return nil, err
	}

	var response MediaUploadResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse media upload response: %w", err)
	}

	return &response, nil
}

// PostTweetWithMedia posts a tweet with media attachments (v2 API)
func (c *Client) PostTweetWithMedia(text string, mediaIDs []string) (*Tweet, error) {
	payload := map[string]interface{}{
		"text": text,
	}

	if len(mediaIDs) > 0 {
		payload["media"] = map[string]interface{}{
			"media_ids": mediaIDs,
		}
	}

	body, err := c.doRequestV2("POST", "/tweets", payload)
	if err != nil {
		return nil, err
	}

	var response TweetResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse tweet response: %w", err)
	}

	return response.Data, nil
}
