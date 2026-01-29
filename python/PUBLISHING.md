# Publishing JIRA Ticket Creator to PyPI

This guide explains how to build and publish the `jira-ticket-creator` Python package to PyPI.

## Prerequisites

1. **Python 3.6+** with pip
2. **Go 1.13+** with CGO enabled
3. **Build tools**: setuptools, wheel, build, twine
4. **PyPI account** at https://pypi.org
5. **API token** from PyPI (stored in `~/.pypirc`)

## Setup

### 1. Create PyPI Account

If you don't have a PyPI account, create one at https://pypi.org/account/register/

### 2. Generate API Token

1. Log in to PyPI
2. Go to Account Settings → API tokens
3. Create a new token for the project
4. Copy the token (format: `pypi-AgE...`)

### 3. Configure Local Environment

Create or update `~/.pypirc`:

```ini
[distutils]
index-servers =
    pypi
    testpypi

[pypi]
repository = https://upload.pypi.org/legacy/
username = __token__
password = pypi-AgE...

[testpypi]
repository = https://test.pypi.org/legacy/
username = __token__
password = pypi-AgE...
```

Or set environment variables:

```bash
export PYPI_API_TOKEN=pypi-AgE...
export TEST_PYPI_API_TOKEN=pypi-AgE...
```

## Building the Package

### Using Makefile (Recommended)

```bash
# Check distribution for issues
make python-publish-check

# Build distribution package
make python-build-dist
```

### Manual Build

```bash
cd python

# Install build tools
pip install --upgrade build twine

# Build C library first
CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api/capi.go

# Build Python distribution
python -m build

# Check for issues
twine check dist/*
```

## Publishing

### To Test PyPI (Recommended First)

Always test in test.pypi.org before publishing to production:

```bash
# Using Makefile
make python-publish-test

# Or manually
cd python
twine upload --repository testpypi dist/*
```

Then test the installation:

```bash
pip install --index-url https://test.pypi.org/simple/ jira-ticket-creator
```

### To Production PyPI

Once testing is successful:

```bash
# Using Makefile
make python-publish

# Or manually
cd python
twine upload dist/*
```

## Continuous Integration / Continuous Deployment

The repository includes a GitHub Actions workflow (`.github/workflows/python-publish.yml`) that:

1. Automatically publishes on GitHub releases
2. Supports manual workflow dispatch for testing
3. Builds for multiple Python versions
4. Tests installation on multiple platforms

### Triggering Automated Publishing

**Option 1: On GitHub Release**

Simply create a release on GitHub, and the package will be automatically published to PyPI.

**Option 2: Manual Workflow Dispatch**

```bash
# Publish to Test PyPI
gh workflow run python-publish.yml -f publish_to=testpypi

# Publish to Production PyPI
gh workflow run python-publish.yml -f publish_to=pypi
```

## Version Management

The version is defined in `pyproject.toml`:

```toml
[project]
version = "1.0.0"
```

Update this before each release:

1. Edit `python/pyproject.toml`
2. Update the `version` field
3. Commit and create a tag: `git tag v1.0.0`
4. Push the tag: `git push origin v1.0.0`
5. Create a GitHub release from the tag

## Package Contents

The published package includes:

- `jira_client.py` - Main Python module
- `__init__.py` - Package initialization
- `README.md` - Documentation
- Pre-compiled C libraries for:
  - Linux (x86_64, ARM64)
  - macOS (Intel, Apple Silicon)
  - Windows (x86_64)

## Post-Publishing Verification

After publishing, verify the package:

```bash
# Install from PyPI
pip install jira-ticket-creator

# Test import
python -c "from jira_client import JiraClient; print(JiraClient.__doc__)"

# View on PyPI
open https://pypi.org/project/jira-ticket-creator/
```

## Troubleshooting

### "Invalid distribution"

Run `twine check dist/*` to see specific issues.

### "Authentication failed"

Ensure your API token is correctly configured in `~/.pypirc` or environment variables.

### "Module not found: libjira"

Ensure the C library was built before packaging:

```bash
cd python
CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api/capi.go
```

### Different library names on different platforms

The build system automatically detects the platform and uses:
- `libjira.so` on Linux
- `libjira.dylib` on macOS
- `libjira.dll` on Windows

## Changelog Format

When creating releases, include detailed changelogs:

```markdown
## [1.0.0] - 2024-01-XX

### Added
- New get_ticket() method to retrieve ticket details
- New search() method with JQL and keyword support
- New update_ticket() method for updating ticket fields
- Full Python type hints
- Comprehensive test suite

### Changed
- Improved error messages
- Enhanced documentation

### Fixed
- Library discovery on Windows

### Security
- Improved token handling
```

## Resources

- [PyPI Help](https://pypi.org/help/)
- [Python Packaging Guide](https://packaging.python.org/)
- [Twine Documentation](https://twine.readthedocs.io/)
- [GitHub Actions: Publish to PyPI](https://github.com/pypa/gh-action-pypi-publish)
