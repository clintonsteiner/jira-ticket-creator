package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProjectMapping represents the structure for mapping ticket prefixes to logical projects
type ProjectMapping struct {
	Mappings map[string]ProjectInfo `json:"mappings"`
}

// ProjectInfo contains information about a logical project
type ProjectInfo struct {
	TicketKeys  []string `json:"ticket_keys"` // Ticket key prefixes (e.g., "PROJ", "BACK")
	Description string   `json:"description"` // Human-readable description
}

// DefaultMappingPath returns the default path for the project mapping file
func DefaultMappingPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.jira/project-mapping.json"
	}
	return filepath.Join(home, ".jira", "project-mapping.json")
}

// LoadMapping loads the project mapping from file
func LoadMapping(path string) (*ProjectMapping, error) {
	if path == "" {
		path = DefaultMappingPath()
	}

	// Expand ~ to home directory
	if path == "~/.jira/project-mapping.json" {
		path = DefaultMappingPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty mapping if file doesn't exist
			return &ProjectMapping{Mappings: make(map[string]ProjectInfo)}, nil
		}
		return nil, fmt.Errorf("failed to read project mapping file: %w", err)
	}

	var pm ProjectMapping
	if err := json.Unmarshal(data, &pm); err != nil {
		return nil, fmt.Errorf("failed to parse project mapping file: %w", err)
	}

	return &pm, nil
}

// SaveMapping saves the project mapping to file
func (pm *ProjectMapping) SaveMapping(path string) error {
	if path == "" {
		path = DefaultMappingPath()
	}

	// Expand ~ to home directory
	if path == "~/.jira/project-mapping.json" {
		path = DefaultMappingPath()
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project mapping: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write project mapping file: %w", err)
	}

	return nil
}

// FindProjectForKey finds the logical project for a given ticket key
func (pm *ProjectMapping) FindProjectForKey(ticketKey string) string {
	for project, info := range pm.Mappings {
		for _, prefix := range info.TicketKeys {
			if ticketKey == prefix || (len(ticketKey) > len(prefix) && ticketKey[:len(prefix)] == prefix) {
				return project
			}
		}
	}
	return ""
}

// AddMapping adds or updates a project mapping
func (pm *ProjectMapping) AddMapping(project string, info ProjectInfo) {
	pm.Mappings[project] = info
}
