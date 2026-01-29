package config

// Defaults contains default configuration values
type Defaults struct {
	IssueType string
	Priority  string
}

// DefaultConfig returns the default configuration values
func DefaultConfig() Defaults {
	return Defaults{
		IssueType: "Task",
		Priority:  "Medium",
	}
}
