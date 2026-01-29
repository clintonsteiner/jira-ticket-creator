package reports

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/clintonsteiner/jira-ticket-creator/internal/jira"
)

// HTMLReporter generates HTML format reports
type HTMLReporter struct{}

// Generate generates an HTML report
func (r *HTMLReporter) Generate(issues []jira.Issue) (string, error) {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>JIRA Tickets Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: #f5f7fa;
            padding: 20px;
            color: #333;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
            padding: 30px;
        }

        h1 {
            color: #0052cc;
            margin-bottom: 10px;
            border-bottom: 3px solid #0052cc;
            padding-bottom: 10px;
        }

        .meta {
            color: #666;
            font-size: 14px;
            margin-bottom: 20px;
        }

        h2 {
            color: #0052cc;
            margin-top: 30px;
            margin-bottom: 15px;
            font-size: 18px;
        }

        .summary-box {
            background-color: #f5f7fa;
            border-left: 4px solid #0052cc;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 4px;
        }

        .summary-box p {
            margin: 5px 0;
        }

        .summary-box strong {
            color: #0052cc;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
        }

        th {
            background-color: #0052cc;
            color: white;
            padding: 12px;
            text-align: left;
            font-weight: 600;
        }

        td {
            padding: 12px;
            border-bottom: 1px solid #eee;
        }

        tr:hover {
            background-color: #f9f9f9;
        }

        .key {
            font-weight: 600;
            color: #0052cc;
            font-family: monospace;
        }

        .type-task { background-color: #deebff; color: #0052cc; }
        .type-story { background-color: #dffcf0; color: #216e4e; }
        .type-bug { background-color: #ffeceb; color: #ae2a19; }
        .type-epic { background-color: #eae6ff; color: #5e4db2; }
        .type-subtask { background-color: #e3f2fd; color: #1565c0; }

        .type-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: 600;
        }

        .priority-critical { color: #ae2a19; font-weight: 600; }
        .priority-high { color: #974f0c; font-weight: 600; }
        .priority-medium { color: #0055cc; }
        .priority-low { color: #216e4e; }

        .detail-section {
            border: 1px solid #eee;
            border-radius: 4px;
            padding: 15px;
            margin-bottom: 15px;
        }

        .detail-section h3 {
            color: #0052cc;
            margin-bottom: 10px;
            font-size: 16px;
        }

        .detail-row {
            margin: 8px 0;
        }

        .detail-row strong {
            display: inline-block;
            width: 120px;
            color: #0052cc;
        }

        .label {
            display: inline-block;
            background-color: #f5f7fa;
            border: 1px solid #ddd;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            margin-right: 5px;
            margin-bottom: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
`)

	// Header
	sb.WriteString("<h1>📋 JIRA Tickets Report</h1>\n")
	sb.WriteString(fmt.Sprintf("<div class=\"meta\">Generated: %s</div>\n", time.Now().Format("2006-01-02 15:04:05")))

	// Summary
	sb.WriteString("<h2>Summary</h2>\n")
	sb.WriteString("<div class=\"summary-box\">\n")
	sb.WriteString(fmt.Sprintf("<p><strong>Total Tickets:</strong> %d</p>\n", len(issues)))

	// Count by type
	typeCount := make(map[string]int)
	priorityCount := make(map[string]int)
	for _, issue := range issues {
		typeCount[issue.Fields.IssueType.Name]++
		if issue.Fields.Priority != nil {
			priorityCount[issue.Fields.Priority.Name]++
		}
	}

	if len(typeCount) > 0 {
		sb.WriteString("<p><strong>By Type:</strong> ")
		for typ, count := range typeCount {
			sb.WriteString(fmt.Sprintf("%s: %d, ", typ, count))
		}
		sb.WriteString("</p>\n")
	}

	if len(priorityCount) > 0 {
		sb.WriteString("<p><strong>By Priority:</strong> ")
		for priority, count := range priorityCount {
			sb.WriteString(fmt.Sprintf("%s: %d, ", priority, count))
		}
		sb.WriteString("</p>\n")
	}

	sb.WriteString("</div>\n")

	// Table
	sb.WriteString("<h2>Tickets</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<thead><tr><th>Key</th><th>Type</th><th>Summary</th><th>Assignee</th><th>Priority</th></tr></thead>\n")
	sb.WriteString("<tbody>\n")

	for _, issue := range issues {
		assignee := "-"
		if issue.Fields.Assignee != nil && issue.Fields.Assignee.Name != "" {
			assignee = html.EscapeString(issue.Fields.Assignee.Name)
		}

		priority := "-"
		priorityClass := ""
		if issue.Fields.Priority != nil {
			priority = html.EscapeString(issue.Fields.Priority.Name)
			priorityClass = fmt.Sprintf(" class=\"priority-%s\"", strings.ToLower(priority))
		}

		typeClass := fmt.Sprintf("type-%s", strings.ToLower(strings.ReplaceAll(issue.Fields.IssueType.Name, " ", "-")))

		sb.WriteString(fmt.Sprintf(
			"<tr><td><span class=\"key\">%s</span></td><td><span class=\"type-badge %s\">%s</span></td><td>%s</td><td>%s</td><td%s>%s</td></tr>\n",
			issue.Key,
			typeClass,
			html.EscapeString(issue.Fields.IssueType.Name),
			html.EscapeString(issue.Fields.Summary),
			assignee,
			priorityClass,
			priority))
	}

	sb.WriteString("</tbody>\n")
	sb.WriteString("</table>\n")

	// Details
	sb.WriteString("<h2>Details</h2>\n")
	for _, issue := range issues {
		sb.WriteString("<div class=\"detail-section\">\n")
		sb.WriteString(fmt.Sprintf("<h3>%s</h3>\n", html.EscapeString(issue.Key)))

		sb.WriteString(fmt.Sprintf("<div class=\"detail-row\"><strong>Summary:</strong> %s</div>\n", html.EscapeString(issue.Fields.Summary)))

		if issue.Fields.Description != "" {
			sb.WriteString(fmt.Sprintf("<div class=\"detail-row\"><strong>Description:</strong> %s</div>\n", html.EscapeString(issue.Fields.Description)))
		}

		sb.WriteString(fmt.Sprintf("<div class=\"detail-row\"><strong>Type:</strong> %s</div>\n", html.EscapeString(issue.Fields.IssueType.Name)))

		if issue.Fields.Priority != nil {
			sb.WriteString(fmt.Sprintf("<div class=\"detail-row\"><strong>Priority:</strong> <span class=\"priority-%s\">%s</span></div>\n",
				strings.ToLower(issue.Fields.Priority.Name),
				html.EscapeString(issue.Fields.Priority.Name)))
		}

		if issue.Fields.Assignee != nil && (issue.Fields.Assignee.Name != "" || issue.Fields.Assignee.EmailAddress != "") {
			assignee := issue.Fields.Assignee.Name
			if assignee == "" {
				assignee = issue.Fields.Assignee.EmailAddress
			}
			sb.WriteString(fmt.Sprintf("<div class=\"detail-row\"><strong>Assignee:</strong> %s</div>\n", html.EscapeString(assignee)))
		}

		if len(issue.Fields.Labels) > 0 {
			sb.WriteString("<div class=\"detail-row\"><strong>Labels:</strong> ")
			for _, label := range issue.Fields.Labels {
				sb.WriteString(fmt.Sprintf("<span class=\"label\">%s</span>", html.EscapeString(label)))
			}
			sb.WriteString("</div>\n")
		}

		sb.WriteString("</div>\n")
	}

	sb.WriteString(`
    </div>
</body>
</html>
`)

	return sb.String(), nil
}
