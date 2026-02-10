package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var taskSubtaskCmd = &cobra.Command{
	Use:   "subtask <parent-task-id> <name>",
	Short: "Create a subtask under an existing task",
	Long:  `Create a subtask under a parent task. The parent's list is used unless --list is specified.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		parentID := args[0]
		name := args[1]
		custom, _ := cmd.Flags().GetBool("custom")
		description, _ := cmd.Flags().GetString("description")
		listID, _ := cmd.Flags().GetString("list")

		// Fetch parent task to get internal ID and list ID
		parentParams := map[string]string{}
		if custom {
			parentParams["custom_task_ids"] = "true"
			parentParams["team_id"] = client.TeamID()
		}

		var parent api.Task
		if err := client.Get(fmt.Sprintf("/task/%s", parentID), parentParams, &parent); err != nil {
			return fmt.Errorf("fetching parent task: %w", err)
		}

		targetListID := listID
		if targetListID == "" {
			targetListID = parent.List.ID
		}
		if targetListID == "" {
			return fmt.Errorf("could not determine list ID for subtask creation")
		}

		subtaskData := map[string]string{
			"name":   name,
			"parent": parent.ID, // must use internal ID
		}
		if description != "" {
			subtaskData["markdown_description"] = description
		}

		body, err := json.Marshal(subtaskData)
		if err != nil {
			return fmt.Errorf("encoding request: %w", err)
		}

		createParams := map[string]string{}
		if custom {
			createParams["custom_task_ids"] = "true"
			createParams["team_id"] = client.TeamID()
		}

		var subtask api.Task
		if err := client.Post(fmt.Sprintf("/list/%s/task", targetListID), bytes.NewReader(body), createParams, &subtask); err != nil {
			return fmt.Errorf("creating subtask: %w", err)
		}

		id := subtask.ID
		if subtask.CustomID != nil && *subtask.CustomID != "" {
			id = *subtask.CustomID
		}

		fmt.Printf("Subtask created: %s %s\n", id, subtask.Name)
		fmt.Printf("Parent: %s (%s)\n", parent.Name, parentID)
		fmt.Printf("List: %s\n", parent.List.Name)
		fmt.Printf("Status: %s\n", subtask.Status.Status)
		fmt.Printf("URL: %s\n", subtask.URL)

		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskSubtaskCmd)
	taskSubtaskCmd.Flags().BoolP("custom", "c", false, "Treat the parent task ID as a custom task ID")
	taskSubtaskCmd.Flags().StringP("description", "d", "", "Subtask description (Markdown)")
	taskSubtaskCmd.Flags().StringP("list", "l", "", "Override the list ID (defaults to parent's list)")
}
