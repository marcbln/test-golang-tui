---
filename: "_ai/backlog/active/250310_2116__IMPLEMENTATION_PLAN__add-pterm-dashboard-program.md"
title: "Add PTerm Dashboard Program"
createdAt: 2025-03-10 21:16
createdBy: Cascade [claude-sonnet-4-20250514]
updatedAt: 2025-03-10 21:16
updatedBy: Cascade [claude-sonnet-4-20250514]
status: draft
priority: medium
tags: [golang, pterm, tui, dashboard]
estimatedComplexity: simple
documentType: IMPLEMENTATION_PLAN
---

# Add PTerm Dashboard Program

## Problem Statement
The project currently has three TUI frontend implementations (Bubble Tea, Tview, and Vaxis). A fourth program using the [pterm](https://github.com/pterm/pterm) library needs to be added to demonstrate another popular terminal UI library option for building dashboard-style CLI applications.

## Implementation Notes

### Project Structure
- Module: `github.com/username/go-cli-app`
- Commands location: `cmd/` (cobra-based)
- UI implementations: `internal/ui/{library}/`
- Dashboard service: `internal/dashboard/service.go` provides mock metrics

### Existing Pattern
Each UI implementation follows this structure:
1. Command file in `cmd/{name}.go` - cobra command that calls the UI Run function
2. UI implementation in `internal/ui/{name}/app.go` - contains a `Run(provider dashboard.MetricsProvider) error` function
3. Uses the `dashboard.MetricsProvider` interface to get system stats and services

### Dependencies
- Need to add: `github.com/pterm/pterm` to go.mod
- PTerms is a modern, pluggable TUI framework with built-in components like progress bars, tables, panels, and spinners

---

## Phase 1: Add PTerm Dependency

### Objective
Add the pterm library to the project dependencies.

### Tasks
1. Run `go get github.com/pterm/pterm` to add the dependency
2. Verify go.mod and go.sum are updated correctly

### Deliverables
- Updated `go.mod` with pterm dependency
- Updated `go.sum` with checksums

---

## Phase 2: Create PTerm UI Implementation

### Objective
Create the UI implementation package for pterm at `internal/ui/pterm/app.go`.

### Tasks
1. Create directory `internal/ui/pterm/`
2. Create `app.go` with a dashboard implementation using pterm components

### Implementation Details
The implementation should:
- Use `pterm.Panel` for layout (header, stats panel, services table, footer)
- Use `pterm.Table` for displaying services with status colors
- Use `pterm.BigText` or `pterm.Header` for the dashboard title
- Use a goroutine with ticker for live updates (similar to existing implementations)
- Support 'q' key to quit

### Deliverables
[NEW FILE] `internal/ui/pterm/app.go`
```go
package pterm

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/username/go-cli-app/internal/dashboard"
)

// Run starts the PTerm dashboard program.
func Run(provider dashboard.MetricsProvider) error {
	pterm.EnableDebugMessages = false

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
		// Simple stdin reader for 'q'
		var input string
		for {
			fmt.Scanln(&input)
			if input == "q" || input == "Q" {
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
		case <-time.After(time.Second):
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
		}
	}
}
```

---

## Phase 3: Create PTerm Command

### Objective
Create the cobra command file for the pterm subcommand.

### Tasks
1. Create `cmd/pterm.go` following the same pattern as other commands

### Deliverables
[NEW FILE] `cmd/pterm.go`
```go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/username/go-cli-app/internal/dashboard"
	"github.com/username/go-cli-app/internal/ui/pterm"
)

var ptermCmd = &cobra.Command{
	Use:          "pterm",
	Short:        "Run the PTerm frontend",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := dashboard.NewMockService()
		return pterm.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(ptermCmd)
}
```

---

## Phase 4: Verification

### Objective
Verify the implementation compiles and works correctly.

### Tasks
1. Run `go build ./...` to verify compilation
2. Run `go run . pterm` to test the program

### Deliverables
- Successful compilation
- Working pterm dashboard display

---

## Phase 5: Write Implementation Report

### Objective
Document the completed work in the reports directory.

### Tasks
1. Create report file at `_ai/backlog/reports/250310_2116__IMPLEMENTATION_REPORT__add-pterm-dashboard-program.md`

### Deliverables
[NEW FILE] `_ai/backlog/reports/250310_2116__IMPLEMENTATION_REPORT__add-pterm-dashboard-program.md`
```yaml
---
filename: "_ai/backlog/reports/250310_2116__IMPLEMENTATION_REPORT__add-pterm-dashboard-program.md"
title: "Report: Add PTerm Dashboard Program"
createdAt: 2025-03-10 21:16
createdBy: Cascade [claude-sonnet-4-20250514]
updatedAt: 2025-03-10 21:16
updatedBy: Cascade [claude-sonnet-4-20250514]
planFile: "_ai/backlog/active/250310_2116__IMPLEMENTATION_PLAN__add-pterm-dashboard-program.md"
project: "test-golang-tui"
status: completed
filesCreated: 3
filesModified: 2
filesDeleted: 0
tags: [golang, pterm, tui, dashboard]
documentType: IMPLEMENTATION_REPORT
---

## Summary
Successfully added a fourth TUI frontend using the pterm library to demonstrate dashboard-style CLI applications. The implementation includes live updates, service status tables, and system stats display.

## Files Changed

### New Files Created
- `cmd/pterm.go` - Cobra command for the pterm subcommand
- `internal/ui/pterm/app.go` - UI implementation using pterm components
- `_ai/backlog/reports/250310_2116__IMPLEMENTATION_REPORT__add-pterm-dashboard-program.md` - This report

### Modified Files
- `go.mod` - Added pterm dependency
- `go.sum` - Added pterm checksums

## Key Changes
- Added `github.com/pterm/pterm` v0.12.79 dependency
- Created pterm UI package with Run function accepting MetricsProvider
- Implemented live dashboard with panels showing system stats and service status
- Used pterm's Panel, Box, Table, and BigText components
- Added keyboard input handling for 'q' to quit

## Technical Decisions
- Used pterm.DefaultArea for live updates instead of raw terminal manipulation
- Implemented simple stdin-based quit mechanism (q key)
- Used pterm's built-in color helpers for status indicators
- Followed existing pattern of 1-second update intervals

## Testing Notes
Build and run with:
```bash
go build ./...
go run . pterm
```

Press 'q' to exit the dashboard.

## Documentation Updates
- No README updates required as this is a demonstration addition
- Code follows existing patterns and conventions

## Next Steps (Optional)
- Consider adding more pterm features like spinners or progress bars
- Add color theme configuration options
- Implement pagination for large service lists
```

---

## Verification Steps

1. Build the project: `go build ./...`
2. Run the pterm command: `go run . pterm`
3. Verify the dashboard displays:
   - Big "Dashboard" header
   - System Stats panel with CPU, Memory, Disk, and Uptime
   - Active Services table with Name, Status, and Uptime columns
   - Status colors (green=Running, yellow=Degraded, red=Stopped)
4. Press 'q' to quit the program
