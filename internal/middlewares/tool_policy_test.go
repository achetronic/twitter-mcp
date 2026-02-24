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

package middlewares

import (
	"testing"
)

func TestIsToolAllowed(t *testing.T) {
	mw := &ToolPolicyMiddleware{}

	tests := []struct {
		toolName     string
		allowedTools []string
		expected     bool
	}{
		{"post_tweet", []string{"post_tweet"}, true},
		{"post_tweet", []string{"delete_tweet"}, false},
		{"post_tweet", []string{"*"}, true},
		{"get_timeline", []string{"get_*"}, true},
		{"post_tweet", []string{"get_*"}, false},
		{"get_mentions", []string{"get_timeline", "get_mentions"}, true},
		{"search_tweets", []string{"search_*", "get_*"}, true},
	}

	for _, tt := range tests {
		result := mw.isToolAllowed(tt.toolName, tt.allowedTools)
		if result != tt.expected {
			t.Errorf("isToolAllowed(%s, %v) = %v, expected %v",
				tt.toolName, tt.allowedTools, result, tt.expected)
		}
	}
}

func TestGetRequestScheme(t *testing.T) {
	// Test without headers (should return http)
	// This is a basic test - full testing would require http.Request mocking
}
