package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage ClickUp tasks",
	Long:  `Search, view, update tasks and manage subtasks and relationships.`,
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
