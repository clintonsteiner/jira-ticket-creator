package storage

import "github.com/clintonsteiner/jira-ticket-creator/internal/jira"

// Repository defines the interface for storing and retrieving ticket records
type Repository interface {
	// Save persists ticket records
	Save(records []jira.TicketRecord) error

	// Load retrieves all ticket records
	Load() ([]jira.TicketRecord, error)

	// Add adds a new ticket record
	Add(record jira.TicketRecord) error

	// GetByKey retrieves a ticket record by key
	GetByKey(key string) (*jira.TicketRecord, error)

	// GetAll retrieves all ticket records
	GetAll() ([]jira.TicketRecord, error)

	// Update updates an existing ticket record
	Update(record jira.TicketRecord) error
}
