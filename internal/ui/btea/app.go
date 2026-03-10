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
	services []dashboard.ServiceStatus
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
