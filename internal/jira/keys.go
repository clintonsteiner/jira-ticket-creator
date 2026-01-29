package jira

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractProjectKey extracts the project key from a ticket key
// Examples: "PROJ-123" -> "PROJ", "ABC-1" -> "ABC"
// Returns the project key and error if invalid format
func ExtractProjectKey(ticketKey string) (string, error) {
	ticketKey = strings.TrimSpace(ticketKey)

	// JIRA ticket keys follow pattern: [PROJECT]-[NUMBER]
	// Project key is uppercase letters/numbers, ticket number is digits
	re := regexp.MustCompile(`^([A-Z]+[A-Z0-9]*)-(\d+)$`)
	matches := re.FindStringSubmatch(ticketKey)

	if len(matches) < 2 {
		return "", fmt.Errorf("invalid ticket key format: '%s' (expected format: PROJECT-123)", ticketKey)
	}

	return matches[1], nil
}

// IsTicketKey checks if a string is a valid ticket key format
func IsTicketKey(key string) bool {
	re := regexp.MustCompile(`^[A-Z]+[A-Z0-9]*-\d+$`)
	return re.MatchString(strings.TrimSpace(key))
}

// IsProjectKey checks if a string is a valid project key format
// Project keys are typically uppercase letters and numbers, no hyphens
func IsProjectKey(key string) bool {
	re := regexp.MustCompile(`^[A-Z]+[A-Z0-9]*$`)
	return re.MatchString(strings.TrimSpace(key))
}

// ResolveProject returns the project key from either a project key or ticket key
// If input is "PROJ", returns "PROJ"
// If input is "PROJ-123", returns "PROJ" (extracted from ticket)
// If input is empty, returns error
func ResolveProject(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", fmt.Errorf("project key or ticket key is required")
	}

	// Check if it's a ticket key (contains hyphen and numbers)
	if IsTicketKey(input) {
		return ExtractProjectKey(input)
	}

	// Check if it's a project key
	if IsProjectKey(input) {
		return input, nil
	}

	return "", fmt.Errorf("invalid project or ticket key format: '%s' (expected PROJECT or PROJECT-123)", input)
}

// ExtractTicketNumber extracts the numeric part of a ticket key
// Examples: "PROJ-123" -> 123
func ExtractTicketNumber(ticketKey string) (string, error) {
	ticketKey = strings.TrimSpace(ticketKey)

	re := regexp.MustCompile(`^([A-Z]+[A-Z0-9]*)-(\d+)$`)
	matches := re.FindStringSubmatch(ticketKey)

	if len(matches) < 3 {
		return "", fmt.Errorf("invalid ticket key format: '%s'", ticketKey)
	}

	return matches[2], nil
}
