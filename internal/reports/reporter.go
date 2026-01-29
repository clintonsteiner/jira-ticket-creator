package reports

import (
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// Reporter defines the interface for generating reports
type Reporter interface {
	Generate(issues []jira.Issue) (string, error)
}

// ReporterFactory creates reporters based on format
func NewReporter(format string) Reporter {
	switch format {
	case "json":
		return &JSONReporter{}
	case "csv":
		return &CSVReporter{}
	case "markdown":
		return &MarkdownReporter{}
	case "html":
		return &HTMLReporter{}
	default:
		return &TableReporter{}
	}
}
