---
filename: "_ai/backlog/active/260310_1834__IMPLEMENTATION_PLAN__tui-evaluation.md"
title: "Evaluate TUI Libraries: Bubbletea, Tview, Vaxis"
createdAt: 2026-03-10 18:34
updatedAt: 2026-03-10 18:34
status: draft
priority: medium
tags:[tui, golang, bubbletea, tview, vaxis]
estimatedComplexity: moderate
documentType: IMPLEMENTATION_PLAN
---

# Problem Statement
We need to evaluate three different Terminal User Interface (TUI) libraries for Go: **Bubbletea**, **Tview**, and **Vaxis**. To compare them effectively, we will build a single CLI application that provides a shared backend service (a mocked System & Service Dashboard) and three separate commands (`bubbletea`, `tview`, `vaxis`) that launch a specialized frontend using each respective library. This allows us to assess the API ergonomics, widget availability, and styling mechanics of each library using the same underlying data constraints.

# Project Environment
- Project Name: Go CLI App
- Module Name: `github.com/username/go-cli-app`
- Go Version: 1.23+
- Frameworks: Cobra (CLI), Bubbletea, Tview, Vaxis

> **Note:** This plan will be implemented by an AI coding agent. The agent should create the directories, write the source files as provided below, resolve dependencies with `go mod tidy`, verify build correctness, and finally write the implementation report.

---

## Phase 1: Project Setup & Shared Services
We will initialize the Go project, define the Cobra root command, and establish the shared business logic following SOLID principles (specifically using a `MetricsProvider` interface so frontends depend on abstractions rather than concrete structs).

```go
// [NEW FILE] go.mod
module github.com/username/go-cli-app

go 1.23
```

```go
// [NEW FILE] main.go
package main

import "github.com/username/go-cli-app/cmd"

func main() {
	cmd.Execute()
}
```

```go
// [NEW FILE] cmd/root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A CLI tool to evaluate TUI libraries",
	Long:  `This application provides 3 frontends (bubbletea, tview, vaxis) for a shared dashboard service to evaluate their respective developer experiences.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be defined here
}
```

```go
// [NEW FILE] internal/dashboard/service.go
package dashboard

import (
	"math/rand"
	"time"
)

// SystemStats holds the current system usage metrics.
type SystemStats struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	Uptime      time.Duration
}

// ServiceStatus holds the health and uptime of an individual service.
type ServiceStatus struct {
	Name   string
	Status string
	Uptime time.Duration
}

// MetricsProvider defines the contract for fetching dashboard data.
type MetricsProvider interface {
	GetSystemStats() SystemStats
	GetServices()[]ServiceStatus
}

// MockService implements MetricsProvider for testing and demonstration.
type MockService struct {
	startTime time.Time
}

// NewMockService creates a new MockService.
func NewMockService() *MockService {
	return &MockService{startTime: time.Now()}
}

// GetSystemStats returns randomized mock system metrics.
func (s *MockService) GetSystemStats() SystemStats {
	return SystemStats{
		CPUUsage:    rand.Float64() * 100,
		MemoryUsage: rand.Float64() * 100,
		DiskUsage:   rand.Float64() * 100,
		Uptime:      time.Since(s.startTime),
	}
}

// GetServices returns a list of mock service health statuses.
func (s *MockService) GetServices() []ServiceStatus {
	return[]ServiceStatus{
		{"Web Server", "Running", time.Since(s.startTime) + time.Hour},
		{"Database", "Running", time.Since(s.startTime) + time.Hour*24},
		{"Cache", "Degraded", time.Since(s.startTime) + time.Minute*30},
		{"Worker Node", "Stopped", 0},
	}
}
```

---

## Phase 2: Bubbletea Frontend Implementation
We'll implement the Bubbletea frontend using The Elm Architecture (Model, Update, View). We'll use `lipgloss` for styling.

```go
// [NEW FILE] cmd/bubbletea.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/username/go-cli-app/internal/dashboard"
	"github.com/username/go-cli-app/internal/ui/btea"
)

var bubbleteaCmd = &cobra.Command{
	Use:          "bubbletea",
	Short:        "Run the Bubble Tea frontend",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args[]string) error {
		provider := dashboard.NewMockService()
		return btea.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(bubbleteaCmd)
}
```

```go
// [NEW FILE] internal/ui/btea/app.go
package btea

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/username/go-cli-app/internal/dashboard"
)

type tickMsg time.Time

type model struct {
	provider dashboard.MetricsProvider
	stats    dashboard.SystemStats
	services[]dashboard.ServiceStatus
}

func initialModel(provider dashboard.MetricsProvider) model {
	return model{
		provider: provider,
		stats:    provider.GetSystemStats(),
		services: provider.GetServices(),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		m.stats = m.provider.GetSystemStats()
		m.services = m.provider.GetServices()
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	statsStyle := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1).MarginBottom(1)

	title := titleStyle.Render("=== System Dashboard (Bubble Tea) ===")

	statsStr := fmt.Sprintf("CPU: %.1f%%\nMem: %.1f%%\nDisk: %.1f%%\nUptime: %s",
		m.stats.CPUUsage, m.stats.MemoryUsage, m.stats.DiskUsage, m.stats.Uptime.Round(time.Second))

	servicesStr := lipgloss.NewStyle().Bold(true).Render("Active Services:\n")
	for _, s := range m.services {
		statusColor := "46" // Green
		if s.Status != "Running" {
			statusColor = "196" // Red
		}
		status := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render(s.Status)
		servicesStr += fmt.Sprintf("  • %s: %s (Uptime: %s)\n", s.Name, status, s.Uptime.Round(time.Second))
	}

	help := lipgloss.NewStyle().Faint(true).Render("\n(Press 'q' to quit)")

	return fmt.Sprintf("%s\n%s\n%s%s\n", title, statsStyle.Render(statsStr), servicesStr, help)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Run starts the Bubble Tea program.
func Run(provider dashboard.MetricsProvider) error {
	p := tea.NewProgram(initialModel(provider))
	_, err := p.Run()
	return err
}
```

---

## Phase 3: Tview Frontend Implementation
We'll implement the Tview frontend utilizing its rich widget ecosystem (`Flex`, `TextView`, `Table`).

```go
// [NEW FILE] cmd/tview.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/username/go-cli-app/internal/dashboard"
	"github.com/username/go-cli-app/internal/ui/tview"
)

var tviewCmd = &cobra.Command{
	Use:          "tview",
	Short:        "Run the Tview frontend",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args[]string) error {
		provider := dashboard.NewMockService()
		return tview.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(tviewCmd)
}
```

```go
//[NEW FILE] internal/ui/tview/app.go
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
```

---

## Phase 4: Vaxis Frontend Implementation
We'll implement the Vaxis frontend utilizing its high-performance event loop and direct cell rendering capabilities.

```go
// [NEW FILE] cmd/vaxis.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/username/go-cli-app/internal/dashboard"
	"github.com/username/go-cli-app/internal/ui/vaxis"
)

var vaxisCmd = &cobra.Command{
	Use:          "vaxis",
	Short:        "Run the Vaxis frontend",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args[]string) error {
		provider := dashboard.NewMockService()
		return vaxis.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(vaxisCmd)
}
```

```go
// [NEW FILE] internal/ui/vaxis/app.go
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
	for i, r := range text {
		win.SetCell(x+i, y, vaxis.Cell{
			Character: vaxis.Character{Rune: r},
		})
	}
}
```

---

## Phase 5: Build, Tidy, and Documentation
The AI agent should resolve the dependencies by running standard Go toolchains, and generate the final documentation.

**Agent Instructions for this Phase:**
1. Run `go mod tidy` in the root directory to fetch all dependencies (`bubbletea`, `tview`, `vaxis`, `cobra`).
2. Run `go build -o myapp main.go` to ensure everything compiles correctly.

```markdown
<!-- [NEW FILE] README.md -->
# Go TUI Evaluator

A sandbox application to evaluate different Go TUI libraries (Bubbletea, Tview, Vaxis) displaying a shared mocked dashboard service.

## Usage

Build the application:
```bash
go build -o myapp main.go
```

Run the different frontends:
```bash
./myapp bubbletea
./myapp tview
./myapp vaxis
```

## Architecture
- `internal/dashboard/` holds the decoupled mock data service.
- `internal/ui/` contains isolated UI implementations for each framework.
- The UI implementations adhere strictly to the shared `MetricsProvider` interface.
```

---

## Phase 6: Report Generation
Once all steps above have been executed successfully, write the final implementation report.

Create the file `_ai/backlog/reports/{YYMMDD_HHmm}__IMPLEMENTATION_REPORT__tui-evaluation.md` using the following structure:

```yaml
---
filename: "_ai/backlog/reports/{YYMMDD_HHmm}__IMPLEMENTATION_REPORT__tui-evaluation.md"
title: "Report: Evaluate TUI Libraries"
createdAt: YYYY-MM-DD HH:mm
updatedAt: YYYY-MM-DD HH:mm
planFile: "_ai/backlog/active/260310_1834__IMPLEMENTATION_PLAN__tui-evaluation.md"
project: "Go CLI App"
status: completed
filesCreated: 9
filesModified: 0
filesDeleted: 0
tags:[tui, golang, bubbletea, tview, vaxis]
documentType: IMPLEMENTATION_REPORT
---

## Summary
Briefly describe how the Cobra application and the three TUI frontends were successfully created sharing the decoupled `MetricsProvider`.

## Files Changed
- (List the files created in `cmd/`, `internal/`, etc.)

## Key Changes
- Scaffolded Cobra root command.
- Decoupled data fetching via `MetricsProvider` to enforce SOLID principles.
- Built Elm-architecture dashboard for `bubbletea`.
- Built Widget-driven dashboard for `tview`.
- Built Direct-cell rendering dashboard for `vaxis`.

## Technical Decisions
- Used `lipgloss` alongside Bubbletea for layout handling.
- Tview utilized `Flex` and `Table` components for out-of-the-box UI elements.
- Vaxis implemented a manual coordinate-based render function `printAt` leveraging its high-performance frame buffer interface.

## Next Steps
- Implement real backend stats fetching using `shirou/gopsutil`.

