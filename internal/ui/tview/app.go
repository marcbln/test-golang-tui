package tview

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/username/go-cli-app/internal/dashboard"
)

// Run starts the Tview program.
func Run(provider dashboard.MetricsProvider) error {
	app := tview.NewApplication()

	statsBox := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statsBox.SetBorder(true).SetTitle(" System Stats ")

	servicesBox := tview.NewTable().
		SetBorders(true)
	servicesBox.SetBorder(true).SetTitle(" Active Services ")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("=== System Dashboard (Tview) ===").SetTextAlign(tview.AlignCenter), 3, 1, false).
		AddItem(tview.NewFlex().
			AddItem(statsBox, 0, 1, false).
			AddItem(servicesBox, 0, 2, false),
			0, 1, false).
		AddItem(tview.NewTextView().SetText("Press 'q' to quit").SetTextAlign(tview.AlignCenter), 1, 1, false)

	updateUI := func() {
		stats := provider.GetSystemStats()
		statsText := fmt.Sprintf("[yellow]CPU:[white] %.1f%%\n[yellow]Mem:[white] %.1f%%\n[yellow]Disk:[white] %.1f%%\n[yellow]Uptime:[white] %s",
			stats.CPUUsage, stats.MemoryUsage, stats.DiskUsage, stats.Uptime.Round(time.Second))
		statsBox.SetText(statsText)

		servicesBox.Clear()
		servicesBox.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		servicesBox.SetCell(0, 1, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		servicesBox.SetCell(0, 2, tview.NewTableCell("Uptime").SetTextColor(tcell.ColorYellow).SetSelectable(false))

		for i, s := range provider.GetServices() {
			color := tcell.ColorGreen
			if s.Status != "Running" {
				color = tcell.ColorRed
			}
			servicesBox.SetCell(i+1, 0, tview.NewTableCell(s.Name).SetTextColor(tcell.ColorWhite))
			servicesBox.SetCell(i+1, 1, tview.NewTableCell(s.Status).SetTextColor(color))
			servicesBox.SetCell(i+1, 2, tview.NewTableCell(s.Uptime.Round(time.Second).String()).SetTextColor(tcell.ColorWhite))
		}
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	updateUI()

	go func() {
		for {
			time.Sleep(time.Second)
			app.QueueUpdateDraw(func() {
				updateUI()
			})
		}
	}()

	return app.SetRoot(flex, true).Run()
}
