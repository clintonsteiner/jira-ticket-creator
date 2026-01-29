# JIRA Ticket Creator Documentation

This directory contains the GitHub Pages documentation for JIRA Ticket Creator.

## Structure

```
docs/
 _config.yml # Jekyll configuration
 Gemfile # Ruby dependencies
 index.md # Homepage
 getting-started.md # Getting started guide
 cli/ # CLI command documentation
 create.md
 query.md
 import.md
 gantt.md
 ...
 api/ # API reference
 go-client.md
 config.md
 ...
 advanced/ # Advanced topics
 project-mapping.md
 ...
 _site/ # Generated site (ignored)
```

## Building Locally

### Prerequisites

- Ruby 3.1+
- Bundler

### Setup

```bash
cd docs
bundle install
```

### Build

```bash
bundle exec jekyll build
```

### Serve Locally

```bash
bundle exec jekyll serve
# Open http://localhost:4000/jira-ticket-creator
```

## Publishing

The documentation is automatically built and published to GitHub Pages when you push to the `master` branch.

The GitHub Actions workflow:
1. Triggers on pushes to `master` (if docs/ or README.md changed)
2. Builds the Jekyll site
3. Deploys to GitHub Pages
4. Available at: https://clintonsteiner.github.io/jira-ticket-creator

## Writing Documentation

### Markdown Format

All pages are standard Markdown with YAML front matter:

```markdown
---
layout: default
title: Page Title
parent: Parent Page Title # Optional, for navigation hierarchy
---

# Page Title

Your content here...
```

### Directory Structure

Organize by feature/topic:
- `docs/cli/` - CLI command documentation
- `docs/api/` - API reference
- `docs/advanced/` - Advanced topics
- `docs/examples/` - Example scripts

### Links

Use relative links:
```markdown
[Link text](./other-page.md)
[Link text](../api/go-client.md)
```

### Code Blocks

Use syntax highlighting:
````markdown
```bash
jira-ticket-creator --help
```

```go
client := jira.NewClient(url, email, token)
```

```json
{"key": "value"}
```
````

## Theme

The documentation uses the **Just the Docs** Jekyll theme, which provides:
- Built-in search functionality
- Left sidebar navigation with hierarchy
- Mobile-responsive design
- Proper semantic HTML

Configuration is in `_config.yml`.

## Navigation

Just the Docs uses YAML front matter to build the navigation tree:

```yaml
---
title: Page Title
parent: Parent Page Title  # Optional, creates hierarchy
nav_order: 2              # Optional, controls order
---
```

For example:
```yaml
---
title: Create Command
parent: CLI Commands
nav_order: 1
---
```

Creates:
```
CLI Commands (parent)
└── Create Command (child, nav_order: 1)
```

## Search

Search is automatically enabled. Users can search by:
- Page titles
- Headings (level 2+)
- Body text

## See Also

- [Jekyll Documentation](https://jekyllrb.com/docs/)
- [Just the Docs Theme](https://just-the-docs.github.io/just-the-docs/)
- [GitHub Pages Docs](https://docs.github.com/en/pages)
