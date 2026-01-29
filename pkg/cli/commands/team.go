package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
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
			ticketFilter, _ := cmd.Flags().GetString("ticket")
			creatorFilter, _ := cmd.Flags().GetString("creator")
			assigneeFilter, _ := cmd.Flags().GetString("assignee")
			return executeTeamSummary(projectFilter, ticketFilter, creatorFilter, assigneeFilter)
		},
	}
	summaryCmd.Flags().String("project", "", "Filter by project")
	summaryCmd.Flags().String("ticket", "", "Filter by ticket key (can be comma-separated)")
	summaryCmd.Flags().String("creator", "", "Filter by creator (can be comma-separated)")
	summaryCmd.Flags().String("assignee", "", "Filter by assignee (can be comma-separated)")

	// Assignments subcommand
	assignCmd := &cobra.Command{
		Use:   "assignments",
		Short: "Show workload and assignments",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFilter, _ := cmd.Flags().GetString("project")
			ticketFilter, _ := cmd.Flags().GetString("ticket")
			creatorFilter, _ := cmd.Flags().GetString("creator")
			assigneeFilter, _ := cmd.Flags().GetString("assignee")
			return executeAssignments(projectFilter, ticketFilter, creatorFilter, assigneeFilter)
		},
	}
	assignCmd.Flags().String("project", "", "Filter by project")
	assignCmd.Flags().String("ticket", "", "Filter by ticket key (can be comma-separated)")
	assignCmd.Flags().String("creator", "", "Filter by creator (can be comma-separated)")
	assignCmd.Flags().String("assignee", "", "Filter by assignee (can be comma-separated)")

	// Timeline subcommand
	timelineCmd := &cobra.Command{
		Use:   "timeline",
		Short: "Show project timeline and progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFilter, _ := cmd.Flags().GetString("project")
			ticketFilter, _ := cmd.Flags().GetString("ticket")
			creatorFilter, _ := cmd.Flags().GetString("creator")
			assigneeFilter, _ := cmd.Flags().GetString("assignee")
			return executeTimeline(projectFilter, ticketFilter, creatorFilter, assigneeFilter)
		},
	}
	timelineCmd.Flags().String("project", "", "Filter by project")
	timelineCmd.Flags().String("ticket", "", "Filter by ticket key (can be comma-separated)")
	timelineCmd.Flags().String("creator", "", "Filter by creator (can be comma-separated)")
	timelineCmd.Flags().String("assignee", "", "Filter by assignee (can be comma-separated)")

	cmd.AddCommand(summaryCmd, assignCmd, timelineCmd)

	return cmd
}

// executeTeamSummary shows tickets grouped by creator
func executeTeamSummary(projectFilter string, ticketFilter string, creatorFilter string, assigneeFilter string) error {
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

	// Filter by ticket key(s) if specified
	if ticketFilter != "" {
		ticketKeys := strings.Split(ticketFilter, ",")
		ticketMap := make(map[string]bool)
		for _, key := range ticketKeys {
			ticketMap[strings.TrimSpace(key)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if ticketMap[r.Key] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for keys: %s\n", ticketFilter)
			return nil
		}
	}

	// Filter by creator(s) if specified
	if creatorFilter != "" {
		creators := strings.Split(creatorFilter, ",")
		creatorMap := make(map[string]bool)
		for _, creator := range creators {
			creatorMap[strings.TrimSpace(creator)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if creatorMap[r.Creator] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for creators: %s\n", creatorFilter)
			return nil
		}
	}

	// Filter by assignee(s) if specified
	if assigneeFilter != "" {
		assignees := strings.Split(assigneeFilter, ",")
		assigneeMap := make(map[string]bool)
		for _, assignee := range assignees {
			assigneeMap[strings.TrimSpace(assignee)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if assigneeMap[r.Assignee] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for assignees: %s\n", assigneeFilter)
			return nil
		}
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
func executeAssignments(projectFilter string, ticketFilter string, creatorFilter string, assigneeFilter string) error {
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

	// Filter by ticket key(s) if specified
	if ticketFilter != "" {
		ticketKeys := strings.Split(ticketFilter, ",")
		ticketMap := make(map[string]bool)
		for _, key := range ticketKeys {
			ticketMap[strings.TrimSpace(key)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if ticketMap[r.Key] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for keys: %s\n", ticketFilter)
			return nil
		}
	}

	// Filter by creator(s) if specified
	if creatorFilter != "" {
		creators := strings.Split(creatorFilter, ",")
		creatorMap := make(map[string]bool)
		for _, creator := range creators {
			creatorMap[strings.TrimSpace(creator)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if creatorMap[r.Creator] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for creators: %s\n", creatorFilter)
			return nil
		}
	}

	// Filter by assignee(s) if specified
	if assigneeFilter != "" {
		assignees := strings.Split(assigneeFilter, ",")
		assigneeMap := make(map[string]bool)
		for _, assignee := range assignees {
			assigneeMap[strings.TrimSpace(assignee)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if assigneeMap[r.Assignee] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for assignees: %s\n", assigneeFilter)
			return nil
		}
	}

	// Filter by project if specified
	if projectFilter != "" {
		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if r.Project == projectFilter {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for project: %s\n", projectFilter)
			return nil
		}
	}

	teamReport := &reports.TeamReport{}
	report := teamReport.GenerateAssignmentMap(records)

	fmt.Println(report)
	return nil
}

// executeTimeline shows project timeline
func executeTimeline(projectFilter string, ticketFilter string, creatorFilter string, assigneeFilter string) error {
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

	// Filter by ticket key(s) if specified
	if ticketFilter != "" {
		ticketKeys := strings.Split(ticketFilter, ",")
		ticketMap := make(map[string]bool)
		for _, key := range ticketKeys {
			ticketMap[strings.TrimSpace(key)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if ticketMap[r.Key] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for keys: %s\n", ticketFilter)
			return nil
		}
	}

	// Filter by creator(s) if specified
	if creatorFilter != "" {
		creators := strings.Split(creatorFilter, ",")
		creatorMap := make(map[string]bool)
		for _, creator := range creators {
			creatorMap[strings.TrimSpace(creator)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if creatorMap[r.Creator] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for creators: %s\n", creatorFilter)
			return nil
		}
	}

	// Filter by assignee(s) if specified
	if assigneeFilter != "" {
		assignees := strings.Split(assigneeFilter, ",")
		assigneeMap := make(map[string]bool)
		for _, assignee := range assignees {
			assigneeMap[strings.TrimSpace(assignee)] = true
		}

		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if assigneeMap[r.Assignee] {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for assignees: %s\n", assigneeFilter)
			return nil
		}
	}

	// Filter by project if specified
	if projectFilter != "" {
		filtered := make([]jira.TicketRecord, 0)
		for _, r := range records {
			if r.Project == projectFilter {
				filtered = append(filtered, r)
			}
		}
		records = filtered

		if len(records) == 0 {
			fmt.Printf("No tickets found for project: %s\n", projectFilter)
			return nil
		}
	}

	teamReport := &reports.TeamReport{}
	report := teamReport.GenerateTimeline(records)

	fmt.Println(report)
	return nil
}
