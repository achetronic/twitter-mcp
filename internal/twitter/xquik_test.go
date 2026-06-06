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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchTweetsUsesXquikBackend(t *testing.T) {
	var gotPath string
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.String()
		gotAuth = r.Header.Get("X-API-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"tweets": [{
				"id": "123",
				"full_text": "hello from xquik",
				"user_id": "u1",
				"created_at": "2026-06-06T00:00:00Z",
				"metrics": {"like_count": 3, "retweet_count": 2, "reply_count": 1, "quote_count": 0}
			}],
			"users": [{"id": "u1", "name": "Ada", "handle": "ada"}],
			"meta": {"next_cursor": "cursor-1"}
		}`))
	}))
	defer server.Close()

	client := NewClient("key", "secret", "token", "tokenSecret", "bearer")
	client.ConfigureXquik(XquikConfig{
		Enabled: true,
		BaseURL: server.URL,
		APIKey:  "xq_test",
	})

	result, err := client.SearchTweets("ai agents", 25)
	if err != nil {
		t.Fatalf("SearchTweets returned error: %v", err)
	}

	if !strings.HasPrefix(gotPath, "/api/v1/x/tweets/search?") {
		t.Fatalf("expected Xquik search endpoint, got %s", gotPath)
	}
	if gotAuth != "xq_test" {
		t.Fatalf("expected X-API-Key header, got %s", gotAuth)
	}
	if len(result.Data) != 1 || result.Data[0].Text != "hello from xquik" {
		t.Fatalf("unexpected tweets: %+v", result.Data)
	}
	if result.Data[0].PublicMetrics.LikeCount != 3 {
		t.Fatalf("unexpected metrics: %+v", result.Data[0].PublicMetrics)
	}
	if len(result.Includes.Users) != 1 || result.Includes.Users[0].Username != "ada" {
		t.Fatalf("unexpected users: %+v", result.Includes.Users)
	}
	if result.Meta.NextToken != "cursor-1" {
		t.Fatalf("expected next cursor, got %s", result.Meta.NextToken)
	}
}

func TestGetUserByUsernameUsesXquikBackend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/x/users/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "ada" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"users": [{"id": "u1", "name": "Ada Lovelace", "screen_name": "ada"}]}`))
	}))
	defer server.Close()

	client := NewClient("key", "secret", "token", "tokenSecret", "bearer")
	client.ConfigureXquik(XquikConfig{
		Enabled: true,
		BaseURL: server.URL,
		APIKey:  "xq_test",
	})

	user, err := client.GetUserByUsername("ada")
	if err != nil {
		t.Fatalf("GetUserByUsername returned error: %v", err)
	}

	if user.ID != "u1" || user.Name != "Ada Lovelace" || user.Username != "ada" {
		t.Fatalf("unexpected user: %+v", user)
	}
}
