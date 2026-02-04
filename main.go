package main

import (
    "context"
    "fmt"
    "os"

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
        containers, err := cmd.ListContainers(ctx)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        cmd.PrintContainers(containers)
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
    fmt.Println("  go run . list             # List containers")
    fmt.Println("  go run . start <id>       # Start container")
    fmt.Println("  go run . stop <id>        # Stop container")
    fmt.Println("  go run . ui               # Launch terminal UI")
}
