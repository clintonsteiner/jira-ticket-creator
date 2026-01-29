package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// JSONRepository implements Repository using JSON files
type JSONRepository struct {
	filepath string
}

// NewJSONRepository creates a new JSON-based repository
func NewJSONRepository(path string) (*JSONRepository, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &JSONRepository{
		filepath: path,
	}, nil
}

// Save persists ticket records to JSON file
func (r *JSONRepository) Save(records []jira.TicketRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal records: %w", err)
	}

	if err := os.WriteFile(r.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// Load retrieves all ticket records from JSON file
func (r *JSONRepository) Load() ([]jira.TicketRecord, error) {
	data, err := os.ReadFile(r.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []jira.TicketRecord{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var records []jira.TicketRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal records: %w", err)
	}

	return records, nil
}

// Add adds a new ticket record
func (r *JSONRepository) Add(record jira.TicketRecord) error {
	records, err := r.Load()
	if err != nil {
		return err
	}

	// Check if record already exists
	for i, existing := range records {
		if existing.Key == record.Key {
			records[i] = record
			return r.Save(records)
		}
	}

	records = append(records, record)
	return r.Save(records)
}

// GetByKey retrieves a ticket record by key
func (r *JSONRepository) GetByKey(key string) (*jira.TicketRecord, error) {
	records, err := r.Load()
	if err != nil {
		return nil, err
	}

	for i := range records {
		if records[i].Key == key {
			return &records[i], nil
		}
	}

	return nil, fmt.Errorf("ticket not found: %s", key)
}

// GetAll retrieves all ticket records
func (r *JSONRepository) GetAll() ([]jira.TicketRecord, error) {
	return r.Load()
}

// Update updates an existing ticket record
func (r *JSONRepository) Update(record jira.TicketRecord) error {
	records, err := r.Load()
	if err != nil {
		return err
	}

	found := false
	for i := range records {
		if records[i].Key == record.Key {
			records[i] = record
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("ticket not found: %s", record.Key)
	}

	return r.Save(records)
}
