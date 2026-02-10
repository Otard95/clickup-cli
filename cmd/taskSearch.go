package cmd

import (
	"fmt"
	"strings"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var taskSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for tasks across workspaces",
	Long: `Search for tasks with optional filters. The query argument performs
client-side text filtering on task names and descriptions.

Note: The assignee flag requires a numeric user ID, not a username.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}

		listID, _ := cmd.Flags().GetString("list")
		spaceID, _ := cmd.Flags().GetString("space")
		assignee, _ := cmd.Flags().GetString("assignee")
		status, _ := cmd.Flags().GetString("status")

		params := map[string]string{}
		if listID != "" {
			params["list_ids[]"] = listID
		}
		if spaceID != "" {
			params["space_ids[]"] = spaceID
		}
		if assignee != "" {
			params["assignees[]"] = assignee
		}
		if status != "" {
			params["statuses[]"] = status
		}

		var resp api.TasksResponse
		if err := client.Get(fmt.Sprintf("/team/%s/task", client.TeamID()), params, &resp); err != nil {
			return fmt.Errorf("searching tasks: %w", err)
		}

		tasks := resp.Tasks

		// Client-side text filter
		if query != "" {
			queryLower := strings.ToLower(query)
			var filtered []api.Task
			for _, t := range tasks {
				if strings.Contains(strings.ToLower(t.Name), queryLower) ||
					strings.Contains(strings.ToLower(t.Description), queryLower) {
					filtered = append(filtered, t)
				}
			}
			tasks = filtered
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found matching your criteria.")
			return nil
		}

		limit := 10
		if len(tasks) < limit {
			limit = len(tasks)
		}

		fmt.Printf("Found %d task(s):\n\n", len(tasks))
		for _, t := range tasks[:limit] {
			fmt.Println(api.FormatTaskSummary(t))
		}
		if len(tasks) > 10 {
			fmt.Printf("... and %d more tasks\n", len(tasks)-10)
		}

		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskSearchCmd)
	taskSearchCmd.Flags().StringP("list", "l", "", "Filter by list ID")
	taskSearchCmd.Flags().StringP("space", "S", "", "Filter by space ID")
	taskSearchCmd.Flags().StringP("assignee", "a", "", "Filter by assignee user ID (numeric)")
	taskSearchCmd.Flags().StringP("status", "s", "", "Filter by status")
}
