package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var listInfoCmd = &cobra.Command{
	Use:   "info <list-id>",
	Short: "Get detailed information about a list",
	Long:  `Show list metadata including statuses, feature flags, and organization.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listID := args[0]

		var list api.ListInfo
		if err := client.Get(fmt.Sprintf("/list/%s", listID), nil, &list); err != nil {
			return fmt.Errorf("getting list info: %w", err)
		}

		fmt.Printf("%s\n", list.Name)
		fmt.Printf("========================================\n\n")
		statusStr := "None"
		if list.Status != nil {
			statusStr = list.Status.Status
		}
		fmt.Printf("ID:               %s\n", list.ID)
		fmt.Printf("Status:           %s\n", statusStr)
		fmt.Printf("Task Count:       %d\n", int(list.TaskCount))
		fmt.Printf("Permission Level: %s\n", list.PermissionLevel)
		fmt.Println()
		fmt.Printf("Space:  %s\n", list.Space.Name)
		fmt.Printf("Folder: %s\n", api.Or(list.Folder.Name, "No folder"))
		fmt.Println()
		fmt.Printf("Due Dates:          %t\n", list.DueDateTime)
		fmt.Printf("Multiple Assignees: %t\n", list.MultipleAssignees)
		fmt.Printf("Time Tracking:      %t\n", list.TimeTracking)
		fmt.Println()

		if len(list.Statuses) > 0 {
			fmt.Println("Statuses:")
			for _, s := range list.Statuses {
				fmt.Printf("  - %s (%s)\n", s.Status, s.Type)
			}
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(listInfoCmd)
}
