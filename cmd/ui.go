package cmd

import (
    "context"
    "fmt"
    "strings"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

func ShowUI(ctx context.Context) error {
    app := tview.NewApplication()
    list := tview.NewList().ShowSecondaryText(false)
    details := tview.NewTextView().SetDynamicColors(true)
    details.SetBorder(true).SetTitle("Details")
    footer := tview.NewTextView().SetTextAlign(tview.AlignCenter)
    footer.SetText("R=refresh | S=start | T=stop | Q=quit")
    status := tview.NewTextView().SetTextAlign(tview.AlignLeft)
    status.SetText("Ready")

    type listItem struct {
        id string
    }
    items := make([]listItem, 0)

    refresh := func() {
        list.Clear()
        items = items[:0]
        details.SetText("Select a container to view details.")
        containers, err := ListContainers(ctx)
        if err != nil {
            list.AddItem("Error", err.Error(), 0, nil)
            return
        }
        if len(containers) == 0 {
            list.AddItem("No containers found", "", 0, nil)
            return
        }

        for _, container := range containers {
            id := container.ID
            if len(id) > 12 {
                id = id[:12]
            }
            title := fmt.Sprintf("%s (%s)", strings.Join(container.Names, ","), id)
            details := fmt.Sprintf("%s | %s", container.Image, container.Status)
            list.AddItem(title, details, 0, nil)
            items = append(items, listItem{id: container.ID})
        }
    }

    refresh()

    list.SetDoneFunc(func() {
        app.Stop()
    })

    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Rune() {
        case 'q', 'Q':
            app.Stop()
            return nil
        case 'r', 'R':
            refresh()
            return nil
        case 's', 'S':
            index := list.GetCurrentItem()
            if index >= 0 && index < len(items) {
                status.SetText("Starting container...")
                if err := StartContainer(ctx, items[index].id); err != nil {
                    status.SetText(fmt.Sprintf("Start failed: %v", err))
                } else {
                    status.SetText("Container started.")
                    refresh()
                }
            }
            return nil
        case 't', 'T':
            index := list.GetCurrentItem()
            if index >= 0 && index < len(items) {
                status.SetText("Stopping container...")
                if err := StopContainer(ctx, items[index].id); err != nil {
                    status.SetText(fmt.Sprintf("Stop failed: %v", err))
                } else {
                    status.SetText("Container stopped.")
                    refresh()
                }
            }
            return nil
        }

        if event.Key() == tcell.KeyEscape {
            app.Stop()
            return nil
        }

        return event
    })

    list.SetChangedFunc(func(index int, _ string, _ string, _ rune) {
        if index >= 0 && index < len(items) {
            info, err := InspectContainer(ctx, items[index].id)
            if err != nil {
                details.SetText(fmt.Sprintf("[red]Error:[-] %v", err))
                return
            }

            stats, err := GetContainerStats(ctx, items[index].id)
            if err != nil {
                stats = ContainerStats{}
            }

            memoryLine := "N/A"
            if stats.MemoryLimit > 0 {
                memoryLine = fmt.Sprintf("%s / %s (%.2f%%)", formatBytes(stats.MemoryUsage), formatBytes(stats.MemoryLimit), stats.MemoryPercent)
            }

            details.SetText(fmt.Sprintf(
                "[yellow]ID:[-] %s\n[yellow]Name:[-] %s\n[yellow]Image:[-] %s\n[yellow]Status:[-] %s\n[yellow]Created:[-] %s\n[yellow]Ports:[-] %v\n[yellow]CPU:[-] %.2f%%\n[yellow]Memory:[-] %s",
                info.ID,
                info.Name,
                info.Config.Image,
                info.State.Status,
                info.Created,
                info.NetworkSettings.Ports,
                stats.CPUPercent,
                memoryLine,
            ))
        }
    })

    body := tview.NewFlex().SetDirection(tview.FlexColumn).
        AddItem(list, 0, 1, true).
        AddItem(details, 0, 2, false)

    layout := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(body, 0, 1, true).
        AddItem(status, 1, 0, false).
        AddItem(footer, 1, 0, false)

    return app.SetRoot(layout, true).Run()
}

func formatBytes(value uint64) string {
    const (
        kb = 1024
        mb = 1024 * kb
        gb = 1024 * mb
    )

    switch {
    case value >= gb:
        return fmt.Sprintf("%.2f GB", float64(value)/float64(gb))
    case value >= mb:
        return fmt.Sprintf("%.2f MB", float64(value)/float64(mb))
    case value >= kb:
        return fmt.Sprintf("%.2f KB", float64(value)/float64(kb))
    default:
        return fmt.Sprintf("%d B", value)
    }
}