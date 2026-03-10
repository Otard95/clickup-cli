---
name: clickup-cli
description: How to use the clickup-cli tool for ClickUp project management. Use when the user wants to interact with ClickUp tasks, lists, spaces, comments, time tracking, or documents. Covers all commands, flags, ID formats, and common workflows.
---

# Using clickup-cli

## Authentication

The CLI requires two environment variables to be set:
- `CLICKUP_API_TOKEN` -- a ClickUp personal API token
- `CLICKUP_TEAM_ID` -- the workspace/team ID

```bash
clickup-cli <command> [flags]
```

## Task IDs

ClickUp has two ID formats:
- **Internal IDs**: short alphanumeric strings like `869b827uy` -- use directly
- **Custom IDs**: prefixed like `MA-25578`, `GQL-456` -- requires the `-c` flag

```bash
# Custom ID (most common)
clickup-cli task get MA-25578 -c

# Internal ID
clickup-cli task get 869b827uy
```

## Commands

### task get -- View task details

```bash
clickup-cli task get <task-id> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--custom` | `-c` | Treat task ID as a custom ID |
| `--subtasks` | `-s` | Include subtasks in output |

```bash
clickup-cli task get MA-25578 -c
clickup-cli task get MA-25578 -c -s   # with subtasks
```

### task search -- Find tasks

```bash
clickup-cli task search [query] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--list` | `-l` | Filter by list ID |
| `--space` | `-S` | Filter by space ID |
| `--assignee` | `-a` | Filter by assignee user ID (numeric) |
| `--status` | `-s` | Filter by status |

The `query` argument performs client-side text filtering on task names and descriptions.

```bash
clickup-cli task search "auth bug"
clickup-cli task search -S 12345 -s "in progress"
clickup-cli task search -l 67890 -a 12345678
```

### task update -- Modify a task

```bash
clickup-cli task update <task-id> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--custom` | `-c` | Treat task ID as a custom ID |
| `--title` | `-t` | New task title |
| `--description` | `-d` | New description (Markdown) |
| `--status` | `-s` | New status |

```bash
clickup-cli task update MA-25578 -c -s "in progress"
clickup-cli task update MA-25578 -c -t "Updated title" -d "New description"
```

### task subtask -- Create a subtask

```bash
clickup-cli task subtask <parent-task-id> <name> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--custom` | `-c` | Treat parent task ID as a custom ID |
| `--description` | `-d` | Subtask description (Markdown) |
| `--list` | `-l` | Override list ID (defaults to parent's list) |

```bash
clickup-cli task subtask MA-25578 "Write API tests" -c
clickup-cli task subtask MA-25578 "Update docs" -c -d "Add endpoint docs"
```

### task rels -- View relationships

```bash
clickup-cli task rels <task-id> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--custom` | `-c` | Treat task ID as a custom ID |

```bash
clickup-cli task rels MA-25578 -c
```

### space search -- Find spaces

```bash
clickup-cli space search [query]
```

The optional `query` filters spaces by name.

```bash
clickup-cli space search
clickup-cli space search "Engineering"
```

### space structure -- View space hierarchy

```bash
clickup-cli space structure <space-id>
```

Shows all folders, lists, and task counts in a tree view.

```bash
clickup-cli space structure 12345
```

### list info -- View list details

```bash
clickup-cli list info <list-id>
```

Shows list metadata including available statuses, feature flags, and organization.

### list tasks -- Get tasks in a list

```bash
clickup-cli list tasks <list-id> [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--archived` | `-a` | Include archived tasks |
| `--assignees` | `-A` | Filter by user IDs (comma-separated) |

```bash
clickup-cli list tasks 67890
clickup-cli list tasks 67890 -A 111,222
```

### comment get -- View task comments

```bash
clickup-cli comment get <task-id>
```

```bash
clickup-cli comment get 869b827uy
```

### time get -- View time entries

```bash
clickup-cli time get [task-id] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--team` | `-t` | Override team ID |

Without a task ID, shows entries for the whole team.

```bash
clickup-cli time get 869b827uy      # entries for a task
clickup-cli time get                  # all team entries
```

### doc pages -- List pages in a document

```bash
clickup-cli doc pages <doc-id>
```

Displays the page tree with IDs and names. Use this to find page IDs for `doc read`.

```bash
clickup-cli doc pages 2g6xm-16448
```

### doc read -- Read a document

```bash
clickup-cli doc read <doc-id> [page-id]
```

Without a page-id, displays all pages. With a page-id, displays only that page.

```bash
clickup-cli doc read 2g6xm-16448              # all pages
clickup-cli doc read 2g6xm-16448 2g6xm-2908   # single page
```

### doc search -- Find documents

```bash
clickup-cli doc search [query] [flags]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--max` | `-n` | Minimum number of results before stopping (default 25) |
| `--cursor` | | Cursor from a previous search to continue paging |

```bash
clickup-cli doc search "onboarding"
clickup-cli doc search -n 5
clickup-cli doc search --cursor <cursor> "onboarding"
```

## Common Workflows

### Explore the workspace

```bash
# List all spaces
clickup-cli space search

# View a space's folder/list structure
clickup-cli space structure <space-id>

# Check available statuses for a list
clickup-cli list info <list-id>

# Get tasks from a list
clickup-cli list tasks <list-id>
```

### Work with a task

```bash
# Get full details
clickup-cli task get MA-25578 -c

# See subtasks too
clickup-cli task get MA-25578 -c -s

# Check dependencies and linked tasks
clickup-cli task rels MA-25578 -c

# Read comments
clickup-cli comment get 869b827uy

# Check time logged
clickup-cli time get 869b827uy

# Update status
clickup-cli task update MA-25578 -c -s "complete"

# Add a subtask
clickup-cli task subtask MA-25578 "Write tests" -c
```

### Work with documents

```bash
# Find a document
clickup-cli doc search "Developer Foundations"

# List its pages
clickup-cli doc pages 2g6xm-16448

# Read a specific page
clickup-cli doc read 2g6xm-16448 2g6xm-2908
```

## Important Notes

- **Assignee filtering** requires numeric user IDs, not usernames
- The API does not support `assignees[]=me`
- Custom IDs always require the `-c` flag
- Description fields support Markdown
