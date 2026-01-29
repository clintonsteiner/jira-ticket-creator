---
layout: default
title: Import Command
parent: CLI Commands
nav_order: 4
has_toc: true
---

# Import Command

Import existing JIRA tickets into local tracking with automatic project mapping.

## Basic Usage

```bash
jira-ticket-creator import --jql "project = PROJ" --map-project backend
```

## Flags

| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `--jql` | string | **Required.** JQL query for tickets to import | `--jql "project = PROJ"` |
| `--map-project` | string | Logical project name for all imported tickets | `--map-project backend` |
| `--map-rule` | strings | Inline mapping rules (repeatable) | `--map-rule "PROJ->backend"` |
| `--dry-run` | boolean | Preview import without saving | `--dry-run` |
| `--update-existing` | boolean | Update already-imported tickets | `--update-existing` |
| `--mapping-path` | string | Project mapping file path | `--mapping-path ~/.jira/mappings.json` |

## Examples

### Basic Import
```bash
# Import all tickets from PROJ project
jira-ticket-creator import --jql "project = PROJ" --map-project backend

# Import in-progress tickets
jira-ticket-creator import --jql "status = 'In Progress'" --map-project backend
```

### Dry Run Preview

```bash
# See what would be imported without saving
jira-ticket-creator import --jql "project = PROJ" \
  --map-project backend \
  --dry-run
```

Output example:
```
Found 15 issue(s) to import

Import Summary:
Total: 15
backend: 15

Dry run completed (no changes saved)
```

### Multiple Project Mapping

```bash
# Map different prefixes to projects
jira-ticket-creator import --jql "project in (PROJ, API, DB)" \
  --map-rule "PROJ->backend" \
  --map-rule "API->backend" \
  --map-rule "DB->devops"
```

### Update Existing Tickets

```bash
# Refresh imported tickets with latest data
jira-ticket-creator import --jql "status = 'In Progress'" \
  --update-existing
```

## Project Mapping

### Using Persistent Mapping File

Create `~/.jira/project-mapping.json`:

```json
{
  "mappings": {
    "backend": {
      "ticket_keys": ["PROJ", "API", "DB"],
      "description": "Backend Team"
    },
    "frontend": {
      "ticket_keys": ["UI", "WEB", "MOBILE"],
      "description": "Frontend Team"
    },
    "devops": {
      "ticket_keys": ["INFRA", "CI", "DEPLOY"],
      "description": "DevOps Team"
    }
  }
}
```

Then import automatically uses these mappings:

```bash
# Will use PROJ->backend mapping from file
jira-ticket-creator import --jql "project = PROJ"
```

### Inline Mapping Rules

For one-off imports, use inline rules:

```bash
jira-ticket-creator import --jql "project = PROJ" \
  --map-rule "PROJ->backend" \
  --map-rule "PROJ->api-team"
```

### No Mapping

If no mapping is specified:

```bash
jira-ticket-creator import --jql "project = PROJ" --map-project backend
# All tickets assigned to 'backend' project
```

## Workflow Examples

### Import Project for Team
```bash
# Step 1: Create mapping file
cat > ~/.jira/project-mapping.json << EOF
{
  "mappings": {
    "team-a": {
      "ticket_keys": ["PROJ"],
      "description": "Team A"
    }
  }
}
EOF

# Step 2: Import all team tickets
jira-ticket-creator import --jql "project = PROJ"

# Step 3: View team report
jira-ticket-creator team summary --project team-a

# Step 4: Generate Gantt chart
jira-ticket-creator gantt --format html --output team-workload.html
```

### Multi-Team Setup
```bash
# Create comprehensive mapping
cat > ~/.jira/project-mapping.json << EOF
{
  "mappings": {
    "backend": {
      "ticket_keys": ["PROJ", "API", "DB"],
      "description": "Backend Team"
    },
    "frontend": {
      "ticket_keys": ["UI", "WEB"],
      "description": "Frontend Team"
    }
  }
}
EOF

# Import from each project
jira-ticket-creator import --jql "project = PROJ" --map-project backend
jira-ticket-creator import --jql "project = UI" --map-project frontend

# View by team
jira-ticket-creator team summary --project backend
jira-ticket-creator team summary --project frontend

# Generate reports
jira-ticket-creator gantt --format mermaid --output gantt.md
jira-ticket-creator pm dashboard
```

### Regular Sync
```bash
# Update all imported tickets weekly
jira-ticket-creator import --jql "updated >= -7d" --update-existing

# Generate fresh report
jira-ticket-creator report --format html --output weekly-report.html
```

### Map Parent Tickets with Children

Map an epic or feature so all child tickets inherit the same project mapping:

```bash
# Import parent epic
jira-ticket-creator import --jql "key = PROJ-500" --map-project feature-checkout

# All subtasks and child stories now map to "feature-checkout"
# Verify with:
jira-ticket-creator search --jql "parent = PROJ-500"

# View all work for this feature
jira-ticket-creator gantt --format html --output feature-checkout.html
```

**Example: Product Feature Breakdown**

```bash
# Mobile Redesign Feature
jira-ticket-creator import --jql "key = PROJ-1000" --map-project mobile-redesign
# Children: Design, Frontend, Backend, QA all map to "mobile-redesign"

# Payment Integration Feature
jira-ticket-creator import --jql "key = PROJ-1001" --map-project payment-integration
# Children: API, UI, Security, Testing all map to "payment-integration"

# View separate work efforts
jira-ticket-creator team summary --project mobile-redesign
jira-ticket-creator team summary --project payment-integration
```

**Example: Sprint-Based Mapping**

```bash
# Map Q1 sprint epic
jira-ticket-creator import --jql "key = PROJ-2000" --map-project q1-sprint

# Map Q2 sprint epic
jira-ticket-creator import --jql "key = PROJ-2001" --map-project q2-sprint

# Compare workload across sprints
jira-ticket-creator gantt --format html --output q1-q2-comparison.html
```

## Output

Successful import:

```
Found 25 issue(s) to import

Import Summary:
Total: 25
backend: 15
frontend: 10

Import completed
Added/Updated: 20
Skipped (existing): 5
```

## Integration with Other Commands

After importing, use other commands for analysis:

```bash
# Search imported tickets
jira-ticket-creator search --jql "project = PROJ"

# Query specific fields
jira-ticket-creator query --jql "project = PROJ AND status = 'In Progress'" \
 --format json

# Generate reports
jira-ticket-creator report --format html

# View team workload
jira-ticket-creator team summary --project backend

# Create Gantt chart
jira-ticket-creator gantt --format mermaid
```

## Tips & Tricks

1. **Import specific status**
 ```bash
 jira-ticket-creator import --jql "status = 'In Progress'" \
 --map-project active
 ```

2. **Import recent changes**
 ```bash
 jira-ticket-creator import --jql "updated >= -24h" \
 --update-existing
 ```

3. **Import by assignee**
 ```bash
 jira-ticket-creator import --jql "assignee = john" \
 --map-project johns-work
 ```

4. **Schedule regular imports**
 ```bash
 # Add to crontab for daily sync
 0 9 * * * /path/to/jira-ticket-creator import --jql "updated >= -1d" --update-existing
 ```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `jql is required` | Add `--jql "your query"` flag |
| `no issues found` | Check JQL query is valid |
| `mapping file not found` | Run without `--mapping-path` or create file at `~/.jira/project-mapping.json` |
| `failed to add ticket` | Check storage directory permissions |

## See Also

- [Query Command](query) - Execute JQL queries with output formatting
- [Project Mapping](../advanced/project-mapping) - Configure project mappings
