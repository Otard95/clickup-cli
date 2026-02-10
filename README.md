# clickup-cli

A CLI for interacting with ClickUp — tasks, lists, spaces, comments, time tracking, and documents.

## Install

```bash
nix profile install github:otard95/clickup-cli
```

Or from source:

```bash
go install github.com/otard95/clickup-cli@latest
```

## Auth

Requires two environment variables:

- `CLICKUP_API_TOKEN` — personal API token from ClickUp Settings > Apps
- `CLICKUP_TEAM_ID` — the numeric ID from your workspace URL (`app.clickup.com/{team_id}/...`)

With [pass-env](https://github.com/otard95/pass-env):

```bash
pass-env cu clickup-cli <command>
```

## Commands

```
clickup-cli task search [query]       Search tasks (--list, --space, --assignee, --status)
clickup-cli task get <id>             Task details (-c custom ID, -s include subtasks)
clickup-cli task update <id>          Update task (--title, --description, --status)
clickup-cli task subtask <parent> <n> Create subtask
clickup-cli task rels <id>            Show dependencies and linked tasks

clickup-cli space search [query]      List/search spaces
clickup-cli space structure <id>      Full folder/list tree

clickup-cli list tasks <id>           Tasks in a list (--assignees, --archived)
clickup-cli list info <id>            List metadata and statuses

clickup-cli comment get <task-id>     Task comments
clickup-cli time get [task-id]        Time entries (task or team)
clickup-cli doc read <id>             Read a document
clickup-cli doc search [query]        Search documents
```

## License

MIT
