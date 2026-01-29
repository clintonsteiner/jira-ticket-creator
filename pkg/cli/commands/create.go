package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/interactive"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// CreateOptions holds the options for the create command
type CreateOptions struct {
	Summary     string
	Description string
	Type        string
	Priority    string
	Assignee    string
	Labels      []string
	Components  []string
	BlockedBy   []string
	Interactive bool
	Template    string
}

// ExecuteCreateCommand executes the create command
func ExecuteCreateCommand(v *viper.Viper, opts CreateOptions) error {
	// Handle interactive mode
	if opts.Interactive {
		wizard := &interactive.TicketWizard{}
		input, err := wizard.Run()
		if err != nil {
			cli.PrintError(err)
			return err
		}

		// Convert wizard input to create options
		opts.Summary = input.Summary
		opts.Description = input.Description
		opts.Type = input.Type
		opts.Priority = input.Priority
		opts.Assignee = input.Assignee
		opts.Labels = input.Labels
		opts.BlockedBy = input.BlockedBy
	}

	// Validate required fields
	if opts.Summary == "" {
		return fmt.Errorf("summary is required")
	}

	// Load configuration with flag overrides
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate required configuration
	if err := cfg.ValidateRequired(); err != nil {
		return err
	}

	// Create JIRA client and services
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	issueService := jira.NewIssueService(client)
	linkService := jira.NewLinkService(client)

	// Build issue fields
	fields := jira.IssueFields{
		Project: jira.Project{
			Key: cfg.JIRA.Project,
		},
		Summary:     opts.Summary,
		Description: opts.Description,
		IssueType: jira.IssueType{
			Name: opts.Type,
		},
	}

	// Add optional fields
	if opts.Priority != "" {
		fields.Priority = &jira.Priority{
			Name: opts.Priority,
		}
	}

	if opts.Assignee != "" {
		fields.Assignee = &jira.User{
			EmailAddress: opts.Assignee,
		}
	}

	if len(opts.Labels) > 0 {
		fields.Labels = opts.Labels
	}

	if len(opts.Components) > 0 {
		components := make([]jira.Component, len(opts.Components))
		for i, comp := range opts.Components {
			components[i] = jira.Component{Name: comp}
		}
		fields.Components = components
	}

	// Create the issue
	resp, err := issueService.CreateIssueWithFields(fields)
	if err != nil {
		cli.PrintError(err)
		return err
	}

	ticketKey := resp.Key

	// Link blocked-by issues
	if len(opts.BlockedBy) > 0 {
		for _, blocker := range opts.BlockedBy {
			if err := linkService.LinkBlocks(blocker, ticketKey); err != nil {
				fmt.Printf("⚠️  Warning: Failed to link %s blocks %s: %v\n", blocker, ticketKey, err)
			}
		}
	}

	// Save ticket record
	err = saveTicketRecord(cfg.JIRA.Project, ticketKey, opts.Summary)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to save ticket record: %v\n", err)
	}

	// Print success message
	fmt.Printf("✅ Ticket created successfully: %s\n", ticketKey)
	if cfg.JIRA.URL != "" {
		// Parse URL to get base
		fmt.Printf("🔗 View at: %s/browse/%s\n", cfg.JIRA.URL, ticketKey)
	}

	return nil
}

// saveTicketRecord saves the created ticket to the local record file
func saveTicketRecord(projectKey, ticketKey, summary string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")

	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		return err
	}

	record := jira.TicketRecord{
		Key:       ticketKey,
		Summary:   summary,
		Status:    "To Do",
		BlockedBy: []string{},
		CreatedAt: time.Now(),
	}

	return repo.Add(record)
}

// NewCreateCommand creates the "create" command with full implementation
func NewCreateCommand() *cobra.Command {
	var opts CreateOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new JIRA ticket",
		Long:  "Create a new JIRA ticket with the specified summary, description, and optional metadata.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Summary, _ = cmd.Flags().GetString("summary")
			opts.Description, _ = cmd.Flags().GetString("description")
			opts.Type, _ = cmd.Flags().GetString("type")
			opts.Priority, _ = cmd.Flags().GetString("priority")
			opts.Assignee, _ = cmd.Flags().GetString("assignee")
			opts.Labels, _ = cmd.Flags().GetStringSlice("labels")
			opts.Components, _ = cmd.Flags().GetStringSlice("component")
			opts.BlockedBy, _ = cmd.Flags().GetStringSlice("blocked-by")
			opts.Interactive, _ = cmd.Flags().GetBool("interactive")
			opts.Template, _ = cmd.Flags().GetString("template")

			return ExecuteCreateCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Summary, "summary", "", "Ticket summary (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Ticket description")
	cmd.Flags().StringVar(&opts.Type, "type", "Task", "Issue type (Task, Story, Bug, Epic, Subtask)")
	cmd.Flags().StringVar(&opts.Priority, "priority", "Medium", "Priority (Low, Medium, High, Critical)")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Assignee email address")
	cmd.Flags().StringSliceVar(&opts.Labels, "labels", []string{}, "Labels to add (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.Components, "component", []string{}, "Components (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.BlockedBy, "blocked-by", []string{}, "Comma-separated list of ticket keys that block this one")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Use interactive mode")
	cmd.Flags().StringVar(&opts.Template, "template", "", "Use a template for creation")

	cmd.MarkFlagRequired("summary")

	return cmd
}
