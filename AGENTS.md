# clickup-cli

## Overview

Go CLI for ClickUp project management. Built with [Cobra](https://github.com/spf13/cobra), it provides commands for tasks, lists, spaces, comments, time tracking, and documents.

### ClickUp Hierarchy

```
Workspace (Team)
└── Spaces
    ├── Folders
    │   └── Lists
    │       └── Tasks
    │           └── Subtasks
    └── Lists (folderless)
        └── Tasks
```

- **Spaces** are top-level organizational units within a workspace
- **Folders** organize lists within a space (optional)
- **Lists** contain tasks (e.g., sprint lists, backlogs)
- **Tasks** can have subtasks and a parent relationship

## Running

The CLI reads `CLICKUP_API_TOKEN` and `CLICKUP_TEAM_ID` from the environment. Use the `cu` pass-env alias to inject secrets:

```bash
pass-env cu clickup-cli <command> [flags]
```

## Project Structure

```
main.go                    # Entrypoint
cmd/                       # Cobra commands (one file per command)
  root.go                  # Root command, PersistentPreRunE loads config + creates API client
  <group>.go               # Parent commands (task, space, list, comment, time, doc)
  <group><Action>.go       # Leaf commands (taskGet, spaceStructure, etc.)
internal/
  api/
    client.go              # HTTP client (Get/Put/Post), auth header injection
    types.go               # ClickUp API response structs
    format.go              # Output formatting helpers (FormatTaskDetail, FormatTaskSummary, etc.)
  config/
    config.go              # Loads CLICKUP_API_TOKEN and CLICKUP_TEAM_ID from env
```

## Adding a New Command

1. Scaffold: `cobra-cli add <cmdName> -p <parentCmd>Cmd --viper=false`
2. Edit the generated file in `cmd/`: set `Use`, `Short`, `Long`, `Args`, change `Run` to `RunE`
3. Use the `client` variable (initialized in `root.go` PersistentPreRunE) for API calls
4. Add types to `internal/api/types.go` if the endpoint returns new shapes
5. Add formatting helpers to `internal/api/format.go` for output
6. Use `RunE` (not `Run`) and return errors -- cobra handles display

## Command Pattern

```go
var exampleCmd = &cobra.Command{
    Use:   "example <required-arg>",
    Short: "One-line description",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        flag, _ := cmd.Flags().GetBool("flag-name")
        params := map[string]string{}
        // ... build params, call client.Get/Put/Post, format output
        var resp api.SomeResponse
        if err := client.Get(endpoint, params, &resp); err != nil {
            return fmt.Errorf("doing thing: %w", err)
        }
        fmt.Print(api.FormatSomething(resp))
        return nil
    },
}

func init() {
    parentCmd.AddCommand(exampleCmd)
    exampleCmd.Flags().BoolP("flag-name", "f", false, "Description")
}
```

## ClickUp API Conventions

- API docs: https://developer.clickup.com/reference
- Base URL v2: `https://api.clickup.com/api/v2`
- Base URL v3: `https://api.clickup.com/api/v3` (used for Docs endpoints)
- Auth: raw token in `Authorization` header (no `Bearer` prefix)
- Custom task IDs (e.g. `MA-123`): require `custom_task_ids=true` and `team_id` query params
- Array query params: use `key[]=value` format (see `api.SetQueryArray`)
- Some fields are inconsistently typed across endpoints:
  - `task_count`: string in folder responses, number elsewhere -- use `FlexInt`
  - `status` on lists: nullable -- use `*Status`
  - `dependencies[].type`: integer, not string

## Task IDs

ClickUp has two ID formats:
- **Internal IDs**: short alphanumeric strings like `869b827uy`
- **Custom IDs**: prefixed with a project code like `MA-25578`, `GQL-456`

Custom IDs are what users typically reference. When using a custom ID with the API, always include `custom_task_ids=true` and `team_id` query params. In the CLI, this is handled by the `-c` / `--custom` flag.

## Testing with curl

When debugging or developing new features, you can test the ClickUp API directly using curl:

### Authentication

The API token is stored in `pass` and must be passed in the `Authorization` header:

```bash
pass show keys/clickup | xargs -I {} curl -s "https://api.clickup.com/api/v2/ENDPOINT" \
  -H "Authorization: {}" \
  -H "Content-Type: application/json"
```

**Important**: Do NOT use `$(pass show keys/clickup)` directly in the Authorization header -- it doesn't work in the Bash tool context. Always use `xargs -I {}` to pass the token.

### Common Test Commands

```bash
# Get tasks from a list
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/list/LIST_ID/task" \
  -H "Authorization: {}" | jq .

# Get tasks filtered by assignee (requires numeric user ID)
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/list/LIST_ID/task?assignees[]=USER_ID" \
  -H "Authorization: {}" | jq .

# Get a specific task
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/task/TASK_ID" \
  -H "Authorization: {}" | jq .

# Get folders in a space
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/space/SPACE_ID/folder" \
  -H "Authorization: {}" | jq .

# Get folder details (including lists)
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/folder/FOLDER_ID" \
  -H "Authorization: {}" | jq .
```

### Finding User IDs

User IDs are required for assignee filtering. To find a user ID:

```bash
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/list/LIST_ID/task" \
  -H "Authorization: {}" | \
  jq '.tasks[].assignees[] | {id, username}'
```

### Important Notes

- **Assignee filtering** requires numeric user IDs, NOT usernames or display names
- The API does NOT support `assignees[]=me` -- you must use actual user IDs
- Use `jq` to parse and filter JSON responses
- The `/user` endpoint may not work with all token types
- Query parameters with arrays use the `param[]=value` format
