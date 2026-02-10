package cmd

import (
    "bytes"
    "context"
    "fmt"
    "time"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
    "github.com/docker/docker/pkg/stdcopy"
)

type ContainerInfo struct {
    ID     string
    Image  string
    Status string
    Names  []string
}

type ImageInfo struct {
    ID       string
    RepoTags []string
    Size     int64
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

func ListImages(ctx context.Context) ([]ImageInfo, error) {
    cli, err := NewDockerClient()
    if err != nil {
        return nil, fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
    if err != nil {
        return nil, fmt.Errorf("list images: %w", err)
    }

    results := make([]ImageInfo, 0, len(images))
    for _, image := range images {
        results = append(results, ImageInfo{
            ID:       image.ID,
            RepoTags: image.RepoTags,
            Size:     image.Size,
        })
    }

    return results, nil
}

func PrintImages(images []ImageInfo) {
    if len(images) == 0 {
        fmt.Println("No images found.")
        return
    }

    for _, image := range images {
        id := image.ID
        if len(id) > 12 {
            id = id[:12]
        }
        fmt.Printf("ID: %s | Tags: %v | Size: %d bytes\n", id, image.RepoTags, image.Size)
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

func FetchContainerLogs(ctx context.Context, containerID string, tail string) (string, error) {
    cli, err := NewDockerClient()
    if err != nil {
        return "", fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    options := types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Tail:       tail,
    }
    reader, err := cli.ContainerLogs(ctx, containerID, options)
    if err != nil {
        return "", fmt.Errorf("fetch logs: %w", err)
    }
    defer reader.Close()

    var stdout bytes.Buffer
    var stderr bytes.Buffer
    if _, err := stdcopy.StdCopy(&stdout, &stderr, reader); err != nil {
        return "", fmt.Errorf("read logs: %w", err)
    }

    combined := append(stdout.Bytes(), stderr.Bytes()...)
    return string(combined), nil
}

func InspectContainer(ctx context.Context, containerID string) (types.ContainerJSON, error) {
    cli, err := NewDockerClient()
    if err != nil {
        return types.ContainerJSON{}, fmt.Errorf("create docker client: %w", err)
    }
    defer cli.Close()

    info, err := cli.ContainerInspect(ctx, containerID)
    if err != nil {
        return types.ContainerJSON{}, fmt.Errorf("inspect container: %w", err)
    }

    return info, nil
}