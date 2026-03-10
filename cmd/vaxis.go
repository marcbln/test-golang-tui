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
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := dashboard.NewMockService()
		return vaxis.Run(provider)
	},
}

func init() {
	rootCmd.AddCommand(vaxisCmd)
}
