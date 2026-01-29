package batch

import (
	"os"
	"testing"
)

func TestParseJSONFile(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		expectError bool
		expectCount int
	}{
		{
			name: "valid JSON with all fields",
			jsonContent: `[
  {
    "summary": "Task 1",
    "description": "Description 1",
    "issue_type": "Task",
    "priority": "High",
    "assignee": "user@email.com",
    "labels": ["label1", "label2"],
    "components": ["Component1"],
    "blocked_by": ["PROJ-100"]
  }
]`,
			expectError: false,
			expectCount: 1,
		},
		{
			name: "valid JSON with multiple tickets",
			jsonContent: `[
  {
    "summary": "Task 1",
    "description": "Description 1"
  },
  {
    "summary": "Task 2",
    "description": "Description 2"
  }
]`,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "empty array",
			jsonContent: "[]",
			expectError: true,
			expectCount: 0,
		},
		{
			name:        "invalid JSON",
			jsonContent: "{invalid json}",
			expectError: true,
			expectCount: 0,
		},
		{
			name: "missing required summary",
			jsonContent: `[
  {
    "description": "Description 1"
  }
]`,
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.jsonContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpfile.Close()

			tickets, err := ParseJSONFile(tmpfile.Name())

			if (err != nil) != tt.expectError {
				t.Errorf("ParseJSONFile() error = %v, expectError %v", err, tt.expectError)
			}

			if len(tickets) != tt.expectCount {
				t.Errorf("ParseJSONFile() returned %d tickets, expected %d", len(tickets), tt.expectCount)
			}
		})
	}
}

func TestParseJSONFileWithDefaults(t *testing.T) {
	jsonContent := `[
  {
    "summary": "Task 1",
    "description": "Description 1"
  }
]`

	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(jsonContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	tickets, err := ParseJSONFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseJSONFile() error = %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	if tickets[0].IssueType != "Task" {
		t.Errorf("Default IssueType = %s, want Task", tickets[0].IssueType)
	}

	if tickets[0].Priority != "Medium" {
		t.Errorf("Default Priority = %s, want Medium", tickets[0].Priority)
	}
}

func TestParseJSONFileWithArraysAndLabels(t *testing.T) {
	jsonContent := `[
  {
    "summary": "Task 1",
    "labels": ["label1", "label2"],
    "components": ["comp1", "comp2"],
    "blocked_by": ["PROJ-100", "PROJ-101"]
  }
]`

	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(jsonContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	tickets, err := ParseJSONFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseJSONFile() error = %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	if len(tickets[0].Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(tickets[0].Labels))
	}

	if len(tickets[0].BlockedBy) != 2 {
		t.Errorf("Expected 2 blocked_by items, got %d", len(tickets[0].BlockedBy))
	}
}
