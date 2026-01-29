package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// TransitionOptions holds the options for the transition command
type TransitionOptions struct {
	Key    string
	Status string
}

// ExecuteTransitionCommand executes the transition command
func ExecuteTransitionCommand(v *viper.Viper, opts TransitionOptions) error {
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

	// Get available transitions
	transitions, err := issueService.GetTransitions(opts.Key)
	if err != nil {
		cli.PrintError(err)
		return err
	}

	// Find matching transition
	var transitionID string
	for _, t := range transitions {
		if strings.EqualFold(t.To.Name, opts.Status) {
			transitionID = t.ID
			break
		}
	}

	if transitionID == "" {
		// List available transitions
		fmt.Printf("❌ Status '%s' not found\n\n", opts.Status)
		fmt.Println("Available statuses:")
		for _, t := range transitions {
			fmt.Printf("  • %s\n", t.To.Name)
		}
		return fmt.Errorf("transition not found: %s", opts.Status)
	}

	// Perform transition
	if err := issueService.TransitionIssue(opts.Key, transitionID); err != nil {
		cli.PrintError(err)
		return err
	}

	// Print success message
	fmt.Printf("✅ Ticket transitioned successfully: %s -> %s\n", opts.Key, opts.Status)

	return nil
}

// NewTransitionCommand creates the "transition" command with full implementation
func NewTransitionCommand() *cobra.Command {
	var opts TransitionOptions

	cmd := &cobra.Command{
		Use:   "transition KEY",
		Short: "Transition a ticket to a new status",
		Long:  "Transition a JIRA ticket to a new workflow status.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Key = args[0]

			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Status, _ = cmd.Flags().GetString("to")

			if opts.Status == "" {
				return fmt.Errorf("--to flag is required")
			}

			return ExecuteTransitionCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Status, "to", "", "Target status (required)")

	return cmd
}
