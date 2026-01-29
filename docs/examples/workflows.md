---
layout: default
title: Common Workflows
parent: Examples
---

# Common Workflows

Real-world examples and scripts for common JIRA workflows.

## Daily Standup Report

```bash
#!/bin/bash
# daily-standup.sh

echo "🚀 DAILY STANDUP - $(date +%A, %B %d, %Y)"
echo "========================================"
echo ""

# Section 1: My Work
echo "📋 MY ASSIGNMENTS"
echo "---"
jira-ticket-creator search --jql "assignee = currentUser() AND status != Done" --format table
echo ""

# Section 2: Created by me that needs attention
echo "👀 I CREATED (IN PROGRESS)"
echo "---"
jira-ticket-creator search --jql "creator = currentUser() AND status = 'In Progress'" --format table
echo ""

# Section 3: Blocked items
echo "🚫 BLOCKED ITEMS"
echo "---"
jira-ticket-creator search --jql "status = Blocked" --format table
echo ""

# Section 4: Critical items
echo "🔴 CRITICAL PRIORITY"
echo "---"
jira-ticket-creator search --jql "priority = Critical AND status != Done" --format table
```

Run it:
```bash
chmod +x daily-standup.sh
./daily-standup.sh
```

## Sprint Planning Report

```bash
#!/bin/bash
# sprint-planning.sh

DATE=$(date +%Y-%m-%d)
REPORT_DIR="reports/$DATE"
mkdir -p "$REPORT_DIR"

echo "🏃 SPRINT PLANNING REPORT - $DATE"
echo "======================================"

# Ready for sprint
echo ""
echo "✅ READY FOR SPRINT"
jira-ticket-creator search --jql "status = 'Ready' OR labels = 'ready-for-sprint'" \
  --format table | tee "$REPORT_DIR/ready.txt"

# Tech debt
echo ""
echo "🔧 TECH DEBT"
jira-ticket-creator search --jql "labels = 'tech-debt' AND status != Done" \
  --format table | tee "$REPORT_DIR/tech-debt.txt"

# Bug backlog
echo ""
echo "🐛 BUG BACKLOG"
jira-ticket-creator search --jql "type = Bug AND status = 'To Do' AND priority >= High" \
  --format table | tee "$REPORT_DIR/bugs.txt"

# Estimation needed
echo ""
echo "❓ NEEDS ESTIMATION"
jira-ticket-creator search --jql "customfield_10000 is EMPTY AND status != Done" \
  --format table | tee "$REPORT_DIR/estimation.txt"

echo ""
echo "✅ Report saved to $REPORT_DIR"
```

## Weekly Status Report

```bash
#!/bin/bash
# weekly-status.sh

WEEK=$(date +%Y-W%V)
REPORT_DIR="reports/weekly/$WEEK"
mkdir -p "$REPORT_DIR"

echo "📊 WEEKLY STATUS REPORT - Week $WEEK"
echo "========================================="

# Completed this week
echo ""
echo "🎉 COMPLETED THIS WEEK"
jira-ticket-creator query --jql "status = Done AND updated >= -7d" \
  --format table > "$REPORT_DIR/completed.txt"

# In progress
echo ""
echo "⏳ IN PROGRESS"
jira-ticket-creator query --jql "status = 'In Progress'" \
  --format table > "$REPORT_DIR/in-progress.txt"

# Blocked
echo ""
echo "🚫 BLOCKED"
jira-ticket-creator query --jql "status = Blocked" \
  --format table > "$REPORT_DIR/blocked.txt"

# Team summary
echo ""
echo "👥 TEAM SUMMARY"
jira-ticket-creator team summary > "$REPORT_DIR/team-summary.txt"

# Gantt chart
jira-ticket-creator gantt --format html --output "$REPORT_DIR/gantt.html"

echo ""
echo "✅ Reports saved to $REPORT_DIR"
ls -la "$REPORT_DIR"
```

## Issue Triage Workflow

```bash
#!/bin/bash
# triage.sh - Organize and categorize issues

echo "🔍 ISSUE TRIAGE"
echo "=================="

# Step 1: Find unassigned high-priority
echo ""
echo "1️⃣  UNASSIGNED HIGH PRIORITY"
UNASSIGNED=$(jira-ticket-creator search --jql "assignee is EMPTY AND priority >= High" --format json)
echo "$UNASSIGNED" | jq -r '.[] | "\(.key): \(.fields.summary)"'

# Step 2: Find critical items without estimates
echo ""
echo "2️⃣  CRITICAL WITHOUT ESTIMATES"
jira-ticket-creator search --jql "priority = Critical AND customfield_10000 is EMPTY" --format table

# Step 3: Find bugs awaiting fixes
echo ""
echo "3️⃣  BUGS AWAITING FIXES"
jira-ticket-creator search --jql "type = Bug AND status = 'Ready for Dev'" --format table

# Step 4: Find items in review
echo ""
echo "4️⃣  IN REVIEW"
jira-ticket-creator search --jql "status = 'In Review'" --format table

echo ""
echo "✅ Triage summary complete"
```

## Backlog Management

```bash
#!/bin/bash
# backlog-management.sh

echo "📚 BACKLOG MANAGEMENT"
echo "======================"

# Current backlog size
BACKLOG_COUNT=$(jira-ticket-creator search --jql "project = PROJ AND status = 'To Do'" | wc -l)
echo "Current backlog: $BACKLOG_COUNT items"

# Items per priority
echo ""
echo "BREAKDOWN BY PRIORITY:"
for PRIORITY in Critical High Medium Low; do
  COUNT=$(jira-ticket-creator search --jql "project = PROJ AND status = 'To Do' AND priority = $PRIORITY" | wc -l)
  echo "  $PRIORITY: $COUNT items"
done

# High priority items
echo ""
echo "HIGH PRIORITY BACKLOG:"
jira-ticket-creator search --jql "project = PROJ AND status = 'To Do' AND priority >= High" \
  --format table

# Oldest items
echo ""
echo "OLDEST ITEMS IN BACKLOG:"
jira-ticket-creator search --jql "project = PROJ AND status = 'To Do' ORDER BY created ASC" \
  --format table | head -10
```

## Multi-Team Project Import

```bash
#!/bin/bash
# multi-team-import.sh

# Create project mapping
cat > ~/.jira/project-mapping.json << 'EOF'
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
EOF

echo "📥 IMPORTING MULTI-TEAM PROJECT"
echo "=================================="

# Import backend tickets
echo ""
echo "1. Importing Backend tickets..."
jira-ticket-creator import --jql "project = PROJ" --map-project backend

# Import frontend tickets
echo ""
echo "2. Importing Frontend tickets..."
jira-ticket-creator import --jql "project = UI" --map-project frontend

# Import DevOps tickets
echo ""
echo "3. Importing DevOps tickets..."
jira-ticket-creator import --jql "project = INFRA" --map-project devops

# Generate reports
echo ""
echo "4. Generating team reports..."
echo ""
echo "BACKEND TEAM:"
jira-ticket-creator team summary --project backend

echo ""
echo "FRONTEND TEAM:"
jira-ticket-creator team summary --project frontend

echo ""
echo "DEVOPS TEAM:"
jira-ticket-creator team summary --project devops

# Generate Gantt
echo ""
echo "5. Generating Gantt chart..."
jira-ticket-creator gantt --format html --output workload.html
echo "✅ Gantt chart: workload.html"
```

## Release Checklist

```bash
#!/bin/bash
# release-checklist.sh - Verify release readiness

echo "🚀 RELEASE CHECKLIST"
echo "===================="

# Check 1: All tasks done
echo ""
echo "1. COMPLETION STATUS"
DONE=$(jira-ticket-creator search --jql "fixVersion = '1.0.0' AND status = Done" | wc -l)
TOTAL=$(jira-ticket-creator search --jql "fixVersion = '1.0.0'" | wc -l)
echo "   Done: $DONE/$TOTAL"

if [ "$DONE" -eq "$TOTAL" ]; then
  echo "   ✅ All items completed"
else
  echo "   ❌ Still $((TOTAL - DONE)) items pending"
  jira-ticket-creator search --jql "fixVersion = '1.0.0' AND status != Done"
fi

# Check 2: No critical bugs
echo ""
echo "2. CRITICAL BUGS"
CRITICAL=$(jira-ticket-creator search --jql "fixVersion = '1.0.0' AND type = Bug AND priority = Critical" | wc -l)
if [ "$CRITICAL" -eq 0 ]; then
  echo "   ✅ No critical bugs"
else
  echo "   ❌ Found $CRITICAL critical bugs"
fi

# Check 3: Documentation done
echo ""
echo "3. DOCUMENTATION"
DOCS=$(jira-ticket-creator search --jql "fixVersion = '1.0.0' AND labels = documentation AND status = Done" | wc -l)
echo "   Documentation items completed: $DOCS"

# Check 4: Testing complete
echo ""
echo "4. TESTING"
TESTING=$(jira-ticket-creator search --jql "fixVersion = '1.0.0' AND labels = testing AND status = Done" | wc -l)
echo "   Testing items completed: $TESTING"

echo ""
echo "✅ Release checklist complete"
```

## On-Call Monitoring

```bash
#!/bin/bash
# on-call-monitor.sh - Real-time monitoring dashboard

while true; do
  clear
  echo "🚨 ON-CALL DASHBOARD - $(date)"
  echo "======================================"

  # Critical production issues
  echo ""
  echo "🔴 CRITICAL PRODUCTION"
  CRITICAL=$(jira-ticket-creator search --jql "labels = production AND priority = Critical AND status != Done")
  if [ -z "$CRITICAL" ]; then
    echo "   ✅ No critical issues"
  else
    echo "$CRITICAL"
  fi

  # Incidents
  echo ""
  echo "⚠️  INCIDENTS"
  INCIDENTS=$(jira-ticket-creator search --jql "type = Incident AND status != Done")
  if [ -z "$INCIDENTS" ]; then
    echo "   ✅ No active incidents"
  else
    echo "$INCIDENTS"
  fi

  # SLA violations
  echo ""
  echo "⏰ SLA APPROACHING"
  jira-ticket-creator search --jql "duedate <= 2h"

  # Update every 5 minutes
  sleep 300
done
```

## Dependency Management

```bash
#!/bin/bash
# dependency-check.sh - Check blocking relationships

echo "🔗 DEPENDENCY ANALYSIS"
echo "======================="

# Blocked items
echo ""
echo "BLOCKED ITEMS:"
BLOCKED=$(jira-ticket-creator search --jql "status = Blocked")
BLOCKED_COUNT=$(echo "$BLOCKED" | wc -l)
echo "Total blocked: $BLOCKED_COUNT"

if [ "$BLOCKED_COUNT" -gt 0 ]; then
  echo ""
  echo "$BLOCKED"
fi

# Visualize dependencies
echo ""
echo "DEPENDENCY TREE:"
jira-ticket-creator visualize --format tree

# Export dependency graph
echo ""
echo "Exporting Mermaid diagram..."
jira-ticket-creator visualize --format mermaid --output dependencies.md

echo ""
echo "✅ Dependency analysis complete"
```

## Sprint Velocity Tracking

```bash
#!/bin/bash
# velocity-tracking.sh

echo "📈 SPRINT VELOCITY TRACKING"
echo "============================"

SPRINT=$(jira-ticket-creator search --jql "sprint = activeSprints()" | head -1)

# Story points completed
echo ""
echo "THIS SPRINT:"
jira-ticket-creator search --jql "sprint = activeSprints() AND status = Done" \
  --format table

# Story points remaining
echo ""
echo "STILL IN PROGRESS:"
jira-ticket-creator search --jql "sprint = activeSprints() AND status != Done" \
  --format table

# Burndown view
echo ""
echo "SPRINT TIMELINE:"
jira-ticket-creator team timeline

echo ""
echo "✅ Velocity report complete"
```

## Bulk Operations

```bash
#!/bin/bash
# bulk-update.sh - Update multiple tickets

echo "🔄 BULK OPERATIONS"
echo "===================="

# Find all tickets to update
TICKETS=$(jira-ticket-creator search --jql "labels = 'needs-review'" --format json)

echo "Processing $(echo "$TICKETS" | jq 'length') tickets..."

# Transition each
echo "$TICKETS" | jq -r '.[] | .key' | while read -r KEY; do
  echo "Updating $KEY..."
  jira-ticket-creator transition "$KEY" --to "In Review"
done

echo ""
echo "✅ Bulk update complete"
```

## Automated Weekly Reports

```bash
#!/bin/bash
# cron-weekly-report.sh
# Add to crontab: 0 9 * * 1 /path/to/cron-weekly-report.sh

WEEK=$(date +%Y-W%V)
REPORT_DIR="$HOME/reports/weekly/$WEEK"
mkdir -p "$REPORT_DIR"

# Generate all reports
jira-ticket-creator query --jql "updated >= -7d" --format csv --output "$REPORT_DIR/weekly.csv"
jira-ticket-creator team summary > "$REPORT_DIR/team.txt"
jira-ticket-creator gantt --format html --output "$REPORT_DIR/gantt.html"
jira-ticket-creator pm dashboard > "$REPORT_DIR/dashboard.txt"

# Email report
mail -s "Weekly JIRA Report - Week $WEEK" \
  manager@company.com \
  < "$REPORT_DIR/dashboard.txt"

echo "✅ Weekly report sent"
```

## See Also

- [Search Examples](../cli/search.md)
- [Query Examples](../cli/query.md)
- [CLI Commands](../cli/)
