package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/clintonsteiner/jira-ticket-creator/internal/templates"
)

// NewTemplateCommand creates the "template" command with full implementation
func NewTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage ticket templates",
		Long:  "List, create, and manage ticket templates.",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeListTemplates()
		},
	}

	cmd.AddCommand(listCmd)

	return cmd
}

// executeListTemplates lists all available templates
func executeListTemplates() error {
	loader := templates.NewLoader()
	templateList := loader.List()

	if len(templateList) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	fmt.Println("Available Templates:")
	fmt.Println("===================")

	for _, t := range templateList {
		fmt.Printf("\n• %s\n", t.Name)
		fmt.Printf("  Description: %s\n", t.DisplayDesc)
		fmt.Printf("  Issue Type: %s\n", t.IssueType)
		fmt.Printf("  Priority: %s\n", t.Priority)
	}

	return nil
}
