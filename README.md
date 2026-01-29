# Jira Ticket Creator CLI

A simple **Go CLI utility** to create Jira tickets, manage dependencies (blocked-by), track their statuses, and generate reports. Fully compatible with **Jira Cloud**.

---

## Features

- Create Jira tickets with summary, description, and blocked-by relationships.
- Validate ticket creation in Jira Cloud and log their status.
- Maintain a persistent report (`created_tickets.json`) of all created tickets.
- View a report of all tickets with blocked-by and current status.
- Update all ticket statuses from Jira Cloud to keep the report in sync.

---

## Requirements

- Go 1.21+
- Jira Cloud account with **API token**
- Jira project key

---

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/jira-ticket-creator.git
cd jira-ticket-creator
```

Build the CLI (optional):

```bash
go build -o jira-create
```

Or run directly:

```bash
go run main.go [flags] <command>
go run github.com/yourusername/jira-ticket-creator@main \
  -url https://yourcompany.atlassian.net \
  -email you@example.com \
  -token your-api-token \
  -project PROJECT1 \
  "Implement login feature" \
  "Login feature blocked until core setup is done" \
  "PROJECT1-12,PROJECT1-15"

```

---

## CLI Usage

### Flags

| Flag       | Description                                | Required |
|------------|--------------------------------------------|----------|
| `-url`     | Jira Cloud URL (`https://yourdomain.atlassian.net`) | ✅        |
| `-email`   | Jira account email                          | ✅        |
| `-token`  | Jira API token                              | ✅        |
| `-project` | Jira project key                             | ✅        |
| `-report`  | Report file path (default `created_tickets.json`) | ❌       |

### Commands

#### 1. Create a ticket

```bash
go run main.go 
  -url https://yourcompany.atlassian.net 
  -email you@example.com 
  -token your-api-token 
  -project PROJECT1 
  "Implement login feature" 
  "Login feature blocked until core setup is done" 
  "PROJECT1-12,PROJECT1-15"
```

- Arguments:
  1. Summary
  2. Description
  3. Optional comma-separated list of ticket keys it is blocked by.

Output:
```
Created issue: PROJECT1-101
Linked PROJECT1-12 -> PROJECT1-101
Linked PROJECT1-15 -> PROJECT1-101
All created tickets logged to created_tickets.json
```

#### 2. View report

```bash
go run main.go report
```

Sample output:
```
Created Jira tickets:
- PROJECT1-101 | Implement login feature | To Do | blocked by: [PROJECT1-12 PROJECT1-15]
- PROJECT1-102 | Setup database | In Progress | blocked by: []
```

#### 3. Update all ticket statuses

```bash
go run main.go update-status
```

Sample output:
```
Updating ticket statuses...
PROJECT1-101 -> In Progress
PROJECT1-102 -> Done
All statuses updated successfully.
```

---

## JSON Report Structure

```json
[
  {
    "key": "PROJECT1-101",
    "summary": "Implement login feature",
    "status": "In Progress",
    "blocked_by": ["PROJECT1-12", "PROJECT1-15"],
    "created_at": "2026-01-29T15:02:00Z"
  }
]
```

---

## Notes

- Works **fully with Jira Cloud** using email + API token.
- Ticket linking uses Jira's **Blocks** relationship.
- The `report` and `update-status` commands operate on the same JSON log file.

---

## Testing

Tests use a mock HTTP server to simulate Jira responses. Run:

```bash
go test -v ./...
```

---

## License

MIT License
EOF

