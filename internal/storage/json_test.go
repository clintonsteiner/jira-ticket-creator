package storage

import (
	"os"
	"testing"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func TestJSONRepository_SaveAndLoad(t *testing.T) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Create repository
	repo, err := NewJSONRepository(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONRepository() error = %v", err)
	}

	// Create test data
	records := []jira.TicketRecord{
		{
			Key:       "PROJ-1",
			Summary:   "Test ticket 1",
			Status:    "To Do",
			BlockedBy: []string{},
			CreatedAt: time.Now(),
			Creator:   "user@email.com",
			Assignee:  "assignee@email.com",
		},
		{
			Key:       "PROJ-2",
			Summary:   "Test ticket 2",
			Status:    "In Progress",
			BlockedBy: []string{"PROJ-1"},
			CreatedAt: time.Now(),
			Creator:   "user@email.com",
		},
	}

	// Save records
	err = repo.Save(records)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load records
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != len(records) {
		t.Errorf("Load() returned %d records, expected %d", len(loaded), len(records))
	}

	if loaded[0].Key != "PROJ-1" {
		t.Errorf("Load() first record key = %s, expected PROJ-1", loaded[0].Key)
	}

	if loaded[1].Key != "PROJ-2" {
		t.Errorf("Load() second record key = %s, expected PROJ-2", loaded[1].Key)
	}
}

func TestJSONRepository_Add(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	repo, err := NewJSONRepository(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONRepository() error = %v", err)
	}

	// Add first record
	record1 := jira.TicketRecord{
		Key:       "PROJ-1",
		Summary:   "Test ticket 1",
		Status:    "To Do",
		CreatedAt: time.Now(),
	}

	err = repo.Add(record1)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Add second record
	record2 := jira.TicketRecord{
		Key:       "PROJ-2",
		Summary:   "Test ticket 2",
		Status:    "In Progress",
		CreatedAt: time.Now(),
	}

	err = repo.Add(record2)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Load and verify
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("Load() returned %d records, expected 2", len(loaded))
	}
}

func TestJSONRepository_GetByKey(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	repo, err := NewJSONRepository(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONRepository() error = %v", err)
	}

	record := jira.TicketRecord{
		Key:       "PROJ-123",
		Summary:   "Test ticket",
		Status:    "To Do",
		CreatedAt: time.Now(),
	}

	err = repo.Add(record)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Get by key
	found, err := repo.GetByKey("PROJ-123")
	if err != nil {
		t.Fatalf("GetByKey() error = %v", err)
	}

	if found == nil {
		t.Fatal("GetByKey() returned nil")
	}

	if found.Key != "PROJ-123" {
		t.Errorf("GetByKey() returned key = %s, expected PROJ-123", found.Key)
	}

	// Try to get non-existent key
	_, err = repo.GetByKey("NONEXISTENT")
	if err == nil {
		t.Error("GetByKey() expected error for non-existent key")
	}
}

func TestJSONRepository_Update(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	repo, err := NewJSONRepository(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONRepository() error = %v", err)
	}

	record := jira.TicketRecord{
		Key:       "PROJ-1",
		Summary:   "Original summary",
		Status:    "To Do",
		CreatedAt: time.Now(),
	}

	err = repo.Add(record)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Update record
	record.Summary = "Updated summary"
	record.Status = "In Progress"

	err = repo.Update(record)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Load and verify
	updated, err := repo.GetByKey("PROJ-1")
	if err != nil {
		t.Fatalf("GetByKey() error = %v", err)
	}

	if updated.Summary != "Updated summary" {
		t.Errorf("Update() summary = %s, expected 'Updated summary'", updated.Summary)
	}

	if updated.Status != "In Progress" {
		t.Errorf("Update() status = %s, expected 'In Progress'", updated.Status)
	}
}

func TestJSONRepository_LoadEmptyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	repo, err := NewJSONRepository(tmpfile.Name())
	if err != nil {
		t.Fatalf("NewJSONRepository() error = %v", err)
	}

	// Load from empty file
	records, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Load() from empty file returned %d records, expected 0", len(records))
	}
}
