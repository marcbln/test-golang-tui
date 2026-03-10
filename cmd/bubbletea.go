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
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := dashboard.NewMockService()
		return btea.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(bubbleteaCmd)
}
