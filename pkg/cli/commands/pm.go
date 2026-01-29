package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

// NewPMCommand creates the "pm" command for project management reporting
func NewPMCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pm",
		Short: "Project management dashboard and reporting",
		Long:  "View project status, hierarchies, risks, and team metrics for executive visibility.",
	}

	// Dashboard subcommand
	dashboardCmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Executive summary dashboard",
		Long:  "High-level overview of project status, team workload, and priorities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePMDashboard()
		},
	}

	// Hierarchy subcommand
	hierarchyCmd := &cobra.Command{
		Use:   "hierarchy",
		Short: "Show ticket hierarchy (parent-child relationships)",
		Long:  "Display tickets organized by epics, stories, and their subtasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePMHierarchy()
		},
	}

	// Risk assessment subcommand
	riskCmd := &cobra.Command{
		Use:   "risk",
		Short: "Risk assessment and blockers",
		Long:  "Identify blocked items, unassigned work, and project risks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePMRisk()
		},
	}

	// Details subcommand
	detailsCmd := &cobra.Command{
		Use:   "details",
		Short: "Detailed ticket inventory",
		Long:  "Complete table of all tickets with status, assignment, and dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePMDetails()
		},
	}

	// Parent ticket creation (will create a sample for now)
	parentCmd := &cobra.Command{
		Use:   "create-parent",
		Short: "Create a parent epic ticket",
		Long:  "Create a parent epic to organize and track related child tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePMCreateParent()
		},
	}

	cmd.AddCommand(dashboardCmd, hierarchyCmd, riskCmd, detailsCmd, parentCmd)

	return cmd
}

// executePMDashboard shows the executive dashboard
func executePMDashboard() error {
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

	pmReport := &reports.PMReport{}
	dashboard := pmReport.GeneratePMDashboard(records)

	fmt.Println(dashboard)
	return nil
}

// executePMHierarchy shows the ticket hierarchy
func executePMHierarchy() error {
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

	pmReport := &reports.PMReport{}
	hierarchy := pmReport.GenerateProjectHierarchy(records)

	fmt.Println(hierarchy)
	return nil
}

// executePMRisk shows risk assessment
func executePMRisk() error {
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

	pmReport := &reports.PMReport{}
	risk := pmReport.GenerateRiskReport(records)

	fmt.Println(risk)
	return nil
}

// executePMDetails shows detailed ticket inventory
func executePMDetails() error {
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

	pmReport := &reports.PMReport{}
	details := pmReport.GenerateTicketDetailsTable(records)

	fmt.Println(details)
	return nil
}

// executePMCreateParent guides user through creating a parent epic
func executePMCreateParent() error {
	fmt.Println("\n CREATE PARENT EPIC\n")
	fmt.Println("A parent epic is used to organize and track related work items.")
	fmt.Println("\nTo create a parent epic:")
	fmt.Println("  1. Use the standard create command with --type Epic")
	fmt.Println("  2. Set it up to block child tickets")
	fmt.Println("\nExample:")
	fmt.Println("  ./jira-ticket-creator create \\")
	fmt.Println("    --summary \"Q1 Platform Upgrade\" \\")
	fmt.Println("    --type Epic \\")
	fmt.Println("    --priority High \\")
	fmt.Println("    --description \"Parent epic for all platform upgrade work\"")
	fmt.Println("\nThen create child tickets with:")
	fmt.Println("  ./jira-ticket-creator create \\")
	fmt.Println("    --summary \"Upgrade Database\" \\")
	fmt.Println("    --type Task \\")
	fmt.Println("    --blocked-by <PARENT-EPIC-KEY>")
	fmt.Println("\nView the hierarchy with:")
	fmt.Println("  ./jira-ticket-creator pm hierarchy")
	fmt.Println()
	return nil
}
