package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatTimestamp converts a ClickUp millisecond Unix timestamp to a human-readable string.
func FormatTimestamp(ms string) string {
	if ms == "" {
		return ""
	}
	n, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return ms
	}
	return time.UnixMilli(n).Format("2006-01-02 15:04")
}

// FormatDurationMs converts milliseconds to a human-friendly duration.
func FormatDurationMs(ms int64) string {
	minutes := ms / 1000 / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	remaining := minutes % 60
	if remaining == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, remaining)
}

// or returns fallback if s is empty.
func Or(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// FormatTaskSummary formats a task for list/search views (compact).
func FormatTaskSummary(t Task) string {
	var b strings.Builder

	id := t.ID
	if t.CustomID != nil && *t.CustomID != "" {
		id = *t.CustomID
	}

	assignees := make([]string, len(t.Assignees))
	for i, a := range t.Assignees {
		assignees[i] = a.Username
	}

	desc := t.Description
	if len(desc) > 200 {
		desc = desc[:200] + "..."
	}

	fmt.Fprintf(&b, "%s  %s  [%s]\n", id, t.Name, t.Status.Status)
	if len(assignees) > 0 {
		fmt.Fprintf(&b, "  Assignees: %s\n", strings.Join(assignees, ", "))
	}
	if t.DueDate != nil {
		fmt.Fprintf(&b, "  Due: %s\n", FormatTimestamp(*t.DueDate))
	}
	if desc != "" {
		fmt.Fprintf(&b, "  %s\n", desc)
	}
	fmt.Fprintf(&b, "  %s\n", t.URL)

	return b.String()
}

// FormatTaskDetail formats a task for the detail view (full info).
func FormatTaskDetail(t Task) string {
	var b strings.Builder

	id := t.ID
	if t.CustomID != nil && *t.CustomID != "" {
		id = fmt.Sprintf("%s (%s)", *t.CustomID, t.ID)
	}

	assignees := make([]string, len(t.Assignees))
	for i, a := range t.Assignees {
		assignees[i] = a.Username
	}
	watchers := make([]string, len(t.Watchers))
	for i, w := range t.Watchers {
		watchers[i] = w.Username
	}
	tags := make([]string, len(t.Tags))
	for i, tg := range t.Tags {
		tags[i] = tg.Name
	}

	priority := "None"
	if t.Priority != nil {
		priority = t.Priority.Priority
	}
	parent := "None (top-level task)"
	if t.Parent != nil && *t.Parent != "" {
		parent = *t.Parent
	}

	fmt.Fprintf(&b, "%s\n", t.Name)
	fmt.Fprintf(&b, "========================================\n\n")

	fmt.Fprintf(&b, "ID:       %s\n", id)
	fmt.Fprintf(&b, "Status:   %s\n", t.Status.Status)
	fmt.Fprintf(&b, "Priority: %s\n", priority)
	fmt.Fprintf(&b, "Created:  %s\n", FormatTimestamp(t.DateCreated))
	if t.DueDate != nil {
		fmt.Fprintf(&b, "Due:      %s\n", FormatTimestamp(*t.DueDate))
	}
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "Assignees: %s\n", Or(strings.Join(assignees, ", "), "Unassigned"))
	fmt.Fprintf(&b, "Watchers:  %s\n", Or(strings.Join(watchers, ", "), "None"))
	fmt.Fprintf(&b, "Creator:   %s\n", t.Creator.Username)
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "List:   %s (ID: %s)\n", t.List.Name, t.List.ID)
	fmt.Fprintf(&b, "Space:  %s\n", t.Space.Name)
	fmt.Fprintf(&b, "Tags:   %s\n", Or(strings.Join(tags, ", "), "None"))
	fmt.Fprintf(&b, "Parent: %s\n", parent)
	fmt.Fprintf(&b, "\n")

	if t.Description != "" {
		fmt.Fprintf(&b, "Description:\n%s\n\n", t.Description)
	}

	fmt.Fprintf(&b, "URL: %s\n", t.URL)
	fmt.Fprintf(&b, "\n")

	// Time tracking
	var est, spent int64
	if t.TimeEstimate != nil {
		est = *t.TimeEstimate
	}
	if t.TimeSpent != nil {
		spent = *t.TimeSpent
	}
	fmt.Fprintf(&b, "Time Estimated: %s\n", FormatDurationMs(est))
	fmt.Fprintf(&b, "Time Spent:     %s\n", FormatDurationMs(spent))

	// Subtasks
	if len(t.Subtasks) > 0 {
		fmt.Fprintf(&b, "\nSubtasks (%d):\n", len(t.Subtasks))
		for _, st := range t.Subtasks {
			stID := st.ID
			if st.CustomID != nil && *st.CustomID != "" {
				stID = *st.CustomID
			}
			stAssignees := make([]string, len(st.Assignees))
			for i, a := range st.Assignees {
				stAssignees[i] = a.Username
			}
			fmt.Fprintf(&b, "  - %s [%s] %s (%s)\n", stID, st.Status.Status, st.Name, Or(strings.Join(stAssignees, ", "), "Unassigned"))
		}
	}

	// Relationships
	if len(t.Dependencies) > 0 || len(t.LinkedTasks) > 0 {
		fmt.Fprintf(&b, "\nRelationships:\n")
		if len(t.Dependencies) > 0 {
			fmt.Fprintf(&b, "  Dependencies (%d):\n", len(t.Dependencies))
			for _, d := range t.Dependencies {
				fmt.Fprintf(&b, "    - %s (type: %d)\n", d.DependsOn, d.Type)
			}
		}
		if len(t.LinkedTasks) > 0 {
			fmt.Fprintf(&b, "  Linked Tasks (%d):\n", len(t.LinkedTasks))
			for _, l := range t.LinkedTasks {
				fmt.Fprintf(&b, "    - %s (created: %s)\n", l.LinkID, FormatTimestamp(l.DateCreated))
			}
		}
	}

	return b.String()
}
