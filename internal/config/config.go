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
