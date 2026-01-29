package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	JIRA struct {
		URL     string
		Email   string
		Token   string
		Project string
	}
	Defaults Defaults
}

// LoadConfig loads configuration with the following priority:
// 1. Command-line flags (passed via viper)
// 2. Environment variables (JIRA_URL, JIRA_EMAIL, JIRA_TOKEN, JIRA_PROJECT)
// 3. Config file (~/.jirarc in YAML format)
// 4. Built-in defaults
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set environment variable binding
	v.BindEnv("jira.url", "JIRA_URL")
	v.BindEnv("jira.email", "JIRA_EMAIL")
	v.BindEnv("jira.token", "JIRA_TOKEN")
	v.BindEnv("jira.project", "JIRA_PROJECT")

	// Set config file paths
	v.SetConfigName(".jirarc")
	v.SetConfigType("yaml")

	// Check home directory for config file
	homeDir, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(homeDir)
	}

	// Try to read config file (non-fatal if not found)
	v.ReadInConfig()

	// Set defaults
	defaults := DefaultConfig()
	v.SetDefault("defaults.issue_type", defaults.IssueType)
	v.SetDefault("defaults.priority", defaults.Priority)

	// Parse configuration into struct
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Store viper instance for later use with command flags
	_viperInstance = v

	return cfg, nil
}

// LoadConfigWithFlags loads configuration and applies command-line flag overrides
// This should be called after flags are parsed
func LoadConfigWithFlags(v *viper.Viper) (*Config, error) {
	// Set environment variable binding
	v.BindEnv("jira.url", "JIRA_URL")
	v.BindEnv("jira.email", "JIRA_EMAIL")
	v.BindEnv("jira.token", "JIRA_TOKEN")
	v.BindEnv("jira.project", "JIRA_PROJECT")

	// Set config file paths
	v.SetConfigName(".jirarc")
	v.SetConfigType("yaml")

	// Check home directory for config file
	homeDir, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(homeDir)
	}

	// Try to read config file (non-fatal if not found)
	v.ReadInConfig()

	// Set defaults
	defaults := DefaultConfig()
	v.SetDefault("defaults.issue_type", defaults.IssueType)
	v.SetDefault("defaults.priority", defaults.Priority)

	// Parse configuration into struct
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return cfg, nil
}

// ValidateRequired checks that required fields are set
func (c *Config) ValidateRequired() error {
	if c.JIRA.URL == "" {
		return fmt.Errorf("JIRA URL is required (set via --url flag, JIRA_URL env var, or ~/.jirarc config file)")
	}
	if c.JIRA.Email == "" {
		return fmt.Errorf("JIRA email is required (set via --email flag, JIRA_EMAIL env var, or ~/.jirarc config file)")
	}
	if c.JIRA.Token == "" {
		return fmt.Errorf("JIRA token is required (set via --token flag, JIRA_TOKEN env var, or ~/.jirarc config file)")
	}
	if c.JIRA.Project == "" {
		return fmt.Errorf("JIRA project is required (set via --project flag, JIRA_PROJECT env var, or ~/.jirarc config file)")
	}
	return nil
}

// GetConfigPath returns the path to the config file, creating it if it doesn't exist
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, ".jirarc")
	return configPath, nil
}

// var to hold viper instance for later flag binding
var _viperInstance *viper.Viper

// GetViperInstance returns the viper instance
func GetViperInstance() *viper.Viper {
	return _viperInstance
}
