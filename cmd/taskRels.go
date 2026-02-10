package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var taskRelsCmd = &cobra.Command{
	Use:   "rels <task-id>",
	Short: "Get task relationships (dependencies and linked tasks)",
	Long:  `Show dependencies and linked tasks for a given task.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		custom, _ := cmd.Flags().GetBool("custom")

		params := map[string]string{}
		if custom {
			params["custom_task_ids"] = "true"
			params["team_id"] = client.TeamID()
		}

		var task api.Task
		if err := client.Get(fmt.Sprintf("/task/%s", taskID), params, &task); err != nil {
			return fmt.Errorf("getting task: %w", err)
		}

		if len(task.Dependencies) == 0 && len(task.LinkedTasks) == 0 {
			fmt.Printf("No relationships found for task: %s\n", task.Name)
			return nil
		}

		fmt.Printf("Relationships for: %s (%s)\n\n", task.Name, taskID)

		if len(task.Dependencies) > 0 {
			fmt.Printf("Dependencies (%d):\n", len(task.Dependencies))
			for _, d := range task.Dependencies {
				fmt.Printf("  - Task %s (type: %d, created: %s)\n",
					d.DependsOn, d.Type, api.FormatTimestamp(d.DateCreated))
			}
			fmt.Println()
		}

		if len(task.LinkedTasks) > 0 {
			fmt.Printf("Linked Tasks (%d):\n", len(task.LinkedTasks))
			for _, l := range task.LinkedTasks {
				fmt.Printf("  - Task %s (created: %s, by user: %s)\n",
					l.LinkID, api.FormatTimestamp(l.DateCreated), l.UserID)
			}
			fmt.Println()
		}

		total := len(task.Dependencies) + len(task.LinkedTasks)
		fmt.Printf("Total: %d relationship(s)\n", total)

		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskRelsCmd)
	taskRelsCmd.Flags().BoolP("custom", "c", false, "Treat the task ID as a custom task ID")
}
