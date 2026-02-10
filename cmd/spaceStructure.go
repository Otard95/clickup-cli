package cmd

import (
	"fmt"

	"github.com/otard95/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var spaceStructureCmd = &cobra.Command{
	Use:   "structure <space-id>",
	Short: "Show full folder/list hierarchy of a space",
	Long:  `Display the complete organizational structure of a space with folders, lists, and task counts.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		spaceID := args[0]

		var foldersResp api.FoldersResponse
		var listsResp api.ListsResponse

		// Fetch folders and folderless lists
		if err := client.Get(fmt.Sprintf("/space/%s/folder", spaceID), nil, &foldersResp); err != nil {
			return fmt.Errorf("getting folders: %w", err)
		}
		if err := client.Get(fmt.Sprintf("/space/%s/list", spaceID), nil, &listsResp); err != nil {
			return fmt.Errorf("getting folderless lists: %w", err)
		}

		if len(foldersResp.Folders) == 0 && len(listsResp.Lists) == 0 {
			fmt.Printf("No folders or lists found in space %s\n", spaceID)
			return nil
		}

		fmt.Printf("Space Structure (ID: %s)\n\n", spaceID)

		// Folders
		if len(foldersResp.Folders) > 0 {
			for fi, folder := range foldersResp.Folders {
				isLastFolder := fi == len(foldersResp.Folders)-1
				prefix := "├──"
				if isLastFolder && len(listsResp.Lists) == 0 {
					prefix = "└──"
				}
				hidden := ""
				if folder.Hidden {
					hidden = " [hidden]"
				}
				fmt.Printf("%s %s (ID: %s)%s\n", prefix, folder.Name, folder.ID, hidden)

				treePrefix := "│"
				if isLastFolder && len(listsResp.Lists) == 0 {
					treePrefix = " "
				}

				if len(folder.Lists) > 0 {
					for li, list := range folder.Lists {
						listPrefix := "├──"
						if li == len(folder.Lists)-1 {
							listPrefix = "└──"
						}
						fmt.Printf("%s   %s %s (ID: %s) - %d tasks\n",
							treePrefix, listPrefix, list.Name, list.ID, int(list.TaskCount))
					}
				} else {
					fmt.Printf("%s   (no lists)\n", treePrefix)
				}

				if !isLastFolder {
					fmt.Printf("│\n")
				}
			}
		}

		// Folderless lists
		if len(listsResp.Lists) > 0 {
			if len(foldersResp.Folders) > 0 {
				fmt.Println()
			}
			fmt.Printf("Folderless Lists (%d):\n", len(listsResp.Lists))
			for li, list := range listsResp.Lists {
				prefix := "├──"
				if li == len(listsResp.Lists)-1 {
					prefix = "└──"
				}
				fmt.Printf("%s %s (ID: %s) - %d tasks\n",
					prefix, list.Name, list.ID, int(list.TaskCount))
			}
		}

		return nil
	},
}

func init() {
	spaceCmd.AddCommand(spaceStructureCmd)
}
