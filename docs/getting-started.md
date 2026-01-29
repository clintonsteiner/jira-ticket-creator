---
layout: default
title: Getting Started
nav_order: 2
has_toc: true
---

# Getting Started

## Installation

**Clone and build:**
```bash
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

**Or install directly:**
```bash
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

Verify it works:
```bash
jira-ticket-creator --help
```

## Set Up Credentials

You need three things: JIRA URL, your email, and an API token.

**Get your API token:**
1. Visit https://id.atlassian.com/manage-profile/security/api-tokens
2. Create a new token
3. Copy it (you won't see it again)

**Option 1: Environment variables** (recommended for most)
```bash
export JIRA_URL=https://company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

**Option 2: Config file** (~/.jirarc)
```yaml
jira:
 url: https://company.atlassian.net
 email: your-email@company.com
 token: your-api-token
 project: PROJ
```

Make it readable only by you:
```bash
chmod 600 ~/.jirarc
```

## Test Your Setup

```bash
jira-ticket-creator search --key PROJ-1
```

If you see a ticket, you're good. If not, double-check:
- JIRA_TOKEN must be an API token, not your password
- User account has access to the project
- URL format is correct (no trailing slash)

## Create Your First Ticket

```bash
jira-ticket-creator create --summary "My first ticket"
```

Done. The key will be printed.

## Next Steps

Now that basics work, try these:

**Search for tickets:**
```bash
jira-ticket-creator search --jql "status = 'To Do'"
jira-ticket-creator search --jql "priority = High"
```

**Query with output:**
```bash
jira-ticket-creator query --jql "project = PROJ" --format json
jira-ticket-creator query --jql "project = PROJ" --format csv --output tickets.csv
```

**Import existing tickets:**
```bash
jira-ticket-creator import --jql "project = PROJ" --map-project backend
```

**View team reports:**
```bash
jira-ticket-creator team summary
jira-ticket-creator gantt --format html --output gantt.html
```

## Useful Aliases

Add these to ~/.bashrc or ~/.zshrc:

```bash
alias my-work='jira-ticket-creator search --jql "assignee = currentUser()"'
alias in-progress='jira-ticket-creator search --jql "status = \"In Progress\""'
alias critical='jira-ticket-creator search --jql "priority = Critical"'
```

Then use:
```bash
my-work
critical
in-progress
```

## Common Issues

**Command not found:** Make sure the binary is in your PATH:
```bash
export PATH=$PATH:$(pwd)
./jira-ticket-creator --help
```

**Authentication failed:** Check credentials:
```bash
echo $JIRA_URL
echo $JIRA_EMAIL
cat ~/.jirarc
```

**Project not found:** Verify the project key and permissions.

## Getting Help

- See [Troubleshooting](troubleshooting) for common problems
- Check command help: `jira-ticket-creator [command] --help`
- Report issues: https://github.com/clintonsteiner/jira-ticket-creator/issues

## What's Next

- [Search Examples](cli/search)
- [Import with Project Mapping](cli/import)
- [Common Workflows](examples/workflows)
- [Complete Command Reference](cli)
