package reports

import (
	"fmt"
	"strings"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// PMReport handles project management reporting with parent-child relationships
type PMReport struct {
	tickets map[string]*jira.TicketRecord
}

// NewPMReport creates a new PM report
func NewPMReport() *PMReport {
	return &PMReport{
		tickets: make(map[string]*jira.TicketRecord),
	}
}

// GenerateProjectHierarchy creates a hierarchical view of tickets
func (r *PMReport) GenerateProjectHierarchy(records []jira.TicketRecord) string {
	// Build ticket map for quick lookup
	ticketMap := make(map[string]*jira.TicketRecord)
	for i := range records {
		ticketMap[records[i].Key] = &records[i]
	}

	output := strings.Builder{}
	output.WriteString("\n PROJECT MANAGEMENT DASHBOARD\n")
	output.WriteString("================================\n\n")

	// Group tickets by status
	byStatus := make(map[string][]jira.TicketRecord)
	for _, ticket := range records {
		byStatus[ticket.Status] = append(byStatus[ticket.Status], ticket)
	}

	// Calculate statistics
	total := len(records)
	completed := len(byStatus["Done"])
	inProgress := len(byStatus["In Progress"])
	todo := len(byStatus["To Do"])
	blocked := 0

	for _, ticket := range records {
		if len(ticket.BlockedBy) > 0 {
			blocked++
		}
	}

	// Overall Progress
	percentComplete := 0
	if total > 0 {
		percentComplete = (completed * 100) / total
	}

	output.WriteString(fmt.Sprintf("Overall Progress: %d%% (%d/%d tickets)\n", percentComplete, completed, total))
	output.WriteString(generateProgressBar(percentComplete))
	output.WriteString(fmt.Sprintf("\n  [DONE] %d\n", completed))
	output.WriteString(fmt.Sprintf("  [IN PROGRESS] %d\n", inProgress))
	output.WriteString(fmt.Sprintf("  [TO DO] %d\n", todo))
	output.WriteString(fmt.Sprintf("  [BLOCKED] %d\n\n", blocked))

	// Ticket hierarchy
	output.WriteString("TICKET HIERARCHY\n")
	output.WriteString("-------------------\n\n")

	// Group by issue type (treating Epics/Stories as parents)
	epics := make(map[string][]jira.TicketRecord)
	standalone := []jira.TicketRecord{}

	for _, ticket := range records {
		if ticket.IssueType == "Epic" || ticket.IssueType == "Story" {
			// Could be parent
			epics[ticket.Key] = []jira.TicketRecord{}
		}
	}

	// Assign children based on blocking relationships
	for _, ticket := range records {
		if len(ticket.BlockedBy) > 0 {
			// This ticket is blocked by others - could be a parent
			for _, blocker := range ticket.BlockedBy {
				if parentTicket, ok := ticketMap[blocker]; ok {
					if parentTicket.IssueType == "Epic" || parentTicket.IssueType == "Story" {
						// Already a parent
						continue
					}
				}
			}
		} else if ticket.IssueType != "Epic" && ticket.IssueType != "Story" {
			standalone = append(standalone, ticket)
		}
	}

	// Display epics and their children
	for _, ticket := range records {
		if ticket.IssueType == "Epic" {
			output.WriteString(fmt.Sprintf("[EPIC] EPIC: %s - %s\n", ticket.Key, ticket.Summary))
			output.WriteString(fmt.Sprintf("   Status: %s | Priority: %s | Owner: %s\n",
				ticket.Status, ticket.Priority, ticket.Assignee))

			// Find child tickets (those blocking this one or blocked by this)
			children := []jira.TicketRecord{}
			for _, candidate := range records {
				if candidate.Key == ticket.Key {
					continue
				}
				// Check if related through blocking
				for _, blocker := range candidate.BlockedBy {
					if blocker == ticket.Key {
						children = append(children, candidate)
						break
					}
				}
			}

			if len(children) > 0 {
				output.WriteString("   ├─ Subtasks:\n")
				for i, child := range children {
					statusIcon := "[TODO]"
					if child.Status == "Done" {
						statusIcon = ""
					} else if child.Status == "In Progress" {
						statusIcon = "[IN_PROGRESS]"
					}
					isLast := i == len(children)-1
					prefix := "   │  ├─ "
					if isLast {
						prefix = "   │  └─ "
					}
					output.WriteString(fmt.Sprintf("%s%s %s (%s) - Assignee: %s\n",
						prefix, statusIcon, child.Key, child.Summary, child.Assignee))
				}
			} else {
				output.WriteString("   └─ No subtasks\n")
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}

// GeneratePMDashboard creates an executive summary for project managers
func (r *PMReport) GeneratePMDashboard(records []jira.TicketRecord) string {
	output := strings.Builder{}

	// Collect statistics
	byStatus := make(map[string]int)
	byPriority := make(map[string]int)
	byAssignee := make(map[string]int)
	blockedCount := 0
	criticalPath := []jira.TicketRecord{}

	for _, ticket := range records {
		byStatus[ticket.Status]++
		byPriority[ticket.Priority]++
		byAssignee[ticket.Assignee]++

		if len(ticket.BlockedBy) > 0 {
			blockedCount++
			if ticket.Priority == "Critical" || ticket.Priority == "High" {
				criticalPath = append(criticalPath, ticket)
			}
		}
	}

	output.WriteString("\n╔════════════════════════════════════════════════════════════════════╗\n")
	output.WriteString("║              PROJECT MANAGEMENT EXECUTIVE SUMMARY                  ║\n")
	output.WriteString("╚════════════════════════════════════════════════════════════════════╝\n\n")

	total := len(records)
	output.WriteString(fmt.Sprintf(" Total Tickets: %d\n\n", total))

	// Status breakdown
	output.WriteString("Status Breakdown:\n")
	output.WriteString("─────────────────\n")
	for status, count := range byStatus {
		pct := (count * 100) / total
		output.WriteString(fmt.Sprintf("  %s: %d (%d%%)\n", status, count, pct))
	}

	// Priority breakdown
	output.WriteString("\nPriority Distribution:\n")
	output.WriteString("──────────────────────\n")
	priorities := []string{"Critical", "High", "Medium", "Low", "Lowest"}
	for _, priority := range priorities {
		if count, ok := byPriority[priority]; ok && count > 0 {
			icon := "[CRITICAL]"
			if priority == "High" {
				icon = "[HIGH]"
			} else if priority == "Medium" {
				icon = "[MEDIUM]"
			} else if priority == "Low" {
				icon = "[LOW]"
			}
			output.WriteString(fmt.Sprintf("  %s %s: %d\n", icon, priority, count))
		}
	}

	// Team workload
	output.WriteString("\nTeam Workload Distribution:\n")
	output.WriteString("──────────────────────────\n")
	for assignee, count := range byAssignee {
		if assignee == "" {
			assignee = "[Unassigned]"
		}
		output.WriteString(fmt.Sprintf("  %s: %d tickets\n", assignee, count))
	}

	// Blockers and risks
	output.WriteString("\n  Critical Items:\n")
	output.WriteString("──────────────────\n")
	if blockedCount > 0 {
		output.WriteString(fmt.Sprintf("  Blocked tickets: %d (dependencies exist)\n", blockedCount))
	}
	if len(criticalPath) > 0 {
		output.WriteString(fmt.Sprintf("  High-priority blocked items: %d\n", len(criticalPath)))
		for _, ticket := range criticalPath {
			output.WriteString(fmt.Sprintf("    • %s: %s (blocked by: %v)\n",
				ticket.Key, ticket.Summary, ticket.BlockedBy))
		}
	}

	return output.String()
}

// GenerateRiskReport identifies project risks
func (r *PMReport) GenerateRiskReport(records []jira.TicketRecord) string {
	output := strings.Builder{}
	output.WriteString("\n PROJECT RISK ASSESSMENT\n")
	output.WriteString("==========================\n\n")

	risks := []string{}

	// Check for critical blocked items
	criticalBlocked := 0
	for _, ticket := range records {
		if (ticket.Priority == "Critical" || ticket.Priority == "High") && len(ticket.BlockedBy) > 0 {
			criticalBlocked++
		}
	}
	if criticalBlocked > 0 {
		risks = append(risks, fmt.Sprintf("HIGH RISK: %d critical/high priority items are blocked", criticalBlocked))
	}

	// Check for unassigned tickets
	unassigned := 0
	for _, ticket := range records {
		if ticket.Assignee == "" && ticket.Status != "Done" {
			unassigned++
		}
	}
	if unassigned > 0 {
		risks = append(risks, fmt.Sprintf("MEDIUM RISK: %d tickets are unassigned", unassigned))
	}

	// Check for stalled items
	inProgress := 0
	for _, ticket := range records {
		if ticket.Status == "In Progress" {
			inProgress++
		}
	}
	if inProgress > 5 {
		risks = append(risks, fmt.Sprintf("MEDIUM RISK: %d items in progress (possible WIP bloat)", inProgress))
	}

	if len(risks) == 0 {
		output.WriteString(" No major risks identified\n")
	} else {
		for i, risk := range risks {
			output.WriteString(fmt.Sprintf("%d. %s\n", i+1, risk))
		}
	}

	output.WriteString("\n Recommendations:\n")
	if criticalBlocked > 0 {
		output.WriteString("  • Review and unblock critical items immediately\n")
		output.WriteString("  • Consider re-prioritizing dependencies\n")
	}
	if unassigned > 0 {
		output.WriteString("  • Assign unassigned tickets to team members\n")
	}
	if inProgress > 5 {
		output.WriteString("  • Focus team on completing in-progress items before starting new work\n")
	}

	return output.String()
}

// generateProgressBar creates a visual progress bar
func generateProgressBar(percent int) string {
	width := 30
	filled := (percent * width) / 100
	empty := width - filled

	bar := "  ["
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := 0; i < empty; i++ {
		bar += " "
	}
	bar += fmt.Sprintf("] %d%%", percent)
	return bar
}

// GenerateTicketDetailsTable creates a detailed table of all tickets
func (r *PMReport) GenerateTicketDetailsTable(records []jira.TicketRecord) string {
	output := strings.Builder{}
	output.WriteString("\n DETAILED TICKET INVENTORY\n")
	output.WriteString("============================\n\n")

	data := [][]string{}
	for _, ticket := range records {
		status := ticket.Status
		assignee := ticket.Assignee
		if assignee == "" {
			assignee = "[Unassigned]"
		}
		blockedBy := strings.Join(ticket.BlockedBy, ", ")
		if blockedBy == "" {
			blockedBy = "—"
		}
		data = append(data, []string{
			ticket.Key,
			ticket.Summary,
			status,
			ticket.Priority,
			assignee,
			blockedBy,
		})
	}

	// Create markdown table
	output.WriteString("| Key | Summary | Status | Priority | Assignee | Blocked By |\n")
	output.WriteString("|-----|---------|--------|----------|----------|------------|\n")
	for _, row := range data {
		// Truncate long summaries for table
		summary := row[1]
		if len(summary) > 30 {
			summary = summary[:27] + "..."
		}
		output.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			row[0], summary, row[2], row[3], row[4], row[5]))
	}

	return output.String()
}
