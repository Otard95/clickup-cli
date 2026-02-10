package cmd

import (
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "View time tracking entries",
	Long:  `Get time tracking entries for tasks or teams.`,
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
