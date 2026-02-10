package cmd

import (
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Manage ClickUp documents",
	Long:  `Read and search ClickUp documents.`,
}

func init() {
	rootCmd.AddCommand(docCmd)
}
