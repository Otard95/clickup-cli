package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var commentGetCmd = &cobra.Command{
	Use:   "get <task-id>",
	Short: "Get all comments for a task",
	Long:  `Retrieve all comments for a specific task showing author, date, and content.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		var resp api.CommentsResponse
		if err := client.Get(fmt.Sprintf("/task/%s/comment", taskID), nil, &resp); err != nil {
			return fmt.Errorf("getting comments: %w", err)
		}

		if len(resp.Comments) == 0 {
			fmt.Println("No comments found for this task.")
			return nil
		}

		fmt.Printf("Found %d comment(s):\n\n", len(resp.Comments))
		for _, c := range resp.Comments {
			fmt.Printf("--- %s  (%s) ---\n", c.User.Username, api.FormatTimestamp(c.Date))
			fmt.Printf("%s\n", c.CommentText)
			fmt.Printf("ID: %s\n\n", c.ID)
		}

		return nil
	},
}

func init() {
	commentCmd.AddCommand(commentGetCmd)
}
