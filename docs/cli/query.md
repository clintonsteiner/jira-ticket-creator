---
layout: default
title: Query Command
parent: CLI Commands
nav_order: 3
has_toc: true
---

# Query Command

Execute JQL (JIRA Query Language) queries with flexible output formatting.

## Basic Usage

```bash
jira-ticket-creator query --jql "project = PROJ"
```

## Flags

| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `--jql` | string | **Required.** JQL query string | `--jql "project = PROJ"` |
| `--format` | string | Output format: table, json, csv, markdown, html | `--format json` |
| `--output` | string | Output file path (default: stdout) | `--output results.csv` |
| `--max-results` | int | Maximum results to fetch (default: 50, max: 1000) | `--max-results 500` |
| `--fields` | string | Comma-separated fields to display | `--fields "key,summary,status,priority"` |

## Available Fields

- `key` - Ticket key (e.g., PROJ-123)
- `type` - Issue type
- `summary` - Ticket summary
- `status` - Current status
- `assignee` - Assigned user
- `priority` - Priority level
- `project` - Project key
- `description` - Full description

## Examples

### Basic Query
```bash
# Get all tickets in PROJ
jira-ticket-creator query --jql "project = PROJ"

# Get in-progress tickets
jira-ticket-creator query --jql "status = 'In Progress'"

# Get high-priority bugs
jira-ticket-creator query --jql "type = Bug AND priority = High"
```

### Field Selection
```bash
# Display specific fields
jira-ticket-creator query --jql "project = PROJ" \
 --fields "key,summary,status,assignee"

# Minimal output
jira-ticket-creator query --jql "project = PROJ" \
 --fields "key,summary"
```

### Export to Different Formats

**JSON**
```bash
jira-ticket-creator query --jql "project = PROJ" \
 --format json \
 --output tickets.json
```

**CSV**
```bash
jira-ticket-creator query --jql "status = 'In Progress'" \
 --format csv \
 --output active.csv
```

**Markdown**
```bash
jira-ticket-creator query --jql "priority = Critical" \
 --format markdown \
 --output critical.md
```

**HTML**
```bash
jira-ticket-creator query --jql "project = PROJ" \
 --format html \
 --output tickets.html
# Open in browser
open tickets.html
```

### Large Result Sets
```bash
# Fetch many results
jira-ticket-creator query --jql "project = PROJ" \
 --max-results 500 \
 --format json \
 --output all-tickets.json
```

### Complex Queries

```bash
# Updated in last 30 days
jira-ticket-creator query --jql "updated >= -30d"

# Created by specific user
jira-ticket-creator query --jql "creator = john"

# Assigned to team
jira-ticket-creator query --jql "assignee in (john, jane, bob)"

# Complex filter
jira-ticket-creator query --jql "project = PROJ AND type = Story AND status != Done AND priority >= High"
```

## JQL Reference

Common JQL operators:

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equals | `status = Done` |
| `!=` | Not equals | `status != Done` |
| `>`, `>=`, `<`, `<=` | Comparison | `priority >= High` |
| `in` | In list | `status in (Done, Closed)` |
| `AND` | Both conditions | `project = PROJ AND type = Bug` |
| `OR` | Either condition | `status = Done OR status = Closed` |
| `NOT` | Negate condition | `NOT status = Done` |

Date operators:

| Operator | Description |
|----------|-------------|
| `-Xd` | X days ago |
| `-Xw` | X weeks ago |
| `-Xm` | X months ago |
| `-Xy` | X years ago |

## Output Examples

### Table Format (default)
```
KEY TYPE SUMMARY STATUS ASSIGNEE PRIORITY
------ ---- -------- ------ --------- --------
PROJ-101 Story Implement OAuth In Progress john@co.com High
PROJ-102 Bug Fix login issue To Do jane@co.com Critical
PROJ-103 Task Update documentation Done Medium
```

### CSV Format
```csv
key,type,summary,status,assignee,priority
PROJ-101,Story,Implement OAuth,In Progress,john@co.com,High
PROJ-102,Bug,Fix login issue,To Do,jane@co.com,Critical
PROJ-103,Task,Update documentation,Done,,Medium
```

### Markdown Format
```markdown
| key | type | summary | status | assignee | priority |
| --- | --- | --- | --- | --- | --- |
| PROJ-101 | Story | Implement OAuth | In Progress | john@co.com | High |
| PROJ-102 | Bug | Fix login issue | To Do | jane@co.com | Critical |
```

### JSON Format
```json
[
 {
 "key": "PROJ-101",
 "fields": {
 "summary": "Implement OAuth",
 "status": {"name": "In Progress"},
 "assignee": {"name": "john@co.com"},
 "priority": {"name": "High"}
 }
 }
]
```

## Tips & Tricks

1. **Save queries for reuse**
 ```bash
 alias my-backlog='jira-ticket-creator query --jql "project = PROJ AND status = \"To Do\""'
 my-backlog
 ```

2. **Pipe to other tools**
 ```bash
 jira-ticket-creator query --jql "project = PROJ" --format csv | column -t -s,
 ```

3. **Create dashboards**
 ```bash
 jira-ticket-creator query --jql "assignee = currentUser()" \
 --format html \
 --output my-work.html
 ```

4. **Monitor critical items**
 ```bash
 watch -n 300 "jira-ticket-creator query --jql 'priority = Critical' --format table"
 ```

## See Also

- [Search Command](search) - Simple JIRA queries
- [Import Command](import) - Import tickets locally
