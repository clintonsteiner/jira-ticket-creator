package jira

import (
	"testing"
)

func TestExtractProjectKey(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{"Valid ticket", "PROJ-123", "PROJ", false},
		{"Single letter project", "A-1", "A", false},
		{"Multi-char project", "ABC-999", "ABC", false},
		{"Mixed alphanumeric project", "TEST123-456", "TEST123", false},
		{"With spaces", "  PROJ-123  ", "PROJ", false},
		{"Invalid: no hyphen", "PROJ123", "", true},
		{"Invalid: wrong format", "proj-123", "", true},
		{"Invalid: no numbers", "PROJ-ABC", "", true},
		{"Invalid: empty", "", "", true},
		{"Invalid: only dash", "-123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractProjectKey(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ExtractProjectKey() error = %v, wantError %v", err, tt.wantError)
			}
			if result != tt.expected {
				t.Errorf("ExtractProjectKey() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestIsTicketKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid ticket", "PROJ-123", true},
		{"Valid single letter", "A-1", true},
		{"Valid multi-char", "ABC-999", true},
		{"Invalid: no hyphen", "PROJ123", false},
		{"Invalid: lowercase", "proj-123", false},
		{"Invalid: no numbers", "PROJ-ABC", false},
		{"Invalid: empty", "", false},
		{"With spaces: trimmed and valid", "  PROJ-123  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTicketKey(tt.input)
			if result != tt.expected {
				t.Errorf("IsTicketKey(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsProjectKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid: single letter", "A", true},
		{"Valid: multiple letters", "PROJ", true},
		{"Valid: alphanumeric", "PROJ123", true},
		{"Invalid: contains hyphen", "PROJ-1", false},
		{"Invalid: lowercase", "proj", false},
		{"Invalid: empty", "", false},
		{"Invalid: contains space", "PR OJ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsProjectKey(tt.input)
			if result != tt.expected {
				t.Errorf("IsProjectKey(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResolveProject(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{"Project key", "PROJ", "PROJ", false},
		{"Ticket key", "PROJ-123", "PROJ", false},
		{"Single letter project", "A", "A", false},
		{"Single letter ticket", "A-1", "A", false},
		{"Invalid format", "invalid", "", true},
		{"Empty string", "", "", true},
		{"With spaces", "  PROJ  ", "PROJ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveProject(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ResolveProject() error = %v, wantError %v", err, tt.wantError)
			}
			if result != tt.expected {
				t.Errorf("ResolveProject() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestExtractTicketNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{"Valid ticket", "PROJ-123", "123", false},
		{"Single digit", "A-1", "1", false},
		{"Large number", "TEST-999999", "999999", false},
		{"Invalid format", "PROJ123", "", true},
		{"Invalid format", "proj-123", "", true},
		{"Empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractTicketNumber(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ExtractTicketNumber() error = %v, wantError %v", err, tt.wantError)
			}
			if result != tt.expected {
				t.Errorf("ExtractTicketNumber() = %s, want %s", result, tt.expected)
			}
		})
	}
}
