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

### Quick Start - Clone from GitHub
```bash
# Clone the repository
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator

# Build the binary
go build -o jira-ticket-creator ./cmd/jira-ticket-creator

# Test it works
./jira-ticket-creator --help
```

### Run Directly from GitHub (No cloning required)
```bash
# Run directly from GitHub repository
go run github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest create --help

# Or compile and install to $GOPATH/bin
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

### Install Globally (Optional)
After building, move the binary to your PATH:
```bash
sudo mv jira-ticket-creator /usr/local/bin/

# Or add current directory to PATH
export PATH=$PATH:$(pwd)

# Then use from anywhere
jira-ticket-creator --help
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

## 🚀 Getting Started

### 1. Setup (Choose One Method)

**Method A: Using Environment Variables (Recommended)**
```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

**Method B: Using Config File**
```bash
# Create ~/.jirarc with your credentials
cat > ~/.jirarc << EOF
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: your-api-token
  project: PROJ
EOF

chmod 600 ~/.jirarc  # Secure the file (readable only by you)
```

### 2. Test Your Setup
```bash
./jira-ticket-creator search --key PROJ-1
```

### 3. Start Creating Tickets
```bash
# Simple ticket
./jira-ticket-creator create --summary "My first ticket"

# Detailed ticket
./jira-ticket-creator create \
  --summary "Add OAuth" \
  --description "Implement OAuth 2.0" \
  --type Story \
  --priority High

# Interactive mode (recommended for first-time users)
./jira-ticket-creator create --interactive
```

### 4. Track Related Work with Parent Tickets (Optional)
For projects with multiple related tickets, use parent epics to organize work:

```bash
# Step 1: Create a parent epic
jira-ticket-creator create \
  --summary "Q1 Platform Upgrade" \
  --type Epic \
  --priority High \
  --description "Parent epic for all platform upgrade work"
# Returns: PROJ-1

# Step 2: Create child tickets that link to the parent
jira-ticket-creator create \
  --summary "Upgrade Database" \
  --type Task \
  --priority High \
  --assignee team-member@company.com \
  --blocked-by PROJ-1

# Step 3: View the project hierarchy and progress
jira-ticket-creator pm hierarchy    # See parent-child relationships
jira-ticket-creator pm dashboard    # See overall progress for your boss
jira-ticket-creator pm risk         # Identify bottlenecks
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

### Project Timeline (2-Week Visualization)
```bash
# ASCII text timeline for terminal
jira-ticket-creator timeline --format ascii --weeks 2

# Mermaid Gantt chart (embed in GitHub, Notion, Markdown)
jira-ticket-creator timeline --format mermaid --weeks 2

# HTML interactive timeline (open in browser)
jira-ticket-creator timeline --format html --weeks 2 > timeline.html
open timeline.html
```

**Features:**
- Shows ticket progress bars over configurable weeks
- Displays status: In Progress (████), Not Started (░░░░), Done
- Highlights critical path items
- Includes progress summary and recommendations
- HTML version with CSS styling and interactive cells

### Project Management Dashboard (For Your Boss)
```bash
# Executive summary - overall progress, team workload, priorities
jira-ticket-creator pm dashboard

# Ticket hierarchy - parent epics with child tasks
jira-ticket-creator pm hierarchy

# Risk assessment - blocked items, bottlenecks, recommendations
jira-ticket-creator pm risk

# Detailed inventory - complete table of all tickets
jira-ticket-creator pm details

# Guidance for creating parent tickets
jira-ticket-creator pm create-parent
```

**What the Dashboard Shows:**
- 📊 Overall completion percentage and burndown
- 🎯 Ticket count by status (Done, In Progress, To Do)
- 🔴 Priority distribution (Critical, High, Medium, Low)
- 👥 Team workload allocation
- ⛔ Blocked items and dependencies
- 🚨 Risk assessment with recommendations
- 📌 Epic/Story hierarchy with subtasks
- 🔗 Blocking relationships and critical path

**Example Output:**
```
╔════════════════════════════════════════════════════════════════════╗
║              PROJECT MANAGEMENT EXECUTIVE SUMMARY                  ║
╚════════════════════════════════════════════════════════════════════╝

📊 Total Tickets: 7
📈 Overall Progress: 14% (1/7 tickets)

Status Breakdown:
  Done: 1 (14%)
  In Progress: 2 (28%)
  To Do: 4 (57%)

Team Workload Distribution:
  bob@company.com: 2 tickets
  charlie@company.com: 2 tickets
  [Unassigned]: 1 tickets

⚠️  Critical Items:
  Blocked tickets: 5 (dependencies exist)
  High-priority blocked items: 4
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

## 🐙 GitHub Repository

**Repository**: https://github.com/clintonsteiner/jira-ticket-creator

### Clone & Use
```bash
# Clone from GitHub
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator

# Build locally
go build -o jira-ticket-creator ./cmd/jira-ticket-creator

# Run without cloning
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

### Contribute
- Report issues: https://github.com/clintonsteiner/jira-ticket-creator/issues
- Submit pull requests with improvements
- Star ⭐ the repo if you find it useful!

## 🆘 Troubleshooting

### Authentication Issues
```bash
# Test your credentials
./jira-ticket-creator search --key PROJ-1

# Common issues:
# - JIRA_TOKEN must be API token, not your password
# - User must have permissions on the project
# - Verify JIRA_URL format: https://domain.atlassian.net (no trailing slash)
```

### Build Issues
```bash
# Update Go modules
go mod tidy
go mod download

# Rebuild
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

### Command Not Found
```bash
# Make sure binary is in PATH
export PATH=$PATH:$(pwd)
./jira-ticket-creator --help

# Or move to /usr/local/bin
sudo mv jira-ticket-creator /usr/local/bin/
```

## 📄 License

MIT License
