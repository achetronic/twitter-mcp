# AGENTS.md

This document helps AI agents work effectively in this repository.

## Project Overview

**Twitter MCP** is a production-ready MCP (Model Context Protocol) server for Twitter/X API integration. Built in Go, it allows AI assistants to read, write, and analyze Twitter content with full OAuth support and JWT-based access control.

### Key Technologies
- **Language**: Go 1.24+
- **MCP Library**: `github.com/mark3labs/mcp-go` v0.43.2
- **Twitter Auth**: `github.com/dghubble/oauth1` for OAuth 1.0a
- **JWT Handling**: `github.com/golang-jwt/jwt/v5`
- **CEL Expressions**: `github.com/google/cel-go` for policy evaluation
- **Configuration**: YAML with environment variable expansion

## Commands

### Development

```bash
# Run the server
make run

# Build binary for current platform
make build

# Clean build artifacts
make clean

# Download dependencies
go mod tidy
```

### Docker

```bash
# Build Docker image
docker build -t twitter-mcp .

# Run with docker-compose
docker-compose up -d
```

## Code Organization

```
.
├── cmd/
│   └── main.go              # Application entrypoint
├── api/
│   └── config_types.go      # Configuration type definitions
├── internal/
│   ├── config/
│   │   └── config.go        # YAML config parsing with env expansion
│   ├── globals/
│   │   └── globals.go       # ApplicationContext (config, logger, context)
│   ├── handlers/
│   │   ├── handlers.go      # HandlersManager for HTTP endpoints
│   │   ├── oauth_authorization_server.go  # /.well-known/oauth-authorization-server
│   │   └── oauth_protected_resource.go    # /.well-known/oauth-protected-resource
│   ├── middlewares/
│   │   ├── interfaces.go            # ToolMiddleware, HttpMiddleware interfaces
│   │   ├── jwt_validation.go        # JWT validation middleware
│   │   ├── jwt_validation_utils.go  # JWKS caching, key conversion
│   │   ├── logging.go               # Access logs middleware
│   │   ├── noop.go                  # No-op middleware
│   │   ├── tool_policy.go           # Tool access control based on JWT claims
│   │   └── utils.go                 # Shared utilities
│   ├── tools/
│   │   ├── tools.go         # ToolsManager - tool registration
│   │   └── handlers.go      # Tool handler implementations
│   └── twitter/
│       └── client.go        # Twitter API client (v1.1 and v2)
├── docs/
│   ├── config-http.yaml     # HTTP transport config example
│   ├── config-stdio.yaml    # Stdio transport config example
│   └── images/              # Logo and banner assets
└── .github/workflows/       # CI/CD pipelines
```

## Architecture Patterns

### Dependency Injection
All managers use a `*Dependencies` struct pattern:
```go
type ToolsManagerDependencies struct {
    AppCtx        *globals.ApplicationContext
    McpServer     *server.MCPServer
    Middlewares   []middlewares.ToolMiddleware
    TwitterClient *twitter.Client
}
```

### Middleware Chain
HTTP middlewares wrap handlers: `accessLogs -> jwtValidation -> handler`
Tool middlewares wrap tool handlers: `toolPolicy -> actualToolHandler`

### Configuration
- Config is loaded from YAML file
- Environment variables are expanded (`$VAR` or `${VAR}`)
- Config is available globally via `appCtx.Config`

## Adding New Tools

1. Add the tool definition in `internal/tools/tools.go`:
```go
tool = mcp.NewTool("my_new_tool",
    mcp.WithDescription("Description of what it does"),
    mcp.WithString("param1",
        mcp.Required(),
        mcp.Description("Parameter description"),
    ),
)
tm.dependencies.McpServer.AddTool(tool, tm.wrapWithMiddlewares(tm.HandleToolMyNewTool))
```

2. Implement the handler in `internal/tools/handlers.go`:
```go
func (tm *ToolsManager) HandleToolMyNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    param1, _ := request.Params.Arguments["param1"].(string)
    // Implementation
    return mcp.NewToolResultText(`{"success": true}`), nil
}
```

3. If needed, add Twitter API methods in `internal/twitter/client.go`

## Available Tools

### Reading
- `get_timeline` - Home timeline
- `get_mentions` - Mentions
- `search_tweets` - Search tweets
- `get_trends` - Trending topics
- `get_me` - Current user info
- `get_user_profile` - User profile by username
- `get_user_tweets` - User's recent tweets
- `get_bookmarks` - Saved bookmarks
- `get_dms` - Direct messages

### Writing
- `post_tweet` - Post a tweet
- `delete_tweet` - Delete a tweet
- `post_thread` - Post a thread
- `like_tweet` / `unlike_tweet` - Like/unlike
- `retweet` / `undo_retweet` - Retweet/undo
- `bookmark_tweet` / `remove_bookmark` - Bookmark management
- `follow_user` / `unfollow_user` - Follow/unfollow
- `send_dm` - Send direct message

### Analysis
- `search_topics` - Search multiple topics
- `get_topics_heat` - Topic popularity analysis

## Tool Policies

Tools can be restricted based on JWT claims using CEL expressions:

```yaml
policies:
  tools:
    - expression: 'payload.groups.exists(g, g == "admins")'
      allowed_tools: ["*"]
    - expression: 'payload.scope.contains("twitter:read")'
      allowed_tools: ["get_*", "search_*"]
```

## Twitter API Notes

- **v1.1 API** (OAuth 1.0a): Used for media upload, trends
- **v2 API** (Bearer token): Used for most operations
- **Free tier**: Very limited (posting only)
- **Basic tier** ($100/mo): Full access to search, timeline, trends

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## Common Issues

### "Rate limit exceeded"
Twitter API has strict rate limits. The client doesn't implement backoff - handle at application level.

### "Could not authenticate you"
Check OAuth credentials. For v1.1 API, you need all four OAuth 1.0a tokens. For v2 API, you need the Bearer token.

### Tool not found in policies
If policies are configured and a tool isn't in any `allowed_tools`, access is denied. Use `"*"` for admin access or `"get_*"` for prefix matching.
