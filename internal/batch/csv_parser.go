package batch

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// TicketData represents a single ticket to create
type TicketData struct {
	Summary     string
	Description string
	IssueType   string
	Priority    string
	Assignee    string
	Labels      []string
	Components  []string
	BlockedBy   []string
}

// ParseCSVFile parses a CSV file and returns ticket data
// Expected columns: summary,description,issue_type,priority,assignee,labels,components,blocked_by
func ParseCSVFile(filepath string) ([]TicketData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Parse header
	header := records[0]
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Validate required columns
	if _, ok := columnMap["summary"]; !ok {
		return nil, fmt.Errorf("CSV must have 'summary' column")
	}

	// Parse data rows
	var tickets []TicketData
	for i, record := range records[1:] {
		ticket := TicketData{
			IssueType: "Task",   // Default
			Priority:  "Medium", // Default
		}

		// Parse summary (required)
		if idx, ok := columnMap["summary"]; ok && idx < len(record) {
			ticket.Summary = strings.TrimSpace(record[idx])
		}

		// Parse other fields
		if idx, ok := columnMap["description"]; ok && idx < len(record) {
			ticket.Description = strings.TrimSpace(record[idx])
		}

		if idx, ok := columnMap["issue_type"]; ok && idx < len(record) {
			val := strings.TrimSpace(record[idx])
			if val != "" {
				ticket.IssueType = val
			}
		}

		if idx, ok := columnMap["priority"]; ok && idx < len(record) {
			val := strings.TrimSpace(record[idx])
			if val != "" {
				ticket.Priority = val
			}
		}

		if idx, ok := columnMap["assignee"]; ok && idx < len(record) {
			ticket.Assignee = strings.TrimSpace(record[idx])
		}

		if idx, ok := columnMap["labels"]; ok && idx < len(record) {
			labels := strings.TrimSpace(record[idx])
			if labels != "" {
				ticket.Labels = strings.Split(labels, ",")
				for i := range ticket.Labels {
					ticket.Labels[i] = strings.TrimSpace(ticket.Labels[i])
				}
			}
		}

		if idx, ok := columnMap["components"]; ok && idx < len(record) {
			components := strings.TrimSpace(record[idx])
			if components != "" {
				ticket.Components = strings.Split(components, ",")
				for i := range ticket.Components {
					ticket.Components[i] = strings.TrimSpace(ticket.Components[i])
				}
			}
		}

		if idx, ok := columnMap["blocked_by"]; ok && idx < len(record) {
			blockedBy := strings.TrimSpace(record[idx])
			if blockedBy != "" {
				ticket.BlockedBy = strings.Split(blockedBy, ",")
				for i := range ticket.BlockedBy {
					ticket.BlockedBy[i] = strings.TrimSpace(ticket.BlockedBy[i])
				}
			}
		}

		// Validate
		if ticket.Summary == "" {
			return nil, fmt.Errorf("row %d: summary is required", i+2)
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
