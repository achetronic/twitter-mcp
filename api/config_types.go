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

package api

import "time"

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

// AccessLogsConfig represents the AccessLogs middleware configuration
type AccessLogsConfig struct {
	ExcludedHeaders []string `yaml:"excluded_headers"`
	RedactedHeaders []string `yaml:"redacted_headers"`
}

// JWTValidationLocalConfig represents the local JWT validation configuration
type JWTValidationLocalConfig struct {
	JWKSUri         string                        `yaml:"jwks_uri"`
	CacheInterval   time.Duration                 `yaml:"cache_interval"`
	AllowConditions []JWTValidationAllowCondition `yaml:"allow_conditions,omitempty"`
}

// JWTValidationAllowCondition represents a condition for allowing a request after the local JWT validation
type JWTValidationAllowCondition struct {
	Expression string `yaml:"expression"`
}

// JWTValidationConfig represents the JWT validation configuration
type JWTValidationConfig struct {
	Strategy        string                   `yaml:"strategy"`
	ForwardedHeader string                   `yaml:"forwarded_header,omitempty"`
	Local           JWTValidationLocalConfig `yaml:"local,omitempty"`
}

// JWTConfig represents the JWT middleware configuration
type JWTConfig struct {
	Enabled    bool                `yaml:"enabled"`
	Validation JWTValidationConfig `yaml:"validation,omitempty"`
}

// MiddlewareConfig represents the middleware configuration section
type MiddlewareConfig struct {
	AccessLogs AccessLogsConfig `yaml:"access_logs"`
	JWT        JWTConfig        `yaml:"jwt,omitempty"`
}

// OAuthAuthorizationServer represents the OAuth Authorization Server configuration
type OAuthAuthorizationServer struct {
	Enabled   bool   `yaml:"enabled"`
	UrlSuffix string `yaml:"url_suffix,omitempty"`
	IssuerUri string `yaml:"issuer_uri"`
}

// OAuthProtectedResourceConfig represents the OAuth Protected Resource configuration
type OAuthProtectedResourceConfig struct {
	Enabled   bool   `yaml:"enabled"`
	UrlSuffix string `yaml:"url_suffix,omitempty"`

	Resource                              string   `yaml:"resource"`
	AuthServers                           []string `yaml:"auth_servers"`
	JWKSUri                               string   `yaml:"jwks_uri"`
	ScopesSupported                       []string `yaml:"scopes_supported"`
	BearerMethodsSupported                []string `yaml:"bearer_methods_supported,omitempty"`
	ResourceSigningAlgValuesSupported     []string `yaml:"resource_signing_alg_values_supported,omitempty"`
	ResourceName                          string   `yaml:"resource_name,omitempty"`
	ResourceDocumentation                 string   `yaml:"resource_documentation,omitempty"`
	ResourcePolicyUri                     string   `yaml:"resource_policy_uri,omitempty"`
	ResourceTosUri                        string   `yaml:"resource_tos_uri,omitempty"`
	TLSClientCertificateBoundAccessTokens bool     `yaml:"tls_client_certificate_bound_access_tokens,omitempty"`
	AuthorizationDetailsTypesSupported    []string `yaml:"authorization_details_types_supported,omitempty"`
	DPoPSigningAlgValuesSupported         []string `yaml:"dpop_signing_alg_values_supported,omitempty"`
	DPoPBoundAccessTokensRequired         bool     `yaml:"dpop_bound_access_tokens_required,omitempty"`
}

// ToolPolicyConfig represents a policy for tool access control
type ToolPolicyConfig struct {
	Expression   string   `yaml:"expression"`
	AllowedTools []string `yaml:"allowed_tools"`
}

// PoliciesConfig represents the policies configuration section
type PoliciesConfig struct {
	Tools []ToolPolicyConfig `yaml:"tools"`
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
	Server                   ServerConfig                 `yaml:"server,omitempty"`
	Middleware               MiddlewareConfig             `yaml:"middleware,omitempty"`
	Policies                 PoliciesConfig               `yaml:"policies,omitempty"`
	OAuthAuthorizationServer OAuthAuthorizationServer     `yaml:"oauth_authorization_server,omitempty"`
	OAuthProtectedResource   OAuthProtectedResourceConfig `yaml:"oauth_protected_resource,omitempty"`
	Twitter                  TwitterConfig                `yaml:"twitter"`
	ScheduleFile             string                       `yaml:"schedule_file,omitempty"`
}
