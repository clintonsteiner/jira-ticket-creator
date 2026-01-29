package reports

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// GanttChart generates Gantt chart visualizations
type GanttChart struct{}

// NewGanttChart creates a new Gantt chart generator
func NewGanttChart() *GanttChart {
	return &GanttChart{}
}

// GenerateASCIIGantt generates ASCII Gantt chart showing tickets by assignee
func (gc *GanttChart) GenerateASCIIGantt(records []jira.TicketRecord, weeks int) string {
	if len(records) == 0 {
		return "No tickets found"
	}

	// Group by assignee
	byAssignee := make(map[string][]jira.TicketRecord)
	unassigned := make([]jira.TicketRecord, 0)

	for _, record := range records {
		if record.Assignee != "" {
			byAssignee[record.Assignee] = append(byAssignee[record.Assignee], record)
		} else {
			unassigned = append(unassigned, record)
		}
	}

	// Sort assignees
	var assignees []string
	for assignee := range byAssignee {
		assignees = append(assignees, assignee)
	}
	sort.Strings(assignees)

	var sb strings.Builder
	sb.WriteString("📊 GANTT CHART - WORKLOAD BY RESOURCE\n")
	sb.WriteString("=====================================\n\n")

	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 0, weeks*7)

	sb.WriteString(fmt.Sprintf("Timeline: %s to %s (%d weeks)\n\n",
		startDate.Format("Jan 02, 2006"),
		endDate.Format("Jan 02, 2006"),
		weeks))

	// Generate week headers
	sb.WriteString("RESOURCE       | ")
	for w := 0; w < weeks; w++ {
		weekStart := startDate.AddDate(0, 0, w*7)
		sb.WriteString(fmt.Sprintf("Week %d (%s) | ",
			w+1,
			weekStart.Format("Jan 02")))
	}
	sb.WriteString("STATUS\n")

	sb.WriteString(strings.Repeat("-", 15))
	sb.WriteString("|")
	for w := 0; w < weeks; w++ {
		sb.WriteString(strings.Repeat("-", 25))
		sb.WriteString("|")
	}
	sb.WriteString(strings.Repeat("-", 15))
	sb.WriteString("\n")

	// For each assignee, show their tickets
	for _, assignee := range assignees {
		tickets := byAssignee[assignee]

		// Show assignee name
		displayName := truncate(assignee, 14)
		sb.WriteString(fmt.Sprintf("%-14s | ", displayName))

		// Create bar representation
		bar := gc.generateTicketBar(tickets, weeks, startDate)
		sb.WriteString(bar)

		// Count by status
		done := 0
		inProgress := 0
		todo := 0
		for _, t := range tickets {
			switch t.Status {
			case "Done", "Closed":
				done++
			case "In Progress":
				inProgress++
			default:
				todo++
			}
		}

		// Show status summary
		statusStr := fmt.Sprintf("%d✓ %d⟳ %d□",
			done, inProgress, todo)
		sb.WriteString(fmt.Sprintf("| %-13s\n", truncate(statusStr, 13)))
	}

	// Unassigned tickets
	if len(unassigned) > 0 {
		displayName := truncate("Unassigned", 14)
		sb.WriteString(fmt.Sprintf("%-14s | ", displayName))
		bar := gc.generateTicketBar(unassigned, weeks, startDate)
		sb.WriteString(bar)

		done := 0
		inProgress := 0
		todo := 0
		for _, t := range unassigned {
			switch t.Status {
			case "Done", "Closed":
				done++
			case "In Progress":
				inProgress++
			default:
				todo++
			}
		}

		statusStr := fmt.Sprintf("%d✓ %d⟳ %d□",
			done, inProgress, todo)
		sb.WriteString(fmt.Sprintf("| %-13s\n", truncate(statusStr, 13)))
	}

	return sb.String()
}

// generateTicketBar generates a visual bar for tickets in a timeline
func (gc *GanttChart) generateTicketBar(tickets []jira.TicketRecord, weeks int, startDate time.Time) string {
	var bar strings.Builder

	for w := 0; w < weeks; w++ {
		weekStart := startDate.AddDate(0, 0, w*7)
		weekEnd := weekStart.AddDate(0, 0, 7)

		// Count tickets in this week
		inWeek := 0
		done := 0
		progress := 0

		for _, t := range tickets {
			// Check if ticket created or due date falls in this week
			ticketInWeek := false

			// Created in this week
			if t.CreatedAt.After(weekStart) && t.CreatedAt.Before(weekEnd) {
				ticketInWeek = true
			}

			// Due date in this week
			if t.EstimatedEndDate != nil &&
				t.EstimatedEndDate.After(weekStart) &&
				t.EstimatedEndDate.Before(weekEnd) {
				ticketInWeek = true
			}

			// Or if created before and not done
			if t.CreatedAt.Before(weekEnd) && t.Status != "Done" && t.Status != "Closed" {
				if t.EstimatedEndDate == nil || t.EstimatedEndDate.After(weekStart) {
					ticketInWeek = true
				}
			}

			if ticketInWeek {
				inWeek++
				if t.Status == "Done" || t.Status == "Closed" {
					done++
				} else if t.Status == "In Progress" {
					progress++
				}
			}
		}

		// Generate visual representation
		if inWeek == 0 {
			bar.WriteString("          | ")
		} else if done > 0 && progress == 0 {
			bar.WriteString("✓✓✓✓✓✓✓✓✓ | ")
		} else if progress > 0 && done == 0 {
			bar.WriteString("⟳⟳⟳⟳⟳⟳⟳⟳⟳ | ")
		} else if progress > 0 && done > 0 {
			filled := int((float64(progress) / float64(inWeek)) * 9)
			visual := strings.Repeat("⟳", filled) + strings.Repeat("✓", 9-filled)
			bar.WriteString(fmt.Sprintf("%-9s | ", visual))
		} else {
			bar.WriteString("□□□□□□□□□ | ")
		}
	}

	return bar.String()
}

// GenerateMermaidGantt generates Mermaid Gantt chart
func (gc *GanttChart) GenerateMermaidGantt(records []jira.TicketRecord) string {
	if len(records) == 0 {
		return "gantt\n    title No tickets found"
	}

	// Group by assignee
	byAssignee := make(map[string][]jira.TicketRecord)
	unassigned := make([]jira.TicketRecord, 0)

	for _, record := range records {
		if record.Assignee != "" {
			byAssignee[record.Assignee] = append(byAssignee[record.Assignee], record)
		} else {
			unassigned = append(unassigned, record)
		}
	}

	// Sort assignees
	var assignees []string
	for assignee := range byAssignee {
		assignees = append(assignees, assignee)
	}
	sort.Strings(assignees)

	var sb strings.Builder
	sb.WriteString("gantt\n")
	sb.WriteString("    title Workload by Resource\n")
	sb.WriteString("    dateFormat YYYY-MM-DD\n\n")

	// Add tickets for each assignee
	for _, assignee := range assignees {
		tickets := byAssignee[assignee]
		sectionName := truncate(assignee, 20)
		sb.WriteString(fmt.Sprintf("    section %s\n", sectionName))

		for _, ticket := range tickets {
			startDate := ticket.CreatedAt.Format("2006-01-02")
			endDate := startDate
			status := "todo"

			if ticket.EstimatedEndDate != nil {
				endDate = ticket.EstimatedEndDate.Format("2006-01-02")
			} else {
				// Default to 7 days from creation
				endDate = ticket.CreatedAt.AddDate(0, 0, 7).Format("2006-01-02")
			}

			// Map status to Gantt status
			if ticket.Status == "Done" || ticket.Status == "Closed" {
				status = "done"
			} else if ticket.Status == "In Progress" {
				status = "active"
			}

			// Ensure valid date range
			if endDate < startDate {
				endDate = ticket.CreatedAt.AddDate(0, 0, 1).Format("2006-01-02")
			}

			sb.WriteString(fmt.Sprintf("    %s :%s, %s, %s\n",
				ticket.Key,
				status,
				startDate,
				endDate))
		}
	}

	// Add unassigned tickets
	if len(unassigned) > 0 {
		sb.WriteString("    section Unassigned\n")
		for _, ticket := range unassigned {
			startDate := ticket.CreatedAt.Format("2006-01-02")
			endDate := startDate
			status := "todo"

			if ticket.EstimatedEndDate != nil {
				endDate = ticket.EstimatedEndDate.Format("2006-01-02")
			} else {
				endDate = ticket.CreatedAt.AddDate(0, 0, 7).Format("2006-01-02")
			}

			if ticket.Status == "Done" || ticket.Status == "Closed" {
				status = "done"
			} else if ticket.Status == "In Progress" {
				status = "active"
			}

			if endDate < startDate {
				endDate = ticket.CreatedAt.AddDate(0, 0, 1).Format("2006-01-02")
			}

			sb.WriteString(fmt.Sprintf("    %s :%s, %s, %s\n",
				ticket.Key,
				status,
				startDate,
				endDate))
		}
	}

	return sb.String()
}

// GenerateHTMLGantt generates HTML Gantt chart
func (gc *GanttChart) GenerateHTMLGantt(records []jira.TicketRecord) string {
	if len(records) == 0 {
		return "<p>No tickets found</p>"
	}

	// Group by assignee
	byAssignee := make(map[string][]jira.TicketRecord)
	unassigned := make([]jira.TicketRecord, 0)

	for _, record := range records {
		if record.Assignee != "" {
			byAssignee[record.Assignee] = append(byAssignee[record.Assignee], record)
		} else {
			unassigned = append(unassigned, record)
		}
	}

	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Gantt Chart - Workload by Resource</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .resource { margin-bottom: 30px; }
        .resource-name { font-weight: bold; font-size: 16px; margin-bottom: 10px; color: #0066cc; }
        .tickets { margin-left: 20px; }
        .ticket { margin-bottom: 8px; padding: 8px; background: #f5f5f5; border-radius: 4px; border-left: 4px solid #0066cc; }
        .ticket-key { font-weight: bold; color: #0066cc; }
        .ticket-summary { color: #666; margin: 5px 0; }
        .status { display: inline-block; padding: 2px 8px; border-radius: 3px; font-size: 12px; margin: 5px 0; }
        .status.done { background: #d4edda; color: #155724; }
        .status.progress { background: #fff3cd; color: #856404; }
        .status.todo { background: #f8d7da; color: #721c24; }
        .priority { font-size: 12px; color: #999; }
        .stats { margin-top: 10px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <h1>📊 Gantt Chart - Workload by Resource</h1>
`)

	// Add tickets for each assignee
	for assignee, tickets := range byAssignee {
		sb.WriteString(`    <div class="resource">
`)
		sb.WriteString(fmt.Sprintf(`        <div class="resource-name">👤 %s (%d tickets)</div>
`, assignee, len(tickets)))
		sb.WriteString(`        <div class="tickets">
`)

		done := 0
		progress := 0
		todo := 0

		for _, ticket := range tickets {
			statusClass := "todo"
			statusText := "To Do"

			if ticket.Status == "Done" || ticket.Status == "Closed" {
				statusClass = "done"
				statusText = "✓ Done"
				done++
			} else if ticket.Status == "In Progress" {
				statusClass = "progress"
				statusText = "⟳ In Progress"
				progress++
			} else {
				todo++
			}

			endDateStr := "No due date"
			if ticket.EstimatedEndDate != nil {
				endDateStr = ticket.EstimatedEndDate.Format("Jan 02, 2006")
			}

			sb.WriteString(fmt.Sprintf(`            <div class="ticket">
                <div><span class="ticket-key">%s</span> <span class="ticket-summary">%s</span></div>
                <div><span class="status %s">%s</span> <span class="priority">[%s] Due: %s</span></div>
            </div>
`, ticket.Key, ticket.Summary, statusClass, statusText, ticket.Priority, endDateStr))
		}

		sb.WriteString(fmt.Sprintf(`            <div class="stats">Summary: %d✓ %d⟳ %d□</div>
`, done, progress, todo))
		sb.WriteString(`        </div>
    </div>
`)
	}

	// Add unassigned tickets
	if len(unassigned) > 0 {
		sb.WriteString(`    <div class="resource">
        <div class="resource-name">⚠️ Unassigned (%d tickets)</div>
        <div class="tickets">
`)

		done := 0
		progress := 0
		todo := 0

		for _, ticket := range unassigned {
			statusClass := "todo"
			statusText := "To Do"

			if ticket.Status == "Done" || ticket.Status == "Closed" {
				statusClass = "done"
				statusText = "✓ Done"
				done++
			} else if ticket.Status == "In Progress" {
				statusClass = "progress"
				statusText = "⟳ In Progress"
				progress++
			} else {
				todo++
			}

			endDateStr := "No due date"
			if ticket.EstimatedEndDate != nil {
				endDateStr = ticket.EstimatedEndDate.Format("Jan 02, 2006")
			}

			sb.WriteString(fmt.Sprintf(`            <div class="ticket">
                <div><span class="ticket-key">%s</span> <span class="ticket-summary">%s</span></div>
                <div><span class="status %s">%s</span> <span class="priority">[%s] Due: %s</span></div>
            </div>
`, ticket.Key, ticket.Summary, statusClass, statusText, ticket.Priority, endDateStr))
		}

		sb.WriteString(fmt.Sprintf(`            <div class="stats">Summary: %d✓ %d⟳ %d□</div>
`, done, progress, todo))
		sb.WriteString(`        </div>
    </div>
`)
	}

	sb.WriteString(`</body>
</html>
`)

	return sb.String()
}

// truncate truncates a string to maxLength
func truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
