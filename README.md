<p align="center">
  <img src="docs/images/header.svg" alt="Twitter MCP" width="800"/>
</p>

<p align="center">
  <em>A Model Context Protocol server that lets AI assistants interact with Twitter/X.<br/>Built in Go, designed to be simple and useful.</em>
</p>

---

## 🎯 What can it do?

This MCP gives your AI assistant the ability to:

- **Read** your timeline, mentions, and search for tweets
- **Write** posts, replies, likes, and retweets
- **Analyze** topic popularity with a heat score system
- **Explore** trending topics by location

The heat score feature is particularly interesting: give it a list of topics and it will tell you which ones are getting more traction right now, based on tweet volume and engagement metrics.

## 🚀 Getting started

### 1. Get your Twitter API credentials

Head to the [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard) and create an app. You'll need:

- **API Key & Secret** (for posting, liking, retweeting)
- **Access Token & Secret** (for acting on behalf of your account)
- **Bearer Token** (for reading public data)

> ⚠️ Heads up: Twitter's free tier is very limited. For full functionality (trends, search, timeline), you'll need at least the Basic tier ($100/month).

### 2. Choose your transport mode

Twitter MCP supports two transport modes:

#### STDIO (simple, for local use)

```yaml
server:
  name: "twitter-mcp"
  version: "0.1.0"
  transport:
    type: "stdio"

twitter:
  api_key: "$TWITTER_API_KEY"
  api_key_secret: "$TWITTER_API_KEY_SECRET"
  access_token: "$TWITTER_ACCESS_TOKEN"
  access_token_secret: "$TWITTER_ACCESS_TOKEN_SECRET"
  bearer_token: "$TWITTER_BEARER_TOKEN"
```

#### HTTP (production, with auth support)

```yaml
server:
  name: "twitter-mcp"
  version: "0.1.0"
  transport:
    type: "http"
    http:
      host: ":8080"

middleware:
  access_logs:
    excluded_headers: ["Accept", "Accept-Encoding"]
    redacted_headers: ["Authorization"]
  jwt:
    enabled: true
    validation:
      strategy: "local"
      local:
        jwks_uri: "https://your-idp.com/.well-known/jwks.json"
        cache_interval: 5m

twitter:
  api_key: "$TWITTER_API_KEY"
  # ... rest of credentials
```

See `docs/config-stdio.yaml` and `docs/config-http.yaml` for full examples.

### 3. Build and run

```bash
go mod tidy
go build -o twitter-mcp ./cmd/main.go
./twitter-mcp -config config.yaml
```

Or use the Makefile:

```bash
make build
make run
```

## 🛠️ Available tools

### Reading

| Tool | What it does |
|------|--------------|
| `get_timeline` | Fetch your home timeline |
| `get_mentions` | See who's mentioning you |
| `search_tweets` | Search tweets with any query |
| `get_trends` | Get trending topics for a location |
| `get_me` | Get your account info |

### Writing

| Tool | What it does |
|------|--------------|
| `post_tweet` | Post a new tweet (supports replies) |
| `delete_tweet` | Delete one of your tweets |
| `like_tweet` | Like a tweet |
| `unlike_tweet` | Remove a like |
| `retweet` | Retweet something |
| `undo_retweet` | Undo a retweet |

### Analysis

| Tool | What it does |
|------|--------------|
| `search_topics` | Search multiple topics at once |
| `get_topics_heat` | Compare topic popularity with heat scores |

## 🔥 The heat score explained

When you call `get_topics_heat` with a list of topics, it returns something like:

```json
[
  {
    "topic": "kubernetes",
    "tweet_count": 20,
    "total_likes": 1250,
    "total_retweets": 340,
    "total_replies": 89,
    "avg_engagement": 86.2,
    "heat_score": 78.5
  },
  {
    "topic": "podman",
    "tweet_count": 15,
    "total_likes": 120,
    "total_retweets": 25,
    "avg_engagement": 10.7,
    "heat_score": 51.2
  }
]
```

The score (0-100) combines tweet volume and engagement. Results come sorted from hottest to coldest, so you can quickly see what's getting attention.

## 🔐 Authentication & Security

When running in HTTP mode, Twitter MCP supports:

- **JWT validation** with JWKS endpoint caching
- **CEL expressions** for fine-grained access control
- **Tool policies** based on JWT claims (groups, scopes, etc.)
- **OAuth 2.0 metadata endpoints** (RFC 9728 compliant)
- **Access logging** with header redaction

This makes it suitable for production deployments behind an identity provider like Keycloak, Auth0, or similar.

### Tool Policies

You can restrict which tools are available based on JWT claims. Policies are evaluated in order, and the first matching policy wins:

```yaml
policies:
  tools:
    # Admins can do everything
    - expression: 'has(payload.groups) && payload.groups.exists(g, g == "admins")'
      allowed_tools: ["*"]
    
    # Writers can post and read
    - expression: 'has(payload.groups) && payload.groups.exists(g, g == "writers")'
      allowed_tools:
        - "post_tweet"
        - "delete_tweet"
        - "get_*"  # Wildcard support
    
    # Readers can only read
    - expression: 'has(payload.scope) && payload.scope.contains("twitter:read")'
      allowed_tools:
        - "get_timeline"
        - "get_mentions"
        - "search_*"
```

The `allowed_tools` field supports:
- Exact matches: `"post_tweet"`
- Wildcards: `"*"` (all tools) or `"get_*"` (prefix match)

## 🌍 Location codes for trends

The `get_trends` tool uses WOEIDs (Where On Earth IDs):

| Location | WOEID |
|----------|-------|
| 🌍 Worldwide | 1 |
| 🇪🇸 Spain | 23424950 |
| 🇺🇸 United States | 23424977 |
| 🇬🇧 United Kingdom | 23424975 |
| Madrid | 766273 |
| Barcelona | 753692 |
| New York | 2459115 |

## 📄 License

Apache 2.0
