package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// SearchOptions holds the options for the search command
type SearchOptions struct {
	Key     string
	Summary string
	JQL     string
	Format  string
}

// ExecuteSearchCommand executes the search command
func ExecuteSearchCommand(v *viper.Viper, opts SearchOptions) error {
	// Load configuration with flag overrides
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate required configuration
	if err := cfg.ValidateRequired(); err != nil {
		return err
	}

	// Create JIRA client
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	issueService := jira.NewIssueService(client)

	var issues []jira.Issue

	// Perform search based on provided options
	if opts.Key != "" {
		// Search by key
		issue, err := issueService.GetIssue(opts.Key)
		if err != nil {
			cli.PrintError(err)
			return err
		}
		issues = []jira.Issue{*issue}
	} else if opts.Summary != "" {
		// Search by summary using JQL
		jql := fmt.Sprintf("text ~ \"%s\"", opts.Summary)
		found, err := issueService.SearchIssues(jql, 0, 50)
		if err != nil {
			cli.PrintError(err)
			return err
		}
		issues = found
	} else if opts.JQL != "" {
		// Search using raw JQL
		found, err := issueService.SearchIssues(opts.JQL, 0, 50)
		if err != nil {
			cli.PrintError(err)
			return err
		}
		issues = found
	} else {
		return fmt.Errorf("must provide one of: --key, --summary, or --jql")
	}

	// Format and output results
	if opts.Format == "json" {
		outputJSON(issues)
	} else {
		outputTable(issues)
	}

	fmt.Printf("\n📊 Found %d issue(s)\n", len(issues))

	return nil
}

// outputTable outputs search results as a table
func outputTable(issues []jira.Issue) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "KEY\tTYPE\tSUMMARY\tSTATUS")
	fmt.Fprintln(w, "---\t----\t-------\t------")

	for _, issue := range issues {
		status := ""
		if issue.Fields.Status != nil && issue.Fields.Status.Name != "" {
			status = issue.Fields.Status.Name
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			issue.Key,
			issue.Fields.IssueType.Name,
			issue.Fields.Summary,
			status)
	}
}

// outputJSON outputs search results as JSON
func outputJSON(issues []jira.Issue) {
	data, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// NewSearchCommand creates the "search" command with full implementation
func NewSearchCommand() *cobra.Command {
	var opts SearchOptions

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for JIRA tickets",
		Long:  "Search for JIRA tickets by key, summary, or JQL query.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Key, _ = cmd.Flags().GetString("key")
			opts.Summary, _ = cmd.Flags().GetString("summary")
			opts.JQL, _ = cmd.Flags().GetString("jql")
			opts.Format, _ = cmd.Flags().GetString("format")

			return ExecuteSearchCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Search by exact ticket key (e.g., PROJ-123)")
	cmd.Flags().StringVar(&opts.Summary, "summary", "", "Search by summary text (partial match supported)")
	cmd.Flags().StringVar(&opts.JQL, "jql", "", "Advanced search using JQL (JIRA Query Language)")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}
