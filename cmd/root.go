package cmd

import (
	"fmt"
	"os"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/otard95/clickup-cli/internal/config"
	"github.com/spf13/cobra"
)

var client *api.Client

var rootCmd = &cobra.Command{
	Use:   "clickup-cli",
	Short: "CLI for interacting with ClickUp",
	Long:  `A command-line interface for ClickUp project management — tasks, lists, spaces, comments, time tracking, and documents.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}
		client = api.NewClient(cfg)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
