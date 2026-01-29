package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// ReportOptions holds options for the report command
type ReportOptions struct {
	Format string
	Output string
}

// ExecuteReportCommand executes the report command
func ExecuteReportCommand(v *viper.Viper, opts ReportOptions) error {
	// Load configuration
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate required configuration
	if err := cfg.ValidateRequired(); err != nil {
		return err
	}

	// Load ticket records from storage
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

	if len(records) == 0 {
		fmt.Println("No tickets found")
		return nil
	}

	// Get ticket details from JIRA
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	issueService := jira.NewIssueService(client)

	var issues []jira.Issue
	for _, record := range records {
		issue, err := issueService.GetIssue(record.Key)
		if err != nil {
			fmt.Printf("⚠️  Failed to fetch details for %s: %v\n", record.Key, err)
			continue
		}
		issues = append(issues, *issue)
	}

	// Generate report
	reporter := reports.NewReporter(opts.Format)
	report, err := reporter.Generate(issues)
	if err != nil {
		cli.PrintError(fmt.Errorf("failed to generate report: %w", err))
		return err
	}

	// Output report
	if opts.Output != "" {
		if err := os.WriteFile(opts.Output, []byte(report), 0644); err != nil {
			cli.PrintError(fmt.Errorf("failed to write report: %w", err))
			return err
		}
		fmt.Printf("✅ Report written to: %s\n", opts.Output)
	} else {
		fmt.Println(report)
	}

	return nil
}

// NewReportCommand creates the "report" command with full implementation
func NewReportCommand() *cobra.Command {
	var opts ReportOptions

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a report of created tickets",
		Long:  "Generate a report of all created tickets in various formats (table, json, csv, markdown, html).",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Format, _ = cmd.Flags().GetString("format")
			opts.Output, _ = cmd.Flags().GetString("output")

			return ExecuteReportCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json, csv, markdown, html (default: table)")
	cmd.Flags().StringVar(&opts.Output, "output", "", "Output file path (optional, default: print to stdout)")

	return cmd
}
