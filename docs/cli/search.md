---
layout: default
title: Search Command
parent: CLI Commands
nav_order: 2
has_toc: true
---

# Search Command

Search for JIRA tickets by key, summary, or JQL query.

## Basic Usage

```bash
jira-ticket-creator search --key PROJ-123
jira-ticket-creator search --summary "login"
jira-ticket-creator search --jql "project = PROJ AND status = 'To Do'"
```

## Flags

| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `--key` | string | Search by ticket key | `--key PROJ-123` |
| `--summary` | string | Search by summary (partial match) | `--summary "authentication"` |
| `--jql` | string | Search using JQL query | `--jql "project = PROJ"` |
| `--format` | string | Output format: table, json | `--format json` |

## Examples

### Search by Key

```bash
# Find a specific ticket
jira-ticket-creator search --key PROJ-123
```

Output:
```
KEY TYPE SUMMARY STATUS
--- ---- ------- ------
PROJ-123 Story Implement OAuth In Progress
```

### Search by Summary (Partial Match)

```bash
# Find tickets mentioning "authentication"
jira-ticket-creator search --summary "authentication"

# Find tickets mentioning "database"
jira-ticket-creator search --summary "database"

# Find tickets mentioning "api"
jira-ticket-creator search --summary "api"
```

### JQL Queries

#### Status-Based

```bash
# Find all to-do items
jira-ticket-creator search --jql "status = 'To Do'"

# Find in-progress items
jira-ticket-creator search --jql "status = 'In Progress'"

# Find completed items
jira-ticket-creator search --jql "status in (Done, Closed)"

# Find items not done
jira-ticket-creator search --jql "status != Done"
```

#### Project-Based

```bash
# Find all tickets in PROJ
jira-ticket-creator search --jql "project = PROJ"

# Find tickets in multiple projects
jira-ticket-creator search --jql "project in (PROJ, BACKEND, API)"

# Find all tickets not in PROJ
jira-ticket-creator search --jql "project != PROJ"
```

#### Type-Based

```bash
# Find all bugs
jira-ticket-creator search --jql "type = Bug"

# Find stories and tasks
jira-ticket-creator search --jql "type in (Story, Task)"

# Find everything except bugs
jira-ticket-creator search --jql "type != Bug"
```

#### Priority-Based

```bash
# Find critical items
jira-ticket-creator search --jql "priority = Critical"

# Find high and critical
jira-ticket-creator search --jql "priority >= High"

# Find non-critical
jira-ticket-creator search --jql "priority < Critical"

# Find items with specific priority
jira-ticket-creator search --jql "priority in (Critical, High)"
```

#### Assignee-Based

```bash
# Find tickets assigned to john
jira-ticket-creator search --jql "assignee = john"

# Find tickets assigned to specific people
jira-ticket-creator search --jql "assignee in (john, jane, bob)"

# Find unassigned tickets
jira-ticket-creator search --jql "assignee is EMPTY"

# Find assigned tickets
jira-ticket-creator search --jql "assignee is not EMPTY"

# Find tickets assigned to current user
jira-ticket-creator search --jql "assignee = currentUser()"
```

#### Date-Based

```bash
# Find tickets created in last 7 days
jira-ticket-creator search --jql "created >= -7d"

# Find tickets updated in last 24 hours
jira-ticket-creator search --jql "updated >= -1d"

# Find tickets created last month
jira-ticket-creator search --jql "created >= -30d"

# Find tickets due soon
jira-ticket-creator search --jql "duedate <= 3d"

# Find overdue tickets
jira-ticket-creator search --jql "duedate < now()"
```

### Complex Queries

#### Combination Filters

```bash
# High priority bugs created recently
jira-ticket-creator search --jql "type = Bug AND priority = High AND created >= -7d"

# In-progress stories in PROJ
jira-ticket-creator search --jql "project = PROJ AND type = Story AND status = 'In Progress'"

# Critical issues not assigned
jira-ticket-creator search --jql "priority = Critical AND assignee is EMPTY"

# Backend work for john
jira-ticket-creator search --jql "project = BACKEND AND assignee = john AND status != Done"

# Ready for review
jira-ticket-creator search --jql "status = 'Review Ready' OR labels = ready-review"
```

#### Using AND/OR

```bash
# Multiple assignees OR creators
jira-ticket-creator search --jql "(assignee = john OR assignee = jane) AND status != Done"

# Multiple projects with high priority
jira-ticket-creator search --jql "(project = PROJ OR project = BACKEND) AND priority >= High"

# Complex logic
jira-ticket-creator search --jql "(type = Bug AND priority = Critical) OR (status = 'Blocked')"
```

#### Text Search

```bash
# Find tickets mentioning OAuth in summary or description
jira-ticket-creator search --jql "text ~ 'OAuth'"

# Exact phrase search
jira-ticket-creator search --jql "text ~ '\"OAuth 2.0\"'"

# Multiple terms
jira-ticket-creator search --jql "summary ~ 'database' OR description ~ 'database'"
```

### Format Examples

#### Table Format (default)

```bash
jira-ticket-creator search --jql "project = PROJ" --format table
```

Output:
```
KEY TYPE SUMMARY STATUS
--- ---- ------- ------
PROJ-100 Story Implement OAuth In Progress
PROJ-101 Task Write documentation To Do
PROJ-102 Bug Fix login issue To Do
PROJ-103 Story Add dark mode Done
```

#### JSON Format

```bash
jira-ticket-creator search --jql "project = PROJ" --format json
```

Output:
```json
[
 {
 "key": "PROJ-100",
 "fields": {
 "summary": "Implement OAuth",
 "issuetype": {"name": "Story"},
 "status": {"name": "In Progress"},
 "priority": {"name": "High"},
 "assignee": {"name": "john"}
 }
 }
]
```

## Real-World Workflows

### Daily Standup

```bash
#!/bin/bash
# daily-standup.sh

echo "=== MY WORK (ASSIGNED TO ME) ==="
jira-ticket-creator search --jql "assignee = currentUser() AND status != Done"

echo ""
echo "=== CREATED BY ME (IN PROGRESS) ==="
jira-ticket-creator search --jql "creator = currentUser() AND status = 'In Progress'"

echo ""
echo "=== HIGH PRIORITY ISSUES ==="
jira-ticket-creator search --jql "priority = Critical AND status != Done"
```

### Sprint Planning

```bash
#!/bin/bash
# Get sprint candidates

echo "=== READY FOR SPRINT ==="
jira-ticket-creator search --jql "status = 'Ready for Sprint' OR labels = ready-sprint"

echo ""
echo "=== TECH DEBT ==="
jira-ticket-creator search --jql "labels = tech-debt AND status != Done"

echo ""
echo "=== BUG FIXES ==="
jira-ticket-creator search --jql "type = Bug AND status = 'To Do' AND priority >= High"
```

### Issue Triage

```bash
#!/bin/bash
# Find items needing triage

echo "=== UNASSIGNED HIGH PRIORITY ==="
jira-ticket-creator search --jql "assignee is EMPTY AND priority >= High"

echo ""
echo "=== UNASSIGNED BUGS ==="
jira-ticket-creator search --jql "type = Bug AND assignee is EMPTY"

echo ""
echo "=== BLOCKED ITEMS ==="
jira-ticket-creator search --jql "status = Blocked"
```

### Monitoring

```bash
#!/bin/bash
# Monitor critical work

# Watch for new critical issues (every 5 minutes)
watch -n 300 'jira-ticket-creator search --jql "priority = Critical AND created >= -1d"'

# Check for overdue items daily
0 9 * * * jira-ticket-creator search --jql "duedate < now()" > ~/critical-overdue.txt
```

## Search Tips

### 1. Save Searches as Aliases

```bash
# Add to ~/.bashrc or ~/.zshrc
alias my-work='jira-ticket-creator search --jql "assignee = currentUser()"'
alias in-progress='jira-ticket-creator search --jql "status = \"In Progress\""'
alias my-backlog='jira-ticket-creator search --jql "assignee = currentUser() AND status = \"To Do\""'
alias critical='jira-ticket-creator search --jql "priority = Critical"'

# Use them
my-work
critical
in-progress
```

### 2. Pipe to Other Tools

```bash
# Count tickets
jira-ticket-creator search --jql "project = PROJ" --format table | wc -l

# Filter results
jira-ticket-creator search --jql "project = PROJ" | grep "High"

# Export to CSV
jira-ticket-creator search --jql "project = PROJ" --format json | jq -r '.[] | [.key, .fields.summary] | @csv'
```

### 3. Monitor Changes

```bash
# Watch for changes (update every 60 seconds)
watch 'jira-ticket-creator search --jql "status = \"In Progress\""'

# Or in tmux
tmux new-session -d -s jira 'watch -n 60 "jira-ticket-creator search --jql \"status = In Progress\""'
```

### 4. Create Smart Dashboards

```bash
#!/bin/bash
# dashboard.sh - Monitor key metrics

while true; do
 clear
 echo "=== JIRA DASHBOARD ==="
 echo "Last updated: $(date)"
 echo ""

 echo "My Work ($(jira-ticket-creator search --jql 'assignee = currentUser() AND status != Done' | wc -l) items)"
 jira-ticket-creator search --jql "assignee = currentUser() AND status != Done"

 echo ""
 echo "Critical Issues ($(jira-ticket-creator search --jql 'priority = Critical' | wc -l) items)"
 jira-ticket-creator search --jql "priority = Critical"

 sleep 300 # Update every 5 minutes
done
```

## Common JQL Functions

| Function | Description | Example |
|----------|-------------|---------|
| `currentUser()` | Current logged-in user | `assignee = currentUser()` |
| `now()` | Current date/time | `duedate < now()` |
| `startOfDay()` | Start of today | `created >= startOfDay()` |
| `endOfDay()` | End of today | `created <= endOfDay()` |

## See Also

- [Query Command](query) - For flexible output formatting and exports
- [Import Command](import) - For importing JIRA tickets locally
