---
layout: default
title: Project Mapping
parent: Advanced
nav_order: 1
has_toc: true
---

# Project Mapping Configuration

Organize tickets into logical projects independent of JIRA project keys.

## What is Project Mapping?

Project mapping allows you to:
- Group tickets from different JIRA projects into logical teams
- Create virtual "projects" based on ticket key prefixes
- Organize multi-team workflows
- Filter reports by logical project

## Configuration File

Default location: `~/.jira/project-mapping.json`

### Basic Structure

```json
{
 "mappings": {
 "backend": {
 "ticket_keys": ["PROJ", "API", "DB"],
 "description": "Backend Team"
 },
 "frontend": {
 "ticket_keys": ["UI", "WEB", "MOBILE"],
 "description": "Frontend Team"
 }
 }
}
```

## Configuration Examples

### Small Team Setup

```json
{
 "mappings": {
 "main": {
 "ticket_keys": ["PROJ"],
 "description": "Main Project"
 }
 }
}
```

### Multi-Team Project

```json
{
 "mappings": {
 "backend": {
 "ticket_keys": ["PROJ", "API", "SERVICE"],
 "description": "Backend Services Team"
 },
 "frontend": {
 "ticket_keys": ["UI", "WEB", "APP"],
 "description": "Frontend Team"
 },
 "mobile": {
 "ticket_keys": ["MOBILE", "IOS", "ANDROID"],
 "description": "Mobile Team"
 },
 "devops": {
 "ticket_keys": ["INFRA", "CI", "DEPLOY"],
 "description": "DevOps/Infrastructure"
 }
 }
}
```

### Enterprise Multi-Product Setup

```json
{
 "mappings": {
 "product-a-backend": {
 "ticket_keys": ["PRODA", "PRODAPI"],
 "description": "Product A Backend"
 },
 "product-a-frontend": {
 "ticket_keys": ["PRODUI"],
 "description": "Product A Frontend"
 },
 "product-b-backend": {
 "ticket_keys": ["PRODB", "PRODBAPI"],
 "description": "Product B Backend"
 },
 "product-b-frontend": {
 "ticket_keys": ["PRODBUI"],
 "description": "Product B Frontend"
 },
 "shared-services": {
 "ticket_keys": ["SHARED", "COMMON", "UTILS"],
 "description": "Shared Services"
 }
 }
}
```

## Using Project Mapping

### Create the Mapping File

```bash
cat > ~/.jira/project-mapping.json << 'EOF'
{
 "mappings": {
 "backend": {
 "ticket_keys": ["PROJ", "API"],
 "description": "Backend Team"
 },
 "frontend": {
 "ticket_keys": ["UI"],
 "description": "Frontend Team"
 }
 }
}
EOF
```

### Import Tickets

```bash
# Tickets with PROJ or API prefix will map to "backend"
jira-ticket-creator import --jql "project in (PROJ, API)"

# Tickets with UI prefix will map to "frontend"
jira-ticket-creator import --jql "project = UI"
```

### View by Project

```bash
# See only backend tickets
jira-ticket-creator team summary --project backend

# See only frontend tickets
jira-ticket-creator team summary --project frontend
```

## Inline Mapping Rules

For one-off imports without modifying the config file:

```bash
jira-ticket-creator import --jql "project in (PROJ, API)" \
 --map-rule "PROJ->backend" \
 --map-rule "API->backend"
```

## Custom Mapping Path

Use a different mapping file:

```bash
jira-ticket-creator import --jql "project = PROJ" \
 --mapping-path ~/my-custom-mappings.json
```

## Auto-Detection

If no mapping is found for a ticket key:

1. Check project mapping file
2. Check inline rules
3. Use `--map-project` value (if provided)
4. Leave project empty (no mapping)

## Examples by Use Case

### Microservices Architecture

Map service names to projects:

```json
{
 "mappings": {
 "auth-service": {
 "ticket_keys": ["AUTH", "AUTHAPI"],
 "description": "Authentication Service"
 },
 "payment-service": {
 "ticket_keys": ["PAY", "PAYAPI"],
 "description": "Payment Service"
 },
 "order-service": {
 "ticket_keys": ["ORDER", "ORDERAPI"],
 "description": "Order Service"
 },
 "user-service": {
 "ticket_keys": ["USER", "USERAPI"],
 "description": "User Service"
 }
 }
}
```

### Agency Multi-Client Setup

```json
{
 "mappings": {
 "client-a": {
 "ticket_keys": ["CA"],
 "description": "Client A Project"
 },
 "client-b": {
 "ticket_keys": ["CB"],
 "description": "Client B Project"
 },
 "client-c": {
 "ticket_keys": ["CC"],
 "description": "Client C Project"
 },
 "internal": {
 "ticket_keys": ["INT", "INFRA"],
 "description": "Internal/Infrastructure"
 }
 }
}
```

### Monorepo Packages

```json
{
 "mappings": {
 "api-package": {
 "ticket_keys": ["API", "APICORE"],
 "description": "API Package"
 },
 "ui-package": {
 "ticket_keys": ["UI", "UICORE"],
 "description": "UI Package"
 },
 "cli-package": {
 "ticket_keys": ["CLI", "CLICORE"],
 "description": "CLI Package"
 },
 "shared-package": {
 "ticket_keys": ["SHARED", "UTILS"],
 "description": "Shared Utilities"
 }
 }
}
```

## Workflows

### Initial Setup

```bash
#!/bin/bash
# setup-mapping.sh

# Step 1: Create mapping file
mkdir -p ~/.jira
cat > ~/.jira/project-mapping.json << 'EOF'
{
 "mappings": {
 "backend": {
 "ticket_keys": ["PROJ", "API", "DB"],
 "description": "Backend Team"
 },
 "frontend": {
 "ticket_keys": ["UI", "WEB"],
 "description": "Frontend Team"
 }
 }
}
EOF

echo " Created project mapping"

# Step 2: Import tickets
echo "Importing backend tickets..."
jira-ticket-creator import --jql "project = PROJ" --map-project backend

echo "Importing frontend tickets..."
jira-ticket-creator import --jql "project = UI" --map-project frontend

# Step 3: Verify
echo ""
echo "Backend summary:"
jira-ticket-creator team summary --project backend

echo ""
echo "Frontend summary:"
jira-ticket-creator team summary --project frontend
```

### Update Mapping

```bash
#!/bin/bash
# update-mapping.sh - Add new team to mapping

# Load existing mapping
MAPPING_FILE="$HOME/.jira/project-mapping.json"

# Add new team (you can edit manually or use jq)
jq '.mappings.devops = {
 "ticket_keys": ["INFRA", "CI", "DEPLOY"],
 "description": "DevOps Team"
}' "$MAPPING_FILE" > "$MAPPING_FILE.tmp" && mv "$MAPPING_FILE.tmp" "$MAPPING_FILE"

echo " Added DevOps team to mapping"
```

## Troubleshooting

### Tickets Not Mapping

**Problem:** Imported tickets have empty project field

**Solution:** Check mapping file and verify ticket key prefixes:

```bash
# View current mapping
cat ~/.jira/project-mapping.json

# Check actual ticket keys
jira-ticket-creator query --jql "project = PROJ" --fields key
```

### Multiple Prefixes Not Working

**Problem:** Only first prefix in list is being used

**Solution:** Ensure all prefixes are in the same project mapping:

```json
{
 "mappings": {
 "backend": {
 "ticket_keys": ["PROJ", "API", "DB"], // ← All in one array
 "description": "Backend Team"
 }
 }
}
```

### Want to Remap Existing Tickets

**Problem:** Already imported tickets don't have correct project

**Solution:** Re-import with `--update-existing` flag:

```bash
# Update existing mapping in file
# Then re-import
jira-ticket-creator import --jql "project = PROJ" \
 --map-project backend \
 --update-existing
```

## Best Practices

1. **Use consistent naming**: Use descriptive, lowercase project names
 ```json
 "backend"
 "Frontend" (inconsistent case)
 "be" (too short)
 ```

2. **Map related tickets together**: Group by team, service, or product
 ```json
 {
 "backend": {
 "ticket_keys": ["PROJ", "API", "DB"],
 "description": "Backend Team"
 }
 }
 ```

3. **Document in description**: Explain what each project contains
 ```json
 "description": "Backend Team - API, database, and services"
 ```

4. **Version your mapping**: Track changes in git
 ```bash
 git add ~/.jira/project-mapping.json
 git commit -m "Update project mapping for new team"
 ```

5. **Backup before changes**: Keep old version
 ```bash
 cp ~/.jira/project-mapping.json ~/.jira/project-mapping.json.bak
 ```

## Parent Ticket Mapping

Map parent tickets (Epics, Features) to work efforts so all child tickets automatically inherit the same project mapping.

### Concept

Parent tickets (Epics, Stories) can own multiple child tickets (Subtasks, Stories). When you map a parent ticket to a project, all its children are automatically mapped to the same project.

```
Epic: "Mobile App Redesign" → maps to "frontend" project
├── Subtask: "Design wireframes"     → inherits "frontend"
├── Subtask: "Implement components"  → inherits "frontend"
└── Subtask: "Write tests"           → inherits "frontend"
```

### Example: Product Feature Mapping

```json
{
  "mappings": {
    "checkout-flow": {
      "description": "Checkout Payment Feature",
      "parent_tickets": ["PROJ-500"],
      "include_children": true
    },
    "backend": {
      "ticket_keys": ["PROJ", "API"],
      "description": "Backend Team"
    }
  }
}
```

### Importing with Parent Mapping

When you import a parent ticket, all its children are automatically included:

```bash
# Import parent and all children
jira-ticket-creator import --jql "key = PROJ-500" --map-project checkout-flow

# Verify children were imported with same mapping
jira-ticket-creator search --jql "parent = PROJ-500"
# All will show project: "checkout-flow"
```

### Workflow Example: Sprint Planning

Map sprint epics to your work efforts:

```bash
# Map Q1 Planning epic to planning project
jira-ticket-creator import --jql "key = PROJ-1000" --map-project q1-planning

# All stories under PROJ-1000 are now part of "q1-planning"
jira-ticket-creator search --jql "parent = PROJ-1000"

# View by work effort
jira-ticket-creator team summary --project q1-planning
jira-ticket-creator gantt --format html --output q1-gantt.html
```

### Workflow Example: Multi-Team Epics

Organize complex features across multiple teams:

```bash
# Frontend epic
jira-ticket-creator import --jql "key = PROJ-200" --map-project frontend-epic

# Backend epic
jira-ticket-creator import --jql "key = PROJ-201" --map-project backend-epic

# DevOps epic
jira-ticket-creator import --jql "key = PROJ-202" --map-project devops-epic

# Each epic's children inherit parent mapping
jira-ticket-creator query --jql "parent in (PROJ-200, PROJ-201, PROJ-202)" --format json
```

### Filtering by Parent-Based Projects

```bash
# View all work in frontend epic
jira-ticket-creator search --jql "parent = PROJ-200"

# Generate report for specific work effort
jira-ticket-creator gantt --format html --output frontend-epic.html
```

### Tips

1. **Use descriptive work effort names** based on business deliverables:
   ```bash
   --map-project "mobile-redesign"      # Good
   --map-project "epic-1"               # Avoid
   ```

2. **Map entire features together**:
   ```bash
   # Create work effort for complete feature
   jira-ticket-creator import --jql "key = PROJ-EPIC" --map-project "feature-name"
   ```

3. **Track by business outcome**:
   ```bash
   # Instead of mapping by team, map by delivery
   --map-project "customer-onboarding"  # Feature name
   --map-project "api-performance"      # Business goal
   --map-project "q1-roadmap"           # Timeline
   ```

## Integration with Other Commands

### Team Reports

```bash
# View specific project
jira-ticket-creator team summary --project backend

# Filter assignments
jira-ticket-creator team assignments --project frontend

# Check timeline
jira-ticket-creator team timeline --project backend
```

### Gantt Charts

```bash
# Generate for specific project
jira-ticket-creator gantt --format html --output backend-gantt.html
# Then filter in the generated chart if needed
```

### Filtering

```bash
# Search with project context
jira-ticket-creator query --jql "project = PROJ" --format table
# Map through project mapping system
```

## See Also

- [Import Command](../cli/import) - Import tickets with mapping
- [Query Command](../cli/query) - Find tickets by criteria
