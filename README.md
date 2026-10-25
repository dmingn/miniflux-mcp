# miniflux-mcp

MCP server for [Miniflux](https://miniflux.app/). Talks to a Miniflux instance over its REST API and exposes a small set of read-only tools over stdio.

## Install

```bash
go install github.com/dmingn/miniflux-mcp/cmd/miniflux-mcp@latest
```

Ensure `$(go env GOBIN)` or `$(go env GOPATH)/bin` is on your `PATH`.

## Configuration

Environment variables only (no config file):

| Variable | Required | Description |
|---|---|---|
| `MINIFLUX_URL` | yes | Base URL of your Miniflux instance |
| `MINIFLUX_API_KEY` | yes | API key from Settings → API Keys |

The server fails fast if either variable is missing or authentication fails at startup.

## Cursor MCP config

```json
{
  "mcpServers": {
    "miniflux": {
      "command": "miniflux-mcp",
      "env": {
        "MINIFLUX_URL": "https://miniflux.example",
        "MINIFLUX_API_KEY": "..."
      }
    }
  }
}
```

## Tools

| Tool | Description |
|---|---|
| `list_categories` | List categories |
| `list_feeds` | List feeds (optional `category_id`) |
| `get_feed` | Get one feed by `feed_id` |
| `get_entries` | List/filter entries (`status`, `feed_id`, `category_id`, `search`, `limit`, `offset`, `order`, `direction`, `starred`) |
| `get_entry` | Get one entry by `entry_id`, including full HTML `content` |
| `get_counters` | Read/unread counters per feed |

### `get_entries` defaults

- `status`: `unread`
- `limit`: `20`

List results include `content_preview` (HTML stripped, first 100 characters) instead of full `content`, and omit the nested `feed` object (`feed_id` is enough; use `get_feed` when needed). Use `get_entry` for the full body (including nested `feed`).

## Scope

Currently read-only. Write operations (mark read, star, create/update feeds, refresh, OPML, etc.) are out of scope for now.

## Release

Distribution is `go install` only (no prebuilt binaries). Pushing a SemVer tag runs CI; when tests pass, a GitHub Release is created with [automatically generated release notes](https://docs.github.com/en/repositories/releasing-projects-on-github/automatically-generated-release-notes).

1. Merge the changes you want to ship to `main`.
2. Create and push an annotated SemVer tag (`vMAJOR.MINOR.PATCH`):

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

3. Wait for CI to pass and the Release to appear, then install that version (or `@latest`):

```bash
go install github.com/dmingn/miniflux-mcp/cmd/miniflux-mcp@v0.1.0
```

The MCP server reports the module version from Go build info (`devel` for a plain local `go build`).

## License

MIT
