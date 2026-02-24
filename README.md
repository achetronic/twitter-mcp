# Twitter MCP

An MCP (Model Context Protocol) server for interacting with Twitter/X API.

## Features

- **Post tweets** - Create new tweets, with optional reply support
- **Delete tweets** - Remove your tweets
- **Get timeline** - View your home timeline
- **Get mentions** - See tweets mentioning you
- **Search tweets** - Search for tweets with query operators
- **Get trends** - View trending topics by location
- **Get tech trends** - Curated search for tech/AI/cloud topics
- **Like/Unlike tweets** - Interact with tweets
- **Retweet/Undo retweet** - Share tweets

## Prerequisites

You need Twitter API credentials. Get them from [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard):

1. Create a project and app
2. Generate OAuth 1.0a credentials (API Key, API Secret, Access Token, Access Token Secret)
3. Generate a Bearer Token for v2 API access

## Configuration

Create a `config.yaml` file:

```yaml
server:
  name: "twitter-mcp"
  version: "0.1.0"

twitter:
  api_key: "your-api-key"
  api_key_secret: "your-api-key-secret"
  access_token: "your-access-token"
  access_token_secret: "your-access-token-secret"
  bearer_token: "your-bearer-token"
```

Or use environment variables (the config supports `$VAR` expansion):

```yaml
twitter:
  api_key: "$TWITTER_API_KEY"
  api_key_secret: "$TWITTER_API_KEY_SECRET"
  access_token: "$TWITTER_ACCESS_TOKEN"
  access_token_secret: "$TWITTER_ACCESS_TOKEN_SECRET"
  bearer_token: "$TWITTER_BEARER_TOKEN"
```

## Building

```bash
go mod tidy
go build -o twitter-mcp ./cmd/main.go
```

## Running

```bash
./twitter-mcp -config config.yaml
```

## Available Tools

| Tool | Description |
|------|-------------|
| `post_tweet` | Post a new tweet (supports replies) |
| `delete_tweet` | Delete a tweet by ID |
| `get_timeline` | Get home timeline |
| `get_mentions` | Get mentions of authenticated user |
| `search_tweets` | Search tweets with query |
| `get_trends` | Get trending topics (WOEID: 1=World, 23424950=Spain) |
| `get_tech_trends` | Search tech/AI/cloud related content |
| `get_me` | Get authenticated user info |
| `like_tweet` | Like a tweet |
| `unlike_tweet` | Remove like from tweet |
| `retweet` | Retweet a tweet |
| `undo_retweet` | Remove retweet |

## WOEIDs for Trends

Common WOEIDs for `get_trends`:
- `1` - Worldwide
- `23424950` - Spain
- `23424977` - United States
- `766273` - Madrid
- `753692` - Barcelona
- `23424975` - United Kingdom
- `2459115` - New York

## License

Apache 2.0
