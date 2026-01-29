package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/batch"
	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// BatchCreateOptions holds options for batch create command
type BatchCreateOptions struct {
	InputFile string
	Format    string
	DryRun    bool
	Verbose   bool
}

// ExecuteBatchCreateCommand executes batch create
func ExecuteBatchCreateCommand(v *viper.Viper, opts BatchCreateOptions) error {
	// Load configuration with flag overrides
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate required configuration
	if err := cfg.ValidateRequired(); err != nil {
		return err
	}

	// Parse input file
	var tickets []batch.TicketData

	if strings.ToLower(opts.Format) == "json" {
		tickets, err = batch.ParseJSONFile(opts.InputFile)
	} else {
		tickets, err = batch.ParseCSVFile(opts.InputFile)
	}

	if err != nil {
		cli.PrintError(fmt.Errorf("failed to parse input file: %w", err))
		return err
	}

	fmt.Printf("📋 Loaded %d ticket(s) from %s\n", len(tickets), opts.InputFile)

	// Create JIRA client and services
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	validator := jira.NewValidator(client)

	// Phase 1: Validation
	fmt.Println("\n🔍 Phase 1: Validation")
	fmt.Println("---------------------")

	processor := batch.NewBatchProcessor(client, cfg.JIRA.Project)
	validationResults := processor.ValidateTickets(tickets, validator)

	validFailures := 0
	for _, result := range validationResults {
		if result.Error != nil {
			validFailures++
			fmt.Printf("❌ Ticket %d: %v\n", result.Index+1, result.Error)
		} else if opts.Verbose {
			fmt.Printf("✓ Ticket %d: %s (valid)\n", result.Index+1, result.TicketData.Summary)
		}
	}

	if validFailures > 0 {
		fmt.Printf("\n❌ %d validation error(s) found. Aborting.\n", validFailures)
		return fmt.Errorf("validation failed")
	}

	fmt.Printf("✅ All %d ticket(s) validated\n", len(tickets))

	// If dry-run, stop here
	if opts.DryRun {
		fmt.Println("\n✨ Dry-run complete - no tickets were created")
		return nil
	}

	// Phase 2: Creation
	fmt.Println("\n📝 Phase 2: Creation")
	fmt.Println("-------------------")

	createResults := processor.CreateTickets(tickets)

	createdCount := 0
	for _, result := range createResults {
		if result.Error == nil {
			createdCount++
			fmt.Printf("✅ [%d] %s -> %s\n", result.Index+1, result.TicketData.Summary, result.CreatedKey)
		} else {
			fmt.Printf("❌ [%d] %s: %v\n", result.Index+1, result.TicketData.Summary, result.Error)
		}
	}

	fmt.Printf("\n✅ Created %d ticket(s)\n", createdCount)

	// Phase 3: Linking
	if createdCount > 0 {
		fmt.Println("\n🔗 Phase 3: Linking")
		fmt.Println("------------------")

		linkResults := processor.LinkTickets(createResults)

		if len(linkResults) > 0 {
			for _, result := range linkResults {
				fmt.Printf("⚠️  %s: %v\n", result.CreatedKey, result.Error)
			}
			fmt.Printf("\n⚠️  %d linking error(s) (tickets were created)\n", len(linkResults))
		} else {
			fmt.Println("✅ All links created successfully")
		}
	}

	// Print summary
	fmt.Println("\n📊 Summary")
	fmt.Println("==========")
	fmt.Printf("Input:   %s\n", opts.InputFile)
	fmt.Printf("Created: %d/%d tickets\n", createdCount, len(tickets))

	return nil
}

// NewBatchCommand creates the "batch" command with full implementation
func NewBatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Batch operations on JIRA tickets",
		Long:  "Perform batch operations like creating multiple tickets from CSV/JSON.",
	}

	var opts BatchCreateOptions

	batchCreate := &cobra.Command{
		Use:   "create",
		Short: "Create tickets from a CSV or JSON file",
		Long: `Create multiple tickets from a CSV or JSON input file.

CSV format (with headers):
  summary,description,issue_type,priority,assignee,labels,components,blocked_by
  "Ticket 1","Description",Task,High,"user@email.com","label1,label2","comp1",

JSON format (array of objects):
  [
    {
      "summary": "Ticket 1",
      "description": "Description",
      "issue_type": "Task",
      "priority": "High",
      "assignee": "user@email.com",
      "labels": ["label1", "label2"],
      "components": ["comp1"],
      "blocked_by": ["PROJ-1"]
    }
  ]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.InputFile, _ = cmd.Flags().GetString("input")
			opts.Format, _ = cmd.Flags().GetString("format")
			opts.DryRun, _ = cmd.Flags().GetBool("dry-run")
			opts.Verbose, _ = cmd.Flags().GetBool("verbose")

			return ExecuteBatchCreateCommand(viper.GetViper(), opts)
		},
	}

	batchCreate.Flags().StringVar(&opts.InputFile, "input", "", "Input file (CSV or JSON) (required)")
	batchCreate.Flags().StringVar(&opts.Format, "format", "csv", "Input format (csv or json)")
	batchCreate.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Validate without creating")
	batchCreate.Flags().BoolVar(&opts.Verbose, "verbose", false, "Verbose output")
	batchCreate.MarkFlagRequired("input")

	cmd.AddCommand(batchCreate)

	return cmd
}
