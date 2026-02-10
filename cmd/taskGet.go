package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var taskGetCmd = &cobra.Command{
	Use:   "get <task-id>",
	Short: "Get detailed information about a task",
	Long: `Get detailed information about a specific task by its ID.

Supports both internal IDs (short alphanumeric) and custom IDs (e.g. MA-123)
when the --custom flag is set.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		custom, _ := cmd.Flags().GetBool("custom")
		subtasks, _ := cmd.Flags().GetBool("subtasks")

		params := map[string]string{}
		if custom {
			params["custom_task_ids"] = "true"
			params["team_id"] = client.TeamID()
		}
		if subtasks {
			params["include_subtasks"] = "true"
		}

		var task api.Task
		if err := client.Get(fmt.Sprintf("/task/%s", taskID), params, &task); err != nil {
			return fmt.Errorf("getting task: %w", err)
		}

		fmt.Print(api.FormatTaskDetail(task))
		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskGetCmd)
	taskGetCmd.Flags().BoolP("custom", "c", false, "Treat the task ID as a custom task ID")
	taskGetCmd.Flags().BoolP("subtasks", "s", false, "Include subtasks in the output")
}
