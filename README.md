# Docker CLI + UI (Go)

Beginner-friendly Go app that lists Docker containers, interacts with them, and provides a simple terminal UI.

## Features

### CLI
- List containers with filters
- Start/stop containers
- View logs (tail)
- Inspect container JSON
- List Docker images
- JSON output for scripting

### UI (Terminal)
- Split view: container list + details pane
- Details include CPU & memory stats (human-readable)
- Actions: refresh, start, stop, quit
- Configurable UI colors via `config.json`

### Build & Release
- `Makefile` with simple commands
- Local binary build
- Cross-platform release binaries (macOS/Linux/Windows)

### Config Defaults
- Save default list filters
- Save default logs tail value
- Save UI theme colors

## Prerequisites

- Go 1.21+
- Docker running locally (Docker Desktop or Docker Engine)

## Install dependencies

```bash
go mod tidy
```

## Configuration (`config.json`)

The app loads defaults from `config.json` (if missing, safe defaults are used).

```json
{
  "list": {
    "running": false,
    "name": "",
    "json": false,
    "logsTail": "100"
  },
  "ui": {
    "footerColor": "white",
    "statusColor": "white",
    "accentColor": "yellow"
  }
}
```

### Config behavior
- `list.running`, `list.name`, `list.json` are used as default values for `list` flags.
- `list.logsTail` is used as default for `logs --tail`.
- `ui.*Color` sets UI colors (footer, status, details border/title accent).
- CLI flags still override config values.

## CLI Usage

```bash
# List containers
go run . list
go run . list --running
go run . list --name redis
go run . list --json

# Start/stop containers
go run . start <container-id>
go run . stop <container-id>

# Logs & inspect
go run . logs <container-id>
go run . logs <container-id> --tail 50
go run . inspect <container-id>

# Images
go run . images
```

## UI Usage

```bash
go run . ui
```

### UI Controls

- `R` refresh containers
- `S` start selected container
- `T` stop selected container
- `Q` or `Esc` quit

### UI Layout

- Left pane: container list
- Right pane: details of selected container (CPU/Memory included)

## Build & Release (KISS)

Use the included `Makefile`:

```bash
# Format code
make fmt

# Build local binary (./cli)
make build

# Build cross-platform binaries in ./dist
make release

# Clean build artifacts
make clean
```

### Release outputs

`make release` creates:
- `dist/cli-darwin-amd64`
- `dist/cli-darwin-arm64`
- `dist/cli-linux-amd64`
- `dist/cli-windows-amd64.exe`

## Project Structure

```
.
├── config.go        # Config models + loader
├── config.json      # Default app config
├── cmd/
│   ├── docker.go   # Docker SDK helpers
│   └── ui.go       # Terminal UI + theme support
├── Makefile         # Build/release helpers
├── main.go         # CLI entrypoint
├── go.mod
└── README.md
```