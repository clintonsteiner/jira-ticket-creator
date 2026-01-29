package cli

import (
	"fmt"
	"strings"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// FormatError formats an error for user-friendly display with suggestions
func FormatError(err error) string {
	var sb strings.Builder

	switch e := err.(type) {
	case *jira.AuthenticationError:
		sb.WriteString("❌ Authentication Failed\n")
		sb.WriteString(fmt.Sprintf("   Error: %s\n", e.Message))
		sb.WriteString("\n")
		sb.WriteString("💡 Suggestions:\n")
		sb.WriteString("   • Check your JIRA email and API token\n")
		sb.WriteString("   • Ensure your API token has the necessary permissions\n")
		sb.WriteString("   • Try setting: export JIRA_EMAIL=your-email\n")
		sb.WriteString("   • Try setting: export JIRA_TOKEN=your-token\n")

	case *jira.NotFoundError:
		sb.WriteString("❌ Resource Not Found\n")
		sb.WriteString(fmt.Sprintf("   Error: %s\n", e.Error()))
		sb.WriteString("\n")
		sb.WriteString("💡 Suggestions:\n")
		sb.WriteString("   • Verify the ticket key exists\n")
		sb.WriteString("   • Check that you're using the correct project key\n")

	case *jira.RateLimitError:
		sb.WriteString("⏱️  Rate Limit Exceeded\n")
		sb.WriteString(fmt.Sprintf("   Error: %s\n", e.Error()))
		sb.WriteString("\n")
		sb.WriteString("💡 Suggestions:\n")
		if e.RetryAfter > 0 {
			sb.WriteString(fmt.Sprintf("   • Please wait %d seconds before retrying\n", e.RetryAfter))
		} else {
			sb.WriteString("   • Please wait a moment and try again\n")
		}
		sb.WriteString("   • Reduce the number of concurrent operations\n")

	case *jira.ValidationError:
		sb.WriteString("❌ Validation Error\n")
		sb.WriteString(fmt.Sprintf("   Field: %s\n", e.Field))
		sb.WriteString(fmt.Sprintf("   Error: %s\n", e.Message))
		if e.Details != "" {
			sb.WriteString(fmt.Sprintf("   Details: %s\n", e.Details))
		}
		sb.WriteString("\n")
		sb.WriteString("💡 Suggestions:\n")
		sb.WriteString("   • Review the field requirements\n")
		sb.WriteString("   • Check the error message for specific guidance\n")

	case *jira.JiraError:
		sb.WriteString("❌ JIRA API Error\n")
		sb.WriteString(fmt.Sprintf("   HTTP Status: %d\n", e.StatusCode))
		sb.WriteString(fmt.Sprintf("   Error: %s\n", e.Error()))
		sb.WriteString("\n")
		sb.WriteString("💡 Suggestions:\n")
		switch e.StatusCode {
		case 400:
			sb.WriteString("   • Check your request format and parameters\n")
			sb.WriteString("   • Ensure all required fields are provided\n")
		case 401:
			sb.WriteString("   • Check your JIRA credentials\n")
		case 403:
			sb.WriteString("   • Verify you have permission for this action\n")
		case 404:
			sb.WriteString("   • Verify the resource exists\n")
		case 500, 502, 503, 504:
			sb.WriteString("   • The JIRA server is experiencing issues\n")
			sb.WriteString("   • Try again in a few moments\n")
		default:
			sb.WriteString("   • Check the error details above\n")
		}

	default:
		sb.WriteString("❌ Error\n")
		sb.WriteString(fmt.Sprintf("   Error: %s\n", err.Error()))
	}

	return sb.String()
}

// PrintError prints an error to stdout in a user-friendly format
func PrintError(err error) {
	fmt.Print(FormatError(err))
}
