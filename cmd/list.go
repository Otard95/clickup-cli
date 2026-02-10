package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage ClickUp lists",
	Long:  `View list details and tasks within lists.`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
