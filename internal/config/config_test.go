package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "loads defaults when no config provided",
			envVars: map[string]string{
				"JIRA_URL":     "",
				"JIRA_EMAIL":   "",
				"JIRA_TOKEN":   "",
				"JIRA_PROJECT": "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, val := range tt.envVars {
				os.Setenv(key, val)
				defer os.Unsetenv(key)
			}

			cfg, err := LoadConfig()

			if (err != nil) != tt.expectError {
				t.Errorf("LoadConfig() error = %v, expectError %v", err, tt.expectError)
			}

			if !tt.expectError && cfg == nil {
				t.Error("LoadConfig() cfg is nil")
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				JIRA: struct {
					URL     string
					Email   string
					Token   string
					Project string
				}{
					URL:     "https://example.atlassian.net",
					Email:   "user@example.com",
					Token:   "token123",
					Project: "PROJ",
				},
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			cfg: &Config{
				JIRA: struct {
					URL     string
					Email   string
					Token   string
					Project string
				}{
					URL:     "",
					Email:   "user@example.com",
					Token:   "token123",
					Project: "PROJ",
				},
			},
			wantErr: true,
		},
		{
			name: "missing email",
			cfg: &Config{
				JIRA: struct {
					URL     string
					Email   string
					Token   string
					Project string
				}{
					URL:     "https://example.atlassian.net",
					Email:   "",
					Token:   "token123",
					Project: "PROJ",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.ValidateRequired()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	defaults := DefaultConfig()

	if defaults.IssueType != "Task" {
		t.Errorf("DefaultConfig().IssueType = %s, want Task", defaults.IssueType)
	}

	if defaults.Priority != "Medium" {
		t.Errorf("DefaultConfig().Priority = %s, want Medium", defaults.Priority)
	}
}
