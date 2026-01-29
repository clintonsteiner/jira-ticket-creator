package interactive

import (
	"fmt"
	"strings"
)

// TicketWizard guides users through creating a ticket interactively
type TicketWizard struct{}

// TicketInput holds the user's ticket input
type TicketInput struct {
	Summary     string
	Description string
	Type        string
	Priority    string
	Assignee    string
	Labels      []string
	Components  []string
	BlockedBy   []string
}

// Run runs the interactive wizard
func (w *TicketWizard) Run() (*TicketInput, error) {
	fmt.Println("\n📝 Creating a new JIRA ticket...")
	fmt.Println("================================\n")

	input := &TicketInput{}

	// Prompt for summary
	summary, err := PromptString("Ticket Summary", true)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}
	input.Summary = summary

	// Prompt for description
	description, err := PromptStringWithDefault("Ticket Description", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}
	input.Description = description

	// Prompt for issue type
	types := []string{"Task", "Story", "Bug", "Epic", "Subtask"}
	issueType, err := PromptSelect("Issue Type", types)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue type: %w", err)
	}
	input.Type = issueType

	// Prompt for priority
	priorities := []string{"Lowest", "Low", "Medium", "High", "Highest"}
	priority, err := PromptSelect("Priority", priorities)
	if err != nil {
		return nil, fmt.Errorf("failed to get priority: %w", err)
	}
	input.Priority = priority

	// Prompt for assignee
	assignee, err := PromptStringWithDefault("Assignee Email", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}
	input.Assignee = assignee

	// Prompt for labels
	fmt.Println("\nLabels (comma-separated):")
	labelsStr, err := PromptStringWithDefault("Labels", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
	}
	if labelsStr != "" {
		input.Labels = strings.Split(labelsStr, ",")
		for i := range input.Labels {
			input.Labels[i] = strings.TrimSpace(input.Labels[i])
		}
	}

	// Prompt for blocked-by
	fmt.Println("\nBlocked By (comma-separated ticket keys):")
	blockedByStr, err := PromptStringWithDefault("Blocked By", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked-by: %w", err)
	}
	if blockedByStr != "" {
		input.BlockedBy = strings.Split(blockedByStr, ",")
		for i := range input.BlockedBy {
			input.BlockedBy[i] = strings.TrimSpace(input.BlockedBy[i])
		}
	}

	// Show summary
	fmt.Println("\n✅ Ticket Summary:")
	fmt.Println("=================")
	fmt.Printf("Summary:    %s\n", input.Summary)
	fmt.Printf("Description: %s\n", input.Description)
	fmt.Printf("Type:       %s\n", input.Type)
	fmt.Printf("Priority:   %s\n", input.Priority)
	if input.Assignee != "" {
		fmt.Printf("Assignee:   %s\n", input.Assignee)
	}
	if len(input.Labels) > 0 {
		fmt.Printf("Labels:     %s\n", strings.Join(input.Labels, ", "))
	}
	if len(input.BlockedBy) > 0 {
		fmt.Printf("Blocked By: %s\n", strings.Join(input.BlockedBy, ", "))
	}

	// Confirm
	confirm, err := PromptConfirm("\nCreate ticket with these details?")
	if err != nil {
		return nil, fmt.Errorf("failed to confirm: %w", err)
	}

	if !confirm {
		fmt.Println("❌ Ticket creation cancelled")
		return nil, fmt.Errorf("user cancelled")
	}

	return input, nil
}
