package cmd

import (
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "View task comments",
	Long:  `Retrieve comments for ClickUp tasks.`,
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
