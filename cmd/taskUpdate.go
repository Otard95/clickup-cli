package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var taskUpdateCmd = &cobra.Command{
	Use:   "update <task-id>",
	Short: "Update a task's title, description, or status",
	Long:  `Update one or more fields on a task. Description supports Markdown formatting.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		custom, _ := cmd.Flags().GetBool("custom")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		if title == "" && description == "" && status == "" {
			return fmt.Errorf("at least one of --title, --description, or --status must be provided")
		}

		params := map[string]string{}
		if custom {
			params["custom_task_ids"] = "true"
			params["team_id"] = client.TeamID()
		}

		data := map[string]string{}
		var updated []string
		if title != "" {
			data["name"] = title
			updated = append(updated, "title")
		}
		if description != "" {
			data["markdown_description"] = description
			updated = append(updated, "description")
		}
		if status != "" {
			data["status"] = status
			updated = append(updated, "status")
		}

		body, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("encoding request: %w", err)
		}

		var task api.Task
		if err := client.Put(fmt.Sprintf("/task/%s", taskID), bytes.NewReader(body), params, &task); err != nil {
			return fmt.Errorf("updating task: %w", err)
		}

		id := task.ID
		if task.CustomID != nil && *task.CustomID != "" {
			id = *task.CustomID
		}

		fmt.Printf("Task updated: %s %s\n", id, task.Name)
		fmt.Printf("Updated fields: %s\n", fmt.Sprintf("%v", updated))
		fmt.Printf("Status: %s\n", task.Status.Status)
		fmt.Printf("URL: %s\n", task.URL)

		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskUpdateCmd)
	taskUpdateCmd.Flags().BoolP("custom", "c", false, "Treat the task ID as a custom task ID")
	taskUpdateCmd.Flags().StringP("title", "t", "", "New task title")
	taskUpdateCmd.Flags().StringP("description", "d", "", "New task description (Markdown)")
	taskUpdateCmd.Flags().StringP("status", "s", "", "New task status")
}
