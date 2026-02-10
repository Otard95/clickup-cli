package api

import (
	"encoding/json"
	"strconv"
)

// FlexInt handles JSON numbers that may arrive as either a number or a string.
type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	// Try number first
	var n int
	if err := json.Unmarshal(b, &n); err == nil {
		*fi = FlexInt(n)
		return nil
	}
	// Try string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if s == "" {
			*fi = 0
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*fi = FlexInt(n)
		return nil
	}
	return nil
}

// Task represents a ClickUp task.
type Task struct {
	ID           string       `json:"id"`
	CustomID     *string      `json:"custom_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Status       Status       `json:"status"`
	Priority     *Priority    `json:"priority"`
	Assignees    []User       `json:"assignees"`
	Watchers     []User       `json:"watchers"`
	Creator      User         `json:"creator"`
	List         ListRef      `json:"list"`
	Space        SpaceRef     `json:"space"`
	Tags         []Tag        `json:"tags"`
	Parent       *string      `json:"parent"`
	DueDate      *string      `json:"due_date"`
	DateCreated  string       `json:"date_created"`
	TimeEstimate *int64       `json:"time_estimate"`
	TimeSpent    *int64       `json:"time_spent"`
	URL          string       `json:"url"`
	Subtasks     []Task       `json:"subtasks"`
	Dependencies []Dependency `json:"dependencies"`
	LinkedTasks  []LinkedTask `json:"linked_tasks"`
}

type Status struct {
	Status string `json:"status"`
	Type   string `json:"type"`
}

type Priority struct {
	Priority string `json:"priority"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Tag struct {
	Name string `json:"name"`
}

type ListRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SpaceRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Dependency struct {
	TaskID      string `json:"task_id"`
	DependsOn   string `json:"depends_on"`
	Type        int    `json:"type"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"userid"`
}

type LinkedTask struct {
	LinkID      string `json:"link_id"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"userid"`
	WorkspaceID string `json:"workspace_id"`
}

// TasksResponse wraps the tasks array from the API.
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// Space represents a ClickUp space.
type Space struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Private  bool     `json:"private"`
	Statuses []Status `json:"statuses"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

// Folder represents a ClickUp folder.
type Folder struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Hidden bool       `json:"hidden"`
	Lists  []ListInfo `json:"lists"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

// ListInfo represents full list metadata.
// Note: the ClickUp API returns task_count as a string in some endpoints
// and status as null for lists inside folders.
type ListInfo struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Status            *Status  `json:"status"`
	TaskCount         FlexInt  `json:"task_count"`
	PermissionLevel   string   `json:"permission_level"`
	Space             SpaceRef `json:"space"`
	Folder            ListRef  `json:"folder"`
	DueDateTime       bool     `json:"due_date_time"`
	MultipleAssignees bool     `json:"multiple_assignees"`
	TimeTracking      bool     `json:"time_tracking"`
	Statuses          []Status `json:"statuses"`
}

type ListsResponse struct {
	Lists []ListInfo `json:"lists"`
}

// Comment represents a ClickUp task comment.
type Comment struct {
	ID          string `json:"id"`
	CommentText string `json:"comment_text"`
	User        User   `json:"user"`
	Date        string `json:"date"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

// TimeEntry represents a ClickUp time entry.
type TimeEntry struct {
	ID          string `json:"id"`
	User        User   `json:"user"`
	Duration    string `json:"duration"`
	Start       string `json:"start"`
	Description string `json:"description"`
}

type TimeEntriesResponse struct {
	Data []TimeEntry `json:"data"`
}

// Document represents a ClickUp document.
type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	DateCreated string `json:"date_created"`
	Creator     User   `json:"creator"`
}

type DocsResponse struct {
	Docs []Document `json:"docs"`
}
