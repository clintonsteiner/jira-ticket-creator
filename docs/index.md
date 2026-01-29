---
layout: default
title: Home
nav_order: 1
has_toc: true
---

# JIRA Ticket Creator

Command-line tool for managing JIRA tickets. Create, search, import, and report on tickets without leaving the terminal.

## Quick Links

- [Get Started](/jira-ticket-creator/getting-started) - Installation and first steps
- [CLI Commands](/jira-ticket-creator/cli) - Create, search, import, gantt
- [Common Workflows](/jira-ticket-creator/examples/workflows) - Real-world examples
- [Troubleshooting](/jira-ticket-creator/troubleshooting) - Common problems and fixes

## What It Does

**Create tickets** individually or in batch:
```bash
jira-ticket-creator create --summary "Fix login bug"
jira-ticket-creator batch create --input tickets.csv
```

**Search and query** with JQL:
```bash
jira-ticket-creator search --jql "status = 'To Do'"
jira-ticket-creator query --jql "project = PROJ" --format csv --output results.csv
```

**Import tickets** with project mapping:
```bash
jira-ticket-creator import --jql "project = PROJ" --map-project backend
```

**View reports and charts:**
```bash
jira-ticket-creator team summary
jira-ticket-creator gantt --format html --output gantt.html
```

## Why Use It

- No browser required - work in your terminal
- Batch operations - create, update, or transition multiple tickets at once
- Multiple output formats - table, JSON, CSV, Markdown, HTML
- Project mapping - organize tickets from multiple JIRA projects into logical teams
- Scripting friendly - pipe commands, write automation

## Installation

```bash
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

Or build from source:
```bash
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

## Setup

Export your JIRA credentials:
```bash
export JIRA_URL=https://company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

Test it:
```bash
jira-ticket-creator search --key PROJ-1
```

## Documentation

- [Getting Started](/jira-ticket-creator/getting-started)
- [CLI Commands](/jira-ticket-creator/cli)
- [Create](/jira-ticket-creator/cli/create)
- [Search](/jira-ticket-creator/cli/search)
- [Query](/jira-ticket-creator/cli/query)
- [Import](/jira-ticket-creator/cli/import)
- [Gantt](/jira-ticket-creator/cli/gantt)
- [Examples](/jira-ticket-creator/examples)
- [API Guide](/jira-ticket-creator/api)
- [Go Client](/jira-ticket-creator/api/go-client)
- [Project Mapping](/jira-ticket-creator/advanced/project-mapping)
- [Troubleshooting](/jira-ticket-creator/troubleshooting)

## Common Tasks

**View your work:**
```bash
jira-ticket-creator search --jql "assignee = currentUser()"
```

**Check critical issues:**
```bash
jira-ticket-creator search --jql "priority = Critical"
```

**Generate team report:**
```bash
jira-ticket-creator team summary
```

**Create Gantt chart:**
```bash
jira-ticket-creator gantt --format html --output workload.html
open workload.html
```

## Project Info

- **Repository**: https://github.com/clintonsteiner/jira-ticket-creator
- **Issues**: https://github.com/clintonsteiner/jira-ticket-creator/issues
- **License**: MIT
