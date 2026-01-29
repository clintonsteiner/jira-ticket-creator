package reports

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// TableReporter generates table format reports
type TableReporter struct{}

// Generate generates a table report
func (r *TableReporter) Generate(issues []jira.Issue) (string, error) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintln(w, "KEY\tTYPE\tSUMMARY\tASSIGNEE\tPRIORITY")
	fmt.Fprintln(w, "---\t----\t-------\t--------\t--------")

	// Write data
	for _, issue := range issues {
		assignee := ""
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.Name
			if assignee == "" {
				assignee = issue.Fields.Assignee.EmailAddress
			}
		}

		priority := ""
		if issue.Fields.Priority != nil {
			priority = issue.Fields.Priority.Name
		}

		// Truncate long summary
		summary := issue.Fields.Summary
		if len(summary) > 40 {
			summary = summary[:37] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			issue.Key,
			issue.Fields.IssueType.Name,
			summary,
			assignee,
			priority)
	}

	w.Flush()
	return buf.String(), nil
}
