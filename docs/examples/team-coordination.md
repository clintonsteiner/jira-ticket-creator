---
layout: default
title: Multi-Team Coordination
parent: Examples & Guides
nav_order: 5
has_toc: true
---

# Multi-Team Coordination & Dependencies

Learn how to coordinate work across multiple teams, manage dependencies, and ensure smooth collaboration on shared efforts.

## The Challenge

When efforts span multiple teams:
- Frontend needs Backend APIs before they can start
- QA can't test until features are code-complete
- Design feedback affects both frontend and backend
- Multiple teams working on same epic
- Unclear who's blocking whom

This guide shows how to use jira-ticket-creator to solve these problems.

---

## Dependency Management

### Identifying Dependencies

```bash
cat > show-dependencies.sh << 'EOF'
#!/bin/bash
# Show dependency graph for an effort

EPIC="PROJ-5000"

echo "=== DEPENDENCY MAP FOR EFFORT ==="
echo ""

# Blocking relationships
echo "1. TICKETS THAT BLOCK OTHERS"
echo "────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks" --format json | \
  jq -r '.[] | "\(.key) blocks: \([.fields.issuelinks[] | select(.type.name=="Blocks") | .outwardIssue.key] | join(", "))"' | \
  sort

echo ""
echo "2. TICKETS BLOCKED BY OTHERS"
echo "────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND linked = 'is blocked by'" --format json | \
  jq -r '.[] | "\(.key) is blocked by: \([.fields.issuelinks[] | select(.type.name=="is blocked by") | .inwardIssue.key] | join(", "))"' | \
  sort

echo ""
echo "3. DEPENDENCY CHAIN (CRITICAL PATH)"
echo "────────────────────────────────────────────────────"
echo "Items that must complete before others can start:"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks AND status != Done" --format table

echo ""
echo "4. BLOCKED BY EXTERNAL TEAMS"
echo "────────────────────────────────────────────────────"
echo "Items waiting on other teams:"
jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" --format table
EOF

chmod +x show-dependencies.sh
./show-dependencies.sh
```

### Creating Dependency Report

```bash
cat > dependency-report.sh << 'EOF'
#!/bin/bash
# Generate dependency report for cross-team sync

EPIC="PROJ-5000"
OUTPUT_DIR="./reports/$(date +%Y-%m-%d)"
mkdir -p "$OUTPUT_DIR"

echo "CROSS-TEAM DEPENDENCY REPORT"
echo "$(date)"
echo "=================================================="
echo ""

echo "FOR: Effort Epic $EPIC"
echo ""

# 1. Team readiness
echo "1. TEAM READINESS"
echo "────────────────────────────────────────────────────────────────"

for TEAM in frontend backend design qa; do
  echo ""
  echo "$TEAM TEAM:"
  TOTAL=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z)" 2>/dev/null | wc -l)
  DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z) AND status = Done" 2>/dev/null | wc -l)
  BLOCKED=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z) AND status = Blocked" 2>/dev/null | wc -l)

  echo "  Total work: $TOTAL items"
  echo "  Completed: $DONE items"
  echo "  Blocked: $BLOCKED items"

  if [ $BLOCKED -gt 0 ]; then
    echo "  Blocked by:"
    jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z) AND status = Blocked" --format table | head -3
  fi
done

echo ""

# 2. Cross-team dependencies
echo "2. CROSS-TEAM DEPENDENCIES"
echo "────────────────────────────────────────────────────────────────"

echo ""
echo "FRONTEND WAITING FOR BACKEND:"
jira-ticket-creator search --jql "parent = $EPIC AND component = Frontend AND status = Blocked AND description ~ API" --format table

echo ""
echo "BACKEND WAITING FOR DESIGN:"
jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND status = Blocked AND description ~ Design" --format table

echo ""
echo "QA WAITING FOR FEATURE COMPLETE:"
jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND status = 'Waiting for Dev'" --format table

echo ""

# 3. Critical path
echo "3. CRITICAL PATH (Must complete first)"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks AND status != Done" --format table | head -10

echo ""

# 4. Risk mitigation
echo "4. AT-RISK HANDOFFS"
echo "────────────────────────────────────────────────────────────────"
echo "Tickets approaching handoff deadline:"
jira-ticket-creator search --jql "parent = $EPIC AND status = 'Code Review' AND duedate <= endOfWeek()" --format table

echo ""
echo "Report saved to: $OUTPUT_DIR"
EOF

chmod +x dependency-report.sh
./dependency-report.sh
```

---

## Cross-Team Sync Template

### Weekly Cross-Team Meeting

```bash
cat > cross-team-sync.sh << 'EOF'
#!/bin/bash
# Generate agenda for cross-team sync meeting

EPIC="PROJ-5000"
TEAMS=("frontend" "backend" "design" "qa")

echo "CROSS-TEAM SYNC AGENDA"
echo "$(date +%A, %B %d, %Y - %H:%M %Z)"
echo "=================================================="
echo ""

echo "EFFORT: $EPIC"
echo "Attendees: Team Leads from @frontend @backend @design @qa"
echo ""

echo "───────────────────────────────────────────────────────────────"
echo "1. CURRENT STATUS (5 min)"
echo "───────────────────────────────────────────────────────────────"
echo ""

TOTAL=$(jira-ticket-creator search --jql "parent = $EPIC" 2>/dev/null | wc -l)
DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Done" 2>/dev/null | wc -l)
PERCENT=$((DONE * 100 / TOTAL))

echo "Overall Progress: $DONE/$TOTAL complete ($PERCENT%)"
echo ""

echo "By Team:"
for TEAM in "${TEAMS[@]}"; do
  TEAM_TOTAL=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z)" 2>/dev/null | wc -l)
  TEAM_DONE=$(jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z) AND status = Done" 2>/dev/null | wc -l)
  if [ $TEAM_TOTAL -gt 0 ]; then
    TEAM_PCT=$((TEAM_DONE * 100 / TEAM_TOTAL))
    echo "  $TEAM: $TEAM_DONE/$TEAM_TOTAL ($TEAM_PCT%)"
  fi
done

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "2. BLOCKERS & RISKS (10 min)"
echo "───────────────────────────────────────────────────────────────"
echo ""

BLOCKERS=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" 2>/dev/null | wc -l)
echo "BLOCKED TICKETS ($BLOCKERS):"
jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked" --format table | head -10
echo ""

CRITICAL=$(jira-ticket-creator search --jql "parent = $EPIC AND priority = Critical" 2>/dev/null | wc -l)
echo "CRITICAL PRIORITY ($CRITICAL):"
jira-ticket-creator search --jql "parent = $EPIC AND priority = Critical AND status != Done" --format table | head -5
echo ""

echo "───────────────────────────────────────────────────────────────"
echo "3. DEPENDENCIES & HANDOFFS (10 min)"
echo "───────────────────────────────────────────────────────────────"
echo ""
echo "Who's waiting on whom:"
echo ""
echo "Frontend waiting for Backend:"
jira-ticket-creator search --jql "parent = $EPIC AND component = Frontend AND linked = 'is blocked by' AND type = Task" --format table | head -5

echo ""
echo "Backend waiting for Design:"
jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND linked = 'is blocked by' AND description ~ Design" --format table | head -5

echo ""
echo "QA waiting for Feature:"
jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND status = Waiting" --format table | head -5

echo ""
echo "───────────────────────────────────────────────────────────────"
echo "4. NEXT WEEK'S PLAN (10 min)"
echo "───────────────────────────────────────────────────────────────"
echo ""
echo "What each team commits to completing next week:"
echo ""

for TEAM in "${TEAMS[@]}"; do
  echo "$TEAM:"
  jira-ticket-creator search --jql "parent = $EPIC AND component = $(echo $TEAM | tr a-z A-Z) AND priority = High AND status != Done" \
    --format table | head -3
  echo ""
done

echo "───────────────────────────────────────────────────────────────"
echo "5. ACTION ITEMS & DECISIONS (10 min)"
echo "───────────────────────────────────────────────────────────────"
echo ""
echo "Items requiring team leads to make decisions or take action:"
echo ""
jira-ticket-creator search --jql "parent = $EPIC AND label = decision-needed" --format table

echo ""
echo "Meeting notes to be updated in JIRA comments"
echo "Next sync: Same time next week"
EOF

chmod +x cross-team-sync.sh
```

---

## Team Handoff Workflow

### Frontend to Backend Handoff

```bash
cat > frontend-backend-handoff.sh << 'EOF'
#!/bin/bash
# Manage Frontend -> Backend handoff

EPIC="PROJ-5000"
HANDOFF_LABEL="frontend-handoff-ready"

echo "FRONTEND -> BACKEND HANDOFF CHECKLIST"
echo "=================================================="
echo ""

# 1. What's ready for handoff
echo "1. FRONTEND WORK READY FOR HANDOFF"
echo "────────────────────────────────────────────────────────────────"
READY=$(jira-ticket-creator search --jql "parent = $EPIC AND component = Frontend AND label = $HANDOFF_LABEL" 2>/dev/null | wc -l)
echo "Frontend stories marked ready: $READY"
echo ""
jira-ticket-creator search --jql "parent = $EPIC AND component = Frontend AND label = $HANDOFF_LABEL" --format table
echo ""

# 2. What Backend is blocked on
echo "2. BACKEND BLOCKED WAITING FOR FRONTEND"
echo "────────────────────────────────────────────────────────────────"
BLOCKED=$(jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND status = Blocked" 2>/dev/null | wc -l)
echo "Backend tasks blocked: $BLOCKED"
echo ""
jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND status = Blocked" --format table
echo ""

# 3. API contracts
echo "3. API CONTRACTS"
echo "────────────────────────────────────────────────────────────────"
echo "API design tickets that must be completed before handoff:"
jira-ticket-creator search --jql "parent = $EPIC AND component = API AND type = Task AND status != Done" --format table
echo ""

# 4. Handoff checklist
echo "4. HANDOFF CHECKLIST"
echo "────────────────────────────────────────────────────────────────"
echo ""
echo "☐ All API contracts documented"
echo "☐ Frontend implementation complete"
echo "☐ Frontend tests passing"
echo "☐ Code review approved"
echo "☐ Demo completed"
echo "☐ Backend team has API specs"
echo "☐ Backend team has access to prototype"
echo ""

echo "5. NEXT STEPS FOR BACKEND"
echo "────────────────────────────────────────────────────────────────"
echo "Backend team should:"
echo "1. Review API specs"
echo "2. Update stories with API implementation tasks"
echo "3. Update dependencies in JIRA"
echo "4. Change status from Blocked to 'Ready to Start'"
echo ""

echo "Once Backend is ready:"
jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND status = 'Ready to Start'" --format table
EOF

chmod +x frontend-backend-handoff.sh
```

### Backend to QA Handoff

```bash
cat > backend-qa-handoff.sh << 'EOF'
#!/bin/bash
# Manage Backend -> QA handoff

EPIC="PROJ-5000"

echo "BACKEND -> QA HANDOFF CHECKLIST"
echo "=================================================="
echo ""

# 1. What's ready for QA
echo "1. BACKEND FEATURES READY FOR TESTING"
echo "────────────────────────────────────────────────────────────────"
READY=$(jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND label = qa-ready" 2>/dev/null | wc -l)
echo "Backend features marked ready for QA: $READY"
echo ""
jira-ticket-creator search --jql "parent = $EPIC AND component = Backend AND label = qa-ready" --format table
echo ""

# 2. QA blocked waiting
echo "2. QA WAITING FOR BACKEND"
echo "────────────────────────────────────────────────────────────────"
BLOCKED=$(jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND status = Blocked" 2>/dev/null | wc -l)
echo "QA test plans blocked: $BLOCKED"
echo ""
jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND status = Blocked" --format table
echo ""

# 3. Test environment
echo "3. TEST ENVIRONMENT READINESS"
echo "────────────────────────────────────────────────────────────────"
echo "Infrastructure ready?"
jira-ticket-creator search --jql "parent = $EPIC AND component = Infra AND label = test-environment" --format table
echo ""

# 4. Test data
echo "4. TEST DATA PREPARED"
echo "────────────────────────────────────────────────────────────────"
jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND type = 'Test' AND status != Done" --format table | head -10
echo ""

# 5. Handoff checklist
echo "5. QA HANDOFF CHECKLIST"
echo "────────────────────────────────────────────────────────────────"
echo ""
echo "☐ All features merged to main branch"
echo "☐ Feature deployed to staging/QA environment"
echo "☐ Deployment scripts validated"
echo "☐ Rollback plan prepared"
echo "☐ Data migration scripts tested"
echo "☐ Test environment is stable"
echo "☐ Test data sets prepared"
echo "☐ QA team has access"
echo ""

echo "6. QA TEST EXECUTION PLAN"
echo "────────────────────────────────────────────────────────────────"
echo "QA team to execute:"
jira-ticket-creator search --jql "parent = $EPIC AND component = QA AND type = Test AND status = Ready" --format table
EOF

chmod +x backend-qa-handoff.sh
```

---

## Dependency Chain Tracking

### Critical Path Analysis

```bash
cat > critical-path.sh << 'EOF'
#!/bin/bash
# Show critical path - items that must complete in order

EPIC="PROJ-5000"

echo "CRITICAL PATH ANALYSIS"
echo "Timeline-blocking items"
echo "=================================================="
echo ""

echo "Items that other work depends on"
echo "--------------------------------------------"
echo ""

# Get items that block others
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks" --format json | \
  jq -r '.[] |
    {
      key: .key,
      summary: .fields.summary,
      status: .fields.status.name,
      duedate: .fields.duedate,
      assignee: .fields.assignee.name // "Unassigned"
    } |
    "[\(.status)] \(.key): \(.summary)\n        Assigned: \(.assignee)\n        Due: \(.duedate // "No due date")"
  ' | column

echo ""
echo "Timeline:"
echo "--------------------------------------------"

# Show critical path timeline
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks ORDER BY duedate ASC" --format json | \
  jq -r '.[] | "\(.fields.duedate // "TBD"): \(.key) - \(.fields.summary)"'

echo ""
echo "Risks:"
echo "--------------------------------------------"

# Items in critical path that are blocked
echo "Critical path items that are BLOCKED:"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks AND status = Blocked" --format table

echo ""
echo "Critical path items OVERDUE:"
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks AND duedate < now()" --format table
EOF

chmod +x critical-path.sh
```

---

## Team Communication

### Automated Notifications

```bash
cat > team-notifications.sh << 'EOF'
#!/bin/bash
# Generate notifications for teams about changes

EPIC="PROJ-5000"

# What changed in last 24 hours
echo "CHANGES IN LAST 24 HOURS:"
jira-ticket-creator search --jql "parent = $EPIC AND updated >= -1d" --format json | \
  jq -r '.[] |
    "[\(.fields.updated | split("T")[0])] \(.fields.status.name): \(.key) - \(.fields.summary)"
  '

echo ""
echo "FOR TEAM LEADS:"
echo "----"

# New blockers
NEW_BLOCKERS=$(jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked AND updated >= -1d" 2>/dev/null | wc -l)
if [ $NEW_BLOCKERS -gt 0 ]; then
  echo "⚠️  NEW BLOCKERS ($NEW_BLOCKERS):"
  jira-ticket-creator search --jql "parent = $EPIC AND status = Blocked AND updated >= -1d" --format table
fi

echo ""

# Unassigned critical
UNASSIGNED=$(jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY AND priority = Critical" 2>/dev/null | wc -l)
if [ $UNASSIGNED -gt 0 ]; then
  echo "🔥 UNASSIGNED CRITICAL ($UNASSIGNED):"
  jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY AND priority = Critical" --format table
fi

echo ""

# Due soon
DUE_SOON=$(jira-ticket-creator search --jql "parent = $EPIC AND duedate >= now() AND duedate <= endOfWeek() AND status != Done" 2>/dev/null | wc -l)
if [ $DUE_SOON -gt 0 ]; then
  echo "📅 DUE THIS WEEK ($DUE_SOON):"
  jira-ticket-creator search --jql "parent = $EPIC AND duedate >= now() AND duedate <= endOfWeek() AND status != Done" --format table
fi
EOF

chmod +x team-notifications.sh
```

---

## Best Practices for Multi-Team Efforts

### 1. Clear Ownership
```bash
# Every ticket should have an assignee
# Every Epic should have a lead
jira-ticket-creator search --jql "parent = $EPIC AND assignee is EMPTY" --format table
# If there are unassigned items, assign them

# Mark parent epic with team
# Add label: "epic-lead:john" "epic-lead:jane"
```

### 2. Dependencies Must Be Explicit
```bash
# Use JIRA linking to show dependencies
# Don't assume implicit order
jira-ticket-creator search --jql "parent = $EPIC AND linked = Blocks" --format table
```

### 3. Regular Syncs
```bash
# Daily: Check blockers
./show-dependencies.sh

# Weekly: Cross-team meeting with dependency report
./cross-team-sync.sh

# Before handoff: Use handoff checklist
./frontend-backend-handoff.sh
```

### 4. Communication
```bash
# Automatic notifications when status changes
# Scheduled sync meeting agendas
# Visible blockers and critical items
```

### 5. Escalation Path
```bash
# Clear who to contact if blocked
# Effort lead has authority to make trade-off decisions
# Technical leads to resolve design conflicts
```

---

## See Also

- [Managing JIRA Efforts](effort-management) - Setting up efforts
- [Building Effort Dashboards](effort-dashboards) - Monitoring progress
- [Common Workflows](workflows) - Automation examples
