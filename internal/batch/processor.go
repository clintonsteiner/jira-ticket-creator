package batch

import (
	"fmt"
	"sync"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// ProcessResult represents the result of processing a single ticket
type ProcessResult struct {
	Index      int
	TicketData TicketData
	CreatedKey string
	Error      error
	Status     string
}

// BatchProcessor handles batch ticket operations
type BatchProcessor struct {
	client        *jira.Client
	projectKey    string
	maxConcurrent int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(client *jira.Client, projectKey string) *BatchProcessor {
	return &BatchProcessor{
		client:        client,
		projectKey:    projectKey,
		maxConcurrent: 3, // Default concurrent operations
	}
}

// ValidateTickets validates all tickets before creation
func (bp *BatchProcessor) ValidateTickets(tickets []TicketData, validator *jira.Validator) []ProcessResult {
	var results []ProcessResult

	for i, ticket := range tickets {
		result := ProcessResult{
			Index:      i,
			TicketData: ticket,
			Status:     "validated",
		}

		// Validate summary
		if ticket.Summary == "" {
			result.Error = fmt.Errorf("summary is required")
			result.Status = "failed"
			results = append(results, result)
			continue
		}

		// Validate issue type
		if err := validator.ValidateIssueType(ticket.IssueType); err != nil {
			result.Error = err
			result.Status = "failed"
			results = append(results, result)
			continue
		}

		// Validate priority if set
		if ticket.Priority != "" {
			if err := validator.ValidatePriority(ticket.Priority); err != nil {
				result.Error = err
				result.Status = "failed"
				results = append(results, result)
				continue
			}
		}

		// Validate blocked-by tickets exist
		if len(ticket.BlockedBy) > 0 {
			if err := validator.ValidateTicketsExist(ticket.BlockedBy); err != nil {
				result.Error = fmt.Errorf("blocked-by validation failed: %w", err)
				result.Status = "failed"
				results = append(results, result)
				continue
			}
		}

		results = append(results, result)
	}

	return results
}

// CreateTickets creates all validated tickets
func (bp *BatchProcessor) CreateTickets(tickets []TicketData) []ProcessResult {
	var results []ProcessResult
	var resultsMutex sync.Mutex

	// Use a semaphore for concurrency control
	semaphore := make(chan struct{}, bp.maxConcurrent)
	var wg sync.WaitGroup

	for i, ticket := range tickets {
		wg.Add(1)
		go func(index int, t TicketData) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			result := bp.createSingleTicket(index, t)

			resultsMutex.Lock()
			results = append(results, result)
			resultsMutex.Unlock()
		}(i, ticket)
	}

	wg.Wait()
	return results
}

// LinkTickets creates links between tickets based on blocked-by relationships
func (bp *BatchProcessor) LinkTickets(createResults []ProcessResult) []ProcessResult {
	linkService := jira.NewLinkService(bp.client)
	var results []ProcessResult

	for _, createResult := range createResults {
		if createResult.Error != nil {
			continue // Skip if creation failed
		}

		for _, blockerKey := range createResult.TicketData.BlockedBy {
			// Create link: blockerKey blocks createdKey
			if err := linkService.LinkBlocks(blockerKey, createResult.CreatedKey); err != nil {
				results = append(results, ProcessResult{
					Index:      createResult.Index,
					TicketData: createResult.TicketData,
					CreatedKey: createResult.CreatedKey,
					Error:      fmt.Errorf("failed to link %s blocks %s: %w", blockerKey, createResult.CreatedKey, err),
					Status:     "partial", // Created but linking failed
				})
				continue
			}
		}
	}

	return results
}

// createSingleTicket creates a single ticket
func (bp *BatchProcessor) createSingleTicket(index int, ticket TicketData) ProcessResult {
	issueService := jira.NewIssueService(bp.client)

	fields := jira.IssueFields{
		Project: jira.Project{
			Key: bp.projectKey,
		},
		Summary:     ticket.Summary,
		Description: ticket.Description,
		IssueType: jira.IssueType{
			Name: ticket.IssueType,
		},
	}

	// Add optional fields
	if ticket.Priority != "" {
		fields.Priority = &jira.Priority{
			Name: ticket.Priority,
		}
	}

	if ticket.Assignee != "" {
		fields.Assignee = &jira.User{
			EmailAddress: ticket.Assignee,
		}
	}

	if len(ticket.Labels) > 0 {
		fields.Labels = ticket.Labels
	}

	if len(ticket.Components) > 0 {
		components := make([]jira.Component, len(ticket.Components))
		for i, comp := range ticket.Components {
			components[i] = jira.Component{Name: comp}
		}
		fields.Components = components
	}

	resp, err := issueService.CreateIssueWithFields(fields)
	if err != nil {
		return ProcessResult{
			Index:      index,
			TicketData: ticket,
			Error:      err,
			Status:     "failed",
		}
	}

	return ProcessResult{
		Index:      index,
		TicketData: ticket,
		CreatedKey: resp.Key,
		Status:     "created",
	}
}

// PrintResults prints batch processing results
func PrintResults(results []ProcessResult, verbose bool) {
	successCount := 0
	failureCount := 0
	partialCount := 0

	fmt.Println("\n📊 Batch Processing Results:")
	fmt.Println("=====================================")

	for _, result := range results {
		switch result.Status {
		case "created":
			successCount++
			if verbose {
				fmt.Printf("✅ [%d] %s -> %s\n", result.Index+1, result.TicketData.Summary, result.CreatedKey)
			}
		case "failed":
			failureCount++
			fmt.Printf("❌ [%d] %s: %v\n", result.Index+1, result.TicketData.Summary, result.Error)
		case "partial":
			partialCount++
			if verbose {
				fmt.Printf("⚠️  [%d] %s (%s): %v\n", result.Index+1, result.TicketData.Summary, result.CreatedKey, result.Error)
			}
		case "validated":
			if verbose {
				fmt.Printf("✓ [%d] %s (validated)\n", result.Index+1, result.TicketData.Summary)
			}
		}
	}

	fmt.Println("=====================================")
	fmt.Printf("Success: %d | Failures: %d | Partial: %d\n", successCount, failureCount, partialCount)

	if failureCount == 0 && partialCount == 0 {
		fmt.Println("🎉 All tickets processed successfully!")
	}
}
