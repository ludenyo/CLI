# Docker CLI + UI (Go)

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
go run . list --running
go run . list --name redis
go run . list --json
go run . start <container-id>
go run . stop <container-id>
go run . logs <container-id>
go run . logs <container-id> --tail 50
go run . inspect <container-id>
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
