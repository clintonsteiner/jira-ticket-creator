package reports

import (
	"encoding/json"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// JSONReporter generates JSON reports
type JSONReporter struct{}

// Generate generates a JSON report
func (r *JSONReporter) Generate(issues []jira.Issue) (string, error) {
	data, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
