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
