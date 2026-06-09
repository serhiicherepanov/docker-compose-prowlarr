# Repository Guidelines

This repository runs Prowlarr behind Traefik and mitmproxy with Docker Compose.

## Project Layout

- `compose.yml` is the main runtime configuration.
- `plugins/lampa-categories` contains a local Traefik middleware plugin.
- `data/` is runtime state and must not be committed.

## Common Checks

Run these before committing changes:

```bash
docker compose config
cd plugins/lampa-categories && go test ./...
```

## Traefik Plugin Notes

The `lampa-categories` plugin rewrites comma-separated `categories` query values into repeated query parameters for Prowlarr search requests.

Expected rewrite:

```text
categories=2000%2C5070 -> categories=2000&categories=5070
```

Keep plugin changes narrow and preserve unrelated query parameters exactly where possible.
