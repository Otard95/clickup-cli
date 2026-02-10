---
name: clickup-cli
description: Go CLI for ClickUp project management. Use when building, running, testing, or extending the clickup CLI. Covers project structure, ClickUp API conventions, authentication via pass-env, and adding new commands.
---

# clickup-cli

## Running

The CLI is installed on the system. Use the `cu` pass-env alias to inject secrets:

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
5. Use `RunE` (not `Run`) and return errors — cobra handles display

## ClickUp API Conventions

- Base URL: `https://api.clickup.com/api/v2`
- Auth: raw token in `Authorization` header (no `Bearer` prefix)
- Custom task IDs (e.g. `MA-123`): require `custom_task_ids=true` and `team_id` query params
- Array query params: use `key[]=value` format (see `api.SetQueryArray`)
- Some fields are inconsistently typed across endpoints:
  - `task_count`: string in folder responses, number elsewhere — use `FlexInt`
  - `status` on lists: nullable — use `*Status`
  - `dependencies[].type`: integer, not string

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
