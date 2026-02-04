# Docker CLI + UI (Go)

Beginner-friendly Go app that lists Docker containers, can start/stop them, and offers a simple terminal UI.

## Prerequisites

- Go 1.21+
- Docker running locally (Docker Desktop or Docker Engine)

## Install dependencies

```bash
go mod tidy
```

## CLI Usage

```bash
go run . list
go run . start <container-id>
go run . stop <container-id>
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
