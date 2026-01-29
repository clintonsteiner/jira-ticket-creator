---
layout: default
title: Troubleshooting & FAQ
nav_order: 7
has_toc: true
---

# Troubleshooting & FAQ

Solutions to common issues and frequently asked questions.

## Installation Issues

### Error: "command not found: jira-ticket-creator"

**Problem:** The binary isn't in your PATH

**Solutions:**

1. **Use full path:**
 ```bash
 /path/to/jira-ticket-creator --help
 ```

2. **Add to PATH:**
 ```bash
 export PATH=$PATH:/path/to/directory
 jira-ticket-creator --help
 ```

3. **Move to /usr/local/bin:**
 ```bash
 sudo mv jira-ticket-creator /usr/local/bin/
 jira-ticket-creator --help
 ```

4. **Use Go directly:**
 ```bash
 go run github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest --help
 ```

### Error: "failed to build"

**Problem:** Go build fails with dependencies

**Solution:**
```bash
go mod tidy
go mod download
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

## Authentication Issues

### Error: "failed to load configuration"

**Problem:** Missing or incorrect JIRA credentials

**Solutions:**

1. **Check environment variables:**
 ```bash
 echo "URL: $JIRA_URL"
 echo "Email: $JIRA_EMAIL"
 echo "Token: $JIRA_TOKEN"
 echo "Project: $JIRA_PROJECT"
 ```

2. **Set credentials:**
 ```bash
 export JIRA_URL=https://company.atlassian.net
 export JIRA_EMAIL=user@company.com
 export JIRA_TOKEN=api-token
 export JIRA_PROJECT=PROJ
 ```

3. **Create config file:**
 ```bash
 cat > ~/.jirarc << 'EOF'
 jira:
 url: https://company.atlassian.net
 email: user@company.com
 token: api-token
 project: PROJ
 EOF
 chmod 600 ~/.jirarc
 ```

### Error: "unauthorized" or "401"

**Problem:** Invalid credentials or token

**Solutions:**

1. **Verify token is API token, not password:**
 - Go to https://id.atlassian.com/manage-profile/security/api-tokens
 - Create or regenerate your API token

2. **Check email address:**
 ```bash
 # Must match JIRA login email, not username
 export JIRA_EMAIL=actual-email@company.com
 ```

3. **Test credentials:**
 ```bash
 jira-ticket-creator search --key PROJ-1
 ```

### Error: "project PROJ not found"

**Problem:** Invalid project key or no permissions

**Solutions:**

1. **Verify project key:**
 - Log into JIRA
 - Check the correct project key (in URL or project settings)

2. **Check permissions:**
 - Ensure your user has access to the project
 - Check role/permissions in project settings

3. **Use different project:**
 ```bash
 jira-ticket-creator create --summary "Test" --project DIFFERENT-PROJ
 ```

## Configuration Issues

### Error: "flag provided but not defined"

**Problem:** Unknown flag used

**Solution:** Check available flags:
```bash
jira-ticket-creator [command] --help
```

### Config file not being read

**Problem:** Config file exists but isn't used

**Solutions:**

1. **Ensure correct location:**
 ```bash
 cat ~/.jirarc
 ```

2. **Check format (must be YAML):**
 ```yaml
 jira:
 url: https://company.atlassian.net
 email: user@company.com
 token: api-token
 ```

3. **Override with env vars:**
 ```bash
 JIRA_URL=https://... JIRA_EMAIL=... jira-ticket-creator create --summary "Test"
 ```

## Command Execution Issues

### Error: "summary is required"

**Problem:** Missing required flag

**Solution:** Add the flag:
```bash
jira-ticket-creator create --summary "Your title"
```

### Error: "ticket not found"

**Problem:** Ticket key doesn't exist

**Solution:** Verify the key:
```bash
jira-ticket-creator search --jql "project = PROJ" | grep KEY
```

### No output from search

**Problem:** Query returned no results

**Solutions:**

1. **Check JQL syntax:**
 ```bash
 # Test in JIRA UI first
 jira-ticket-creator search --jql "project = PROJ"
 ```

2. **Verify project has tickets:**
 ```bash
 jira-ticket-creator search --jql "project = PROJ"
 ```

3. **Try simpler query:**
 ```bash
 jira-ticket-creator search --jql "project = PROJ" | head -5
 ```

## Storage Issues

### Error: "failed to initialize storage"

**Problem:** Can't create or access storage directory

**Solution:** Check directory permissions:
```bash
# Create with correct permissions
mkdir -p ~/.jira
chmod 755 ~/.jira

# Verify
ls -la ~/.jira
```

### Error: "failed to load tickets"

**Problem:** Corrupted or missing tickets.json

**Solution:**

1. **Backup old file:**
 ```bash
 cp ~/.jira/tickets.json ~/.jira/tickets.json.bak
 ```

2. **Start fresh:**
 ```bash
 rm ~/.jira/tickets.json
 jira-ticket-creator create --summary "First ticket"
 ```

### Tickets not persisting

**Problem:** Tickets created but don't appear when running reports

**Solution:** Verify storage location:
```bash
ls -la ~/.jira/tickets.json
cat ~/.jira/tickets.json | jq . | head -20
```

## Report Issues

### Error: "failed to generate report"

**Problem:** Report generation failed

**Solution:** Check for corrupted data:
```bash
# Validate JSON
jq . ~/.jira/tickets.json > /dev/null

# Check file
head -5 ~/.jira/tickets.json
```

### Empty or missing data in reports

**Problem:** Reports show no tickets

**Solutions:**

1. **Verify tickets exist:**
 ```bash
 ls -la ~/.jira/tickets.json
 wc -l ~/.jira/tickets.json
 ```

2. **Check team summary:**
 ```bash
 jira-ticket-creator team summary
 ```

3. **Query directly:**
 ```bash
 jira-ticket-creator query --jql "project = PROJ"
 ```

## Performance Issues

### Slow queries

**Problem:** Searches are slow

**Solutions:**

1. **Reduce result set:**
 ```bash
 jira-ticket-creator query --jql "project = PROJ AND status = 'To Do'" --max-results 50
 ```

2. **Use specific filters:**
 ```bash
 # Bad (broad)
 jira-ticket-creator query --jql "summary ~ 'test'"

 # Better (specific)
 jira-ticket-creator query --jql "project = PROJ AND summary ~ 'test'"
 ```

3. **Check JIRA performance:**
 - Try the query in JIRA UI first
 - Check if JIRA instance is slow

### Memory issues with large imports

**Problem:** Out of memory during import

**Solutions:**

1. **Import in batches:**
 ```bash
 # Instead of all at once
 jira-ticket-creator import --jql "project = PROJ AND created >= -1d"
 jira-ticket-creator import --jql "project = PROJ AND created <= -1d" --update-existing
 ```

2. **Limit results:**
 ```bash
 jira-ticket-creator query --jql "project = PROJ" --max-results 500
 ```

## Advanced Debugging

### Enable verbose output

```bash
VERBOSE=1 jira-ticket-creator create --summary "Test"
```

### Check API calls

```bash
# Use curl to test API directly
curl -u "email:token" \
 "https://company.atlassian.net/rest/api/2/issue/PROJ-1"
```

### Validate configuration

```bash
#!/bin/bash
# validate-config.sh

echo "Checking configuration..."
echo "JIRA_URL: $JIRA_URL"
echo "JIRA_EMAIL: $JIRA_EMAIL"
echo "JIRA_PROJECT: $JIRA_PROJECT"

# Test connection
if curl -s -u "$JIRA_EMAIL:$JIRA_TOKEN" \
 "$JIRA_URL/rest/api/2/myself" > /dev/null; then
 echo " Authentication successful"
else
 echo " Authentication failed"
fi

# Check for config file
if [ -f ~/.jirarc ]; then
 echo " Config file found"
else
 echo " No config file found"
fi
```

## FAQ

### Q: Can I use the same instance for multiple JIRA accounts?

**A:** Use environment variables or config files:
```bash
# Account 1
export JIRA_EMAIL=user1@company.com
export JIRA_TOKEN=token1
jira-ticket-creator create --summary "From account 1"

# Account 2
export JIRA_EMAIL=user2@company.com
export JIRA_TOKEN=token2
jira-ticket-creator create --summary "From account 2"
```

### Q: How do I update multiple tickets at once?

**A:** Use batch operations:
```bash
jira-ticket-creator search --jql "status = 'To Do'" --format json | \
 jq -r '.[] | .key' | \
 while read KEY; do
 jira-ticket-creator update "$KEY" --priority High
 done
```

### Q: Can I schedule automatic imports?

**A:** Yes, use cron:
```bash
# Run daily at 9 AM
0 9 * * * /path/to/jira-ticket-creator import --jql "updated >= -1d" --update-existing
```

### Q: How do I backup my tickets?

**A:** Backup the storage file:
```bash
# Daily backup
0 0 * * * cp ~/.jira/tickets.json ~/backups/tickets-$(date +%Y-%m-%d).json

# Or version control
git add ~/.jira/tickets.json
git commit -m "Backup tickets"
```

### Q: What if I want to migrate to a different JIRA instance?

**A:** Export and import:
```bash
# Export from old instance
JIRA_URL=old.atlassian.net JIRA_EMAIL=... JIRA_TOKEN=... \
 jira-ticket-creator import --jql "project = PROJ"

# Tickets now in ~/.jira/tickets.json

# Import to new instance
JIRA_URL=new.atlassian.net JIRA_EMAIL=... JIRA_TOKEN=... \
 jira-ticket-creator create ...
```

### Q: Can I use custom fields?

**A:** Not directly in CLI, but you can use the Go library to set custom fields when creating issues.

### Q: How do I contribute?

**A:** See CONTRIBUTING.md in the repository or open an issue on GitHub.

## Getting Help

### 1. Check the Docs
- [Getting Started](/jira-ticket-creator/getting-started)
- [CLI Commands](/jira-ticket-creator/cli)
- [API Reference](/jira-ticket-creator/api)

### 2. Search Issues
- GitHub Issues: https://github.com/clintonsteiner/jira-ticket-creator/issues

### 3. Run Help Command
```bash
jira-ticket-creator --help
jira-ticket-creator [command] --help
```

### 4. Check Configuration
```bash
# Verify all settings
echo $JIRA_URL
echo $JIRA_EMAIL
cat ~/.jirarc
```

### 5. Test with Simple Command
```bash
jira-ticket-creator search --key PROJ-1
```

## Report a Bug

If you find a bug:

1. **Reproduce it with minimal example:**
 ```bash
 jira-ticket-creator [command] --flag value
 ```

2. **Check error message:**
 ```bash
 jira-ticket-creator [command] 2>&1 | cat
 ```

3. **Open GitHub issue** with:
 - Command that failed
 - Error message
 - Expected behavior
 - System info (OS, Go version)

## See Also

- [Getting Started](/jira-ticket-creator/getting-started) - Initial setup guide
- [CLI Commands](/jira-ticket-creator/cli) - Command reference
