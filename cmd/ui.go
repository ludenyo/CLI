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

    layout := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(list, 0, 1, true).
        AddItem(status, 1, 0, false).
        AddItem(footer, 1, 0, false)

    return app.SetRoot(layout, true).Run()
}