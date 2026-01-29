package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// TeamReport generates team-based reports with creator, timeline, and assignee info
type TeamReport struct{}

// GenerateTeamSummary generates a summary of tickets by creator
func (tr *TeamReport) GenerateTeamSummary(records []jira.TicketRecord) string {
	return tr.GenerateTeamSummaryWithFilter(records, "")
}

// GenerateTeamSummaryWithFilter generates a summary of tickets by creator, optionally filtered by project
func (tr *TeamReport) GenerateTeamSummaryWithFilter(records []jira.TicketRecord, projectFilter string) string {
	if len(records) == 0 {
		return "No tickets found"
	}

	// Filter by project if specified
	if projectFilter != "" {
		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if r.Project == projectFilter {
				filtered = append(filtered, r)
			}
		}
		records = filtered
		if len(records) == 0 {
			return fmt.Sprintf("No tickets found for project: %s", projectFilter)
		}
	}

	// Group by creator
	byCreator := make(map[string][]jira.TicketRecord)
	for _, record := range records {
		byCreator[record.Creator] = append(byCreator[record.Creator], record)
	}

	var sb strings.Builder
	sb.WriteString("📊 TEAM TICKET SUMMARY\n")
	sb.WriteString("======================\n\n")

	// Summary stats
	sb.WriteString(fmt.Sprintf("Total Tickets Created: %d\n", len(records)))
	sb.WriteString(fmt.Sprintf("Number of Team Members: %d\n", len(byCreator)))
	if projectFilter != "" {
		sb.WriteString(fmt.Sprintf("Project Filter: %s\n", projectFilter))
	}
	sb.WriteString("\n")

	// By creator
	sb.WriteString("Tickets by Creator:\n")
	sb.WriteString("-------------------\n")
	for creator, tickets := range byCreator {
		completed := 0
		pending := 0
		overdue := 0

		for _, t := range tickets {
			if t.Status == "Done" || t.Status == "Closed" {
				completed++
			} else if t.EstimatedEndDate != nil && t.EstimatedEndDate.Before(time.Now()) {
				overdue++
			} else {
				pending++
			}
		}

		sb.WriteString(fmt.Sprintf("\n👤 %s\n", creator))
		sb.WriteString(fmt.Sprintf("   Total: %d | Completed: %d | Pending: %d | Overdue: %d\n",
			len(tickets), completed, pending, overdue))

		// Group by project within creator
		byProject := make(map[string][]jira.TicketRecord)
		noProject := make([]jira.TicketRecord, 0)
		for _, t := range tickets {
			if t.Project != "" {
				byProject[t.Project] = append(byProject[t.Project], t)
			} else {
				noProject = append(noProject, t)
			}
		}

		// List tickets by project
		for project, projectTickets := range byProject {
			sb.WriteString(fmt.Sprintf("   [%s] %d ticket(s)\n", project, len(projectTickets)))
			for _, t := range projectTickets {
				status := "⏳"
				if t.Status == "Done" || t.Status == "Closed" {
					status = "✅"
				} else if t.EstimatedEndDate != nil && t.EstimatedEndDate.Before(time.Now()) {
					status = "⚠️ "
				}

				sb.WriteString(fmt.Sprintf("      %s %s [%s] → %s\n", status, t.Key, t.Priority, t.Summary))

				if t.Assignee != "" {
					sb.WriteString(fmt.Sprintf("         Assigned to: %s\n", t.Assignee))
				}
			}
		}

		// List unassigned to projects
		if len(noProject) > 0 {
			sb.WriteString(fmt.Sprintf("   [unassigned] %d ticket(s)\n", len(noProject)))
			for _, t := range noProject {
				status := "⏳"
				if t.Status == "Done" || t.Status == "Closed" {
					status = "✅"
				}

				sb.WriteString(fmt.Sprintf("      %s %s [%s] → %s\n", status, t.Key, t.Priority, t.Summary))
			}
		}
	}

	return sb.String()
}

// GenerateAssignmentMap generates a map of who is assigned to what
func (tr *TeamReport) GenerateAssignmentMap(records []jira.TicketRecord) string {
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

	var sb strings.Builder
	sb.WriteString("👥 WORKLOAD ASSIGNMENT MAP\n")
	sb.WriteString("==========================\n\n")

	// Assigned work
	sb.WriteString("Assigned Work:\n")
	sb.WriteString("---------------\n")

	for assignee, tickets := range byAssignee {
		completed := 0
		pending := 0
		critical := 0

		for _, t := range tickets {
			if t.Status == "Done" || t.Status == "Closed" {
				completed++
			} else if t.Priority == "Critical" || t.Priority == "High" {
				critical++
			} else {
				pending++
			}
		}

		sb.WriteString(fmt.Sprintf("\n👤 %s\n", assignee))
		sb.WriteString(fmt.Sprintf("   Total Work: %d | Completed: %d | Pending: %d | Critical: %d\n",
			len(tickets), completed, pending, critical))

		// List critical items
		for _, t := range tickets {
			if (t.Priority == "Critical" || t.Priority == "High") &&
				(t.Status != "Done" && t.Status != "Closed") {
				sb.WriteString(fmt.Sprintf("   🚨 %s [%s] - %s\n", t.Key, t.Priority, t.Summary))
			}
		}
	}

	// Unassigned work
	if len(unassigned) > 0 {
		sb.WriteString("\n\n⚠️  Unassigned Work:\n")
		sb.WriteString("-------------------\n")
		for _, t := range unassigned {
			sb.WriteString(fmt.Sprintf("   %s [%s] - %s\n", t.Key, t.Priority, t.Summary))
		}
	}

	return sb.String()
}

// GenerateTimeline generates a timeline view of work completion
func (tr *TeamReport) GenerateTimeline(records []jira.TicketRecord) string {
	if len(records) == 0 {
		return "No tickets found"
	}

	var sb strings.Builder
	sb.WriteString("📅 PROJECT TIMELINE\n")
	sb.WriteString("===================\n\n")

	// Calculate metrics
	total := len(records)
	completed := 0
	dueWithin7Days := 0
	overdue := 0
	noDueDate := 0

	for _, t := range records {
		if t.Status == "Done" || t.Status == "Closed" {
			completed++
		} else if t.EstimatedEndDate == nil {
			noDueDate++
		} else {
			daysLeft := t.EstimatedEndDate.Sub(time.Now())
			if daysLeft.Hours()/24 <= 7 && daysLeft.Hours()/24 > 0 {
				dueWithin7Days++
			} else if daysLeft < 0 {
				overdue++
			}
		}
	}

	// Progress bar
	progress := float64(completed) / float64(total) * 100
	filledBars := int(progress / 10)
	sb.WriteString(fmt.Sprintf("Overall Progress: ["))
	for i := 0; i < 10; i++ {
		if i < filledBars {
			sb.WriteString("█")
		} else {
			sb.WriteString("░")
		}
	}
	sb.WriteString(fmt.Sprintf("] %.0f%%\n\n", progress))

	// Status breakdown
	sb.WriteString("Status Breakdown:\n")
	sb.WriteString(fmt.Sprintf("  ✅ Completed: %d/%d\n", completed, total))
	sb.WriteString(fmt.Sprintf("  ⏳ In Progress: %d\n", total-completed-overdue))
	sb.WriteString(fmt.Sprintf("  🚨 Overdue: %d\n", overdue))
	sb.WriteString(fmt.Sprintf("  ⏰ Due Within 7 Days: %d\n", dueWithin7Days))
	sb.WriteString(fmt.Sprintf("  ❓ No Due Date: %d\n", noDueDate))

	// Critical issues
	critical := 0
	for _, t := range records {
		if (t.Priority == "Critical" || t.Priority == "Highest") &&
			(t.Status != "Done" && t.Status != "Closed") {
			critical++
		}
	}

	if critical > 0 {
		sb.WriteString(fmt.Sprintf("\n⚠️  Critical Items to Address: %d\n", critical))
	}

	// Recommended actions
	sb.WriteString("\n📋 Recommended Actions:\n")
	if overdue > 0 {
		sb.WriteString(fmt.Sprintf("  1. Address %d overdue ticket(s)\n", overdue))
	}
	if dueWithin7Days > 0 {
		sb.WriteString(fmt.Sprintf("  2. Accelerate work on %d ticket(s) due within 7 days\n", dueWithin7Days))
	}
	if noDueDate > 0 {
		sb.WriteString(fmt.Sprintf("  3. Set due dates for %d unscheduled ticket(s)\n", noDueDate))
	}

	return sb.String()
}
