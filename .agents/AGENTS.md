# AGENTS.md

This document helps AI agents work effectively in this repository.

## Project Overview

**Twitter MCP** is a production-ready MCP (Model Context Protocol) server for Twitter/X API integration. Built in Go, it allows AI assistants to read, write, analyze Twitter content, and schedule tweets with full OAuth support and JWT-based access control.

### Key Technologies
- **Language**: Go 1.24+
- **MCP Library**: `github.com/mark3labs/mcp-go` v0.44.0
- **Twitter Auth**: `github.com/dghubble/oauth1` for OAuth 1.0a
- **JWT Handling**: `github.com/golang-jwt/jwt/v5`
- **CEL Expressions**: `github.com/google/cel-go` for policy evaluation
- **UUID**: `github.com/google/uuid` for scheduled tweet IDs
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
│   ├── config_types.go      # Configuration type definitions
│   └── schedule_types.go    # ScheduledTweet and ScheduleStore types
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
│   ├── schedule/
│   │   └── store.go         # YAML-backed persistent store for scheduled tweets
│   ├── tools/
│   │   ├── tools.go                 # ToolsManager - tool registration
│   │   ├── handlers.go              # Twitter tool handler implementations
│   │   ├── schedule_handlers.go     # Schedule tool handler implementations
│   │   └── helpers.go               # getArgs, getString, getInt, getStringSlice
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
    ScheduleStore *schedule.Store
}
```

### Middleware Chain
HTTP middlewares wrap handlers: `accessLogs -> jwtValidation -> handler`
Tool middlewares wrap tool handlers: `toolPolicy -> actualToolHandler`

### Configuration
- Config is loaded from YAML file
- Environment variables are expanded (`$VAR` or `${VAR}`)
- Config is available globally via `appCtx.Config`
- Schedule file path configured via `schedule_file` (default: `schedule.yaml`)

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

2. Implement the handler in `internal/tools/handlers.go` (or a new file):
```go
func (tm *ToolsManager) HandleToolMyNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := getArgs(request)
    param1 := getString(args, "param1", "")
    // Implementation
    return mcp.NewToolResultText(`{"success": true}`), nil
}
```

3. If needed, add Twitter API methods in `internal/twitter/client.go`

## Available Tools

### Reading
- `get_me` - Current user info
- `get_timeline` - Home timeline
- `get_mentions` - Mentions
- `search_tweets` - Search tweets (last 24h, sorted by recency)
- `get_trends` - Trending topics by location (requires v1.1 API access)
- `get_user_profile` - User profile by username
- `get_user_tweets` - User's recent tweets
- `get_bookmarks` - Saved bookmarks

### Writing
- `post_tweet` - Post a tweet (supports replies)
- `post_thread` - Post a thread
- `delete_tweet` - Delete a tweet
- `like_tweet` / `unlike_tweet` - Like/unlike
- `retweet` / `undo_retweet` - Retweet/undo
- `bookmark_tweet` / `remove_bookmark` - Bookmark management
- `follow_user` / `unfollow_user` - Follow/unfollow

### Analysis
- `search_topics` - Search multiple topics at once (last 24h)
- `get_topics_heat` - Topic popularity heat score (last 24h)

### Scheduling
- `schedule_tweet` - Add a tweet or thread to the scheduling queue
- `schedule_update` - Modify a scheduled tweet (content, date, reviewed status)
- `schedule_delete` - Remove a scheduled tweet from the queue
- `schedule_list` - List scheduled tweets, optionally filtered by status
- `schedule_get_publishable` - Get tweets ready to publish (reviewed + scheduled_at past + cooldown respected)
- `schedule_publish` - Publish a specific scheduled tweet by ID

## Scheduling System

Tweets are stored in a YAML file (`schedule.yaml` by default, configurable via `schedule_file`).

### Statuses
- `pending` - Added but not reviewed yet
- `reviewed` - Approved and ready to publish when scheduled_at arrives
- `published` - Successfully published
- `failed` - Publishing failed (see `fail_reason`)

### Content format
Content is always `[]string`. One element for a tweet, multiple for a thread. This keeps the code simple and consistent.

### Publishing flow
The AI is always in the loop. There is no background worker. The recommended flow is:
1. Call `schedule_get_publishable` to check what's ready
2. Decide which one to publish (respect timing, don't publish multiple at once)
3. Call `schedule_publish` with the chosen ID

## Twitter API Notes

- **v1.1 API** (OAuth 1.0a): Used for media upload, trends
- **v2 API** (Bearer token): Used for most read operations
- **v2 API** (OAuth 1.0a User Context): Used for all write operations
- **DMs removed**: Requires OAuth 2.0 with redirect callback — not suitable for self-hosted setups
- **Free tier**: Very limited (posting only)
- **Basic tier** ($100/mo): Full access to search, timeline, trends
- **Search results**: Limited to last 24 hours, sorted by recency

## Tool Policies

Tools can be restricted based on JWT claims using CEL expressions:

```yaml
policies:
  tools:
    - expression: 'payload.groups.exists(g, g == "admins")'
      allowed_tools: ["*"]
    - expression: 'payload.scope.contains("twitter:read")'
      allowed_tools: ["get_*", "search_*", "schedule_list", "schedule_get_publishable"]
```

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
Twitter API has strict rate limits. The client doesn't implement backoff — handle at application level.

### "Could not authenticate you"
Check OAuth credentials. For write operations, you need all four OAuth 1.0a tokens. For read operations, you need the Bearer token.

### "CreditsDepleted"
You've run out of API credits. Check your Twitter Developer Portal to top up or wait for the monthly reset.

### Tool not found in policies
If policies are configured and a tool isn't in any `allowed_tools`, access is denied. Use `"*"` for admin access or `"get_*"` for prefix matching.

### Scheduled tweet not appearing in get_publishable
Check that: (1) the tweet has `reviewed: true`, (2) `scheduled_at` is in the past, (3) enough time has passed since the last published tweet (`min_hours_since_last`).
