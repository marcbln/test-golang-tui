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
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := dashboard.NewMockService()
		return tview.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(tviewCmd)
}
