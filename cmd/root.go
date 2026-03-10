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
