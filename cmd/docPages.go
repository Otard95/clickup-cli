package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var docPagesCmd = &cobra.Command{
	Use:   "pages <doc-id>",
	Short: "List pages in a ClickUp document",
	Long:  `Display the page tree of a ClickUp document showing page names and IDs.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		docID := args[0]
		endpoint := fmt.Sprintf("/workspaces/%s/docs/%s/page_listing", client.TeamID(), docID)

		var pages []api.DocPage
		if err := client.GetV3(endpoint, nil, &pages); err != nil {
			return fmt.Errorf("listing pages: %w", err)
		}

		if len(pages) == 0 {
			fmt.Println("No pages found.")
			return nil
		}

		printPageTree(pages, 0)
		return nil
	},
}

func printPageTree(pages []api.DocPage, depth int) {
	for _, p := range pages {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("%s%s  %s\n", indent, p.ID, p.Name)
		if len(p.Pages) > 0 {
			printPageTree(p.Pages, depth+1)
		}
	}
}

func init() {
	docCmd.AddCommand(docPagesCmd)
}
