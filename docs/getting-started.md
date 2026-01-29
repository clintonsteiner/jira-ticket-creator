---
layout: default
title: Getting Started
---

# Getting Started with JIRA Ticket Creator

## Installation

### Prerequisites
- Go 1.21 or later
- JIRA Cloud account with API access

### Build from Source

```bash
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

### Install from GitHub

```bash
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

### Verify Installation

```bash
jira-ticket-creator --version
jira-ticket-creator --help
```

## Setup Your Credentials

You need three pieces of information:
1. **JIRA URL** - Your Atlassian instance (e.g., https://company.atlassian.net)
2. **Email** - Your JIRA email address
3. **API Token** - Generate from Account Settings

### Get Your API Token

1. Go to [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click "Create API token"
3. Give it a name (e.g., "jira-ticket-creator")
4. Copy the token

### Option 1: Environment Variables (Recommended)

```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ  # Your default project key
```

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make permanent:

```bash
cat >> ~/.zshrc << 'EOF'
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
EOF

source ~/.zshrc
```

### Option 2: Configuration File

Create `~/.jirarc`:

```yaml
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: your-api-token
  project: PROJ
```

Make it readable only by you:

```bash
chmod 600 ~/.jirarc
```

### Option 3: Command-Line Flags

Pass credentials directly:

```bash
jira-ticket-creator create --summary "Task" \
  --url https://company.atlassian.net \
  --email user@company.com \
  --token api-token \
  --project PROJ
```

## Test Your Setup

```bash
# Search for an existing ticket to verify credentials work
jira-ticket-creator search --key PROJ-1
```

Expected output:
```
KEY         TYPE   SUMMARY                  STATUS
---         ----   -------                  ------
PROJ-1      Task   Your Existing Ticket     Done
```

## Create Your First Ticket

```bash
jira-ticket-creator create --summary "Hello JIRA Ticket Creator"
```

Expected output:
```
✅ Ticket created successfully!
PROJ-123
```

## Next Steps

### 1. Create a More Detailed Ticket

```bash
jira-ticket-creator create \
  --summary "Implement user authentication" \
  --description "Add OAuth 2.0 support to the API" \
  --type Story \
  --priority High \
  --assignee john@company.com
```

### 2. Search for Tickets

```bash
# Find all in-progress tickets
jira-ticket-creator search --jql "status = 'In Progress'"

# Find high-priority bugs
jira-ticket-creator search --jql "type = Bug AND priority = High"
```

### 3. Query with Flexible Output

```bash
# Export to JSON
jira-ticket-creator query --jql "project = PROJ" --format json --output tickets.json

# Export to CSV
jira-ticket-creator query --jql "status = 'To Do'" --format csv --output todo.csv

# Export to HTML
jira-ticket-creator query --jql "priority = Critical" --format html --output critical.html
```

### 4. Batch Create Tickets

Create `tickets.csv`:
```csv
summary,description,issue_type,priority,assignee
"Task 1","Description 1",Task,High,john@company.com
"Task 2","Description 2",Task,Medium,jane@company.com
"Task 3","Description 3",Bug,Critical,bob@company.com
```

Import them:
```bash
jira-ticket-creator batch create --input tickets.csv
```

### 5. View Team Reports

```bash
# Summary by creator
jira-ticket-creator team summary

# Workload by assignee
jira-ticket-creator team assignments

# Project timeline
jira-ticket-creator team timeline
```

### 6. Generate Gantt Chart

```bash
# View in terminal
jira-ticket-creator gantt --format ascii

# Generate for documentation
jira-ticket-creator gantt --format mermaid --output gantt.md

# Create interactive chart
jira-ticket-creator gantt --format html --output gantt.html
open gantt.html
```

## Common Workflows

### Setup for Multiple Teams

```bash
# Create project mapping
cat > ~/.jira/project-mapping.json << 'EOF'
{
  "mappings": {
    "backend": {
      "ticket_keys": ["PROJ", "API"],
      "description": "Backend Team"
    },
    "frontend": {
      "ticket_keys": ["UI", "WEB"],
      "description": "Frontend Team"
    }
  }
}
EOF

# Import tickets for each team
jira-ticket-creator import --jql "project = PROJ" --map-project backend
jira-ticket-creator import --jql "project = UI" --map-project frontend

# View team reports
jira-ticket-creator team summary --project backend
```

### Weekly Status Report

```bash
#!/bin/bash
# Create weekly report

DATE=$(date +%Y-%m-%d)
REPORT_DIR="reports/$DATE"
mkdir -p "$REPORT_DIR"

jira-ticket-creator gantt --format html --output "$REPORT_DIR/gantt.html"
jira-ticket-creator pm dashboard > "$REPORT_DIR/dashboard.txt"
jira-ticket-creator team summary > "$REPORT_DIR/team.txt"

echo "Reports available in $REPORT_DIR"
```

### Monitor Critical Issues

```bash
# Create alias
alias critical='jira-ticket-creator query --jql "priority = Critical" --format table'

# Use it
critical

# Or set up monitoring
watch -n 300 critical  # Update every 5 minutes
```

## Troubleshooting

### Test Your Configuration

```bash
# Check if environment variables are set
echo "URL: $JIRA_URL"
echo "Email: $JIRA_EMAIL"
echo "Project: $JIRA_PROJECT"

# Test connection
jira-ticket-creator search --key PROJ-1
```

### Common Issues

| Problem | Solution |
|---------|----------|
| Command not found | Ensure binary is in PATH or use full path `./jira-ticket-creator` |
| Authentication failed | Check API token (not password), user permissions, URL format |
| Project not found | Verify project key and that you have access to it |
| File not found | Use absolute paths or check working directory |

### Enable Debug Mode

Set verbose environment variable:
```bash
VERBOSE=1 jira-ticket-creator create --summary "Test"
```

## Get Help

### View Command Help

```bash
jira-ticket-creator --help
jira-ticket-creator create --help
jira-ticket-creator search --help
```

### Documentation

- Full [API Documentation](./docs/api/)
- [CLI Command Reference](./docs/cli/)
- [Advanced Topics](./docs/advanced/)

### Report Issues

Found a bug? Report it on [GitHub Issues](https://github.com/clintonsteiner/jira-ticket-creator/issues)

## Next: Explore More Features

1. **[Create Tickets](./docs/cli/create/)** - Detailed creation options
2. **[Search & Query](./docs/cli/query/)** - Advanced searching
3. **[Batch Operations](./docs/cli/batch/)** - Bulk importing
4. **[Reports & Analytics](./docs/cli/reports/)** - Generate insights
5. **[Gantt Charts](./docs/cli/gantt/)** - Visualize workload
6. **[Project Management](./docs/cli/pm/)** - Executive dashboards

---

**You're all set! Start creating and managing tickets like a pro.** 🚀
