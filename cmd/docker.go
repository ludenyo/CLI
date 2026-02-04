package cmd

import (
    "context"
    "fmt"
    "time"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
)

type ContainerInfo struct {
    ID     string
    Image  string
    Status string
    Names  []string
}

func NewDockerClient() (*client.Client, error) {
    return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func ListContainers(ctx context.Context) ([]ContainerInfo, error) {
    cli, err := NewDockerClient()
    if err != nil {
        return nil, fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
    if err != nil {
        return nil, fmt.Errorf("list containers: %w", err)
    }

    results := make([]ContainerInfo, 0, len(containers))
    for _, container := range containers {
        results = append(results, ContainerInfo{
            ID:     container.ID,
            Image:  container.Image,
            Status: container.Status,
            Names:  container.Names,
        })
    }

    return results, nil
}

func PrintContainers(containers []ContainerInfo) {
    if len(containers) == 0 {
        fmt.Println("No containers found.")
        return
    }

    for _, container := range containers {
        id := container.ID
        if len(id) > 12 {
            id = id[:12]
        }
        fmt.Printf("ID: %s | Image: %s | Status: %s | Names: %v\n", id, container.Image, container.Status, container.Names)
    }
}

func StartContainer(ctx context.Context, containerID string) error {
    cli, err := NewDockerClient()
    if err != nil {
        return fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
        return fmt.Errorf("start container: %w", err)
    }

    return nil
}

func StopContainer(ctx context.Context, containerID string) error {
    cli, err := NewDockerClient()
    if err != nil {
        return fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    timeout := int((10 * time.Second).Seconds())
    if err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
        return fmt.Errorf("stop container: %w", err)
    }

    return nil
}