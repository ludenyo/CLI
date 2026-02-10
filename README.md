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

## Prerequisites

- Go 1.21+
- Docker running locally (Docker Desktop or Docker Engine)

## Install dependencies

```bash
go mod tidy
```

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

## Project Structure

```
.
├── cmd/
│   ├── docker.go   # Docker SDK helpers
│   └── ui.go       # Terminal UI
├── main.go         # CLI entrypoint
├── go.mod
└── README.md
```

## Ideas for Next Steps

- Auto-refresh stats in UI
- UI search/filter box
- Remove containers/images with confirmation
- Export output to CSV
- Build binary with Makefile