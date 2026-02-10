package cmd

import (
	"fmt"
	"strconv"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var timeGetCmd = &cobra.Command{
	Use:   "get [task-id]",
	Short: "Get time tracking entries",
	Long: `Get time tracking entries for a specific task or the whole team.

If a task ID is given, shows time entries for that task.
Otherwise, shows entries for the configured team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var resp api.TimeEntriesResponse
		var context string

		if len(args) > 0 {
			taskID := args[0]
			if err := client.Get(fmt.Sprintf("/task/%s/time", taskID), nil, &resp); err != nil {
				return fmt.Errorf("getting time entries: %w", err)
			}
			context = fmt.Sprintf("task %s", taskID)
		} else {
			teamID, _ := cmd.Flags().GetString("team")
			if teamID == "" {
				teamID = client.TeamID()
			}
			if err := client.Get(fmt.Sprintf("/team/%s/time_entries", teamID), nil, &resp); err != nil {
				return fmt.Errorf("getting time entries: %w", err)
			}
			context = fmt.Sprintf("team %s", teamID)
		}

		if len(resp.Data) == 0 {
			fmt.Printf("No time entries found for %s.\n", context)
			return nil
		}

		fmt.Printf("Found %d time entry/entries for %s:\n\n", len(resp.Data), context)
		for _, e := range resp.Data {
			durationMs, _ := strconv.ParseInt(e.Duration, 10, 64)
			fmt.Printf("  %s  %s  %s\n", api.FormatTimestamp(e.Start), api.FormatDurationMs(durationMs), e.User.Username)
			if e.Description != "" {
				fmt.Printf("    %s\n", e.Description)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	timeCmd.AddCommand(timeGetCmd)
	timeGetCmd.Flags().StringP("team", "t", "", "Override team ID (defaults to CLICKUP_TEAM_ID)")
}
