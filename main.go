package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	jiraURL    string
	email      string
	apiToken   string
	projectKey string
	reportFile string
)

func init() {
	flag.StringVar(&jiraURL, "url", "", "Jira instance URL")
	flag.StringVar(&email, "email", "", "Jira account email")
	flag.StringVar(&apiToken, "token", "", "Jira API token")
	flag.StringVar(&projectKey, "project", "", "Jira project key")
	flag.StringVar(&reportFile, "report", "created_tickets.json", "File to log created tickets")
}

// TicketRecord is used to log created tickets
type TicketRecord struct {
	Key       string    `json:"key"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

func loadReport(filename string) ([]TicketRecord, error) {
	var records []TicketRecord
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []TicketRecord{}, nil
		}
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func saveReport(filename string, records []TicketRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(records)
}

func main() {
	flag.Parse()

	if jiraURL == "" || email == "" || apiToken == "" || projectKey == "" {
		fmt.Println("Please provide all required flags: -url, -email, -token, -project")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: create_ticket <summary> <description> [blocked_by_issue_keys(comma-separated)]")
		fmt.Println("Or: create_ticket report  -> to show logged tickets")
		os.Exit(1)
	}

	if args[0] == "report" {
		records, err := loadReport(reportFile)
		if err != nil {
			fmt.Println("Error loading report:", err)
			os.Exit(1)
		}
		if len(records) == 0 {
			fmt.Println("No tickets created yet.")
			return
		}
		fmt.Println("Created Jira tickets:")
		for _, r := range records {
			fmt.Printf("- %s | %s | %s\n", r.Key, r.Summary, r.CreatedAt.Format(time.RFC3339))
		}
		return
	}

	summary := args[0]
	description := args[1]
	blockedBy := []string{}
	if len(args) > 2 {
		blockedBy = strings.Split(args[2], ",")
	}

	createdTickets := []TicketRecord{}

	newIssueKey, err := CreateIssue(jiraURL, email, apiToken, projectKey, summary, description, "Task")
	if err != nil {
		fmt.Println("Error creating issue:", err)
		os.Exit(1)
	}
	fmt.Println("Created issue:", newIssueKey)

	createdTickets = append(createdTickets, TicketRecord{
		Key:       newIssueKey,
		Summary:   summary,
		CreatedAt: time.Now(),
	})

	for _, blocker := range blockedBy {
		err := LinkIssue(jiraURL, email, apiToken, blocker, newIssueKey)
		if err != nil {
			fmt.Printf("Failed to link %s -> %s: %v\n", blocker, newIssueKey, err)
		} else {
			fmt.Printf("Linked %s -> %s\n", blocker, newIssueKey)
		}
	}

	// Load existing report
	existingRecords, err := loadReport(reportFile)
	if err != nil {
		fmt.Println("Error loading existing report:", err)
		os.Exit(1)
	}

	// Append new tickets
	existingRecords = append(existingRecords, createdTickets...)

	// Save report
	if err := saveReport(reportFile, existingRecords); err != nil {
		fmt.Println("Error saving report:", err)
		os.Exit(1)
	}

	fmt.Printf("All created tickets logged to %s\n", reportFile)
}

