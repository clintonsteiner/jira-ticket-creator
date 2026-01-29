# JIRA Ticket Creator - Python Integration

Direct Python integration with JIRA Ticket Creator using C bindings (CGO). No subprocess calls or REST API - pure direct function calls.

## Features

- ✓ Direct function calls (no subprocess or REST API)
- ✓ No external dependencies beyond standard library
- ✓ Fast and efficient C bindings
- ✓ Type hints and documentation
- ✓ Cross-platform (Linux, macOS, Windows)

## Installation

### From PyPI (Recommended)

```bash
pip install jira-ticket-creator
```

This includes the pre-compiled C library for your platform.

### From Source (Development)

#### 1. Build C Library

```bash
# Build the C library
make python-build

# Or manually
cd python
CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api/capi.go
```

#### 2. Install Python Package

```bash
# Install locally for development
make python-install

# Or manually
cd python
pip install -e .

# Or using modern build tools
pip install -e ".[dev]"
```

## Quick Start

```python
from jira_client import JiraClient

# Initialize client
client = JiraClient(
    url="https://company.atlassian.net",
    email="user@company.com",
    token="api-token",
    project="PROJ"
)

# Create a ticket
ticket = client.create_ticket(
    summary="Fix login bug",
    description="Login button not responding",
    issue_type="Bug",
    priority="High"
)

print(f"Created: {ticket['key']}")
```

## API Reference

### JiraClient

```python
from jira_client import JiraClient

client = JiraClient(
    url: str,              # JIRA base URL
    email: str,            # JIRA email
    token: str,            # JIRA API token
    project: str = None,   # Project key (e.g., "PROJ")
    ticket: str = None,    # Or ticket key (e.g., "PROJ-123")
    lib_path: str = None   # Path to libjira.so (auto-detected)
)
```

### Methods

#### create_ticket()

```python
ticket = client.create_ticket(
    summary: str,                    # Required: Ticket summary
    description: str = "",           # Ticket description
    issue_type: str = "Task",       # Task, Story, Bug, Epic, Subtask
    priority: str = "Medium",       # Lowest, Low, Medium, High, Highest
    assignee: str = "",             # Assignee email
    labels: List[str] = None,       # List of labels
    blocked_by: List[str] = None    # Blocking ticket keys
) -> Dict[str, str]
```

Returns dictionary with: `key`, `id`, `url`

#### get_ticket()

```python
ticket = client.get_ticket(ticket_key: str) -> Dict[str, Any]
```

Retrieve ticket details by key.

Returns dictionary with: `key`, `id`, `summary`, `description`, `status`, `issue_type`, `priority`, `assignee`, `labels`, `url`

#### search()

```python
# Search using JQL
results = client.search(jql: str = "") -> Dict[str, Any]

# Search using keyword arguments
results = client.search(
    status="In Progress",           # Filter by status
    key="PROJ-123",                 # Search by ticket key
    summary="feature",              # Search by summary text
    assignee="user@company.com",    # Filter by assignee
    issue_type="Bug"                # Filter by issue type
) -> Dict[str, Any]
```

Returns dictionary with: `total`, `count`, `tickets`

#### update_ticket()

```python
response = client.update_ticket(
    ticket_key: str,                # Ticket to update
    summary: str = None,            # New summary
    description: str = None,        # New description
    priority: str = None,           # New priority
    assignee: str = None            # New assignee email
) -> Dict[str, Any]
```

#### extract_project_key()

```python
project = client.extract_project_key("PROJ-123")  # Returns "PROJ"
```

#### get_version()

```python
version = client.get_version()  # Returns library version
```

## Examples

### Create a simple task

```python
from jira_client import JiraClient

client = JiraClient(
    url="https://company.atlassian.net",
    email="user@company.com",
    token="your-api-token",
    project="PROJ"
)

ticket = client.create_ticket(summary="Implement new feature")
print(f"Created: {ticket['key']}")
```

### Create a bug with high priority

```python
ticket = client.create_ticket(
    summary="Fix critical bug",
    issue_type="Bug",
    priority="Highest",
    labels=["critical", "urgent"]
)
```

### Create with assignee and blockers

```python
ticket = client.create_ticket(
    summary="Update API endpoints",
    assignee="john@company.com",
    blocked_by=["PROJ-123", "PROJ-124"]
)
```

### Get ticket details

```python
ticket = client.get_ticket("PROJ-123")
print(f"Status: {ticket['status']}")
print(f"Assignee: {ticket['assignee']}")
print(f"Labels: {ticket['labels']}")
```

### Search tickets

```python
# Search by JQL
results = client.search(jql='project = PROJ AND status = "In Progress"')
print(f"Found {results['count']} tickets")

# Search by keyword
results = client.search(status="In Progress")
for ticket in results['tickets']:
    print(f"  - {ticket['key']}: {ticket['summary']}")

# Search by multiple criteria
results = client.search(
    status="To Do",
    priority="High",
    assignee="john@company.com"
)
```

### Update ticket

```python
response = client.update_ticket(
    "PROJ-123",
    priority="High",
    description="Updated description",
    assignee="jane@company.com"
)
```

### Use ticket key instead of project key

```python
client = JiraClient(
    url="https://company.atlassian.net",
    email="user@company.com",
    token="api-token",
    ticket="PROJ-100"  # Project auto-extracted as "PROJ"
)
```

## Testing

```bash
# Run Python example
make python-example

# Run Python tests
make python-test
```

## Performance

- **No subprocess overhead**: Direct C function calls
- **No REST API latency**: Compiled Go binary
- **Minimal memory**: Python ctypes wrapper is lightweight
- **Fast execution**: ~1-5ms per operation

## Supported Platforms

- Linux (x86_64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (64-bit)

## Troubleshooting

### "Could not find libjira library"

Build the library:
```bash
make python-build
```

### "CGO_ENABLED=1 required"

Install Go with CGO support:
```bash
go install -tag cgo
```

### Library not found in custom location

Specify lib_path:
```python
client = JiraClient(
    url="...",
    email="...",
    token="...",
    project="...",
    lib_path="/path/to/libjira.so"
)
```

## Development

```bash
# Setup development environment
make dev-setup

# Format code
make fmt

# Lint
make lint

# Build everything
make all
```

## License

MIT License - Same as JIRA Ticket Creator
