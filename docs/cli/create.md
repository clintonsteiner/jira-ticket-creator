---
layout: default
title: Create Command
parent: CLI Commands
nav_order: 1
---

# Create Command

Create new JIRA tickets with full field control.

## Basic Usage

```bash
jira-ticket-creator create --summary "My new ticket"
```

## Flags

| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `--summary` | string | **Required.** Ticket summary/title | `--summary "Implement OAuth"` |
| `--description` | string | Detailed description | `--description "Add OAuth 2.0 support"` |
| `--type` | string | Issue type (Task, Story, Bug, Epic) | `--type Story` |
| `--priority` | string | Priority (Critical, High, Medium, Low) | `--priority High` |
| `--assignee` | string | Assignee email address | `--assignee john@company.com` |
| `--labels` | string | Comma-separated labels | `--labels "auth,security,backend"` |
| `--components` | string | Comma-separated components | `--components "API,Auth"` |
| `--blocked-by` | string | Comma-separated blocking ticket keys | `--blocked-by "PROJ-100,PROJ-101"` |
| `--interactive` | boolean | Interactive prompt mode | `--interactive` |
| `--template` | string | Use a template | `--template bug` |

## Examples

### Simple Ticket
```bash
jira-ticket-creator create --summary "Fix login bug"
```

### Detailed Ticket
```bash
jira-ticket-creator create \
 --summary "Implement OAuth 2.0" \
 --description "Add OAuth authentication to API" \
 --type Story \
 --priority High \
 --assignee john@company.com \
 --labels "auth,security" \
 --components "API"
```

### Ticket with Dependencies
```bash
jira-ticket-creator create \
 --summary "Update API docs" \
 --priority Medium \
 --blocked-by "PROJ-100"
```

### Interactive Mode
```bash
jira-ticket-creator create --interactive
```

This will prompt you for:
- Summary
- Description
- Issue type
- Priority
- Assignee
- Labels
- Components
- Dependencies

### Using a Template
```bash
# List available templates
jira-ticket-creator template list

# Create from template
jira-ticket-creator create --template bug
```

## Global Flags

These flags work with all commands:

| Flag | Type | Description |
|------|------|-------------|
| `--url` | string | JIRA base URL (env: JIRA_URL) |
| `--email` | string | JIRA email (env: JIRA_EMAIL) |
| `--token` | string | JIRA API token (env: JIRA_TOKEN) |
| `--project` | string | JIRA project key (env: JIRA_PROJECT) |
| `--ticket` | string | JIRA ticket key for auto-extracting project (env: JIRA_TICKET) |
| `--config` | string | Path to config file (default: ~/.jirarc) |

## Configuration

You can provide credentials via:

**1. Command-line flags (highest priority)**
```bash
jira-ticket-creator create --summary "Task" \
 --url https://company.atlassian.net \
 --email user@company.com \
 --token api-token \
 --project PROJ
```

**2. Environment variables**
```bash
export JIRA_URL=https://company.atlassian.net
export JIRA_EMAIL=user@company.com
export JIRA_TOKEN=api-token
export JIRA_PROJECT=PROJ
jira-ticket-creator create --summary "Task"
```

**3. Configuration file (~/.jirarc)**
```yaml
jira:
 url: https://company.atlassian.net
 email: user@company.com
 token: api-token
 project: PROJ
```

## Return Value

The command returns the newly created ticket key:

```
 Ticket created successfully!
PROJ-123
```

Use this key in subsequent commands:

```bash
# Assign the ticket
jira-ticket-creator update PROJ-123 --assignee jane@company.com

# Transition it
jira-ticket-creator transition PROJ-123 --to "In Progress"

# View it
jira-ticket-creator search --key PROJ-123
```

## Error Handling

Common errors and solutions:

| Error | Cause | Solution |
|-------|-------|----------|
| `failed to load configuration` | Missing credentials | Set JIRA_URL, JIRA_EMAIL, JIRA_TOKEN env vars or use config file |
| `project PROJ not found` | Invalid project | Verify project key and that you have permissions |
| `summary is required` | Missing --summary flag | Add `--summary "Your title"` |
| `unauthorized` | Invalid credentials | Check JIRA_TOKEN is an API token, not a password |
| `issue type not valid` | Invalid issue type | Use valid type: Task, Story, Bug, Epic, etc. |

## See Also

- [Search Command](./search.md)
- [Update Command](./update.md)
- [Batch Command](./batch.md)
- [Templates](../advanced/templates.md)
