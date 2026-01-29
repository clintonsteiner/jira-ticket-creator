package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// UpdateOptions holds the options for the update command
type UpdateOptions struct {
	Key         string
	Summary     string
	Description string
	Priority    string
	Assignee    string
	Labels      []string
}

// ExecuteUpdateCommand executes the update command
func ExecuteUpdateCommand(v *viper.Viper, opts UpdateOptions) error {
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

	// Build update fields - only include non-empty values
	fields := jira.IssueFields{}

	if opts.Summary != "" {
		fields.Summary = opts.Summary
	}

	if opts.Description != "" {
		fields.Description = opts.Description
	}

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

	// Update the issue
	if err := issueService.UpdateIssue(opts.Key, fields); err != nil {
		cli.PrintError(err)
		return err
	}

	// Print success message
	fmt.Printf("✅ Ticket updated successfully: %s\n", opts.Key)

	return nil
}

// NewUpdateCommand creates the "update" command with full implementation
func NewUpdateCommand() *cobra.Command {
	var opts UpdateOptions

	cmd := &cobra.Command{
		Use:   "update KEY",
		Short: "Update a JIRA ticket",
		Long:  "Update an existing JIRA ticket with new values.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]

			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Summary, _ = cmd.Flags().GetString("summary")
			opts.Description, _ = cmd.Flags().GetString("description")
			opts.Priority, _ = cmd.Flags().GetString("priority")
			opts.Assignee, _ = cmd.Flags().GetString("assignee")
			opts.Labels, _ = cmd.Flags().GetStringSlice("labels")

			return ExecuteUpdateCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Summary, "summary", "", "New ticket summary")
	cmd.Flags().StringVar(&opts.Description, "description", "", "New ticket description")
	cmd.Flags().StringVar(&opts.Priority, "priority", "", "New priority level (Lowest, Low, Medium, High, Highest)")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "New assignee email (e.g., user@company.com)")
	cmd.Flags().StringSliceVar(&opts.Labels, "labels", []string{}, "New labels (comma-separated, e.g., --labels bug,review)")

	return cmd
}
