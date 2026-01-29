package interactive

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptString prompts for a string value
func PromptString(label string, required bool) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	if required {
		prompt.Validate = func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("this field is required")
			}
			return nil
		}
	}

	return prompt.Run()
}

// PromptStringWithDefault prompts for a string with a default value
func PromptStringWithDefault(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	return prompt.Run()
}

// PromptSelect prompts to select from a list
func PromptSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	return result, err
}

// PromptMultiSelect prompts to select multiple items
func PromptMultiSelect(label string, items []string) ([]string, error) {
	selected := []string{}

	for {
		prompt := promptui.Select{
			Label: label + " (selected: " + strings.Join(selected, ", ") + ") [space to select, enter to finish]",
			Items: items,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			return selected, err
		}

		item := items[idx]
		found := false
		for _, s := range selected {
			if s == item {
				found = true
				break
			}
		}

		if found {
			// Remove item
			newSelected := []string{}
			for _, s := range selected {
				if s != item {
					newSelected = append(newSelected, s)
				}
			}
			selected = newSelected
		} else {
			// Add item
			selected = append(selected, item)
		}

		// Check if user wants to continue
		continuePrompt := promptui.Prompt{
			Label: "Add more (y/n)",
		}

		cont, _ := continuePrompt.Run()
		if strings.ToLower(cont) != "y" {
			break
		}
	}

	return selected, nil
}

// PromptConfirm prompts for yes/no confirmation
func PromptConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	return result == "y", err
}
