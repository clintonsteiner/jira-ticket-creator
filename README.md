# JIRA Ticket Creator CLI

A powerful, feature-rich command-line tool for managing JIRA tickets with advanced features like batch operations, dependency visualization, multiple reporting formats, and team collaboration tools.

## ✨ Features

### Core Features (13 Major Enhancements)
1. **Better Credentials Handling** - Config hierarchy: Flags > Env Vars > Config File > Defaults
2. **Refactored Architecture** - Clean package structure with proper separation of concerns
3. **Better Error Handling** - Custom error types with retry logic and user-friendly messages
4. **Richer Ticket Configuration** - Priority, assignee, labels, components, and issue types
5. **Ticket Updates & Transitions** - Update fields and move tickets through workflow
6. **Validation Before Creation** - Pre-flight checks to prevent invalid tickets
7. **Batch Ticket Creation** - Create multiple tickets from CSV/JSON files
8. **Search/Query Capability** - Search by key, summary, or JQL queries
9. **Enhanced Reporting** - Multiple formats: table, JSON, CSV, Markdown, HTML
10. **Dependency Visualization** - ASCII tree, Mermaid diagrams, DOT (Graphviz)
11. **Interactive Mode** - Guided ticket creation with prompts
12. **Templates/Workflows** - Built-in and custom ticket templates
13. **Shell Completions** - Bash, Zsh, Fish, PowerShell support

### Additional Features
- **Team Collaboration** - Track creator, assignee, and timeline information
- **Project Timeline** - View progress bars and estimated completion dates
- **Workload Management** - See who is assigned what and critical items

## 📦 Installation

### Prerequisites
- Go 1.21 or later
- JIRA Cloud account with API token

### Build from Source
```bash
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

### Move to PATH
```bash
sudo mv jira-ticket-creator /usr/local/bin/
```

## ⚙️ Configuration

### Configuration Hierarchy (Highest to Lowest Priority)
1. Command-line flags (e.g., `--url`, `--email`)
2. Environment variables (e.g., `JIRA_URL`, `JIRA_EMAIL`)
3. Configuration file (`~/.jirarc`)
4. Default values

### Environment Variables
```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

### Configuration File (~/.jirarc)
```yaml
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: ${JIRA_TOKEN}
  project: PROJ

defaults:
  issue_type: Task
  priority: Medium
  labels:
    - auto-created
```

## 📖 Command Reference

### Create Tickets
```bash
# Basic
jira-ticket-creator create --summary "New task"

# Advanced
jira-ticket-creator create \
  --summary "OAuth 2.0" \
  --type Story \
  --priority High \
  --assignee john@company.com \
  --labels "auth,security"

# Interactive
jira-ticket-creator create --interactive
```

### Batch Operations
```bash
# From CSV
jira-ticket-creator batch create --input tickets.csv

# From JSON
jira-ticket-creator batch create --input tickets.json --format json

# Validation only
jira-ticket-creator batch create --input tickets.csv --dry-run
```

### Update & Transition
```bash
jira-ticket-creator update PROJ-123 --priority Critical
jira-ticket-creator transition PROJ-123 --to "In Progress"
```

### Search
```bash
jira-ticket-creator search --key PROJ-123
jira-ticket-creator search --summary "login"
jira-ticket-creator search --jql "project = PROJ AND status = 'To Do'"
```

### Reports
```bash
jira-ticket-creator report --format markdown --output report.md
jira-ticket-creator report --format html --output report.html
jira-ticket-creator report --format csv --output report.csv
```

### Team Reports
```bash
jira-ticket-creator team summary      # By creator
jira-ticket-creator team assignments  # Workload
jira-ticket-creator team timeline     # Progress
```

### Visualization
```bash
jira-ticket-creator visualize --format tree
jira-ticket-creator visualize --format mermaid --output diagram.md
jira-ticket-creator visualize --format dot --output diagram.dot
```

### Templates
```bash
jira-ticket-creator template list
jira-ticket-creator create --template bug
```

### Shell Completions
```bash
# Bash
jira-ticket-creator completion bash | sudo tee /etc/bash_completion.d/jira-ticket-creator

# Zsh
jira-ticket-creator completion zsh > "${fpath[1]}/_jira-ticket-creator"
```

## 📋 Batch File Formats

### CSV Format
```csv
summary,description,issue_type,priority,assignee,labels,components,blocked_by
"Task 1","Description",Task,High,"user@email.com","label1,label2","Component1","PROJ-100"
```

### JSON Format
```json
[
  {
    "summary": "Task 1",
    "description": "Description",
    "issue_type": "Task",
    "priority": "High",
    "assignee": "user@email.com",
    "labels": ["label1", "label2"],
    "blocked_by": ["PROJ-100"]
  }
]
```

## 🏗️ Architecture

- **internal/config/** - Configuration management (Viper + Cobra)
- **internal/jira/** - JIRA API client with retry logic
- **internal/storage/** - JSON-based ticket storage
- **internal/batch/** - CSV/JSON parsers and batch processor
- **internal/reports/** - Multiple report formats
- **internal/templates/** - Template engine
- **internal/interactive/** - Interactive prompt system
- **pkg/cli/commands/** - CLI command implementations

## 🧪 Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📄 License

MIT License

