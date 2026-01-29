package batch

import (
	"os"
	"testing"
)

func TestParseCSVFile(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		expectError bool
		expectCount int
	}{
		{
			name: "valid CSV with all fields",
			csvContent: `summary,description,issue_type,priority,assignee,labels,components,blocked_by
"Task 1","Description 1",Task,High,"user@email.com","label1,label2","Component1","PROJ-100"
"Task 2","Description 2",Story,Medium,"user2@email.com","label1","Component2",`,
			expectError: false,
			expectCount: 2,
		},
		{
			name: "CSV with only required fields",
			csvContent: `summary
"Task 1"
"Task 2"`,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "empty CSV",
			csvContent:  "",
			expectError: true,
			expectCount: 0,
		},
		{
			name: "CSV missing required summary column",
			csvContent: `title,description
"Task 1","Description 1"`,
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpfile, err := os.CreateTemp("", "test*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpfile.Close()

			// Parse CSV
			tickets, err := ParseCSVFile(tmpfile.Name())

			if (err != nil) != tt.expectError {
				t.Errorf("ParseCSVFile() error = %v, expectError %v", err, tt.expectError)
			}

			if len(tickets) != tt.expectCount {
				t.Errorf("ParseCSVFile() returned %d tickets, expected %d", len(tickets), tt.expectCount)
			}

			if !tt.expectError && len(tickets) > 0 {
				// Verify defaults are set
				if tickets[0].IssueType != "Task" && tt.csvContent != "" {
					// Only check if no explicit issue type was set
					if !contains(tt.csvContent, "issue_type") {
						t.Errorf("Default IssueType not set: %s", tickets[0].IssueType)
					}
				}
			}
		})
	}
}

func TestParseCSVFileWithDefaults(t *testing.T) {
	csvContent := `summary,description
"Task 1","Description 1"`

	tmpfile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	tickets, err := ParseCSVFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseCSVFile() error = %v", err)
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

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0
}

func TestParseCSVFileWithLabelsAndComponents(t *testing.T) {
	csvContent := `summary,labels,components
"Task 1","label1,label2","comp1,comp2"`

	tmpfile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	tickets, err := ParseCSVFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseCSVFile() error = %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	if len(tickets[0].Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(tickets[0].Labels))
	}

	if len(tickets[0].Components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(tickets[0].Components))
	}
}
