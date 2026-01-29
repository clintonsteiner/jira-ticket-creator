---
layout: default
title: Go Client Library
parent: API Reference
---

# Go Client Library

The JIRA Ticket Creator provides a Go client library for programmatic access to JIRA operations.

## Installation

```bash
go get github.com/clintonsteiner/jira-ticket-creator
```

## Basic Usage

### Initialize Client

```go
package main

import (
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(
		"https://company.atlassian.net",
		"user@company.com",
		"api-token",
	)

	service := jira.NewIssueService(client)
	// Now use the service
}
```

### Create an Issue

```go
package main

import (
	"fmt"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(
		"https://company.atlassian.net",
		"user@company.com",
		"api-token",
	)

	service := jira.NewIssueService(client)

	resp, err := service.CreateIssue(
		"PROJ",
		"My Ticket Summary",
		"Detailed description",
		"Task",
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created ticket: %s\n", resp.Key)
}
```

### Create Issue with Full Fields

```go
fields := jira.IssueFields{
	Project: jira.Project{Key: "PROJ"},
	Summary: "Implement OAuth 2.0",
	Description: "Add OAuth authentication to API",
	IssueType: jira.IssueType{Name: "Story"},
	Priority: &jira.Priority{Name: "High"},
	Assignee: &jira.User{EmailAddress: "john@company.com"},
	Labels: []string{"auth", "security"},
	Components: []jira.Component{
		{Name: "API"},
		{Name: "Security"},
	},
}

resp, err := service.CreateIssueWithFields(fields)
if err != nil {
	panic(err)
}

fmt.Printf("Created: %s\n", resp.Key)
```

### Get an Issue

```go
issue, err := service.GetIssue("PROJ-123")
if err != nil {
	panic(err)
}

fmt.Printf("Summary: %s\n", issue.Fields.Summary)
fmt.Printf("Status: %s\n", issue.Fields.IssueType.Name)
fmt.Printf("Assignee: %s\n", issue.Fields.Assignee.Name)
```

### Search Issues

```go
issues, err := service.SearchIssues("project = PROJ AND status = 'To Do'", 0, 50)
if err != nil {
	panic(err)
}

for _, issue := range issues {
	fmt.Printf("%s: %s\n", issue.Key, issue.Fields.Summary)
}
```

### Update an Issue

```go
fields := jira.IssueFields{
	Summary: "Updated summary",
	Priority: &jira.Priority{Name: "High"},
}

err := service.UpdateIssue("PROJ-123", fields)
if err != nil {
	panic(err)
}
```

### Transition an Issue

```go
// First get available transitions
transitions, err := service.GetTransitions("PROJ-123")
if err != nil {
	panic(err)
}

// Find the transition ID for "In Progress"
for _, t := range transitions {
	if t.Name == "In Progress" {
		err = service.TransitionIssue("PROJ-123", t.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Transitioned to In Progress")
		break
	}
}
```

## Working with Storage

### Load Ticket Records

```go
package main

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		panic(err)
	}

	// Load all records
	records, err := repo.GetAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		fmt.Printf("%s: %s (%s)\n",
			record.Key,
			record.Summary,
			record.Status,
		)
	}
}
```

### Save Tickets

```go
records := []jira.TicketRecord{
	{
		Key: "PROJ-123",
		Summary: "My ticket",
		Status: "To Do",
		Priority: "Medium",
		IssueType: "Task",
	},
}

err := repo.Save(records)
if err != nil {
	panic(err)
}
```

### Add a Ticket

```go
record := jira.TicketRecord{
	Key: "PROJ-456",
	Summary: "New ticket",
	Status: "To Do",
	Priority: "High",
	IssueType: "Bug",
}

err := repo.Add(record)
if err != nil {
	panic(err)
}
```

## Working with Reports

### Generate Team Report

```go
package main

import (
	"fmt"
	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	// Load records (see storage example above)
	records, _ := repo.GetAll()

	// Generate team report
	teamReport := &reports.TeamReport{}
	output := teamReport.GenerateTeamSummary(records)

	fmt.Println(output)
}
```

### Generate Team Report with Filter

```go
output := teamReport.GenerateTeamSummaryWithFilter(records, "backend")
fmt.Println(output)
```

### Generate Gantt Chart

```go
gantt := reports.NewGanttChart()

// ASCII format
ascii := gantt.GenerateASCIIGantt(records, 2)
fmt.Println(ascii)

// Mermaid format
mermaid := gantt.GenerateMermaidGantt(records)
fmt.Println(mermaid)

// HTML format
html := gantt.GenerateHTMLGantt(records)
// Write to file or serve
```

## Error Handling

### Retry Logic

The client has built-in retry logic for transient failures:

```go
client := jira.NewClient(url, email, token)
// Automatically retries on network errors

// If error occurs after retries:
issue, err := service.GetIssue("PROJ-123")
if err != nil {
	// Handle error
	fmt.Printf("Failed to get issue: %v\n", err)
}
```

### Common Errors

```go
if err != nil {
	// Type assertions for specific errors
	if errors.Is(err, context.Canceled) {
		// Request was canceled
	} else if errors.Is(err, context.DeadlineExceeded) {
		// Request timeout
	} else {
		// Generic error
		fmt.Printf("Error: %v\n", err)
	}
}
```

## Configuration

### Load from Config File

```go
package main

import (
	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	service := jira.NewIssueService(client)
	// Use service
}
```

### Load with Environment Overrides

```go
import "github.com/spf13/viper"

func main() {
	v := viper.New()
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		panic(err)
	}
	// Use cfg
}
```

## Example: Bulk Operations

```go
package main

import (
	"fmt"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func createMultipleTickets(service *jira.IssueService, summaries []string) {
	for _, summary := range summaries {
		resp, err := service.CreateIssue("PROJ", summary, "", "Task")
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", summary, err)
			continue
		}
		fmt.Printf("Created: %s\n", resp.Key)
	}
}

func updateMultipleTickets(service *jira.IssueService, keys []string, priority string) {
	fields := jira.IssueFields{
		Priority: &jira.Priority{Name: priority},
	}

	for _, key := range keys {
		err := service.UpdateIssue(key, fields)
		if err != nil {
			fmt.Printf("Failed to update %s: %v\n", key, err)
			continue
		}
		fmt.Printf("Updated: %s\n", key)
	}
}

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Create multiple tickets
	createMultipleTickets(service, []string{
		"Task 1",
		"Task 2",
		"Task 3",
	})

	// Update multiple tickets
	updateMultipleTickets(service, []string{
		"PROJ-100",
		"PROJ-101",
		"PROJ-102",
	}, "High")
}
```

## Concurrency

For large operations, use goroutines:

```go
package main

import (
	"sync"
	"fmt"
)

func processTicketsConc(service *jira.IssueService, keys []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrency to 5

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			issue, err := service.GetIssue(k)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Processed: %s\n", issue.Key)
		}(key)
	}

	wg.Wait()
}
```

## See Also

- [Configuration API](./config.md)
- [Storage API](./storage.md)
- [Report Generators](./reports.md)
