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

package config

import (
	"os"
	"twitter-mcp/api"

	"gopkg.in/yaml.v3"
)

// Marshal converts a Configuration to YAML bytes
func Marshal(config api.Configuration) (bytes []byte, err error) {
	bytes, err = yaml.Marshal(config)
	return bytes, err
}

// Unmarshal converts YAML bytes to a Configuration
func Unmarshal(bytes []byte) (config api.Configuration, err error) {
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

// ReadFile reads and parses a configuration file
func ReadFile(filepath string) (config api.Configuration, err error) {
	var fileBytes []byte
	fileBytes, err = os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	// Expand environment variables present in the config
	fileExpandedEnv := os.ExpandEnv(string(fileBytes))

	config, err = Unmarshal([]byte(fileExpandedEnv))

	return config, err
}
