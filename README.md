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

**Option A: Using Project Key (explicit project)**
```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

**Option B: Using Ticket Key (auto-extract project)**
```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_TICKET=PROJ-123  # Project key extracted automatically
```

### Configuration File (~/.jirarc)

**Method 1: With explicit project key**
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

**Method 2: With ticket key (project auto-extracted)**
```yaml
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: ${JIRA_TOKEN}
  ticket: PROJ-123  # Project key (PROJ) extracted automatically

defaults:
  issue_type: Task
  priority: Medium
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

**Method B: Using Config File with Project Key**
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

**Method C: Using Ticket Key (Auto-Extract Project)**
Instead of specifying a project key, you can specify a ticket key and the project key will be automatically extracted:

```bash
# Environment variable
export JIRA_TICKET=PROJ-123

# Or in config file
cat > ~/.jirarc << EOF
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: your-api-token
  ticket: PROJ-123  # Project extracted as "PROJ"
EOF
```

This is useful when you want to work with a specific ticket rather than a default project.

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

### Project & Ticket Configuration

You can specify the project using either `--project` (project key) or `--ticket` (ticket key):

```bash
# Using project key
./jira-ticket-creator create --summary "Task" \
  --url https://company.atlassian.net \
  --email user@company.com \
  --token api-token \
  --project PROJ

# Using ticket key (project auto-extracted)
./jira-ticket-creator create --summary "Task" \
  --url https://company.atlassian.net \
  --email user@company.com \
  --token api-token \
  --ticket PROJ-123

# Using environment variables
export JIRA_URL=https://company.atlassian.net
export JIRA_EMAIL=user@company.com
export JIRA_TOKEN=api-token
export JIRA_TICKET=PROJ-123
./jira-ticket-creator create --summary "Task"

# Using config file
./jira-ticket-creator create --summary "Task"  # Uses ~/.jirarc with ticket or project
```

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

### Query (Advanced JQL Search)
Execute JQL queries with flexible output formatting and field selection:

```bash
# Basic query with table output
jira-ticket-creator query --jql "project = PROJ"

# Specify fields to display
jira-ticket-creator query --jql "project = PROJ" --fields "key,summary,status,priority"

# Export to CSV
jira-ticket-creator query --jql "status = 'In Progress'" --format csv --output active.csv

# Export to JSON
jira-ticket-creator query --jql "project = PROJ" --format json --output tickets.json

# Export to Markdown (for documentation)
jira-ticket-creator query --jql "type = Story" --format markdown --output stories.md

# Export to HTML (for browser viewing)
jira-ticket-creator query --jql "priority = Critical" --format html --output critical.html

# Fetch more results
jira-ticket-creator query --jql "project = PROJ" --max-results 500
```

**Flags:**
- `--jql <query>` (required) - JQL query string
- `--format <format>` - Output format: `table` (default), `json`, `csv`, `markdown`, `html`
- `--output <path>` - Write results to file (default: stdout)
- `--max-results <n>` - Maximum results to fetch (default: 50, max: 1000)
- `--fields <list>` - Comma-separated fields to display (default: key,type,summary,status,assignee,priority)

**Available Fields:** key, type, summary, status, assignee, priority, project, description

### Import (Bulk Import with Project Mapping)
Import existing JIRA tickets into local tracking with automatic project mapping:

```bash
# Basic import (assign all to one project)
jira-ticket-creator import --jql "project = PROJ" --map-project backend

# Dry run to preview
jira-ticket-creator import --jql "project = PROJ" --map-project backend --dry-run

# Map multiple prefixes to projects
jira-ticket-creator import --jql "project in (PROJ, BACK)" \
  --map-rule "PROJ->backend" --map-rule "BACK->backend"

# Update existing imported tickets
jira-ticket-creator import --jql "status = 'In Progress'" --update-existing

# Use custom mapping file
jira-ticket-creator import --jql "project = PROJ" \
  --map-project backend --mapping-path ~/my-mappings.json
```

**Flags:**
- `--jql <query>` (required) - JQL query to select tickets to import
- `--map-project <name>` - Logical project name to assign to all tickets
- `--map-rule <rule>` (repeatable) - Inline mapping rules (format: `PREFIX->project`)
- `--dry-run` - Preview import without saving
- `--update-existing` - Update already-imported tickets
- `--mapping-path <path>` - Project mapping file location (default: ~/.jira/project-mapping.json)

### Project Mapping Configuration
Store persistent project mappings for consistent ticket organization:

**Default location:** `~/.jira/project-mapping.json`

```json
{
  "mappings": {
    "backend": {
      "ticket_keys": ["PROJ", "BACK", "API"],
      "description": "Backend Team"
    },
    "frontend": {
      "ticket_keys": ["UI", "FRONT", "WEB"],
      "description": "Frontend Team"
    },
    "devops": {
      "ticket_keys": ["INFRA", "DEPLOY", "OPS"],
      "description": "DevOps Team"
    }
  }
}
```

The import command will automatically use these mappings to assign projects based on ticket key prefixes. Use `--map-rule` for one-off mappings or create the config file for persistent mappings.

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

# Filter by project
jira-ticket-creator team summary --project backend
jira-ticket-creator team summary --project frontend

# Filter by specific tickets
jira-ticket-creator team summary --ticket "PROJ-123,PROJ-456"
jira-ticket-creator team assignments --ticket "PROJ-123"

# Filter by creator/assignee
jira-ticket-creator team summary --creator "John Smith"
jira-ticket-creator team assignments --assignee "jane@company.com"

# Combine multiple filters
jira-ticket-creator team summary --project backend --creator "John Smith"
jira-ticket-creator team timeline --assignee "jane@company.com" --project frontend
```

**Team Report Flags (available on all subcommands):**
- `--project <name>` - Filter tickets by logical project name
- `--ticket <keys>` - Filter by ticket key(s), comma-separated (e.g., "PROJ-1,PROJ-2")
- `--creator <name>` - Filter by creator, comma-separated for multiple users
- `--assignee <name>` - Filter by assignee, comma-separated for multiple users

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

## 📚 Complete Command Reference

### Global Flags (Available on all commands)
- `--url <url>` - JIRA base URL (env: JIRA_URL)
- `--email <email>` - JIRA email address (env: JIRA_EMAIL)
- `--token <token>` - JIRA API token (env: JIRA_TOKEN)
- `--project <key>` - JIRA project key (env: JIRA_PROJECT)
- `--ticket <key>` - JIRA ticket key for auto-extracting project (env: JIRA_TICKET)
- `--config <path>` - Path to config file (default: ~/.jirarc)
- `--help` - Show help for command
- `--version` - Show version information

### create
Create new JIRA tickets

**Flags:**
- `--summary <text>` (required) - Ticket summary/title
- `--description <text>` - Detailed description
- `--type <type>` - Issue type (Task, Story, Bug, Epic, etc.)
- `--priority <level>` - Priority (Critical, High, Medium, Low)
- `--assignee <email>` - Assignee email address
- `--labels <labels>` - Comma-separated labels
- `--components <components>` - Comma-separated components
- `--blocked-by <keys>` - Comma-separated blocking ticket keys
- `--interactive` - Interactive prompt mode
- `--template <name>` - Use a template

### search
Search for JIRA tickets by key, summary, or JQL

**Flags:**
- `--key <key>` - Search by ticket key
- `--summary <text>` - Search by summary (partial match)
- `--jql <query>` - Raw JQL query
- `--format <format>` - Output format (table, json)

### query
Execute JQL queries with multiple output formats

**Flags:**
- `--jql <query>` (required) - JQL query string
- `--format <format>` - Output format: table, json, csv, markdown, html
- `--output <path>` - Write results to file
- `--max-results <n>` - Maximum results (default: 50, max: 1000)
- `--fields <list>` - Comma-separated fields to display

### import
Import existing tickets from JIRA with project mapping

**Flags:**
- `--jql <query>` (required) - JQL query for tickets to import
- `--map-project <name>` - Logical project name for all imported tickets
- `--map-rule <rule>` (repeatable) - Inline mapping rules (PREFIX->project)
- `--dry-run` - Preview without saving
- `--update-existing` - Update already-imported tickets
- `--mapping-path <path>` - Project mapping file location

### update
Update existing JIRA tickets

**Flags:**
- `--priority <level>` - New priority
- `--assignee <email>` - New assignee
- `--labels <labels>` - New labels
- `--summary <text>` - New summary
- `--description <text>` - New description

### transition
Move tickets through workflow states

**Flags:**
- `--to <state>` (required) - Target workflow state
- `--comment <text>` - Optional comment

### batch
Create multiple tickets from CSV or JSON file

**Flags:**
- `--input <file>` (required) - Input file (CSV or JSON)
- `--format <format>` - File format (auto-detected from extension)
- `--dry-run` - Validate without creating

### report
Generate ticket reports in multiple formats

**Flags:**
- `--format <format>` - Output format: table, json, csv, markdown, html
- `--output <path>` - Output file path
- `--filter <query>` - JQL filter for report

### team
Team-based reporting

**Subcommands:**
- `team summary [--project <name>]` - Tickets by creator
- `team assignments` - Workload assignments
- `team timeline` - Project timeline

### pm
Project management dashboard for executives

**Subcommands:**
- `pm dashboard` - Overall progress summary
- `pm hierarchy` - Parent-child ticket relationships
- `pm risk` - Risk assessment and bottlenecks
- `pm details` - Complete inventory
- `pm create-parent` - Guidance for parent tickets

### timeline
Project timeline visualization

**Flags:**
- `--format <format>` - Output format (ascii, mermaid, html)
- `--weeks <n>` - Number of weeks to display (default: 2)
- `--output <path>` - Output file path

### visualize
Dependency visualization

**Flags:**
- `--format <format>` - Output format (tree, mermaid, dot)
- `--output <path>` - Output file path

### template
Template management

**Subcommands:**
- `template list` - List available templates
- `template create --name <name>` - Create new template
- `template delete --name <name>` - Delete template

### completion
Generate shell completion scripts

**Subcommands:**
- `completion bash` - Bash completion
- `completion zsh` - Zsh completion
- `completion fish` - Fish completion
- `completion powershell` - PowerShell completion

## 📄 License

MIT License

```bash
# Import tickets by JQL search
  ./jira-ticket-creator import --jql "project = PROJ AND updated >= -30d" --map-project backend

  # Import all tickets from project
  ./jira-ticket-creator import --project PROJ --all --map-project backend

  # Import specific tickets
  ./jira-ticket-creator import --keys "PROJ-1,PROJ-2,PROJ-3" --map-project backend
```
