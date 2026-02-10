package cmd

import (
	"fmt"
	"strings"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var listTasksCmd = &cobra.Command{
	Use:   "tasks <list-id>",
	Short: "Get all tasks in a list",
	Long: `Get all tasks within a specific list. Supports filtering by assignee
using numeric user IDs (comma-separated for multiple).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listID := args[0]
		archived, _ := cmd.Flags().GetBool("archived")
		assignees, _ := cmd.Flags().GetString("assignees")

		params := map[string]string{}
		if archived {
			params["archived"] = "true"
		}

		// Build the endpoint; assignees need array params
		endpoint := fmt.Sprintf("/list/%s/task", listID)
		if assignees != "" {
			ids := strings.Split(assignees, ",")
			for i := range ids {
				ids[i] = strings.TrimSpace(ids[i])
			}
			endpoint = api.SetQueryArray(endpoint, "assignees[]", ids)
		}

		var resp api.TasksResponse
		if err := client.Get(endpoint, params, &resp); err != nil {
			return fmt.Errorf("getting tasks: %w", err)
		}

		if len(resp.Tasks) == 0 {
			fmt.Println("No tasks found in this list.")
			return nil
		}

		fmt.Printf("Found %d task(s) in list:\n\n", len(resp.Tasks))
		for _, t := range resp.Tasks {
			fmt.Println(api.FormatTaskSummary(t))
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listTasksCmd)
	listTasksCmd.Flags().BoolP("archived", "a", false, "Include archived tasks")
	listTasksCmd.Flags().StringP("assignees", "A", "", "Filter by assignee user IDs (comma-separated)")
}
