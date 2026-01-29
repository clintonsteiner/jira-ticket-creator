package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// QueryOptions holds the options for the query command
type QueryOptions struct {
	JQL        string
	Format     string
	Output     string
	MaxResults int
	Fields     string
}

// ExecuteQueryCommand executes the query command
func ExecuteQueryCommand(v *viper.Viper, opts QueryOptions) error {
	// Load configuration with flag overrides
	cfg, err := config.LoadConfigWithFlags(v)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate required configuration
	if err := cfg.ValidateRequired(); err != nil {
		return err
	}

	if opts.JQL == "" {
		return fmt.Errorf("--jql flag is required")
	}

	// Create JIRA client
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	issueService := jira.NewIssueService(client)

	var allIssues []jira.Issue
	startAt := 0
	fetchSize := 50
	if opts.MaxResults < fetchSize {
		fetchSize = opts.MaxResults
	}

	// Fetch all results with pagination
	for {
		issues, err := issueService.SearchIssues(opts.JQL, startAt, fetchSize)
		if err != nil {
			cli.PrintError(err)
			return err
		}

		allIssues = append(allIssues, issues...)

		if len(allIssues) >= opts.MaxResults || len(issues) < fetchSize {
			break
		}

		startAt += fetchSize
	}

	// Trim to max results
	if len(allIssues) > opts.MaxResults {
		allIssues = allIssues[:opts.MaxResults]
	}

	// Format and output results
	var output string
	switch opts.Format {
	case "json":
		output, err = formatJSON(allIssues)
	case "csv":
		output, err = formatCSV(allIssues, opts.Fields)
	case "markdown":
		output, err = formatMarkdown(allIssues, opts.Fields)
	case "html":
		output, err = formatHTML(allIssues, opts.Fields)
	case "table":
		fallthrough
	default:
		output, err = formatTable(allIssues, opts.Fields)
	}

	if err != nil {
		return err
	}

	// Write output
	if opts.Output != "" {
		if err := os.WriteFile(opts.Output, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("✅ Results written to %s (%d issues)\n", opts.Output, len(allIssues))
	} else {
		fmt.Println(output)
		fmt.Printf("\n📊 Found %d issue(s)\n", len(allIssues))
	}

	return nil
}

// parseFields parses the fields string into a slice
func parseFields(fieldsStr string) []string {
	if fieldsStr == "" {
		return []string{"key", "type", "summary", "status", "assignee", "priority"}
	}
	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(strings.ToLower(fields[i]))
	}
	return fields
}

// getFieldValue retrieves a field value from an issue
func getFieldValue(issue jira.Issue, field string) string {
	switch field {
	case "key":
		return issue.Key
	case "type":
		return issue.Fields.IssueType.Name
	case "summary":
		summary := issue.Fields.Summary
		if len(summary) > 50 {
			return summary[:47] + "..."
		}
		return summary
	case "status":
		return issue.Fields.IssueType.Name // Placeholder, status not in Issue struct
	case "assignee":
		if issue.Fields.Assignee != nil {
			return issue.Fields.Assignee.Name
		}
		return "Unassigned"
	case "priority":
		if issue.Fields.Priority != nil {
			return issue.Fields.Priority.Name
		}
		return "None"
	case "project":
		return issue.Fields.Project.Key
	case "description":
		desc := issue.Fields.Description
		if len(desc) > 50 {
			return desc[:47] + "..."
		}
		return desc
	default:
		return ""
	}
}

// formatTable formats issues as a table
func formatTable(issues []jira.Issue, fieldsStr string) (string, error) {
	fields := parseFields(fieldsStr)

	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)

	// Write header
	header := strings.Join(fields, "\t")
	fmt.Fprintln(w, strings.ToUpper(header))
	fmt.Fprintln(w, strings.Repeat("-\t", len(fields)))

	// Write rows
	for _, issue := range issues {
		var row []string
		for _, field := range fields {
			row = append(row, getFieldValue(issue, field))
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
	return sb.String(), nil
}

// formatJSON formats issues as JSON
func formatJSON(issues []jira.Issue) (string, error) {
	data, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// formatCSV formats issues as CSV
func formatCSV(issues []jira.Issue, fieldsStr string) (string, error) {
	fields := parseFields(fieldsStr)

	var sb strings.Builder
	w := csv.NewWriter(&sb)

	// Write header
	w.Write(fields)

	// Write rows
	for _, issue := range issues {
		var row []string
		for _, field := range fields {
			row = append(row, getFieldValue(issue, field))
		}
		w.Write(row)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// formatMarkdown formats issues as Markdown
func formatMarkdown(issues []jira.Issue, fieldsStr string) (string, error) {
	fields := parseFields(fieldsStr)

	var sb strings.Builder
	sb.WriteString("| " + strings.Join(fields, " | ") + " |\n")
	sb.WriteString("|" + strings.Repeat(" --- |", len(fields)) + "\n")

	for _, issue := range issues {
		var row []string
		for _, field := range fields {
			row = append(row, getFieldValue(issue, field))
		}
		sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}

	return sb.String(), nil
}

// formatHTML formats issues as HTML
func formatHTML(issues []jira.Issue, fieldsStr string) (string, error) {
	fields := parseFields(fieldsStr)

	var sb strings.Builder
	sb.WriteString("<table>\n<thead>\n<tr>\n")

	for _, field := range fields {
		sb.WriteString("<th>" + strings.ToUpper(field) + "</th>\n")
	}

	sb.WriteString("</tr>\n</thead>\n<tbody>\n")

	for _, issue := range issues {
		sb.WriteString("<tr>\n")
		for _, field := range fields {
			sb.WriteString("<td>" + getFieldValue(issue, field) + "</td>\n")
		}
		sb.WriteString("</tr>\n")
	}

	sb.WriteString("</tbody>\n</table>")
	return sb.String(), nil
}

// NewQueryCommand creates the "query" command
func NewQueryCommand() *cobra.Command {
	var opts QueryOptions

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Execute JQL queries and display results",
		Long:  "Execute JQL (JIRA Query Language) queries and display results in multiple formats (table, json, csv, markdown, html).",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.JQL, _ = cmd.Flags().GetString("jql")
			opts.Format, _ = cmd.Flags().GetString("format")
			opts.Output, _ = cmd.Flags().GetString("output")
			opts.MaxResults, _ = cmd.Flags().GetInt("max-results")
			opts.Fields, _ = cmd.Flags().GetString("fields")

			return ExecuteQueryCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.JQL, "jql", "", "JQL query string (required)")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json, csv, markdown, html")
	cmd.Flags().StringVar(&opts.Output, "output", "", "Output file path (default: stdout)")
	cmd.Flags().IntVar(&opts.MaxResults, "max-results", 50, "Maximum results to fetch (max: 1000)")
	cmd.Flags().StringVar(&opts.Fields, "fields", "", "Comma-separated fields to display (default: key,type,summary,status,assignee,priority)")

	cmd.MarkFlagRequired("jql")

	return cmd
}
