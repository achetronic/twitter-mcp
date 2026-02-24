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
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("key", "secret", "token", "tokenSecret", "bearer")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.bearerToken != "bearer" {
		t.Errorf("expected bearerToken to be 'bearer', got '%s'", client.bearerToken)
	}
}

func TestLogBase10(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{1, 0},
		{10, 1},
		{100, 2},
		{1000, 3},
		{0, 0},
		{-1, 0},
	}

	for _, tt := range tests {
		result := logBase10(tt.input)
		// Allow some tolerance for the approximation
		if result < tt.expected-0.5 || result > tt.expected+0.5 {
			t.Errorf("logBase10(%f) = %f, expected ~%f", tt.input, result, tt.expected)
		}
	}
}

func TestSortTopicsByHeat(t *testing.T) {
	topics := []TopicHeat{
		{Topic: "low", HeatScore: 10},
		{Topic: "high", HeatScore: 90},
		{Topic: "mid", HeatScore: 50},
	}

	sortTopicsByHeat(topics)

	if topics[0].Topic != "high" {
		t.Errorf("expected first topic to be 'high', got '%s'", topics[0].Topic)
	}
	if topics[1].Topic != "mid" {
		t.Errorf("expected second topic to be 'mid', got '%s'", topics[1].Topic)
	}
	if topics[2].Topic != "low" {
		t.Errorf("expected third topic to be 'low', got '%s'", topics[2].Topic)
	}
}
