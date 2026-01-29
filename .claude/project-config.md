# JIRA Ticket Creator - Project Configuration

## Repository
GitHub: https://github.com/clintonsteiner/jira-ticket-creator

## Quick Start from GitHub

### Clone & Build
```bash
# Clone the repository
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator

# Build the CLI tool
go build -o jira-ticket-creator ./cmd/jira-ticket-creator

# Verify installation
./jira-ticket-creator --help
```

### Configure Credentials

**Option 1: Environment Variables (Recommended for automation)**
```bash
export JIRA_URL=https://your-domain.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

**Option 2: Config File (~/.jirarc)**
```yaml
jira:
  url: https://your-domain.atlassian.net
  email: your-email@company.com
  token: your-api-token
  project: PROJ

defaults:
  issue_type: Task
  priority: Medium
```

**Option 3: Using Ticket Key (auto-extract project)**
Instead of a project key, you can use a ticket key and the project will be automatically extracted:

```bash
export JIRA_URL=https://your-domain.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_TICKET=PROJ-123  # Project extracted as "PROJ"
```

Or in config file:
```yaml
jira:
  url: https://your-domain.atlassian.net
  email: your-email@company.com
  token: your-api-token
  ticket: PROJ-123  # Project key extracted automatically
```

**Option 4: Command-line Flags (Highest priority)**
```bash
# Using project key
./jira-ticket-creator create \
  --url https://your-domain.atlassian.net \
  --email your-email@company.com \
  --token your-api-token \
  --project PROJ \
  --summary "My ticket"

# Or using ticket key
./jira-ticket-creator create \
  --url https://your-domain.atlassian.net \
  --email your-email@company.com \
  --token your-api-token \
  --ticket PROJ-123 \
  --summary "My ticket"
```

## Key Development Files

### Core Logic
- `internal/config/config.go` - Configuration system with hierarchy
- `internal/jira/client.go` - HTTP client with exponential backoff retry
- `internal/jira/types.go` - API types and TicketRecord struct
- `internal/jira/validation.go` - Input validation functions
- `internal/storage/json.go` - Persistence layer

### CLI Commands
- `pkg/cli/commands/root.go` - Root command with persistent flags
- `pkg/cli/commands/create.go` - Create single tickets with interactive mode
- `pkg/cli/commands/batch.go` - Batch operations (CSV/JSON)
- `pkg/cli/commands/team.go` - Team reporting (summary, assignments, timeline)
- `pkg/cli/commands/timeline.go` - Project timeline visualization

### Tests (30%+ coverage)
- `internal/config/config_test.go`
- `internal/jira/validation_test.go`
- `internal/batch/*_test.go`
- `internal/storage/*_test.go`

## Build & Test Commands

```bash
# Build
go build -o jira-ticket-creator ./cmd/jira-ticket-creator

# Run all tests
go test ./...

# Run specific test
go test -run TestValidatePriority ./internal/jira -v

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Example Usage

### Create Ticket
```bash
./jira-ticket-creator create \
  --summary "Implement OAuth" \
  --description "Add OAuth 2.0 support" \
  --type Story \
  --priority High \
  --assignee john@company.com \
  --labels "auth,security" \
  --blocked-by "PROJ-1,PROJ-2"
```

### Batch Create from CSV
```bash
./jira-ticket-creator batch create --input tickets.csv --dry-run
./jira-ticket-creator batch create --input tickets.csv
```

### Team Reports
```bash
./jira-ticket-creator team summary        # Tickets by creator
./jira-ticket-creator team assignments    # Work allocation
./jira-ticket-creator team timeline       # Project timeline
```

### Timeline Visualization
```bash
# Three output formats:
./jira-ticket-creator timeline --format ascii --weeks 2    # Terminal output
./jira-ticket-creator timeline --format mermaid --weeks 2  # Gantt chart for Markdown
./jira-ticket-creator timeline --format html --weeks 2 > timeline.html  # Standalone HTML
```

### Search & Report
```bash
./jira-ticket-creator search --key PROJ-123
./jira-ticket-creator search --jql "project = PROJ AND status = 'In Progress'"
./jira-ticket-creator report --format csv --output report.csv
```

## Architecture Notes

- **Config Hierarchy**: Flags > Environment > Config File > Defaults
- **HTTP Client**: Retries on 429/5xx with exponential backoff (1s, 2s, 4s)
- **Storage**: JSON-based, extensible to databases
- **Batch Processing**: 3-phase pipeline (validation → creation → linking)
- **CLI Framework**: Cobra with Viper for configuration

## Common Issues

**Build Error**: `go mod tidy` then `go build`
**Auth Error**: Check JIRA_TOKEN is API token, not password
**Test Failure**: Run `go test -v` for details; ensure Go 1.19+

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Config management
- `github.com/manifoldco/promptui` - Interactive mode
- `github.com/olekukonko/tablewriter` - Table formatting
