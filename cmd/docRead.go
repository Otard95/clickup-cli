package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var docReadCmd = &cobra.Command{
	Use:   "read <doc-id> [page-id]",
	Short: "Read a ClickUp document or a specific page",
	Long: `Retrieve and display a ClickUp document's content.

Without a page-id, displays all pages. With a page-id, displays only
that page. Use "doc pages" to list available page IDs.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		docID := args[0]
		docEndpoint := fmt.Sprintf("/workspaces/%s/docs/%s", client.TeamID(), docID)

		var doc api.Document
		if err := client.GetV3(docEndpoint, nil, &doc); err != nil {
			return fmt.Errorf("reading document: %w", err)
		}

		fmt.Printf("%s\n", doc.Name)
		fmt.Printf("========================================\n\n")
		fmt.Printf("ID:      %s\n", doc.ID)
		fmt.Printf("Created: %s\n", api.FormatTimestampInt(doc.DateCreated))
		fmt.Printf("Creator: %d\n", doc.Creator)
		fmt.Println()

		if len(args) == 2 {
			pageID := args[1]
			var page api.DocPage
			if err := client.GetV3(fmt.Sprintf("%s/pages/%s", docEndpoint, pageID), nil, &page); err != nil {
				return fmt.Errorf("reading page: %w", err)
			}
			printPages([]api.DocPage{page}, 0)
		} else {
			var pages []api.DocPage
			if err := client.GetV3(docEndpoint+"/pages", nil, &pages); err != nil {
				return fmt.Errorf("reading document pages: %w", err)
			}
			if len(pages) == 0 {
				fmt.Println("(No pages)")
				return nil
			}
			printPages(pages, 0)
		}

		return nil
	},
}

func printPages(pages []api.DocPage, depth int) {
	for _, p := range pages {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}

		fmt.Printf("%s## %s\n", indent, p.Name)
		if p.SubTitle != "" {
			fmt.Printf("%s*%s*\n", indent, p.SubTitle)
		}
		if p.Content != "" {
			fmt.Printf("%s%s\n", indent, p.Content)
		}
		fmt.Println()

		if len(p.Pages) > 0 {
			printPages(p.Pages, depth+1)
		}
	}
}

func init() {
	docCmd.AddCommand(docReadCmd)
}
