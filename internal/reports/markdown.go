package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// MarkdownReporter generates Markdown format reports
type MarkdownReporter struct{}

// Generate generates a Markdown report
func (r *MarkdownReporter) Generate(issues []jira.Issue) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# JIRA Tickets Report\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Tickets**: %d\n", len(issues)))

	// Count by type
	typeCount := make(map[string]int)
	priorityCount := make(map[string]int)
	for _, issue := range issues {
		typeCount[issue.Fields.IssueType.Name]++
		if issue.Fields.Priority != nil {
			priorityCount[issue.Fields.Priority.Name]++
		}
	}

	if len(typeCount) > 0 {
		sb.WriteString("- **By Type**:\n")
		for typ, count := range typeCount {
			sb.WriteString(fmt.Sprintf("  - %s: %d\n", typ, count))
		}
	}

	if len(priorityCount) > 0 {
		sb.WriteString("- **By Priority**:\n")
		for priority, count := range priorityCount {
			sb.WriteString(fmt.Sprintf("  - %s: %d\n", priority, count))
		}
	}

	sb.WriteString("\n")

	// Table
	sb.WriteString("## Tickets\n\n")
	sb.WriteString("| Key | Type | Summary | Assignee | Priority |\n")
	sb.WriteString("|-----|------|---------|----------|----------|\n")

	for _, issue := range issues {
		assignee := "-"
		if issue.Fields.Assignee != nil && issue.Fields.Assignee.Name != "" {
			assignee = issue.Fields.Assignee.Name
		} else if issue.Fields.Assignee != nil && issue.Fields.Assignee.EmailAddress != "" {
			assignee = issue.Fields.Assignee.EmailAddress
		}

		priority := "-"
		if issue.Fields.Priority != nil {
			priority = issue.Fields.Priority.Name
		}

		summary := issue.Fields.Summary
		if len(summary) > 50 {
			summary = summary[:47] + "..."
		}

		sb.WriteString(fmt.Sprintf("| [%s](#%s) | %s | %s | %s | %s |\n",
			issue.Key,
			strings.ToLower(strings.ReplaceAll(issue.Key, "-", "")),
			issue.Fields.IssueType.Name,
			summary,
			assignee,
			priority))
	}

	sb.WriteString("\n")

	// Details
	sb.WriteString("## Details\n\n")
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("### %s\n\n", issue.Key))
		sb.WriteString(fmt.Sprintf("**Summary**: %s\n\n", issue.Fields.Summary))

		if issue.Fields.Description != "" {
			sb.WriteString(fmt.Sprintf("**Description**: %s\n\n", issue.Fields.Description))
		}

		sb.WriteString(fmt.Sprintf("**Type**: %s\n\n", issue.Fields.IssueType.Name))

		if issue.Fields.Priority != nil {
			sb.WriteString(fmt.Sprintf("**Priority**: %s\n\n", issue.Fields.Priority.Name))
		}

		if issue.Fields.Assignee != nil && (issue.Fields.Assignee.Name != "" || issue.Fields.Assignee.EmailAddress != "") {
			name := issue.Fields.Assignee.Name
			if name == "" {
				name = issue.Fields.Assignee.EmailAddress
			}
			sb.WriteString(fmt.Sprintf("**Assignee**: %s\n\n", name))
		}

		if len(issue.Fields.Labels) > 0 {
			sb.WriteString(fmt.Sprintf("**Labels**: %s\n\n", strings.Join(issue.Fields.Labels, ", ")))
		}

		sb.WriteString("---\n\n")
	}

	return sb.String(), nil
}
