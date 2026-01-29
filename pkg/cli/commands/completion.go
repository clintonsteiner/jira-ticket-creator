package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCompletionCommand creates the "completion" command
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completions",
		Long: `Generate shell completion scripts for your shell.

To load completions:

Bash:
  $ source <(jira-ticket-creator completion bash)
  # To load completions for each session, execute once:
  $ jira-ticket-creator completion bash > /etc/bash_completion.d/jira-ticket-creator

Zsh:
  $ source <(jira-ticket-creator completion zsh)
  # To load completions for each session, execute once:
  $ jira-ticket-creator completion zsh > "${fpath[1]}/_jira-ticket-creator"

Fish:
  $ jira-ticket-creator completion fish | source
  # To load completions for each session, execute once:
  $ jira-ticket-creator completion fish > ~/.config/fish/completions/jira-ticket-creator.fish
`,
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			case "fish":
				return root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return root.GenPowerShellCompletion(os.Stdout)
			}
			return fmt.Errorf("unknown shell: %s", args[0])
		},
	}

	return cmd
}
