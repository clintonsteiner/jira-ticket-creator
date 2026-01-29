package jira

import "testing"

func TestValidatePriority(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		priority  string
		wantError bool
	}{
		{"Valid: Lowest", "Lowest", false},
		{"Valid: Low", "Low", false},
		{"Valid: Medium", "Medium", false},
		{"Valid: High", "High", false},
		{"Valid: Highest", "Highest", false},
		{"Case insensitive: MEDIUM", "MEDIUM", false},
		{"Case insensitive: medium", "medium", false},
		{"Invalid: Critical", "Critical", true},
		{"Invalid: Urgent", "Urgent", true},
		{"Invalid: empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePriority(tt.priority)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePriority(%s) error = %v, wantError %v", tt.priority, err, tt.wantError)
			}
		})
	}
}

func TestValidateIssueType(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		issueType string
		wantError bool
	}{
		{"Valid: Task", "Task", false},
		{"Valid: Story", "Story", false},
		{"Valid: Bug", "Bug", false},
		{"Valid: Epic", "Epic", false},
		{"Valid: Subtask", "Subtask", false},
		{"Case insensitive: TASK", "TASK", false},
		{"Case insensitive: story", "story", false},
		{"Invalid: Feature", "Feature", true},
		{"Invalid: Improvement", "Improvement", true},
		{"Invalid: empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateIssueType(tt.issueType)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateIssueType(%s) error = %v, wantError %v", tt.issueType, err, tt.wantError)
			}
		})
	}
}

func TestValidateSummary(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		summary   string
		wantError bool
	}{
		{"Valid: short summary", "Add new feature", false},
		{"Valid: 255 chars", make255String(), false},
		{"Invalid: empty", "", true},
		{"Invalid: too long", make256String(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSummary(tt.summary)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateSummary() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateLabel(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		label     string
		wantError bool
	}{
		{"Valid: simple label", "feature", false},
		{"Valid: label with dash", "feature-new", false},
		{"Valid: label with underscore", "feature_new", false},
		{"Invalid: empty label", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateLabel(tt.label)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLabel(%s) error = %v, wantError %v", tt.label, err, tt.wantError)
			}
		})
	}
}

func TestValidateLabels(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		labels    []string
		wantError bool
	}{
		{"Valid: single label", []string{"feature"}, false},
		{"Valid: multiple labels", []string{"feature", "urgent", "ui"}, false},
		{"Invalid: contains empty", []string{"feature", "", "ui"}, true},
		{"Valid: empty list", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateLabels(tt.labels)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateLabels() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateCreateIssueRequest(t *testing.T) {
	validator := &Validator{client: nil}

	tests := []struct {
		name      string
		fields    IssueFields
		wantError bool
	}{
		{
			name: "Valid request",
			fields: IssueFields{
				Project:   Project{Key: "PROJ"},
				Summary:   "Valid ticket",
				IssueType: IssueType{Name: "Task"},
			},
			wantError: false,
		},
		{
			name: "Missing summary",
			fields: IssueFields{
				Project:   Project{Key: "PROJ"},
				Summary:   "",
				IssueType: IssueType{Name: "Task"},
			},
			wantError: true,
		},
		{
			name: "Invalid issue type",
			fields: IssueFields{
				Project:   Project{Key: "PROJ"},
				Summary:   "Valid ticket",
				IssueType: IssueType{Name: "InvalidType"},
			},
			wantError: true,
		},
		{
			name: "Missing project",
			fields: IssueFields{
				Project:   Project{Key: ""},
				Summary:   "Valid ticket",
				IssueType: IssueType{Name: "Task"},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateIssueRequest(tt.fields)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCreateIssueRequest() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func make255String() string {
	s := ""
	for i := 0; i < 255; i++ {
		s += "a"
	}
	return s
}

func make256String() string {
	s := ""
	for i := 0; i < 256; i++ {
		s += "a"
	}
	return s
}
