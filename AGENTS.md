# ClickUp MCP Server

## Overview

This MCP server provides comprehensive access to ClickUp's project management features through a secure, containerized interface. It focuses on the most commonly used ClickUp functionality: tasks, lists, spaces, comments, documents, and time tracking.

### ClickUp Hierarchy

Understanding the ClickUp structure is essential for navigating the workspace:

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

**Key points**:
- **Spaces** are top-level organizational units within a workspace
- **Folders** organize lists within a space (optional - lists can exist directly in spaces)
- **Lists** contain tasks (e.g., sprint lists, backlogs)
- **Tasks** can have subtasks and a parent relationship
- Use `get_space_structure` to view the complete hierarchy in one call, then `get_tasks_in_list` to work with specific lists

## Available Tools

### Task Management
- **search_tasks**: Search tasks with filters (query, list_id, space_id, assignee, status, custom_task_ids). Note: assignee requires user ID, not username
- **get_task_by_id**: Get complete task details including metadata, assignees, parent task ID, time tracking, and optionally subtasks with `include_subtasks="true"` (supports custom_task_ids)
- **get_tasks_in_list**: List all tasks in a specific list (supports archived tasks, assignee filtering by user ID, and custom_task_ids)
- **update_task**: Update a task's title, description (Markdown supported), and/or status (supports custom_task_ids)
- **create_subtask**: Create a subtask under an existing parent task (supports custom_task_ids)
- **get_task_relationships**: Get detailed information about task dependencies and linked tasks (supports custom_task_ids)

### Project Organization
- **search_spaces**: Find and list available spaces across teams
- **get_space_structure**: Get complete space hierarchy showing all folders and lists with IDs and task counts in a tree view
- **get_list_info**: Detailed information about a list including statuses and settings

### Collaboration
- **get_task_comments**: Retrieve all comments for a task with authors and timestamps
- **read_document**: Read ClickUp document content by ID
- **search_documents**: Find documents across workspaces

### Time Tracking
- **get_time_entries**: Get time entries for tasks or teams with duration and descriptions

## Authentication

The server uses personal API token authentication with the ClickUp API v2. Two environment variables are required:
- `CLICKUP_API_TOKEN`: Your ClickUp API token (stored as Docker secret)
- `CLICKUP_TEAM_ID`: Your ClickUp team/workspace ID (stored as Docker secret)

## Common Usage Patterns

### Finding Tasks
```
# Search by keyword
"Find all tasks related to 'user authentication'"

# Filter by assignee  
"Show me tasks assigned to john.doe"

# Filter by status
"List all tasks in 'In Progress' status"

# Combine filters
"Find high priority tasks in the Development space assigned to Sarah"
```

### Project Discovery
```
# Explore workspace structure
"What spaces are available in my ClickUp workspace?"

# View complete space hierarchy
"Show me the complete structure of space 10230170"
"Get the folder and list structure for the MA space"

# Get project details
"Tell me about the 'Q4 Campaign' list"
```

### Task Details
```
# Complete task information
"Get full details for task ID 12345"

# Using custom task IDs (requires team context)
"Get task details for custom task ID MA-23587 with custom task IDs enabled"

# Using internal task IDs (direct access)
"Get task details for internal task ID 12345"

# Get task with subtasks included
"Get task MA-25189 with subtasks included"

# Include collaboration info
"Show me all comments on the bug report task"

# Time tracking
"How much time was logged on task 67890?"

# Status management
"What are the available statuses for list 12345?" (use get_list_info)
"Update task CU-abc123 to 'In Progress' status" (use update_task with status parameter)

# Subtask management
"Get task MA-25189 with subtasks included" (use get_task_by_id with include_subtasks="true")
"Create a subtask under MA-25189 called 'Write API tests'"
"Create subtask 'Update documentation' for task 12345 with description 'Add new endpoint docs'"

# Task relationships
"Show me the relationships and dependencies for task MA-23048"
"Get all linked tasks for MA-25189"
```

## Data Format

### Task Information
Tasks include comprehensive metadata:
- Basic info (ID, name, status, priority, dates)
- People (assignees, watchers, creator)
- Organization (list, space, tags)
- Time tracking (estimated vs actual time)
- URLs for direct access

### Comments
Comments show full conversation threads with:
- Author information
- Timestamps
- Complete message content
- Comment IDs for reference

### Time Entries
Time tracking data includes:
- User who logged time
- Duration in minutes
- Start timestamps
- Descriptions of work performed

## Error Handling

The server provides clear error messages for:
- Missing API tokens
- Invalid task/list/space IDs
- Network connectivity issues
- API rate limiting
- Permission restrictions

## Performance Considerations

- API responses are limited to reasonable sizes (10 tasks max in search results)
- Requests include proper timeouts (30 seconds)
- Error handling prevents cascading failures
- Logging helps with debugging issues

## Integration Tips

1. **Start with search_spaces** to understand workspace structure
2. **Use get_list_info** to understand project organization
3. **Combine search_tasks with filters** for targeted results
4. **Always get full task details** with get_task_by_id for important tasks
5. **Check comments and time entries** for complete context

## Security Notes

- API tokens never appear in logs or responses
- All communication uses HTTPS
- Container runs as non-root user
- No persistent storage of sensitive data

## Limitations

- Some ClickUp features (goals, guests, roles) are not implemented
- Document API endpoints may vary based on ClickUp plan level
- Rate limiting applies (100 requests/minute typical)
- Requires appropriate ClickUp permissions for data access

## Testing with curl

When debugging or developing new features, you can test the ClickUp API directly using curl:

### Authentication

The API token is stored in `pass` and must be passed in the `Authorization` header:

```bash
# Basic pattern
pass show keys/clickup | xargs -I {} curl -s "https://api.clickup.com/api/v2/ENDPOINT" \
  -H "Authorization: {}" \
  -H "Content-Type: application/json"
```

**Important**: Do NOT use `$(pass show keys/clickup)` directly in the Authorization header - it doesn't work in the Bash tool context. Always use `xargs -I {}` to pass the token.

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
# Get tasks and extract assignee information
pass show keys/clickup | xargs -I {} curl -s \
  "https://api.clickup.com/api/v2/list/LIST_ID/task" \
  -H "Authorization: {}" | \
  jq '.tasks[].assignees[] | {id, username}'
```

### Important Notes

- **Assignee filtering** requires numeric user IDs, NOT usernames or display names
- The API does NOT support `assignees[]=me` - you must use actual user IDs
- Use `jq` to parse and filter JSON responses
- The `/user` endpoint may not work with all token types
- Query parameters with arrays use the `param[]=value` format
