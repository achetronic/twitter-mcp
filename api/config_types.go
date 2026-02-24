package api

// ServerTransportHTTPConfig represents the HTTP transport configuration
type ServerTransportHTTPConfig struct {
	Host string `yaml:"host"`
}

// ServerTransportConfig represents the transport configuration
type ServerTransportConfig struct {
	Type string                    `yaml:"type"`
	HTTP ServerTransportHTTPConfig `yaml:"http,omitempty"`
}

// ServerConfig represents the server configuration section
type ServerConfig struct {
	Name      string                `yaml:"name"`
	Version   string                `yaml:"version"`
	Transport ServerTransportConfig `yaml:"transport,omitempty"`
}

// TwitterConfig represents the Twitter/X API configuration
type TwitterConfig struct {
	// OAuth 1.0a credentials (for v1.1 API - posting tweets, etc.)
	APIKey            string `yaml:"api_key"`
	APIKeySecret      string `yaml:"api_key_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`

	// OAuth 2.0 Bearer Token (for v2 API - read operations)
	BearerToken string `yaml:"bearer_token"`
}

// Configuration represents the complete configuration structure
type Configuration struct {
	Server  ServerConfig  `yaml:"server,omitempty"`
	Twitter TwitterConfig `yaml:"twitter"`
}
