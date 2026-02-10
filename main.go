package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/ludenyo/cli/cmd"
)

func main() {
    ctx := context.Background()
    if len(os.Args) < 2 {
        printUsage()
        return
    }

    switch os.Args[1] {
    case "list":
        listFlags := flag.NewFlagSet("list", flag.ExitOnError)
        runningOnly := listFlags.Bool("running", false, "Show running containers only")
        nameFilter := listFlags.String("name", "", "Filter containers by name or image")
        jsonOutput := listFlags.Bool("json", false, "Output JSON")
        _ = listFlags.Parse(os.Args[2:])

        containers, err := cmd.ListContainers(ctx)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        filtered := filterContainers(containers, *runningOnly, *nameFilter)
        if *jsonOutput {
            outputJSON(filtered)
            return
        }
        cmd.PrintContainers(filtered)
    case "start":
        if len(os.Args) < 3 {
            fmt.Println("Please provide a container ID.")
            return
        }
        if err := cmd.StartContainer(ctx, os.Args[2]); err != nil {
            fmt.Println("Error:", err)
        }
    case "stop":
        if len(os.Args) < 3 {
            fmt.Println("Please provide a container ID.")
            return
        }
        if err := cmd.StopContainer(ctx, os.Args[2]); err != nil {
            fmt.Println("Error:", err)
        }
	case "logs":
		logsFlags := flag.NewFlagSet("logs", flag.ExitOnError)
		tail := logsFlags.String("tail", "100", "Number of lines to show")
		_ = logsFlags.Parse(os.Args[2:])
		args := logsFlags.Args()
		if len(args) < 1 {
			fmt.Println("Please provide a container ID.")
			return
		}
		output, err := cmd.FetchContainerLogs(ctx, args[0], *tail)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(output)
	case "inspect":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a container ID.")
			return
		}
		info, err := cmd.InspectContainer(ctx, os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(data))
	case "images":
		images, err := cmd.ListImages(ctx)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		cmd.PrintImages(images)
    case "ui":
        if err := cmd.ShowUI(ctx); err != nil {
            fmt.Println("Error:", err)
        }
    default:
        printUsage()
    }
}

func printUsage() {
    fmt.Println("Usage:")
    fmt.Println("  go run . list [flags]     # List containers")
    fmt.Println("  go run . start <id>       # Start container")
    fmt.Println("  go run . stop <id>        # Stop container")
	fmt.Println("  go run . logs <id> [--tail N]  # Show container logs")
	fmt.Println("  go run . inspect <id>     # Inspect container JSON")
	fmt.Println("  go run . images           # List Docker images")
    fmt.Println("  go run . ui               # Launch terminal UI")
    fmt.Println("\nList flags:")
    fmt.Println("  --running                 Show only running containers")
    fmt.Println("  --name <text>             Filter by name or image")
    fmt.Println("  --json                    Output JSON")
}

func filterContainers(containers []cmd.ContainerInfo, runningOnly bool, nameFilter string) []cmd.ContainerInfo {
    if !runningOnly && nameFilter == "" {
        return containers
    }

    filtered := make([]cmd.ContainerInfo, 0, len(containers))
    normalized := strings.ToLower(nameFilter)
    for _, container := range containers {
        if runningOnly && !strings.HasPrefix(strings.ToLower(container.Status), "up") {
            continue
        }

        if normalized != "" {
            haystack := strings.ToLower(container.Image + " " + strings.Join(container.Names, " "))
            if !strings.Contains(haystack, normalized) {
                continue
            }
        }

        filtered = append(filtered, container)
    }

    return filtered
}

func outputJSON(containers []cmd.ContainerInfo) {
    data, err := json.MarshalIndent(containers, "", "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data))
}
