# httpServer

Minimal Go HTTP server that serves static files and a few simple API endpoints.

## Overview

- Serves files from the project root under the `/app/` prefix.
- Exposes simple endpoints under `/api/` for health, metrics and reset.
- Default port: `8080`.

## Prerequisites

- Go (1.18+ recommended).

## Build & Run

From the project root:

```bash
# Run without building
go run .

# Build and run
go build -o httpserver .
./httpserver    # Unix
.\httpserver.exe # Windows
```

The server listens on port `8080` by default.

## Endpoints

- Static files: `GET /app/<path>` — serves files from the repository root. Example: `/app/index.html`.
- Health: `GET /api/healthz` — readiness/health endpoint.
- Metrics: `GET /api/metrics` — returns a simple text metric (file server hit count).
- Reset: `POST /api/reset` — resets metrics/state (if implemented).

## Examples

Fetch index:

```bash
curl http://localhost:8080/app/index.html
```

Check health:

```bash
curl http://localhost:8080/api/healthz
```

View metrics:

```bash
curl http://localhost:8080/api/metrics
```

Reset metrics:

```bash
curl -X POST http://localhost:8080/api/reset
```

## Notes

- The server maps `/app/` to the current working directory (`.`) so any files (for example `index.html` or the `assets/` folder) will be served there.
- If you want a different port or file root, update `main.go` constants.

## License

This project has no license; add one if you plan to publish.
