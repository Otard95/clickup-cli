package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var docReadCmd = &cobra.Command{
	Use:   "read <doc-id>",
	Short: "Read a ClickUp document by ID",
	Long:  `Retrieve and display a ClickUp document's content, metadata, and creator info.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		docID := args[0]

		var doc api.Document
		if err := client.Get(fmt.Sprintf("/doc/%s", docID), nil, &doc); err != nil {
			return fmt.Errorf("reading document: %w", err)
		}

		fmt.Printf("%s\n", doc.Name)
		fmt.Printf("========================================\n\n")
		fmt.Printf("ID:      %s\n", doc.ID)
		fmt.Printf("Created: %s\n", api.FormatTimestamp(doc.DateCreated))
		fmt.Printf("Creator: %s\n", doc.Creator.Username)
		fmt.Println()
		if doc.Content != "" {
			fmt.Println(doc.Content)
		} else {
			fmt.Println("(No content)")
		}

		return nil
	},
}

func init() {
	docCmd.AddCommand(docReadCmd)
}
