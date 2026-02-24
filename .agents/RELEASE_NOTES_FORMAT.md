# Release Notes Format

Use this format when writing release notes. Keep it friendly, concise, and useful.

---

## Template

```markdown
# Twitter MCP vX.Y.Z

One or two sentences describing what this release is about. For major releases, highlight the big picture. For patches, be specific about what's fixed.

## ✨ What's new

Short paragraphs or bullet points grouped by theme. Don't list every single commit—focus on what matters to users.

**Theme or area** — Description of what changed and why it's useful. If there are multiple related changes, group them together.

## 💥 Breaking changes

Only if applicable. Be clear about what breaks and how to migrate.

## 🐛 Bug fixes

Only if applicable. Brief descriptions, no need to reference issue numbers unless they add context.

## 📝 Notes

Any caveats, known issues, or things users should be aware of.

---

## 🔐 Checksums

SHA256 checksums are provided for all binaries.
```

---

## Example: v0.1.0

```markdown
# Twitter MCP v0.1.0

First release. A Model Context Protocol server that lets AI assistants interact with Twitter/X—read timelines, post tweets, analyze trends, and more.

## ✨ What's new

**Reading tools** — Fetch your home timeline, mentions, search tweets, explore trending topics by location, view user profiles and their tweets, access bookmarks and DMs.

**Writing tools** — Post tweets and threads, reply to conversations, like/unlike, retweet/undo, manage bookmarks, follow/unfollow users, send direct messages.

**Topic analysis** — Compare multiple topics at once with a heat score (0-100) that combines tweet volume and engagement. Useful for spotting what's trending in your niche.

**Two transport modes** — STDIO for simple local setups, HTTP for production with full auth support.

**Security (HTTP mode)** — JWT validation with JWKS caching, tool-level access control using CEL expressions, policies based on JWT claims (groups, scopes), OAuth 2.0 metadata endpoints (RFC 9728), and access logging with header redaction.

## 📦 Platforms

Binaries available for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64).

Docker images at `ghcr.io/achetronic/twitter-mcp:v0.1.0` for linux/amd64 and linux/arm64.

---

## 🔐 Checksums

SHA256 checksums are provided for all binaries.
```

---

## Guidelines

1. **Be human** — Write like you're explaining to a friend, not documenting for a compliance audit.

2. **Group by impact** — Users care about what they can do now, not which file you changed.

3. **Skip the noise** — Internal refactors, dependency bumps, and CI tweaks don't need to be in release notes unless they affect users.

4. **One emoji per header** — Keep it minimal. Just the section headers, not every line.

5. **Keep it scannable** — Bold the theme/area, then explain. People skim.

6. **Breaking changes up front** — If something breaks, say it clearly and early.
