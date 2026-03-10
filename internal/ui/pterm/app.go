package pterm

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/username/go-cli-app/internal/dashboard"
)

// Run starts the PTerm dashboard program.
func Run(provider dashboard.MetricsProvider) error {

	// Create header
	header, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString("Dashboard"),
	).Srender()

	// Create a live area for updates
	area, _ := pterm.DefaultArea.Start()
	defer area.Stop()

	// Keyboard listener for quit
	quit := make(chan bool, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			ch, _, err := reader.ReadRune()
			if err != nil {
				continue
			}
			if ch == 'q' || ch == 'Q' {
				quit <- true
				return
			}
		}
	}()

	for {
		select {
		case <-quit:
			pterm.Println("\nGoodbye!")
			return nil
		default:
			// Get current data
			stats := provider.GetSystemStats()
			services := provider.GetServices()

			// Build stats panel content
			statsContent := fmt.Sprintf(`CPU Usage:    %.1f%%
Memory Usage: %.1f%%
Disk Usage:   %.1f%%
Uptime:       %s`,
				stats.CPUUsage,
				stats.MemoryUsage,
				stats.DiskUsage,
				stats.Uptime.Round(time.Second))

			// Build services table
			tableData := pterm.TableData{
				{"Name", "Status", "Uptime"},
			}
			for _, svc := range services {
				status := svc.Status
				if status == "Running" {
					status = pterm.Green(status)
				} else if status == "Degraded" {
					status = pterm.Yellow(status)
				} else {
					status = pterm.Red(status)
				}
				tableData = append(tableData, []string{
					svc.Name,
					status,
					svc.Uptime.Round(time.Second).String(),
				})
			}

			table, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()

			// Create panels
			leftPanel := pterm.DefaultBox.WithTitle("System Stats").Sprint(statsContent)
			rightPanel := pterm.DefaultBox.WithTitle("Active Services").Sprint(table)

			panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{
				{{Data: leftPanel}}, {{Data: rightPanel}},
			}).Srender()

			// Combine all content
			content := header + "\n" + panels + "\n" + pterm.Gray("Press 'q' to quit")

			area.Update(content)
			time.Sleep(time.Second)
		}
	}
}
