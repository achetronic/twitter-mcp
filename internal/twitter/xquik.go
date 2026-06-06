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
	"encoding/json"
	"fmt"
)

type xquikTweet struct {
	ID            string         `json:"id"`
	Text          string         `json:"text"`
	FullText      string         `json:"full_text"`
	AuthorID      string         `json:"author_id"`
	UserID        string         `json:"user_id"`
	CreatedAt     string         `json:"created_at"`
	PublicMetrics *PublicMetrics `json:"public_metrics"`
	Metrics       *PublicMetrics `json:"metrics"`
}

type xquikUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Handle     string `json:"handle"`
	ScreenName string `json:"screen_name"`
}

type xquikTweetsEnvelope struct {
	Data     []xquikTweet `json:"data"`
	Tweets   []xquikTweet `json:"tweets"`
	Items    []xquikTweet `json:"items"`
	Users    []xquikUser  `json:"users"`
	Includes struct {
		Users []xquikUser `json:"users"`
	} `json:"includes"`
	Meta struct {
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
		NextCursor  string `json:"next_cursor"`
	} `json:"meta"`
}

type xquikUsersEnvelope struct {
	Data  []xquikUser `json:"data"`
	Users []xquikUser `json:"users"`
	Items []xquikUser `json:"items"`
	User  *xquikUser  `json:"user"`
}

func parseXquikTweetsResponse(body []byte) (*TweetsResponse, error) {
	var envelope xquikTweetsEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse Xquik search response: %w", err)
	}

	source := firstTweets(envelope)
	response := &TweetsResponse{}
	response.Data = make([]Tweet, 0, len(source))
	for _, item := range source {
		response.Data = append(response.Data, mapXquikTweet(item))
	}

	users := append([]xquikUser{}, envelope.Users...)
	users = append(users, envelope.Includes.Users...)
	response.Includes.Users = make([]User, 0, len(users))
	for _, user := range users {
		response.Includes.Users = append(response.Includes.Users, mapXquikUser(user))
	}

	response.Meta.ResultCount = envelope.Meta.ResultCount
	if response.Meta.ResultCount == 0 {
		response.Meta.ResultCount = len(response.Data)
	}
	if envelope.Meta.NextToken != "" {
		response.Meta.NextToken = envelope.Meta.NextToken
	} else {
		response.Meta.NextToken = envelope.Meta.NextCursor
	}

	return response, nil
}

func parseXquikUserResponse(body []byte) (*User, error) {
	var envelope xquikUsersEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse Xquik user response: %w", err)
	}

	if envelope.User != nil {
		user := mapXquikUser(*envelope.User)
		return &user, nil
	}

	users := firstUsers(envelope)
	if len(users) == 0 {
		return nil, fmt.Errorf("Xquik user response contained no users")
	}

	user := mapXquikUser(users[0])
	return &user, nil
}

func firstTweets(envelope xquikTweetsEnvelope) []xquikTweet {
	if len(envelope.Data) > 0 {
		return envelope.Data
	}
	if len(envelope.Tweets) > 0 {
		return envelope.Tweets
	}
	return envelope.Items
}

func firstUsers(envelope xquikUsersEnvelope) []xquikUser {
	if len(envelope.Data) > 0 {
		return envelope.Data
	}
	if len(envelope.Users) > 0 {
		return envelope.Users
	}
	return envelope.Items
}

func mapXquikTweet(item xquikTweet) Tweet {
	text := item.Text
	if text == "" {
		text = item.FullText
	}

	authorID := item.AuthorID
	if authorID == "" {
		authorID = item.UserID
	}

	metrics := item.PublicMetrics
	if metrics == nil {
		metrics = item.Metrics
	}

	return Tweet{
		ID:            item.ID,
		Text:          text,
		AuthorID:      authorID,
		CreatedAt:     item.CreatedAt,
		PublicMetrics: metrics,
	}
}

func mapXquikUser(item xquikUser) User {
	username := item.Username
	if username == "" {
		username = item.Handle
	}
	if username == "" {
		username = item.ScreenName
	}

	return User{
		ID:       item.ID,
		Name:     item.Name,
		Username: username,
	}
}
