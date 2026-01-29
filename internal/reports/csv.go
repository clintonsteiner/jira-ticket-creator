package reports

import (
	"bytes"
	"encoding/csv"
	"strings"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// CSVReporter generates CSV reports
type CSVReporter struct{}

// Generate generates a CSV report
func (r *CSVReporter) Generate(issues []jira.Issue) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Key", "Type", "Summary", "Status", "Assignee", "Priority", "Labels"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

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

		labels := strings.Join(issue.Fields.Labels, ";")

		row := []string{
			issue.Key,
			issue.Fields.IssueType.Name,
			issue.Fields.Summary,
			"", // Status would come from issue.Fields.Status if available
			assignee,
			priority,
			labels,
		}

		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
