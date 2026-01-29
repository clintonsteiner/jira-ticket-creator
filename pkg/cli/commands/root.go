package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCommand creates the root command for the CLI
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira-ticket-creator",
		Short: "A CLI tool to manage JIRA tickets",
		Long: `jira-ticket-creator is a command-line tool that helps you manage JIRA tickets efficiently.
It supports creating, updating, searching, and reporting on tickets with features like:
- Batch ticket creation
- Dependency visualization
- Multiple report formats
- Interactive mode`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind viper to flags for all commands
			return viper.BindPFlags(cmd.PersistentFlags())
		},
	}

	// Add persistent flags
	cmd.PersistentFlags().String("url", "", "JIRA base URL (or set JIRA_URL env var)")
	cmd.PersistentFlags().String("email", "", "JIRA email address (or set JIRA_EMAIL env var)")
	cmd.PersistentFlags().String("token", "", "JIRA API token (or set JIRA_TOKEN env var)")
	cmd.PersistentFlags().String("project", "", "JIRA project key (or set JIRA_PROJECT env var)")
	cmd.PersistentFlags().String("ticket", "", "JIRA ticket key to extract project (or set JIRA_TICKET env var, e.g., PROJ-123)")
	cmd.PersistentFlags().String("config", "", "Path to config file (default: ~/.jirarc)")

	// Bind to viper
	viper.BindPFlag("jira.url", cmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("jira.email", cmd.PersistentFlags().Lookup("email"))
	viper.BindPFlag("jira.token", cmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("jira.project", cmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("jira.ticket", cmd.PersistentFlags().Lookup("ticket"))

	// Add subcommands
	cmd.AddCommand(NewCreateCommand())
	cmd.AddCommand(NewReportCommand())
	cmd.AddCommand(NewUpdateCommand())
	cmd.AddCommand(NewTransitionCommand())
	cmd.AddCommand(NewSearchCommand())
	cmd.AddCommand(NewQueryCommand())
	cmd.AddCommand(NewImportCommand())
	cmd.AddCommand(NewBatchCommand())
	cmd.AddCommand(NewVisualizeCommand())
	cmd.AddCommand(NewTemplateCommand())
	cmd.AddCommand(NewTeamCommand())
	cmd.AddCommand(NewTimelineCommand())
	cmd.AddCommand(NewPMCommand())
	cmd.AddCommand(NewCompletionCommand())

	return cmd
}
