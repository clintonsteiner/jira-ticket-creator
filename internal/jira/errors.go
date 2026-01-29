package jira

import (
	"fmt"
)

// JiraError represents a JIRA API error
type JiraError struct {
	StatusCode    int
	Message       string
	Errors        map[string]interface{}
	ErrorMessages []string
}

// Error implements the error interface
func (e *JiraError) Error() string {
	if len(e.ErrorMessages) > 0 {
		return fmt.Sprintf("JIRA error (HTTP %d): %s", e.StatusCode, e.ErrorMessages[0])
	}
	if e.Message != "" {
		return fmt.Sprintf("JIRA error (HTTP %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("JIRA error (HTTP %d)", e.StatusCode)
}

// IsRetryable determines if a request should be retried based on the status code
func (e *JiraError) IsRetryable() bool {
	switch e.StatusCode {
	case 429: // Rate limit
		return true
	case 500, 502, 503, 504: // Server errors
		return true
	default:
		return false
	}
}

// ValidationError represents a validation error before API call
type ValidationError struct {
	Field   string
	Message string
	Details string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("validation error in %s: %s (%s)", e.Field, e.Message, e.Details)
	}
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Message string
}

// Error implements the error interface
func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	Key      string
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("%s not found: %s", e.Resource, e.Key)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	RetryAfter int
	Message    string
}

// Error implements the error interface
func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limited: please retry after %d seconds (%s)", e.RetryAfter, e.Message)
	}
	return fmt.Sprintf("rate limited: %s", e.Message)
}

// IsRetryable determines if this error should be retried
func (e *RateLimitError) IsRetryable() bool {
	return true
}
