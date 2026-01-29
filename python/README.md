# JIRA Ticket Creator - Python Integration

Direct Python integration with JIRA Ticket Creator using C bindings (CGO). No subprocess calls or REST API - pure direct function calls.

## Features

- ✓ Direct function calls (no subprocess or REST API)
- ✓ No external dependencies beyond standard library
- ✓ Fast and efficient C bindings
- ✓ Type hints and documentation
- ✓ Cross-platform (Linux, macOS, Windows)

## Installation

### Build C Library

```bash
# Build the C library
make python-build

# Or manually
cd python
CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api/capi.go
```

### Install Python Package

```bash
# Install locally for development
make python-install

# Or manually
cd python
pip install -e .
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
ticket = client.create_ticket(
    summary="Implement new feature"
)
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
