package batch

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONTicketData represents ticket data in JSON format
type JSONTicketData struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	IssueType   string   `json:"issue_type,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Components  []string `json:"components,omitempty"`
	BlockedBy   []string `json:"blocked_by,omitempty"`
}

// ParseJSONFile parses a JSON file and returns ticket data
// Expected format: array of ticket objects
func ParseJSONFile(filepath string) ([]TicketData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var jsonTickets []JSONTicketData
	if err := json.Unmarshal(data, &jsonTickets); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(jsonTickets) == 0 {
		return nil, fmt.Errorf("JSON file contains no tickets")
	}

	// Convert to TicketData
	var tickets []TicketData
	for i, jt := range jsonTickets {
		if jt.Summary == "" {
			return nil, fmt.Errorf("ticket %d: summary is required", i)
		}

		ticket := TicketData{
			Summary:     jt.Summary,
			Description: jt.Description,
			IssueType:   jt.IssueType,
			Priority:    jt.Priority,
			Assignee:    jt.Assignee,
			Labels:      jt.Labels,
			Components:  jt.Components,
			BlockedBy:   jt.BlockedBy,
		}

		// Set defaults
		if ticket.IssueType == "" {
			ticket.IssueType = "Task"
		}
		if ticket.Priority == "" {
			ticket.Priority = "Medium"
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
