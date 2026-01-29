package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clintonsteiner/jira-ticket-creator/internal/reports"
	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
	"github.com/clintonsteiner/jira-ticket-creator/pkg/cli"
)

// VisualizeOptions holds options for the visualize command
type VisualizeOptions struct {
	Format string
	Output string
}

// ExecuteVisualizeCommand executes the visualize command
func ExecuteVisualizeCommand(v *viper.Viper, opts VisualizeOptions) error {
	// No configuration needed for visualization, just load local storage

	// Load ticket records from storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	recordFile := filepath.Join(homeDir, ".jira", "tickets.json")
	repo, err := storage.NewJSONRepository(recordFile)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create visualizer
	visualizer := reports.NewVisualizer(repo)

	// Check for circular dependencies
	cycles, err := visualizer.DetectCircularDependencies()
	if err != nil {
		cli.PrintError(fmt.Errorf("failed to detect circular dependencies: %w", err))
		return err
	}

	if len(cycles) > 0 {
		fmt.Println("⚠️  Circular dependencies detected:")
		for _, cycle := range cycles {
			fmt.Printf("   - %v\n", cycle)
		}
		fmt.Println()
	}

	// Generate visualization
	var output string

	switch opts.Format {
	case "mermaid":
		var err error
		output, err = visualizer.GenerateMermaid()
		if err != nil {
			cli.PrintError(err)
			return err
		}

		// Add Mermaid wrapper for easy integration
		if opts.Output == "" {
			output = "```mermaid\n" + output + "```"
		}

	case "dot":
		var err error
		output, err = visualizer.GenerateDOT()
		if err != nil {
			cli.PrintError(err)
			return err
		}

	default: // tree
		var err error
		output, err = visualizer.GenerateTree()
		if err != nil {
			cli.PrintError(err)
			return err
		}
	}

	// Output visualization
	if opts.Output != "" {
		if err := os.WriteFile(opts.Output, []byte(output), 0644); err != nil {
			cli.PrintError(fmt.Errorf("failed to write visualization: %w", err))
			return err
		}
		fmt.Printf("✅ Visualization written to: %s\n", opts.Output)
	} else {
		fmt.Println(output)
	}

	return nil
}

// NewVisualizeCommand creates the "visualize" command with full implementation
func NewVisualizeCommand() *cobra.Command {
	var opts VisualizeOptions

	cmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize ticket dependencies",
		Long:  "Create a visualization of ticket dependencies and relationships. Formats: tree (ASCII), mermaid (diagram), dot (Graphviz).",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind flags to viper
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// Read values from flags
			opts.Format, _ = cmd.Flags().GetString("format")
			opts.Output, _ = cmd.Flags().GetString("output")

			return ExecuteVisualizeCommand(viper.GetViper(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "tree", "Output format: tree (ASCII tree), mermaid (Mermaid diagram), dot (Graphviz DOT format)")
	cmd.Flags().StringVar(&opts.Output, "output", "", "Output file path (optional, default: print to stdout)")

	return cmd
}
