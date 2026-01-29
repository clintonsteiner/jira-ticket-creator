---
layout: default
title: Managing JIRA Efforts
parent: Examples
nav_order: 3
---

# Managing JIRA Efforts - Step-by-Step Guides

Learn how to organize, track, and report on work efforts using jira-ticket-creator. An "effort" is any large body of work that may span multiple teams: products, features, initiatives, roadmap items, or quarterly goals.

## Table of Contents

- [Concepts](#concepts)
- [Setting Up an Effort](#setting-up-an-effort)
- [Tracking Effort Progress](#tracking-effort-progress)
- [Multi-Team Efforts](#multi-team-efforts)
- [Effort Reporting](#effort-reporting)
- [Real-World Scenarios](#real-world-scenarios)

## Concepts

### What is an Effort?

An effort is a body of work that:
- Typically owns an Epic or Feature in JIRA
- May have subtasks, stories, and tasks under it
- Involves one or more teams
- Has business value or goals
- Takes weeks or months to complete

### Examples of Efforts

- **Product Feature**: "Mobile App Redesign" (design, frontend, backend, QA)
- **Initiative**: "Payment System Modernization" (backend team focus)
- **Quarterly Goal**: "Q1 2025 Platform Reliability" (infrastructure, operations)
- **Customer Delivery**: "Enterprise Onboarding" (sales, backend, support)
- **Technical Debt**: "API v2 Migration" (backend, operations)

## Setting Up an Effort

### Step 1: Create or Identify Parent Epic in JIRA

First, create an Epic in JIRA or identify an existing one:

```bash
# Example: Create epic in JIRA web UI
# Title: "Mobile App Redesign"
# Key: PROJ-1000
# Link all work to this epic
```

### Step 2: Define Your Effort Locally

Create a local directory for effort tracking:

```bash
mkdir -p ~/jira-efforts/2025
cd ~/jira-efforts/2025
```

### Step 3: Create Effort Configuration File

Create a file documenting the effort:

```bash
cat > mobile-redesign.sh << 'EOF'
#!/bin/bash
# Mobile App Redesign Effort Configuration

EFFORT_NAME="mobile-redesign"
EFFORT_EPIC="PROJ-1000"
EFFORT_GOAL="Redesign mobile app for iOS and Android"
START_DATE="2025-02-01"
TARGET_DATE="2025-04-30"

# Teams involved
TEAMS=("frontend" "backend" "qa")

# Key metrics
echo "Setting up effort: $EFFORT_NAME"
echo "Epic: $EFFORT_EPIC"
echo "Target completion: $TARGET_DATE"
EOF

chmod +x mobile-redesign.sh
```

### Step 4: Import Effort and Its Children

Import the epic and all child tickets:

```bash
jira-ticket-creator import --jql "key = PROJ-1000" \
  --map-project mobile-redesign
```

Verify all children were imported:

```bash
jira-ticket-creator search --jql "parent = PROJ-1000" --format table
```

### Step 5: View Initial Status

Get a snapshot of the effort:

```bash
# See all tickets in the effort
jira-ticket-creator search --jql "parent = PROJ-1000"

# View by ticket type
jira-ticket-creator search --jql "parent = PROJ-1000 AND type = Task"

# View by assignee
jira-ticket-creator search --jql "parent = PROJ-1000 AND assignee is not EMPTY"
```

---

## Tracking Effort Progress

### Daily Progress Check

```bash
cat > check-effort-daily.sh << 'EOF'
#!/bin/bash
# Check daily progress on effort

EFFORT="mobile-redesign"
EPIC="PROJ-1000"
DATE=$(date +%Y-%m-%d)

echo "=== EFFORT PROGRESS: $EFFORT ==="
echo "Date: $DATE"
echo ""

# Status breakdown
echo "STATUS BREAKDOWN:"
jira-ticket-creator search --jql "parent = $EPIC" | awk '{
  if ($3 == "To") print "  To Do: " ++to_do
  else if ($3 == "In") print "  In Progress: " ++in_progress
  else if ($3 == "Done") print "  Done: " ++done
}
END {
  print "  Total: " to_do + in_progress + done
}'

echo ""
echo "RECENTLY UPDATED:"
jira-ticket-creator search --jql "parent = $EPIC AND updated >= -1d" --format table

echo ""
echo "BLOCKED ITEMS:"
jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" --format table
EOF

chmod +x check-effort-daily.sh
./check-effort-daily.sh
```

### Weekly Progress Report

```bash
cat > effort-weekly-report.sh << 'EOF'
#!/bin/bash
# Generate weekly effort report

EFFORT="mobile-redesign"
EPIC="PROJ-1000"
WEEK=$(date +%Y-W%V)
REPORT_DIR="reports/$EFFORT/$WEEK"

mkdir -p "$REPORT_DIR"

echo "=== WEEKLY PROGRESS REPORT ==="
echo "Effort: $EFFORT"
echo "Week: $WEEK"
echo ""

# 1. Completed this week
echo "1. COMPLETED THIS WEEK"
jira-ticket-creator search --jql "parent = $EPIC AND status = Done AND updated >= -7d" \
  --format table | tee "$REPORT_DIR/completed.txt"
COMPLETED=$(wc -l < "$REPORT_DIR/completed.txt")
echo "Total: $COMPLETED tickets"
echo ""

# 2. Currently in progress
echo "2. IN PROGRESS ($(date))"
jira-ticket-creator search --jql "parent = $EPIC AND status = 'In Progress'" \
  --format table | tee "$REPORT_DIR/in-progress.txt"
echo ""

# 3. Blocked or at risk
echo "3. BLOCKED/AT RISK"
jira-ticket-creator search --jql "parent = $EPIC AND (status = Blocked OR priority = Critical)" \
  --format table | tee "$REPORT_DIR/blocked.txt"
echo ""

# 4. By assignee (who's working on what)
echo "4. WORKLOAD BY PERSON"
jira-ticket-creator search --jql "parent = $EPIC AND assignee is not EMPTY" \
  --format json | jq -r '.[] | .fields.assignee.name' | \
  sort | uniq -c | tee "$REPORT_DIR/assignee-summary.txt"
echo ""

# 5. Visual gantt chart
echo "Generating visual gantt chart..."
jira-ticket-creator gantt --format html --output "$REPORT_DIR/gantt.html"

echo "Report saved to: $REPORT_DIR"
ls -lh "$REPORT_DIR"
EOF

chmod +x effort-weekly-report.sh
./effort-weekly-report.sh
```

---

## Multi-Team Efforts

When an effort spans multiple teams, organize by team within the effort.

### Step 1: Import Effort for Each Team's Perspective

```bash
# Import same epic but organize by team
jira-ticket-creator import --jql "key = PROJ-1000" --map-project "mobile-redesign"

# This maps parent + all children to "mobile-redesign"
# All team members see same tickets under this project
```

### Step 2: Track Team Allocations

```bash
cat > effort-team-allocation.sh << 'EOF'
#!/bin/bash
# Show how each team is allocated to this effort

EFFORT="mobile-redesign"
EPIC="PROJ-1000"

echo "=== EFFORT ALLOCATION BY TEAM ==="
echo ""

# Frontend team allocation
echo "FRONTEND TEAM:"
jira-ticket-creator search --jql "parent = $EPIC AND labels = frontend" \
  --format table
echo ""

# Backend team allocation
echo "BACKEND TEAM:"
jira-ticket-creator search --jql "parent = $EPIC AND labels = backend" \
  --format table
echo ""

# QA team allocation
echo "QA TEAM:"
jira-ticket-creator search --jql "parent = $EPIC AND labels = qa" \
  --format table
echo ""

# Unassigned team work
echo "UNASSIGNED (needs assignment):"
jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY" \
  --format table
EOF

chmod +x effort-team-allocation.sh
./effort-team-allocation.sh
```

### Step 3: Review Team Dependencies

```bash
# Show tickets that depend on other tickets
cat > effort-dependencies.sh << 'EOF'
#!/bin/bash

EPIC="PROJ-1000"

echo "=== EFFORT DEPENDENCIES ==="
echo ""
echo "BLOCKING OTHER TICKETS:"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks" --format table

echo ""
echo "BLOCKED BY OTHER TICKETS:"
jira-ticket-creator search --jql "parent = $EPIC AND linked = 'is blocked by'" --format table
EOF

chmod +x effort-dependencies.sh
./effort-dependencies.sh
```

---

## Effort Reporting

### Health Check Dashboard

```bash
cat > effort-health-check.sh << 'EOF'
#!/bin/bash
# Quick health check of effort status

EFFORT="mobile-redesign"
EPIC="PROJ-1000"

echo "╔════════════════════════════════════╗"
echo "║   EFFORT HEALTH CHECK               ║"
echo "║   $EFFORT"
echo "║   $(date +%Y-%m-%d)"
echo "╚════════════════════════════════════╝"
echo ""

# Get raw counts
TOTAL=$(jira-ticket-creator search --jql "parent = $EPIC" | wc -l)
DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Done" | wc -l)
IN_PROGRESS=$(jira-ticket-creator search --jql "parent = $EPIC AND status = 'In Progress'" | wc -l)
TO_DO=$(jira-ticket-creator search --jql "parent = $EPIC AND status = 'To Do'" | wc -l)
BLOCKED=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" | wc -l)

# Calculate percentages
PERCENT_DONE=$((DONE * 100 / TOTAL))
PERCENT_IN_PROGRESS=$((IN_PROGRESS * 100 / TOTAL))

# Progress bar
FILLED=$((PERCENT_DONE / 10))
EMPTY=$((10 - FILLED))
PROGRESS_BAR=$(printf '█%.0s' $(seq 1 $FILLED))$(printf '░%.0s' $(seq 1 $EMPTY))

echo "OVERALL PROGRESS: $PERCENT_DONE%"
echo "[$PROGRESS_BAR]"
echo ""

echo "TICKET STATUS:"
echo "  Done:         $DONE/$TOTAL ($PERCENT_DONE%)"
echo "  In Progress:  $IN_PROGRESS/$TOTAL ($PERCENT_IN_PROGRESS%)"
echo "  To Do:        $TO_DO/$TOTAL"
echo "  Blocked:      $BLOCKED"
echo ""

# Risk indicators
echo "RISKS & CONCERNS:"
if [ $BLOCKED -gt 0 ]; then
  echo "  ⚠ $BLOCKED blocked tickets"
  jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" --format table
fi

UNASSIGNED=$(jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY" | wc -l)
if [ $UNASSIGNED -gt 0 ]; then
  echo "  ⚠ $UNASSIGNED unassigned tickets"
fi

CRITICAL=$(jira-ticket-creator search --jql "parent = $EPIC AND priority = Critical" | wc -l)
if [ $CRITICAL -gt 0 ]; then
  echo "  🔥 $CRITICAL critical priority tickets"
fi

echo ""
echo "TEAM STATUS:"
jira-ticket-creator team summary --ticket $EPIC
EOF

chmod +x effort-health-check.sh
./effort-health-check.sh
```

### Visual Gantt Chart for Effort

```bash
# Generate interactive HTML gantt chart for effort
jira-ticket-creator gantt --format html --output mobile-redesign-gantt.html
open mobile-redesign-gantt.html
```

### Export Effort Data for Stakeholders

```bash
# Export to JSON for analysis
jira-ticket-creator search --jql "parent = PROJ-1000" --format json > effort-data.json

# Export to CSV for Excel
jira-ticket-creator search --jql "parent = PROJ-1000" --format csv > effort-data.csv

# Export to Markdown for GitHub/Confluence
jira-ticket-creator query --jql "parent = PROJ-1000" --format markdown > effort-summary.md
```

---

## Real-World Scenarios

### Scenario 1: Quarterly OKR (Goal Tracking)

```bash
#!/bin/bash
# Track quarterly OKR as effort

# Setup
mkdir -p ~/efforts/Q1-2025/okr-payment-systems
cd ~/efforts/Q1-2025/okr-payment-systems

# Quarterly OKR Epic in JIRA: PROJ-2000
# Goal: Modernize payment processing system
# Leads: John (backend), Jane (QA)

# Step 1: Import the goal/epic
jira-ticket-creator import --jql "key = PROJ-2000" --map-project "q1-payment-modernization"

# Step 2: Daily sync (runs each morning)
cat > daily-sync.sh << 'EOF'
#!/bin/bash
jira-ticket-creator search --jql "parent = PROJ-2000 AND updated >= -1d" --format table
echo ""
echo "URGENT BLOCKERS:"
jira-ticket-creator search --jql "parent = PROJ-2000 AND status = Blocked"
EOF

# Step 3: Weekly leadership review
cat > weekly-executive-summary.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y-%m-%d)
echo "Q1 PAYMENT SYSTEMS - Weekly Update - $DATE"

# Progress
echo "Progress:"
jira-ticket-creator search --jql "parent = PROJ-2000" | awk '
  {
    if ($3 ~ /Done/) done++
    if ($3 ~ /Progress/) progress++
    if ($3 ~ /To/) todo++
  }
  END { print "  Done: " done, "  In Progress: " progress, "  To Do: " todo }
'

# Key metrics
echo ""
echo "Key Metrics:"
echo "  On track items: $(jira-ticket-creator search --jql 'parent = PROJ-2000 AND priority != Critical' | wc -l)"
echo "  At risk items: $(jira-ticket-creator search --jql 'parent = PROJ-2000 AND (priority = Critical OR status = Blocked)' | wc -l)"

# Next steps
echo ""
echo "Next Steps:"
jira-ticket-creator search --jql "parent = PROJ-2000 AND status = 'To Do' AND priority = High" --format table
EOF

chmod +x daily-sync.sh weekly-executive-summary.sh

# Step 4: Run reviews
./daily-sync.sh
./weekly-executive-summary.sh
```

### Scenario 2: Cross-Functional Product Feature

```bash
#!/bin/bash
# Coordinate effort across frontend, backend, design teams

# Feature Epic: PROJ-5000 "Dark Mode Implementation"

# Step 1: Setup per-team views
jira-ticket-creator import --jql "key = PROJ-5000" --map-project "dark-mode-feature"

# Step 2: Track by team discipline
cat > feature-by-team.sh << 'EOF'
#!/bin/bash

FEATURE_EPIC="PROJ-5000"

echo "=== DARK MODE - BY TEAM ==="
echo ""

# Design tasks
echo "DESIGN TASKS:"
jira-ticket-creator search --jql "parent = $FEATURE_EPIC AND type = Task AND assignee in (designer1, designer2)" \
  --format table
echo ""

# Frontend implementation
echo "FRONTEND DEVELOPMENT:"
jira-ticket-creator search --jql "parent = $FEATURE_EPIC AND component = Frontend" \
  --format table
echo ""

# Backend support
echo "BACKEND SUPPORT:"
jira-ticket-creator search --jql "parent = $FEATURE_EPIC AND component = Backend" \
  --format table
echo ""

# QA testing
echo "QA TESTING:"
jira-ticket-creator search --jql "parent = $FEATURE_EPIC AND type = 'Test'" \
  --format table
EOF

chmod +x feature-by-team.sh
./feature-by-team.sh

# Step 3: Identify blockers between teams
echo "TEAM DEPENDENCIES:"
jira-ticket-creator search --jql "parent = $FEATURE_EPIC AND status = Blocked" --format table
```

### Scenario 3: Technical Debt Initiative with Multiple Repos

```bash
#!/bin/bash
# Track technical debt across multiple service repos

# Effort Epic: PROJ-3000 "API v2 Migration"
# Involves: api-service, auth-service, notification-service repos

cat > tech-debt-initiative.sh << 'EOF'
#!/bin/bash

INITIATIVE_EPIC="PROJ-3000"
TARGET_DATE="2025-06-30"

# Import
jira-ticket-creator import --jql "key = $INITIATIVE_EPIC" --map-project "api-v2-migration"

# Track effort across repos
echo "=== API v2 MIGRATION EFFORT ==="
echo "Target: $TARGET_DATE"
echo ""

# By service
for SERVICE in api-service auth-service notification-service; do
  echo "SERVICE: $SERVICE"
  jira-ticket-creator search --jql "parent = $INITIATIVE_EPIC AND labels = $SERVICE" --format table
  echo ""
done

# Progress summary
echo "OVERALL PROGRESS:"
jira-ticket-creator gantt --format ascii --weeks 8

# Risk report
echo ""
echo "RISKS:"
jira-ticket-creator search --jql "parent = $INITIATIVE_EPIC AND (status = Blocked OR priority = Critical OR updated <= -30d)" \
  --format table
EOF

chmod +x tech-debt-initiative.sh
./tech-debt-initiative.sh
```

---

## Effort Lifecycle

### Phase 1: Planning (Week 1)

```bash
# Create effort epic in JIRA
# Import it locally
jira-ticket-creator import --jql "key = PROJ-1000" --map-project my-effort

# Define all stories and tasks
# Estimate (add story points in JIRA)
# Create script for daily tracking
cp ~/templates/daily-check.sh ./effort-check.sh
```

### Phase 2: Execution (Weeks 2-8)

```bash
# Daily: Check blockers
./effort-check.sh

# Weekly: Generate report
./effort-weekly-report.sh

# Monitor: Use gantt chart
jira-ticket-creator gantt --format html --output gantt.html
open gantt.html
```

### Phase 3: Completion (Final Week)

```bash
# Verify all done
jira-ticket-creator search --jql "parent = PROJ-1000 AND status != Done" --format table

# Generate final report
jira-ticket-creator search --jql "parent = PROJ-1000" --format json > final-effort-report.json

# Archive effort
mv ~/efforts/my-effort ~/efforts/completed/my-effort-2025-04-30
```

---

## Tips & Best Practices

1. **Use consistent naming**: `project-effort-name` not `proj_effort` or `PROJECT-EFFORT`
2. **Tag tickets clearly**: Use labels and components for team/service
3. **Automate reports**: Create shell scripts for recurring reports
4. **Check weekly**: Don't wait for blockers to pile up
5. **Involve teams early**: Import and share dashboard with stakeholders
6. **Archive efforts**: Keep organized history of completed efforts

## See Also

- [Import Command](../cli/import) - Import and organize efforts
- [Project Mapping](../advanced/project-mapping) - Configure effort organization
- [Common Workflows](workflows) - Other automation examples
