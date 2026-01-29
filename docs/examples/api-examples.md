---
layout: default
title: API Examples
parent: Examples
nav_order: 2
---

# API Usage Examples

Programmatic examples using the Go client library.

## Creating Tickets Programmatically

### Simple Ticket Creation

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	// Create client
	client := jira.NewClient(
		"https://company.atlassian.net",
		"user@company.com",
		"api-token",
	)

	service := jira.NewIssueService(client)

	// Create a simple ticket
	resp, err := service.CreateIssue(
		"PROJ",
		"Build new feature",
		"Implement user authentication system",
		"Story",
	)
	if err != nil {
		log.Fatalf("Failed to create ticket: %v", err)
	}

	fmt.Printf("Created ticket: %s\n", resp.Key)
}
```

### Bulk Ticket Creation

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	tickets := []struct {
		Summary string
		Description string
		Type string
	}{
		{"Task 1", "First task", "Task"},
		{"Task 2", "Second task", "Task"},
		{"Bug 1", "Critical bug", "Bug"},
	}

	for _, t := range tickets {
		resp, err := service.CreateIssue("PROJ", t.Summary, t.Description, t.Type)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", t.Summary, err)
			continue
		}
		fmt.Printf(" Created: %s\n", resp.Key)
	}
}
```

### Create with Full Fields

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	fields := jira.IssueFields{
		Project: jira.Project{Key: "PROJ"},
		Summary: "Implement OAuth 2.0 authentication",
		Description: "Add OAuth 2.0 support to API with scope-based permissions",
		IssueType: jira.IssueType{Name: "Story"},
		Priority: &jira.Priority{Name: "High"},
		Assignee: &jira.User{EmailAddress: "john@company.com"},
		Labels: []string{"auth", "security", "backend"},
		Components: []jira.Component{
			{Name: "API"},
			{Name: "Security"},
		},
	}

	resp, err := service.CreateIssueWithFields(fields)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("Created with fields: %s\n", resp.Key)
}
```

## Searching and Querying

### Search with Pagination

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Search with pagination
	var allIssues []jira.Issue
	startAt := 0
	pageSize := 50

	for {
		issues, err := service.SearchIssues("project = PROJ", startAt, pageSize)
		if err != nil {
			log.Fatalf("Search failed: %v", err)
		}

		allIssues = append(allIssues, issues...)

		if len(issues) < pageSize {
			break
		}
		startAt += pageSize
	}

	fmt.Printf("Found %d total issues\n", len(allIssues))

	for _, issue := range allIssues {
		fmt.Printf("%s: %s\n", issue.Key, issue.Fields.Summary)
	}
}
```

### Complex Search Query

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Complex JQL
	jql := `project = PROJ
		AND type = Story
		AND status = "In Progress"
		AND priority >= High
		ORDER BY created DESC`

	issues, err := service.SearchIssues(jql, 0, 100)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("Found %d high-priority in-progress stories\n", len(issues))

	for _, issue := range issues {
		fmt.Printf(" %s: %s (%s)\n",
			issue.Key,
			issue.Fields.Summary,
			issue.Fields.Priority.Name,
		)
	}
}
```

## Updating Issues

### Update Single Issue

```go
package main

import (
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Update fields
	fields := jira.IssueFields{
		Summary: "Updated summary",
		Priority: &jira.Priority{Name: "Critical"},
	}

	err := service.UpdateIssue("PROJ-123", fields)
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}

	fmt.Println(" Updated PROJ-123")
}
```

### Bulk Update Issues

```go
package main

import (
	"fmt"
	"sync"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	keys := []string{"PROJ-100", "PROJ-101", "PROJ-102"}

	fields := jira.IssueFields{
		Priority: &jira.Priority{Name: "High"},
	}

	// Update in parallel
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Max 5 concurrent

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := service.UpdateIssue(k, fields)
			if err != nil {
				fmt.Printf(" Failed to update %s: %v\n", k, err)
			} else {
				fmt.Printf(" Updated %s\n", k)
			}
		}(key)
	}

	wg.Wait()
	fmt.Println("Done updating all tickets")
}
```

## Working with Local Storage

### Load and Process Tickets

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		log.Fatalf("Failed to initialize repo: %v", err)
	}

	// Load all records
	records, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Failed to load records: %v", err)
	}

	// Process by status
	statuses := make(map[string]int)
	for _, record := range records {
		statuses[record.Status]++
	}

	fmt.Println("Ticket count by status:")
	for status, count := range statuses {
		fmt.Printf(" %s: %d\n", status, count)
	}
}
```

### Save Imported Tickets

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	// Create records from issues
	records := []jira.TicketRecord{
		{
			Key: "PROJ-123",
			Summary: "Sample ticket",
			Status: "To Do",
			CreatedAt: time.Now(),
			Creator: "imported",
			Priority: "Medium",
			IssueType: "Task",
			Project: "backend",
		},
	}

	// Save them
	for _, record := range records {
		err := repo.Add(record)
		if err != nil {
			fmt.Printf("Failed to add %s: %v\n", record.Key, err)
			continue
		}
		fmt.Printf(" Saved %s\n", record.Key)
	}
}
```

## Generating Reports

### Team Report with Filter

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	records, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	// Generate team report filtered by project
	teamReport := &reports.TeamReport{}
	output := teamReport.GenerateTeamSummaryWithFilter(records, "backend")

	fmt.Println(output)
}
```

### Generate Gantt Chart

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	records, err := repo.GetAll()
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	// Generate different formats
	gantt := reports.NewGanttChart()

	// ASCII
	ascii := gantt.GenerateASCIIGantt(records, 2)
	fmt.Println(ascii)

	// Mermaid (save to file)
	mermaid := gantt.GenerateMermaidGantt(records)
	err = os.WriteFile("gantt.md", []byte("```mermaid\n"+mermaid+"\n```"), 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	// HTML (save to file)
	html := gantt.GenerateHTMLGantt(records)
	err = os.WriteFile("gantt.html", []byte(html), 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println(" Generated gantt.md and gantt.html")
}
```

## Import Workflow

### Import with Project Mapping

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

func main() {
	// 1. Load project mapping
	mapping, err := config.LoadMapping("")
	if err != nil {
		log.Printf("Warning: %v", err)
		mapping = &config.ProjectMapping{Mappings: make(map[string]config.ProjectInfo)}
	}

	// 2. Search for issues
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	issues, err := service.SearchIssues("project = PROJ", 0, 100)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// 3. Convert to TicketRecords
	records := make([]jira.TicketRecord, 0)
	for _, issue := range issues {
		project := mapping.FindProjectForKey(issue.Key)

		record := jira.TicketRecord{
			Key: issue.Key,
			Summary: issue.Fields.Summary,
			Status: "Open",
			CreatedAt: time.Now(),
			Creator: "imported",
			Priority: "Medium",
			IssueType: issue.Fields.IssueType.Name,
			Project: project,
		}

		if issue.Fields.Priority != nil {
			record.Priority = issue.Fields.Priority.Name
		}

		records = append(records, record)
	}

	// 4. Save to local storage
	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	for _, record := range records {
		err := repo.Add(record)
		if err != nil {
			fmt.Printf("Failed to add %s: %v\n", record.Key, err)
		}
	}

	fmt.Printf(" Imported %d tickets\n", len(records))
}
```

## Error Handling Best Practices

```go
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Try to get issue with error handling
	issue, err := service.GetIssue("PROJ-123")

	if err != nil {
		// Check error type
		if errors.Is(err, context.Canceled) {
			log.Println("Request was canceled")
		} else if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Request timeout")
		} else {
			log.Printf("Error: %v\n", err)

			// Specific error messages
			if err.Error() == "not found" {
				fmt.Println("Ticket not found")
			} else if err.Error() == "unauthorized" {
				fmt.Println("Authentication failed - check credentials")
			} else {
				fmt.Println("Unexpected error")
			}
		}
		return
	}

	fmt.Printf("Found: %s\n", issue.Key)
}
```

## Concurrent Operations

### Process Large Batches in Parallel

```go
package main

import (
	"fmt"
	"sync"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

func main() {
	client := jira.NewClient(url, email, token)
	service := jira.NewIssueService(client)

	// Simulate batch of keys
	keys := generateKeys(1000)

	// Process with controlled concurrency
	results := processConcurrent(service, keys, 10)

	// Print results
	success := 0
	failed := 0
	for result := range results {
		if result.err == nil {
			success++
		} else {
			failed++
		}
	}

	fmt.Printf("Processed: %d success, %d failed\n", success, failed)
}

type result struct {
	key string
	err error
}

func processConcurrent(service *jira.IssueService, keys []string, workers int) <-chan result {
	results := make(chan result, len(keys))
	semaphore := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			issue, err := service.GetIssue(k)
			results <- result{
				key: k,
				err: err,
			}

			if err == nil {
				fmt.Printf("Processed %s: %s\n", k, issue.Fields.Summary)
			}
		}(key)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func generateKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("PROJ-%d", i+1)
	}
	return keys
}
```

## See Also

- [Go Client Library Reference](../api/go-client.md)
- [Configuration API](../api/config.md)
- [API Examples](./api-examples.md)
