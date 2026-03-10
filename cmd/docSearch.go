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
	Long: `Find documents across the ClickUp workspace, optionally filtered by name.

The ClickUp API does not support server-side search, so this command
paginates through all workspace docs and filters by name client-side.

Use --cursor to continue from a previous search result.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		maxResults, _ := cmd.Flags().GetInt("max")
		cursor, _ := cmd.Flags().GetString("cursor")

		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}
		lowerQuery := strings.ToLower(query)

		endpoint := fmt.Sprintf("/workspaces/%s/docs", client.TeamID())
		var matches []api.Document
		nextCursor := ""

		for {
			params := map[string]string{}
			if cursor != "" {
				// The v3 API only respects "next_cursor" as the query param
				// despite documenting "cursor" as the non-deprecated option.
				params["next_cursor"] = cursor
			}

			var resp api.DocsResponse
			if err := client.GetV3(endpoint, params, &resp); err != nil {
				return fmt.Errorf("searching documents: %w", err)
			}

			for _, d := range resp.Docs {
				if query == "" || strings.Contains(strings.ToLower(d.Name), lowerQuery) {
					matches = append(matches, d)
				}
			}

			nextCursor = resp.NextCursor
			if len(matches) >= maxResults || nextCursor == "" {
				break
			}
			cursor = nextCursor
		}

		if len(matches) == 0 {
			fmt.Println("No documents found.")
			return nil
		}

		fmt.Printf("Found %d document(s):\n\n", len(matches))
		for _, d := range matches {
			fmt.Printf("%s  %s  (created: %s)\n",
				d.ID, d.Name, api.FormatTimestampInt(d.DateCreated))
		}

		if nextCursor != "" {
			fmt.Printf("\nMore results available. Continue with:\n  clickup-cli doc search --cursor %s", nextCursor)
			if query != "" {
				fmt.Printf(" %q", query)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	docCmd.AddCommand(docSearchCmd)
	docSearchCmd.Flags().IntP("max", "n", 25, "Minimum number of results before stopping")
	docSearchCmd.Flags().String("cursor", "", "Cursor from a previous search to continue paging")
}
