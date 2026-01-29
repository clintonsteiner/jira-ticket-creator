package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

// NewTeamCommand creates the "team" command for team-based reporting
func NewTeamCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "View team-based ticket reports",
		Long:  "View tickets organized by creator, assignee, and timeline information.",
	}

	// Team summary subcommand
	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Show ticket summary by creator",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFilter, _ := cmd.Flags().GetString("project")
			return executeTeamSummary(projectFilter)
		},
	}
	summaryCmd.Flags().String("project", "", "Filter by project")

	// Assignments subcommand
	assignCmd := &cobra.Command{
		Use:   "assignments",
		Short: "Show workload and assignments",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAssignments()
		},
	}

	// Timeline subcommand
	timelineCmd := &cobra.Command{
		Use:   "timeline",
		Short: "Show project timeline and progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeTimeline()
		},
	}

	cmd.AddCommand(summaryCmd, assignCmd, timelineCmd)

	return cmd
}

// executeTeamSummary shows tickets grouped by creator
func executeTeamSummary(projectFilter string) error {
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

	teamReport := &reports.TeamReport{}
	var report string
	if projectFilter != "" {
		report = teamReport.GenerateTeamSummaryWithFilter(records, projectFilter)
	} else {
		report = teamReport.GenerateTeamSummary(records)
	}

	fmt.Println(report)
	return nil
}

// executeAssignments shows workload assignments
func executeAssignments() error {
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

	teamReport := &reports.TeamReport{}
	report := teamReport.GenerateAssignmentMap(records)

	fmt.Println(report)
	return nil
}

// executeTimeline shows project timeline
func executeTimeline() error {
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

	teamReport := &reports.TeamReport{}
	report := teamReport.GenerateTimeline(records)

	fmt.Println(report)
	return nil
}
