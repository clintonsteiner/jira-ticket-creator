package jira

import (
	"fmt"
	"strings"
)

// Validator handles validation of JIRA data
type Validator struct {
	client *Client
}

// NewValidator creates a new validator
func NewValidator(client *Client) *Validator {
	return &Validator{client: client}
}

// ValidatePriority checks if a priority value is valid
func (v *Validator) ValidatePriority(priority string) error {
	validPriorities := []string{"Lowest", "Low", "Medium", "High", "Highest"}
	for _, valid := range validPriorities {
		if strings.EqualFold(priority, valid) {
			return nil
		}
	}
	return &ValidationError{
		Field:   "priority",
		Message: "invalid priority",
		Details: fmt.Sprintf("must be one of: %s", strings.Join(validPriorities, ", ")),
	}
}

// ValidateIssueType checks if an issue type is valid
func (v *Validator) ValidateIssueType(issueType string) error {
	validTypes := []string{"Task", "Story", "Bug", "Epic", "Subtask"}
	for _, valid := range validTypes {
		if strings.EqualFold(issueType, valid) {
			return nil
		}
	}
	return &ValidationError{
		Field:   "issue_type",
		Message: "invalid issue type",
		Details: fmt.Sprintf("must be one of: %s", strings.Join(validTypes, ", ")),
	}
}

// ValidateTicketExists checks if a ticket with the given key exists
func (v *Validator) ValidateTicketExists(key string) error {
	_, err := v.client.GetIssue(key)
	if _, ok := err.(*NotFoundError); ok {
		return &ValidationError{
			Field:   "ticket_key",
			Message: fmt.Sprintf("ticket not found: %s", key),
			Details: "verify the ticket key is correct",
		}
	}
	return err
}

// ValidateTicketsExist checks if multiple tickets with the given keys exist
func (v *Validator) ValidateTicketsExist(keys []string) error {
	for _, key := range keys {
		if err := v.ValidateTicketExists(key); err != nil {
			return err
		}
	}
	return nil
}

// ValidateSummary checks if a summary is valid
func (v *Validator) ValidateSummary(summary string) error {
	if summary == "" {
		return &ValidationError{
			Field:   "summary",
			Message: "summary is required",
		}
	}
	if len(summary) > 255 {
		return &ValidationError{
			Field:   "summary",
			Message: "summary is too long",
			Details: fmt.Sprintf("maximum 255 characters, got %d", len(summary)),
		}
	}
	return nil
}

// ValidateDescription checks if a description is valid
func (v *Validator) ValidateDescription(description string) error {
	// No strict limits on description in Jira Cloud
	return nil
}

// ValidateCreateIssueRequest validates a create issue request
func (v *Validator) ValidateCreateIssueRequest(fields IssueFields) error {
	// Validate summary
	if err := v.ValidateSummary(fields.Summary); err != nil {
		return err
	}

	// Validate issue type
	if err := v.ValidateIssueType(fields.IssueType.Name); err != nil {
		return err
	}

	// Validate priority if set
	if fields.Priority != nil && fields.Priority.Name != "" {
		if err := v.ValidatePriority(fields.Priority.Name); err != nil {
			return err
		}
	}

	// Validate project key
	if fields.Project.Key == "" {
		return &ValidationError{
			Field:   "project",
			Message: "project key is required",
		}
	}

	return nil
}

// ValidateLabel checks if a label is valid
func (v *Validator) ValidateLabel(label string) error {
	if label == "" {
		return &ValidationError{
			Field:   "labels",
			Message: "label cannot be empty",
		}
	}
	return nil
}

// ValidateLabels checks if multiple labels are valid
func (v *Validator) ValidateLabels(labels []string) error {
	for _, label := range labels {
		if err := v.ValidateLabel(label); err != nil {
			return err
		}
	}
	return nil
}
