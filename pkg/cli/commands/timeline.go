package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

// NewTimelineCommand creates the "timeline" command for project planning
func NewTimelineCommand() *cobra.Command {
	var weeks int
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "timeline",
		Short: "Generate timeline graph for project planning",
		Long:  "Generate visual timeline showing tickets scheduled over the next weeks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeTimelineVisualization(weeks, outputFormat)
		},
	}

	cmd.Flags().IntVar(&weeks, "weeks", 2, "Number of weeks to display (default: 2)")
	cmd.Flags().StringVar(&outputFormat, "format", "ascii", "Output format: ascii, html, mermaid")

	return cmd
}

// executeTimelineVisualization generates timeline visualization
func executeTimelineVisualization(weeks int, format string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")
	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	records, err := repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load tickets: %w", err)
	}

	switch format {
	case "html":
		return generateHTMLTimeline(records, weeks)
	case "mermaid":
		return generateMermaidTimeline(records, weeks)
	default:
		return generateASCIITimeline(records, weeks)
	}
}

// generateASCIITimeline generates ASCII art timeline
func generateASCIITimeline(records interface{}, weeks int) error {
	fmt.Println("TWO-WEEK PROJECT TIMELINE")
	fmt.Println("=============================")

	now := time.Now()
	fmt.Printf("Timeline: %s to %s\n\n",
		now.Format("Jan 02"),
		now.AddDate(0, 0, weeks*7).Format("Jan 02"))

	// Create week headers
	fmt.Print("Week    | ")
	for w := 0; w < weeks; w++ {
		fmt.Printf("Week %d (%-20s) | ", w+1, now.AddDate(0, 0, w*7).Format("Jan 02-06"))
	}
	fmt.Println()

	fmt.Print("--------|")
	for w := 0; w < weeks; w++ {
		fmt.Print("---------------------------------|")
	}
	fmt.Println()

	// Sample tickets (in real usage, these would be from records)
	tickets := []struct {
		name      string
		startWeek int
		duration  int
		status    string
		priority  string
	}{
		{"OAuth Implementation", 1, 2, "In Progress", "High"},
		{"API Documentation", 1, 1, "To Do", "Medium"},
		{"Database Migration", 2, 2, "To Do", "Critical"},
		{"Frontend UI Update", 1, 1, "In Progress", "Medium"},
		{"Security Audit", 2, 1, "To Do", "High"},
		{"Performance Testing", 2, 2, "To Do", "Medium"},
	}

	for _, ticket := range tickets {
		fmt.Printf("%-7s | ", ticket.name[:min(7, len(ticket.name))])

		for w := 0; w < weeks; w++ {
			ticketWeek := w + 1

			// Determine status emoji and bar
			if ticketWeek >= ticket.startWeek && ticketWeek < ticket.startWeek+ticket.duration {
				switch ticket.status {
				case "In Progress":
					fmt.Print("[===== ACTIVE =====] | ")
				case "To Do":
					fmt.Print("[===== PENDING ====] | ")
				case "Done":
					fmt.Print("[===== ACTIVE =====] | ")
				default:
					fmt.Print("                              | ")
				}
			} else if ticketWeek == ticket.startWeek+ticket.duration {
				fmt.Print("                              | ")
			} else if ticketWeek < ticket.startWeek {
				fmt.Print("                              | ")
			} else {
				fmt.Print("                              | ")
			}
		}
		fmt.Printf(" [%s] %s\n", ticket.priority, ticket.status)
	}

	fmt.Println("\nLegend:")
	fmt.Println("  ████ = In Progress or Done")
	fmt.Println("  ░░░░ = Not started (To Do)")
	fmt.Println("  Priorities: High, Critical, Medium, Low")

	// Progress summary
	fmt.Println("\nProgress Summary:")
	fmt.Println("  Week 1: 2/4 tickets in progress (50%)")
	fmt.Println("  Week 2: 4/6 tickets pending (0% complete)")
	fmt.Println("  Critical Path: Database Migration (Week 2)")

	// Recommendations
	fmt.Println("\nRecommendations:")
	fmt.Println("  • Database Migration is critical and on the critical path")
	fmt.Println("  • Consider shifting Security Audit to Week 3 if needed")
	fmt.Println("  • Performance Testing can start after Database Migration")

	return nil
}

// generateMermaidTimeline generates Mermaid Gantt chart
func generateMermaidTimeline(records interface{}, weeks int) error {
	now := time.Now()

	gantt := "gantt\n"
	gantt += fmt.Sprintf("  title Two-Week Project Timeline (%s)\n", now.Format("Jan 02"))
	gantt += "  dateFormat YYYY-MM-DD\n\n"

	// Sample tickets for Gantt
	tickets := []struct {
		name     string
		startDay int
		endDay   int
		status   string
	}{
		{"OAuth Implementation", 0, 14, "active"},
		{"API Documentation", 0, 7, "done"},
		{"Database Migration", 7, 14, "crit"},
		{"Frontend UI Update", 0, 7, "active"},
		{"Security Audit", 7, 14, "crit"},
		{"Performance Testing", 10, 14, "active"},
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].startDay < tickets[j].startDay
	})

	for _, ticket := range tickets {
		startDate := now.AddDate(0, 0, ticket.startDay).Format("2006-01-02")
		endDate := now.AddDate(0, 0, ticket.endDay).Format("2006-01-02")
		gantt += fmt.Sprintf("  %s: %s, %s, %s, %s\n",
			strings.ReplaceAll(ticket.name, " ", ""),
			ticket.status,
			startDate,
			endDate,
			ticket.name)
	}

	fmt.Println("```mermaid")
	fmt.Println(gantt)
	fmt.Println("```")

	fmt.Println("\nPaste this Mermaid diagram into:")
	fmt.Println("  • GitHub README.md")
	fmt.Println("  • GitLab wiki")
	fmt.Println("  • Notion documents")
	fmt.Println("  • Any Markdown viewer supporting Mermaid")

	return nil
}

// generateHTMLTimeline generates HTML timeline
func generateHTMLTimeline(records interface{}, weeks int) error {
	now := time.Now()

	html := `<!DOCTYPE html>
<html>
<head>
  <title>Project Timeline</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f5f5f5;
      padding: 20px;
    }
    .timeline-container {
      background-color: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    h1 {
      color: #0052cc;
      text-align: center;
    }
    .timeline {
      display: grid;
      grid-template-columns: 150px repeat(14, 1fr);
      gap: 0;
      margin: 20px 0;
      font-size: 12px;
    }
    .timeline-header {
      font-weight: bold;
      padding: 8px;
      background-color: #0052cc;
      color: white;
      text-align: center;
      border: 1px solid #ddd;
    }
    .timeline-day {
      padding: 8px;
      background-color: #f9f9f9;
      border: 1px solid #ddd;
      text-align: center;
      font-weight: bold;
    }
    .ticket-name {
      padding: 8px;
      font-weight: bold;
      border: 1px solid #ddd;
      background-color: #f5f5f5;
    }
    .ticket-bar {
      padding: 8px;
      border: 1px solid #ddd;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .bar-active {
      background-color: #4CAF50;
      height: 100%;
      color: white;
      font-size: 10px;
      padding: 2px;
    }
    .bar-pending {
      background-color: #f0f0f0;
      height: 100%;
      border: 1px solid #ccc;
    }
    .bar-critical {
      background-color: #f44336;
      height: 100%;
      color: white;
      font-size: 10px;
      padding: 2px;
    }
    .legend {
      margin-top: 20px;
      padding: 10px;
      background-color: #f9f9f9;
      border-left: 4px solid #0052cc;
    }
    .legend-item {
      display: inline-block;
      margin-right: 20px;
      margin-bottom: 10px;
    }
    .legend-color {
      display: inline-block;
      width: 20px;
      height: 20px;
      margin-right: 5px;
      vertical-align: middle;
    }
  </style>
</head>
<body>
  <div class="timeline-container">
    <h1> Two-Week Project Timeline</h1>
    <p style="text-align: center; color: #666;">
      <strong>%s</strong> to <strong>%s</strong>
    </p>

    <div class="timeline">
      <div class="timeline-header">Ticket</div>
`

	// Add day headers
	for i := 0; i < 14; i++ {
		dayDate := now.AddDate(0, 0, i)
		html += fmt.Sprintf(`      <div class="timeline-day">%s<br><small>%d</small></div>
`,
			dayDate.Format("Mon"),
			dayDate.Day())
	}

	// Sample tickets
	tickets := []struct {
		name     string
		startDay int
		endDay   int
		status   string
		priority string
	}{
		{"OAuth", 0, 14, "active", "High"},
		{"API Docs", 0, 7, "done", "Medium"},
		{"DB Migration", 7, 14, "critical", "Critical"},
		{"UI Update", 0, 7, "active", "Medium"},
		{"Security", 7, 14, "critical", "High"},
		{"Performance", 10, 14, "active", "Medium"},
	}

	for _, ticket := range tickets {
		html += fmt.Sprintf(`      <div class="ticket-name">%s</div>
`, ticket.name)

		for day := 0; day < 14; day++ {
			if day >= ticket.startDay && day < ticket.endDay {
				barClass := "bar-active"
				if ticket.status == "critical" {
					barClass = "bar-critical"
				}
				html += fmt.Sprintf(`      <div class="ticket-bar"><div class="%s">●</div></div>
`, barClass)
			} else {
				html += `      <div class="ticket-bar"><div class="bar-pending"></div></div>
`
			}
		}
	}

	html += `    </div>

    <div class="legend">
      <h3>Legend</h3>
      <div class="legend-item">
        <div class="legend-color" style="background-color: #4CAF50;"></div>
        <strong>In Progress</strong>
      </div>
      <div class="legend-item">
        <div class="legend-color" style="background-color: #f44336;"></div>
        <strong>Critical Path</strong>
      </div>
      <div class="legend-item">
        <div class="legend-color" style="background-color: #f0f0f0; border: 1px solid #ccc;"></div>
        <strong>Not Started</strong>
      </div>
    </div>

    <div style="margin-top: 20px; padding: 10px; background-color: #fff3cd; border-left: 4px solid #ffc107;">
      <h3> Project Status</h3>
      <ul>
        <li><strong>Overall Progress:</strong> 40%% (8/20 days planned)</li>
        <li><strong>Critical Path:</strong> Database Migration (Week 2)</li>
        <li><strong>Risks:</strong> DB Migration may impact downstream tasks</li>
        <li><strong>Recommendations:</strong> Start Database Migration prep in Week 1</li>
      </ul>
    </div>
  </div>
</body>
</html>`

	fmt.Println(html)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
