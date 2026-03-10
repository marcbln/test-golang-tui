package vaxis

import (
	"fmt"
	"time"

	"git.sr.ht/~rockorager/vaxis"
	"github.com/username/go-cli-app/internal/dashboard"
)

// Run starts the Vaxis program.
func Run(provider dashboard.MetricsProvider) error {
	vx, err := vaxis.New(vaxis.Options{})
	if err != nil {
		return err
	}
	defer vx.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Initial render
	draw(vx, provider)

	for {
		select {
		case ev := <-vx.Events():
			switch ev := ev.(type) {
			case vaxis.Key:
				if ev.String() == "q" || ev.String() == "Ctrl+c" {
					return nil
				}
			case vaxis.Resize:
				draw(vx, provider)
			}
		case <-ticker.C:
			draw(vx, provider)
		}
	}
}

func draw(vx *vaxis.Vaxis, provider dashboard.MetricsProvider) {
	win := vx.Window()
	win.Clear()

	stats := provider.GetSystemStats()
	services := provider.GetServices()

	printAt(win, 0, 0, "=== System Dashboard (Vaxis) ===")

	statsText := fmt.Sprintf("CPU: %.1f%% | Mem: %.1f%% | Disk: %.1f%% | Uptime: %s",
		stats.CPUUsage, stats.MemoryUsage, stats.DiskUsage, stats.Uptime.Round(time.Second))
	printAt(win, 0, 2, statsText)

	printAt(win, 0, 4, "Active Services:")
	for i, s := range services {
		svcText := fmt.Sprintf("  - %-15s | %-8s | %s", s.Name, s.Status, s.Uptime.Round(time.Second))
		printAt(win, 0, 5+i, svcText)
	}

	printAt(win, 0, 7+len(services), "(Press 'q' to quit)")

	vx.Render()
}

// printAt is a helper to draw strings to the vaxis window at specific coordinates.
func printAt(win vaxis.Window, x, y int, text string) {
	chars := vaxis.Characters(text)
	for i, c := range chars {
		win.SetCell(x+i, y, vaxis.Cell{
			Character: c,
		})
	}
}
