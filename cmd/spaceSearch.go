package cmd

import (
	"fmt"
	"strings"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var spaceSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for spaces in the workspace",
	Long:  `List all spaces, optionally filtered by name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}

		var resp api.SpacesResponse
		if err := client.Get(fmt.Sprintf("/team/%s/space", client.TeamID()), nil, &resp); err != nil {
			return fmt.Errorf("searching spaces: %w", err)
		}

		if len(resp.Spaces) == 0 {
			fmt.Println("No spaces found.")
			return nil
		}

		queryLower := strings.ToLower(query)
		count := 0
		for _, s := range resp.Spaces {
			if query != "" && !strings.Contains(strings.ToLower(s.Name), queryLower) {
				continue
			}
			count++
			private := ""
			if s.Private {
				private = " [private]"
			}
			fmt.Printf("%s  %s%s  (%d statuses)\n", s.ID, s.Name, private, len(s.Statuses))
		}

		if count == 0 {
			fmt.Println("No spaces matching your query.")
		}

		return nil
	},
}

func init() {
	spaceCmd.AddCommand(spaceSearchCmd)
}
