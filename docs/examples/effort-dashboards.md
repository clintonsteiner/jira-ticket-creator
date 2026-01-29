---
layout: default
title: Building Effort Dashboards
parent: Examples
nav_order: 4
has_toc: true
---

# Building Effort Dashboards

This guide shows how to create comprehensive dashboards for monitoring multiple efforts, team workloads, and project health.

## Overview

Build automated dashboards that answer key questions:
- What efforts are we running?
- What's the status of each effort?
- Which teams are most loaded?
- What are the risks and blockers?
- Are we on track for deadlines?

## Dashboard Types

1. **Executive Dashboard** - High-level status for leadership
2. **Team Dashboard** - Work status for individual teams
3. **Effort Dashboard** - Detailed tracking for effort leads
4. **Risk Dashboard** - Blockers, dependencies, red flags

---

## Executive Dashboard

For leadership visibility into all efforts.

### Setup

```bash
mkdir -p ~/dashboards/executive
cd ~/dashboards/executive

cat > executive-dashboard.sh << 'EOF'
#!/bin/bash
# Executive-level dashboard for all efforts

OUTPUT_DIR="./reports/$(date +%Y-%m-%d)"
mkdir -p "$OUTPUT_DIR"

echo "EXECUTIVE DASHBOARD - $(date +%Y-%m-%d)"
echo "Organization Status Report"
echo "=================================================="
echo ""

# Define all efforts in organization
EFFORTS=(
  "PROJ-1000:mobile-redesign:2025-04-30"
  "PROJ-2000:q1-payment-systems:2025-03-31"
  "PROJ-3000:api-v2-migration:2025-06-30"
  "PROJ-4000:infrastructure-upgrade:2025-05-15"
)

echo "EFFORT STATUS SUMMARY"
echo "────────────────────────────────────────────────────────────────"
printf "%-30s | %5s | %10s | %10s | STATUS\n" "Effort" "Done" "Progress" "To Do"
echo "────────────────────────────────────────────────────────────────"

for EFFORT_INFO in "${EFFORTS[@]}"; do
  IFS=':' read -r EPIC PROJECT_NAME TARGET_DATE <<< "$EFFORT_INFO"

  # Get counts
  TOTAL=$(jira-ticket-creator search --jql "key = $EPIC OR parent = $EPIC" 2>/dev/null | wc -l)
  DONE=$(jira-ticket-creator search --jql "(key = $EPIC OR parent = $EPIC) AND status = Done" 2>/dev/null | wc -l)
  IN_PROGRESS=$(jira-ticket-creator search --jql "(key = $EPIC OR parent = $EPIC) AND status = 'In Progress'" 2>/dev/null | wc -l)
  TODO=$((TOTAL - DONE - IN_PROGRESS))

  # Calculate health
  if [ $TOTAL -gt 0 ]; then
    PERCENT=$((DONE * 100 / TOTAL))
    if [ $PERCENT -ge 75 ]; then STATUS="✅ ON TRACK";
    elif [ $PERCENT -ge 50 ]; then STATUS="IN PROGRESS";
    else STATUS="AT RISK"; fi
  else
    PERCENT=0
    STATUS="PLANNING"
  fi

  printf "%-30s | %5d | %10d | %10d | %s\n" "$PROJECT_NAME" "$DONE" "$IN_PROGRESS" "$TODO" "$STATUS"
done

echo "────────────────────────────────────────────────────────────────"
echo ""

# Critical issues across all efforts
echo "CRITICAL ISSUES REQUIRING ATTENTION"
echo "────────────────────────────────────────────────────────────────"

# Blockers
BLOCKERS=$(jira-ticket-creator search --jql "status = Blocked" 2>/dev/null | wc -l)
echo "Blocked Tickets: $BLOCKERS"
if [ $BLOCKERS -gt 0 ]; then
  jira-ticket-creator search --jql "status = Blocked" --format table | head -5
fi
echo ""

# Overdue
OVERDUE=$(jira-ticket-creator search --jql "duedate < now() AND status != Done" 2>/dev/null | wc -l)
echo "Overdue Tickets: $OVERDUE"
if [ $OVERDUE -gt 0 ]; then
  jira-ticket-creator search --jql "duedate < now() AND status != Done" --format table | head -5
fi
echo ""

# High priority
CRITICAL=$(jira-ticket-creator search --jql "priority = Critical AND status != Done" 2>/dev/null | wc -l)
echo "Critical Priority: $CRITICAL"
if [ $CRITICAL -gt 0 ]; then
  jira-ticket-creator search --jql "priority = Critical AND status != Done" --format table | head -5
fi
echo ""

# Resource utilization
echo "RESOURCE UTILIZATION"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator team summary 2>/dev/null | head -20
echo ""

# Team workload
echo "TEAM WORKLOAD VISUALIZATION"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator gantt --format ascii --weeks 2

echo ""
echo "Report generated: $OUTPUT_DIR"
echo "Full HTML report: $OUTPUT_DIR/dashboard.html"

# Generate HTML version
jira-ticket-creator gantt --format html --output "$OUTPUT_DIR/dashboard.html"
EOF

chmod +x executive-dashboard.sh
```

### Run It

```bash
./executive-dashboard.sh
# Output shows all efforts, blockers, risks, team status

# Schedule it to run daily
# Add to crontab:
# 0 9 * * * /path/to/executive-dashboard.sh
```

---

## Team Dashboard

For individual team leads to monitor their team's work.

### Setup

```bash
mkdir -p ~/dashboards/teams/backend
cd ~/dashboards/teams/backend

cat > team-dashboard.sh << 'EOF'
#!/bin/bash
# Team-level dashboard for backend team

TEAM="backend"
OUTPUT_DIR="./reports/$(date +%Y-%m-%d)"
mkdir -p "$OUTPUT_DIR"

echo "$TEAM TEAM DASHBOARD"
echo "$(date +%Y-%m-%d %H:%M:%S)"
echo "=================================================="
echo ""

# 1. Team workload
echo "1. TEAM WORKLOAD"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "project in (PROJ, API, DB)" --format table
echo ""

# 2. My assignments
echo "2. MY ASSIGNMENTS"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "assignee = currentUser() AND status != Done" --format table
echo ""

# 3. Effort breakdown for team
echo "3. EFFORT BREAKDOWN (What I'm working on)"
echo "────────────────────────────────────────────────────────────────"

EFFORTS=("PROJ-1000" "PROJ-2000" "PROJ-3000")
for EPIC in "${EFFORTS[@]}"; do
  COUNT=$(jira-ticket-creator search --jql "parent = $EPIC AND assignee = currentUser()" 2>/dev/null | wc -l)
  if [ $COUNT -gt 0 ]; then
    echo "Effort $EPIC: $COUNT tickets"
  fi
done
echo ""

# 4. Blockers
echo "4. BLOCKERS (What's stuck)"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "project in (PROJ, API, DB) AND status = Blocked" --format table
echo ""

# 5. Today's focus
echo "5. THIS WEEK'S FOCUS"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "project in (PROJ, API, DB) AND priority = High AND status != Done" --format table | head -10
echo ""

# 6. Code review needed
echo "6. IN REVIEW / READY TO MERGE"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "project in (PROJ, API, DB) AND status = 'In Review'" --format table
echo ""

# Generate gantt for team
echo "Generating team gantt chart..."
jira-ticket-creator gantt --format html --output "$OUTPUT_DIR/team-gantt.html"

echo ""
echo "Report saved to: $OUTPUT_DIR"
echo "Open in browser: open $OUTPUT_DIR/team-gantt.html"
EOF

chmod +x team-dashboard.sh
```

### Run It

```bash
./team-dashboard.sh
# Shows personal work, team blockers, efforts, priorities
```

---

## Effort Lead Dashboard

For effort owners to track progress in detail.

### Setup

```bash
mkdir -p ~/dashboards/efforts/mobile-redesign
cd ~/dashboards/efforts/mobile-redesign

cat > effort-dashboard.sh << 'EOF'
#!/bin/bash
# Effort lead dashboard for a single effort/epic

EPIC="PROJ-1000"
EFFORT_NAME="Mobile App Redesign"
OUTPUT_DIR="./reports/$(date +%Y-%m-%d)"
mkdir -p "$OUTPUT_DIR"

echo "EFFORT DASHBOARD: $EFFORT_NAME"
echo "$(date +%Y-%m-%d)"
echo "=================================================="
echo ""

# 1. Overall progress
echo "1. OVERALL PROGRESS"
echo "────────────────────────────────────────────────────────────────"

TOTAL=$(jira-ticket-creator search --jql "parent = $EPIC" 2>/dev/null | wc -l)
DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Done" 2>/dev/null | wc -l)
IN_PROGRESS=$(jira-ticket-creator search --jql "parent = $EPIC AND status = 'In Progress'" 2>/dev/null | wc -l)
TODO=$((TOTAL - DONE - IN_PROGRESS))

if [ $TOTAL -gt 0 ]; then
  PERCENT=$((DONE * 100 / TOTAL))
  FILLED=$((PERCENT / 10))
  EMPTY=$((10 - FILLED))

  PROGRESS_BAR=$(printf '█%.0s' $(seq 1 $FILLED))$(printf '░%.0s' $(seq 1 $EMPTY))
  echo "Progress: [$PROGRESS_BAR] $PERCENT%"
  echo ""
  echo "Done:        $DONE/$TOTAL"
  echo "In Progress: $IN_PROGRESS/$TOTAL"
  echo "To Do:       $TODO/$TOTAL"
else
  echo "No tickets found for this effort"
fi
echo ""

# 2. By component
echo "2. PROGRESS BY COMPONENT"
echo "────────────────────────────────────────────────────────────────"

for COMPONENT in Frontend Backend QA Design; do
  COUNT=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $COMPONENT" 2>/dev/null | wc -l)
  DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $COMPONENT AND status = Done" 2>/dev/null | wc -l)

  if [ $COUNT -gt 0 ]; then
    PCT=$((DONE * 100 / COUNT))
    echo "$COMPONENT: $DONE/$COUNT (${PCT}%)"
  fi
done
echo ""

# 3. Risk analysis
echo "3. RISK ANALYSIS"
echo "────────────────────────────────────────────────────────────────"

BLOCKED=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" 2>/dev/null | wc -l)
CRITICAL=$(jira-ticket-creator search --jql "parent = $EPIC AND priority = Critical" 2>/dev/null | wc -l)
UNASSIGNED=$(jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY" 2>/dev/null | wc -l)

echo "Blocked: $BLOCKED"
echo "Critical Priority: $CRITICAL"
echo "Unassigned: $UNASSIGNED"

if [ $BLOCKED -gt 0 ] || [ $CRITICAL -gt 0 ]; then
  echo ""
  echo "Items needing attention:"
  if [ $BLOCKED -gt 0 ]; then
    jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" --format table | head -5
  fi
fi
echo ""

# 4. Team assignments
echo "4. TEAM ASSIGNMENTS"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND assignee is not EMPTY" --format json 2>/dev/null | \
  jq -r '.[] | .fields.assignee.name' | sort | uniq -c | sort -rn

echo ""

# 5. Upcoming deadlines
echo "5. UPCOMING DEADLINES"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND duedate is not EMPTY AND duedate <= endOfMonth()" \
  --format table | head -10
echo ""

# 6. Visual gantt
jira-ticket-creator gantt --format html --output "$OUTPUT_DIR/gantt.html"

echo "Report saved to: $OUTPUT_DIR"
echo "View gantt: open $OUTPUT_DIR/gantt.html"
EOF

chmod +x effort-dashboard.sh
```

### Run It

```bash
./effort-dashboard.sh
# Detailed view of one effort: progress, risks, deadlines, team load
```

---

## Risk Dashboard

Monitor risks, blockers, and red flags across all work.

### Setup

```bash
cat > risk-dashboard.sh << 'EOF'
#!/bin/bash
# Risk/blocker monitoring dashboard

OUTPUT_DIR="./reports/$(date +%Y-%m-%d)"
mkdir -p "$OUTPUT_DIR"

echo "RISK & BLOCKER DASHBOARD"
echo "$(date +%Y-%m-%d)"
echo "=================================================="
echo ""

# 1. Currently blocked
echo "1. CURRENTLY BLOCKED TICKETS"
echo "────────────────────────────────────────────────────────────────"
BLOCKED=$(jira-ticket-creator search --jql "status = Blocked" 2>/dev/null | wc -l)
echo "Total blocked: $BLOCKED"
echo ""
jira-ticket-creator search --jql "status = Blocked" --format table | tee "$OUTPUT_DIR/blocked.txt"
echo ""

# 2. Overdue items
echo "2. OVERDUE ITEMS"
echo "────────────────────────────────────────────────────────────────"
OVERDUE=$(jira-ticket-creator search --jql "duedate < now() AND status != Done" 2>/dev/null | wc -l)
echo "Total overdue: $OVERDUE"
echo ""
jira-ticket-creator search --jql "duedate < now() AND status != Done" --format table | tee "$OUTPUT_DIR/overdue.txt"
echo ""

# 3. High priority without assignee
echo "3. HIGH PRIORITY UNASSIGNED"
echo "────────────────────────────────────────────────────────────────"
UNASSIGNED=$(jira-ticket-creator search --jql "priority >= High AND assignee is EMPTY" 2>/dev/null | wc -l)
echo "Total unassigned high priority: $UNASSIGNED"
echo ""
jira-ticket-creator search --jql "priority >= High AND assignee is EMPTY" --format table | tee "$OUTPUT_DIR/unassigned.txt"
echo ""

# 4. No activity in 2 weeks
echo "4. STALE TICKETS (No update in 2 weeks)"
echo "────────────────────────────────────────────────────────────────"
STALE=$(jira-ticket-creator search --jql "status != Done AND updated <= -14d" 2>/dev/null | wc -l)
echo "Total stale: $STALE"
echo ""
jira-ticket-creator search --jql "status != Done AND updated <= -14d" --format table | head -10 | tee "$OUTPUT_DIR/stale.txt"
echo ""

# 5. Critical priority
echo "5. CRITICAL PRIORITY ITEMS"
echo "────────────────────────────────────────────────────────────────"
CRITICAL=$(jira-ticket-creator search --jql "priority = Critical" 2>/dev/null | wc -l)
echo "Total critical: $CRITICAL"
echo ""
jira-ticket-creator search --jql "priority = Critical AND status != Done" --format table | tee "$OUTPUT_DIR/critical.txt"
echo ""

# 6. Dependencies at risk
echo "6. DEPENDENCY ISSUES"
echo "────────────────────────────────────────────────────────────────"
echo "Blocking other tickets:"
jira-ticket-creator search --jql "linked = Blocks AND status != Done" --format table | head -10
echo ""

# Generate summary report
TOTAL_ISSUES=$(BLOCKED + OVERDUE + UNASSIGNED + STALE + CRITICAL)
echo ""
echo "RISK SUMMARY"
echo "────────────────────────────────────────────────────────────────"
echo "Blocked: $BLOCKED"
echo "Overdue: $OVERDUE"
echo "High Priority Unassigned: $UNASSIGNED"
echo "Stale: $STALE"
echo "Critical: $CRITICAL"
echo ""
echo "Total items needing attention: $((BLOCKED + OVERDUE + UNASSIGNED + CRITICAL))"

echo ""
echo "Reports saved to: $OUTPUT_DIR"
EOF

chmod +x risk-dashboard.sh
```

### Run It

```bash
./risk-dashboard.sh
# Identifies all issues that need immediate attention
```

---

## Automated Scheduling

Schedule dashboards to run automatically and send reports.

### Using Cron

```bash
# Edit crontab
crontab -e

# Add these lines:

# Executive dashboard - 9am daily
0 9 * * * /path/to/executive-dashboard.sh

# Team dashboards - 8am daily
0 8 * * * /path/to/team-dashboard.sh

# Effort dashboards - 2pm daily
0 14 * * * /path/to/effort-dashboard.sh

# Risk dashboard - 7am daily (early warning)
0 7 * * * /path/to/risk-dashboard.sh

# Weekly summary - Friday 5pm
0 17 * * 5 /path/to/executive-dashboard.sh > ~/weekly-summary.txt
```

### Using GitHub Actions (if repo-based)

```yaml
# .github/workflows/daily-dashboard.yml
name: Daily Dashboard Reports

on:
  schedule:
    - cron: '0 9 * * *'  # 9am UTC daily

jobs:
  dashboard:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: go build -o jira-ticket-creator ./cmd/jira-ticket-creator
      - run: ./dashboards/executive-dashboard.sh
      - name: Commit report
        run: |
          git config user.name "Dashboard Bot"
          git add reports/
          git commit -m "Daily dashboard report"
          git push
```

---

## Dashboard Best Practices

1. **Keep it simple** - Executive dashboard should fit on one screen
2. **Use clear metrics** - % complete, number of blockers, days until deadline
3. **Highlight risks** - Blockers, overdue, critical items need visual prominence
4. **Run consistently** - Same time each day so people expect it
5. **Archive reports** - Keep history for trend analysis
6. **Share widely** - Send to all stakeholders, not just leads
7. **Make actionable** - Each item should show who to contact or what to do

---

## Example Report Output Structure

```
DASHBOARD NAME
  Date/Time

KEY METRICS
  Overall progress: X%
  Team workload: Y items
  Blockers: Z

STATUS BY EFFORT
  Effort 1: 50% complete, 2 blockers
  Effort 2: 75% complete, on track

RISKS & BLOCKERS
  List of critical items

TEAM STATUS
  Resource utilization

UPCOMING
  Deadlines this week

ACTIONS NEEDED
  Specific decisions/escalations
```

---

## See Also

- [Managing JIRA Efforts](effort-management) - How to organize efforts
- [Common Workflows](workflows) - Other automation scripts
- [Gantt Command](../cli/gantt) - Generate visual charts
