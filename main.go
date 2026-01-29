package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
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
	flag.StringVar(&jiraURL, "url", "", "Jira Cloud URL (https://yourdomain.atlassian.net)")
	flag.StringVar(&email, "email", "", "Jira account email")
	flag.StringVar(&apiToken, "token", "", "Jira API token")
	flag.StringVar(&projectKey, "project", "", "Jira project key")
	flag.StringVar(&reportFile, "report", "created_tickets.json", "File to log created tickets")
}

type TicketRecord struct {
	Key       string    `json:"key"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	BlockedBy []string  `json:"blocked_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Load report from JSON
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

// Save report to JSON
func saveReport(filename string, records []TicketRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(records)
}

// Validate ticket exists and get status
func ValidateTicket(key string) (string, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/issue/%s", jiraURL, key), nil)
	req.SetBasicAuth(email, apiToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ticket %s not found (status %d)", key, resp.StatusCode)
	}

	var data struct {
		Fields struct {
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Fields.Status.Name, nil
}

// ------------------- Main CLI -------------------
func main() {
	flag.Parse()

	if jiraURL == "" || email == "" || apiToken == "" || projectKey == "" {
		fmt.Println("Provide all flags: -url, -email, -token, -project")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: create_ticket <summary> <description> [blocked_by_keys]")
		fmt.Println("   OR create_ticket report")
		fmt.Println("   OR create_ticket update-status")
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "report":
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
			fmt.Printf("- %s | %s | %s | blocked by: %v\n", r.Key, r.Summary, r.Status, r.BlockedBy)
		}

	case "update-status":
		records, err := loadReport(reportFile)
		if err != nil {
			fmt.Println("Error loading report:", err)
			os.Exit(1)
		}
		if len(records) == 0 {
			fmt.Println("No tickets to update.")
			return
		}

		fmt.Println("Updating ticket statuses...")
		for i, r := range records {
			status, err := ValidateTicket(r.Key)
			if err != nil {
				fmt.Printf("Failed to get status for %s: %v\n", r.Key, err)
				continue
			}
			records[i].Status = status
			fmt.Printf("%s -> %s\n", r.Key, status)
		}

		if err := saveReport(reportFile, records); err != nil {
			fmt.Println("Error saving updated report:", err)
			os.Exit(1)
		}
		fmt.Println("All statuses updated successfully.")

	default:
		if len(args) < 2 {
			fmt.Println("Usage: create_ticket <summary> <description> [blocked_by_keys]")
			os.Exit(1)
		}
		summary := args[0]
		description := args[1]
		blockedBy := []string{}
		if len(args) > 2 {
			blockedBy = strings.Split(args[2], ",")
		}

		newIssueKey, err := CreateIssue(jiraURL, email, apiToken, projectKey, summary, description, "Task")
		if err != nil {
			fmt.Println("Error creating issue:", err)
			os.Exit(1)
		}
		fmt.Println("Created issue:", newIssueKey)

		for _, blocker := range blockedBy {
			err := LinkIssue(jiraURL, email, apiToken, blocker, newIssueKey)
			if err != nil {
				fmt.Printf("Failed to link %s -> %s: %v\n", blocker, newIssueKey, err)
			} else {
				fmt.Printf("Linked %s -> %s\n", blocker, newIssueKey)
			}
		}

		status, err := ValidateTicket(newIssueKey)
		if err != nil {
			fmt.Println("Validation failed:", err)
			status = "Unknown"
		}

		existingRecords, _ := loadReport(reportFile)
		newRecord := TicketRecord{
			Key:       newIssueKey,
			Summary:   summary,
			Status:    status,
			BlockedBy: blockedBy,
			CreatedAt: time.Now(),
		}
		existingRecords = append(existingRecords, newRecord)
		_ = saveReport(reportFile, existingRecords)
		fmt.Printf("All created tickets logged to %s\n", reportFile)
	}
}
