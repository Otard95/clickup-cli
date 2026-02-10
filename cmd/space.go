package cmd

import (
	"github.com/spf13/cobra"
)

var spaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Manage ClickUp spaces",
	Long:  `Search spaces and view their folder/list structure.`,
}

func init() {
	rootCmd.AddCommand(spaceCmd)
}
