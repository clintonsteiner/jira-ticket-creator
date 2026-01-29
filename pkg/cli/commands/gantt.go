package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

// NewGanttCommand creates the "gantt" command for Gantt chart visualization
func NewGanttCommand() *cobra.Command {
	var outputFormat string
	var outputFile string
	var weeks int

	cmd := &cobra.Command{
		Use:   "gantt",
		Short: "Generate Gantt chart showing workload by resource",
		Long: `Generate Gantt chart visualization showing tickets scheduled by assigned resource.
Displays ticket status, timeline, and workload distribution across team members.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGanttCommand(outputFormat, outputFile, weeks)
		},
	}

	cmd.Flags().StringVar(&outputFormat, "format", "ascii", "Output format: ascii, mermaid, html")
	cmd.Flags().StringVar(&outputFile, "output", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&weeks, "weeks", 2, "Number of weeks to display (for ascii format)")

	return cmd
}

// executeGanttCommand executes the gantt command
func executeGanttCommand(format string, outputFile string, weeks int) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
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

	if len(records) == 0 {
		fmt.Println("No tickets found")
		return nil
	}

	// Generate Gantt chart
	gantt := reports.NewGanttChart()
	var output string

	switch format {
	case "mermaid":
		output = gantt.GenerateMermaidGantt(records)
		// Add markdown code block for easy embedding
		if outputFile == "" {
			output = "```mermaid\n" + output + "```"
		}

	case "html":
		output = gantt.GenerateHTMLGantt(records)

	case "ascii":
		fallthrough
	default:
		output = gantt.GenerateASCIIGantt(records, weeks)
	}

	// Output results
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("✅ Gantt chart written to %s\n", outputFile)
	} else {
		fmt.Println(output)
	}

	return nil
}
