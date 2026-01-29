package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/config"
	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// ImportOptions holds the options for the import command
type ImportOptions struct {
	JQL            string
	MapProject     string
	MapRules       []string
	DryRun         bool
	UpdateExisting bool
	MappingPath    string
}

// ExecuteImportCommand executes the import command
func ExecuteImportCommand(v *viper.Viper, opts ImportOptions) error {
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

	// Load project mapping
	mapping, err := config.LoadMapping(opts.MappingPath)
	if err != nil {
		fmt.Printf("⚠️  Warning: could not load project mapping: %v\n", err)
		mapping = &config.ProjectMapping{Mappings: make(map[string]config.ProjectInfo)}
	}

	// Parse inline mapping rules
	inlineRules := make(map[string]string)
	for _, rule := range opts.MapRules {
		parts := strings.Split(rule, "->")
		if len(parts) != 2 {
			return fmt.Errorf("invalid map rule format: %s (use format: PREFIX->project)", rule)
		}
		prefix := strings.TrimSpace(parts[0])
		project := strings.TrimSpace(parts[1])
		inlineRules[prefix] = project
	}

	// Create JIRA client
	client := jira.NewClient(cfg.JIRA.URL, cfg.JIRA.Email, cfg.JIRA.Token)
	issueService := jira.NewIssueService(client)

	// Execute JQL query
	issues, err := issueService.SearchIssues(opts.JQL, 0, 1000)
	if err != nil {
		cli.PrintError(err)
		return err
	}

	if len(issues) == 0 {
		fmt.Println("ℹ️  No issues found matching the query")
		return nil
	}

	fmt.Printf("📥 Found %d issue(s) to import\n", len(issues))

	// Create ticket records from issues
	records := make([]jira.TicketRecord, 0, len(issues))
	for _, issue := range issues {
		project := opts.MapProject
		if project == "" {
			// Try to find project from key prefix
			keyPrefix := extractKeyPrefix(issue.Key)
			if rule, exists := inlineRules[keyPrefix]; exists {
				project = rule
			} else if mapped := mapping.FindProjectForKey(issue.Key); mapped != "" {
				project = mapped
			}
		}

		record := jira.TicketRecord{
			Key:       issue.Key,
			Summary:   issue.Fields.Summary,
			Status:    "Open",
			CreatedAt: time.Now(),
			Creator:   "imported",
			Assignee:  "",
			Priority:  "Medium",
			IssueType: issue.Fields.IssueType.Name,
			Project:   project,
		}

		if issue.Fields.Assignee != nil && issue.Fields.Assignee.Name != "" {
			record.Assignee = issue.Fields.Assignee.Name
		}

		if issue.Fields.Priority != nil && issue.Fields.Priority.Name != "" {
			record.Priority = issue.Fields.Priority.Name
		}

		records = append(records, record)
	}

	// Display summary
	fmt.Printf("\n📋 Import Summary:\n")
	fmt.Printf("   Total: %d\n", len(records))

	// Group by project
	byProject := make(map[string]int)
	noProject := 0
	for _, r := range records {
		if r.Project == "" {
			noProject++
		} else {
			byProject[r.Project]++
		}
	}

	for project, count := range byProject {
		fmt.Printf("   %s: %d\n", project, count)
	}

	if noProject > 0 {
		fmt.Printf("   (unmapped): %d\n", noProject)
	}

	// Dry run?
	if opts.DryRun {
		fmt.Println("\n✅ Dry run completed (no changes saved)")
		return nil
	}

	// Load existing records
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")
	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	existing, _ := repo.GetAll()

	// Update or add records
	processed := 0
	skipped := 0
	for _, record := range records {
		if existing != nil {
			foundExisting := false
			for _, e := range existing {
				if e.Key == record.Key {
					foundExisting = true
					if opts.UpdateExisting {
						if err := repo.Update(record); err != nil {
							fmt.Printf("⚠️  Failed to update %s: %v\n", record.Key, err)
						} else {
							processed++
						}
					} else {
						skipped++
					}
					break
				}
			}
			if !foundExisting {
				if err := repo.Add(record); err != nil {
					fmt.Printf("⚠️  Failed to add %s: %v\n", record.Key, err)
				} else {
					processed++
				}
			}
		} else {
			if err := repo.Add(record); err != nil {
				fmt.Printf("⚠️  Failed to add %s: %v\n", record.Key, err)
			} else {
				processed++
			}
		}
	}

	fmt.Printf("\n✅ Import completed\n")
	fmt.Printf("   Added/Updated: %d\n", processed)
	if skipped > 0 {
		fmt.Printf("   Skipped (existing): %d\n", skipped)
	}

	return nil
}

// extractKeyPrefix extracts the prefix from a ticket key (e.g., "PROJ" from "PROJ-123")
func extractKeyPrefix(key string) string {
	parts := strings.Split(key, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// NewImportCommand creates the "import" command
func NewImportCommand() *cobra.Command {
	var opts ImportOptions

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import existing JIRA tickets into local tracking",
		Long: `Import existing JIRA tickets into local tracking with optional project mapping.

The import command executes a JQL query and saves the results to local storage.
You can map tickets to logical projects for better organization.

Examples:
  # Import all tickets from a project
  import --jql "project = PROJ" --map-project backend

  # Dry run to preview what would be imported
  import --jql "project = PROJ" --map-project backend --dry-run

  # Use inline mapping rules
  import --jql "project in (PROJ, BACK)" \
    --map-rule "PROJ->backend" --map-rule "BACK->backend"

  # Update existing imported tickets
  import --jql "status = 'In Progress'" --update-existing`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.JQL, _ = cmd.Flags().GetString("jql")
			opts.MapProject, _ = cmd.Flags().GetString("map-project")
			opts.MapRules, _ = cmd.Flags().GetStringSlice("map-rule")
			opts.DryRun, _ = cmd.Flags().GetBool("dry-run")
			opts.UpdateExisting, _ = cmd.Flags().GetBool("update-existing")
			opts.MappingPath, _ = cmd.Flags().GetString("mapping-path")

			return ExecuteImportCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.JQL, "jql", "", "JQL query for ticket selection (required)")
	cmd.Flags().StringVar(&opts.MapProject, "map-project", "", "Logical project name to assign to all imported tickets")
	cmd.Flags().StringSliceVar(&opts.MapRules, "map-rule", []string{}, "Inline mapping rules (repeatable, format: PREFIX->project)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be imported without saving")
	cmd.Flags().BoolVar(&opts.UpdateExisting, "update-existing", false, "Update existing local tickets")
	cmd.Flags().StringVar(&opts.MappingPath, "mapping-path", "", "Path to project mapping file (default: ~/.jira/project-mapping.json)")

	cmd.MarkFlagRequired("jql")

	return cmd
}
