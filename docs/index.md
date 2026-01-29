---
layout: default
title: JIRA Ticket Creator - CLI Documentation
---

# 🎫 JIRA Ticket Creator

A powerful, feature-rich command-line tool for managing JIRA tickets with advanced features like batch operations, dependency visualization, multiple reporting formats, and team collaboration tools.

## ✨ Key Features

- **🚀 Fast Ticket Creation** - Create tickets individually or in batch from CSV/JSON
- **📊 Multiple Report Formats** - Table, JSON, CSV, Markdown, HTML
- **🔗 Dependency Visualization** - ASCII trees, Mermaid diagrams, Graphviz DOT format
- **👥 Team Collaboration** - Track creator, assignee, timeline, and workload
- **📈 Gantt Charts** - Workload visualization by resource
- **🔍 Advanced Search** - JQL queries with flexible output
- **📥 Bulk Import** - Import existing tickets with project mapping
- **⚙️ Project Management** - Executive dashboards and risk assessment
- **🔄 Batch Operations** - Transition, update, and manage multiple tickets
- **📋 Templates** - Built-in and custom ticket templates
- **🐚 Shell Completions** - Bash, Zsh, Fish, PowerShell support

## 🚀 Quick Start

### Installation

```bash
# Build from source
git clone https://github.com/clintonsteiner/jira-ticket-creator.git
cd jira-ticket-creator
go build -o jira-ticket-creator ./cmd/jira-ticket-creator

# Or install from GitHub
go install github.com/clintonsteiner/jira-ticket-creator/cmd/jira-ticket-creator@latest
```

### Setup (Choose One Method)

**Method A: Environment Variables (Recommended)**
```bash
export JIRA_URL=https://your-company.atlassian.net
export JIRA_EMAIL=your-email@company.com
export JIRA_TOKEN=your-api-token
export JIRA_PROJECT=PROJ
```

**Method B: Config File**
```bash
cat > ~/.jirarc << EOF
jira:
  url: https://your-company.atlassian.net
  email: your-email@company.com
  token: your-api-token
  project: PROJ
EOF
chmod 600 ~/.jirarc
```

### Create Your First Ticket

```bash
jira-ticket-creator create --summary "My first ticket"
```

## 📚 Documentation

### Getting Started
- [Installation Guide](./docs/getting-started/)
- [Configuration](./docs/configuration/)
- [First Ticket](./docs/first-ticket/)

### CLI Commands
- [Create Tickets](./docs/cli/create/)
- [Search & Query](./docs/cli/search/)
- [Update & Transition](./docs/cli/update/)
- [Batch Operations](./docs/cli/batch/)
- [Reports](./docs/cli/reports/)
- [Gantt Charts](./docs/cli/gantt/)
- [Team Reports](./docs/cli/team/)
- [Project Management](./docs/cli/pm/)
- [Visualizations](./docs/cli/visualize/)

### Advanced Topics
- [Project Mapping](./docs/advanced/project-mapping/)
- [Dependency Management](./docs/advanced/dependencies/)
- [Templates](./docs/advanced/templates/)
- [Shell Completions](./docs/advanced/completions/)

### API Reference
- [Go Client Library](./docs/api/go-client/)
- [Configuration API](./docs/api/config/)
- [Storage API](./docs/api/storage/)
- [Report Generators](./docs/api/reports/)

## 💡 Common Use Cases

### Create a Ticket
```bash
jira-ticket-creator create \
  --summary "Implement OAuth 2.0" \
  --description "Add OAuth authentication" \
  --type Story \
  --priority High \
  --assignee john@company.com
```

### Search for Tickets
```bash
jira-ticket-creator search --jql "project = PROJ AND status = 'To Do'"
jira-ticket-creator query --jql "priority = Critical" --format json
```

### Import Tickets with Project Mapping
```bash
jira-ticket-creator import --jql "project = PROJ" \
  --map-project backend \
  --map-rule "PROJ->backend" \
  --map-rule "API->backend"
```

### Generate Reports
```bash
jira-ticket-creator report --format html --output report.html
jira-ticket-creator team summary --project backend
jira-ticket-creator gantt --format mermaid --output gantt.md
```

### Batch Create Tickets
```bash
jira-ticket-creator batch create --input tickets.csv --format csv
```

### View Project Management Dashboard
```bash
jira-ticket-creator pm dashboard    # Executive summary
jira-ticket-creator pm hierarchy    # Ticket relationships
jira-ticket-creator pm risk         # Risk assessment
```

## 🔧 Architecture

The tool is built with a clean, modular architecture:

```
├── cmd/
│   └── jira-ticket-creator/      # CLI entry point
├── internal/
│   ├── config/                    # Configuration management
│   ├── jira/                      # JIRA API client
│   ├── storage/                   # JSON-based storage
│   ├── batch/                     # Batch processing
│   ├── reports/                   # Report generators
│   ├── templates/                 # Template engine
│   └── interactive/               # Interactive prompts
└── pkg/
    └── cli/
        ├── commands/              # CLI command implementations
        └── errors.go              # Error handling
```

## 📊 Example Workflows

### Multi-Team Project Setup
```bash
# Create project mapping
cat > ~/.jira/project-mapping.json << EOF
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

# Import tickets by team
jira-ticket-creator import --jql "project = PROJ" --map-project backend
jira-ticket-creator import --jql "project = UI" --map-project frontend

# View team reports
jira-ticket-creator team summary --project backend
jira-ticket-creator gantt --format html --output workload.html
```

### Executive Reporting
```bash
# Generate comprehensive dashboard
jira-ticket-creator pm dashboard > dashboard.txt

# Create Gantt chart
jira-ticket-creator gantt --format mermaid --output gantt.md

# Export to multiple formats
jira-ticket-creator report --format html --output report.html
jira-ticket-creator report --format csv --output report.csv
```

### Dependency Management
```bash
# Visualize ticket dependencies
jira-ticket-creator visualize --format tree
jira-ticket-creator visualize --format mermaid --output dependencies.md

# Check for blocked items
jira-ticket-creator pm risk
```

## 🆘 Troubleshooting

### Authentication Issues
```bash
# Test your credentials
jira-ticket-creator search --key PROJ-1

# Common issues:
# - Token must be API token, not password
# - User needs project permissions
# - URL format: https://domain.atlassian.net (no trailing slash)
```

### Build Issues
```bash
# Update dependencies
go mod tidy
go mod download

# Rebuild
go build -o jira-ticket-creator ./cmd/jira-ticket-creator
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - See LICENSE file for details

## 🐙 Repository

**GitHub**: [github.com/clintonsteiner/jira-ticket-creator](https://github.com/clintonsteiner/jira-ticket-creator)

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/clintonsteiner/jira-ticket-creator/issues)
- **Documentation**: This website
- **README**: [Project README](https://github.com/clintonsteiner/jira-ticket-creator#readme)

---

**Made with ❤️ for teams managing JIRA at scale**
