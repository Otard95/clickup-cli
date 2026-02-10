package cmd

import (
	"fmt"
	"strings"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var docSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for documents in the workspace",
	Long:  `Find documents across the ClickUp workspace, optionally filtered by query.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}

		params := map[string]string{}
		if query != "" {
			params["query"] = query
		}

		var resp api.DocsResponse
		if err := client.Get(fmt.Sprintf("/team/%s/docs", client.TeamID()), params, &resp); err != nil {
			return fmt.Errorf("searching documents: %w", err)
		}

		if len(resp.Docs) == 0 {
			fmt.Println("No documents found.")
			return nil
		}

		fmt.Printf("Found %d document(s):\n\n", len(resp.Docs))
		for _, d := range resp.Docs {
			fmt.Printf("%s  %s  (created: %s, by: %s)\n",
				d.ID, d.Name, api.FormatTimestamp(d.DateCreated), d.Creator.Username)
		}

		return nil
	},
}

func init() {
	docCmd.AddCommand(docSearchCmd)
}
