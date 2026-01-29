# GitHub Actions Workflows

Comprehensive CI/CD pipeline for the JIRA Ticket Creator project.

## Workflows Overview

### 1. **CI (.github/workflows/ci.yml)**
Main continuous integration workflow triggered on push and pull requests.

**Jobs:**
- **go-tests**: Run Go tests with race detector and coverage
  - Tests: `go test -v -race -coverprofile=go-coverage.out ./...`
  - Coverage uploaded to Codecov

- **python-tests**: Run Python tests on multiple Python versions (3.9, 3.11, 3.12)
  - Builds C library (CGO)
  - Runs pytest with coverage
  - Uploads coverage to Codecov

- **lint**: Linting and formatting checks
  - golangci-lint for Go
  - go fmt verification
  - black, flake8, isort for Python

- **build**: Build CLI binaries and C library artifacts
  - Builds for Linux and macOS
  - Uploads artifacts (retention: 5 days)

- **status**: Final status check

**Triggers:**
- `push` to `master`, `main`, `develop`
- `pull_request` to `master`, `main`

**Features:**
- Concurrency control (cancels older runs)
- Module caching (Go modules, pip packages)
- Cross-platform builds
- Coverage reporting

---

### 2. **Python Tests (.github/workflows/python-tests.yml)**
Dedicated Python testing workflow with comprehensive coverage.

**Jobs:**
- **test**: Multi-platform, multi-version testing
  - Platforms: Linux, macOS, Windows
  - Python versions: 3.8 - 3.12
  - Builds C library
  - Runs pytest with verbose output
  - Coverage reporting

- **lint**: Python linting
  - black formatting
  - isort import sorting
  - flake8 analysis

- **build-dist**: Build distribution package
  - Builds wheel and sdist
  - Validates with twine
  - Uploads artifacts

- **summary**: Test summary

**Triggers:**
- Changes to `python/**` or `internal/api/**`
- Direct trigger on workflow file changes

**Features:**
- Matrix testing (3 OS × 5 Python versions)
- Smart exclusions to save CI time
- Distribution validation
- Artifact preservation

---

### 3. **Go Tests (.github/workflows/go-tests.yml)**
Dedicated Go testing workflow.

**Jobs:**
- **test**: Multi-version Go testing
  - Go versions: 1.21 - 1.25
  - Platforms: Linux, macOS, Windows
  - Race detector enabled
  - Coverage collection

- **build**: Binary builds
  - CLI binary compilation
  - Platform-specific naming
  - Binary verification

- **lint**: Go linting
  - golangci-lint
  - go vet
  - go fmt check

- **summary**: Test summary

**Triggers:**
- Changes to `**.go`, `go.mod`, `go.sum`
- Direct trigger on workflow file changes

**Features:**
- Module caching
- Race condition detection
- Comprehensive linting
- Cross-version compatibility

---

### 4. **Python Publish (.github/workflows/python-publish.yml)**
PyPI publishing pipeline.

**Jobs:**
- **build**: Build distribution packages
  - Multi-version Python compilation
  - C library building
  - Distribution validation

- **test-install**: Verify installation
  - Tests install from Test PyPI
  - Cross-platform verification

**Triggers:**
- GitHub releases (publishes to PyPI)
- Manual workflow dispatch

**Features:**
- Automatic release publishing
- Test PyPI staging support
- Multi-platform testing
- Installation verification

---

## Workflow Triggers

| Workflow | Trigger | Branches |
|----------|---------|----------|
| **CI** | push, PR | master, main, develop |
| **Python Tests** | push, PR | Any (path filtered) |
| **Go Tests** | push, PR | Any (path filtered) |
| **Python Publish** | Release, Manual | N/A |

## Coverage Integration

- **Codecov**: Automatically uploads coverage from CI
  - Go coverage: `go-coverage.out`
  - Python coverage: `coverage.xml`
  - Flags: `go`, `python` for separate tracking

## Caching Strategy

### Go Modules
```yaml
cache:
  path: ~/go/pkg/mod
  key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### Python Packages
```yaml
cache:
  path: ~/.cache/pip
  key: ${{ runner.os }}-pip-${{ matrix.python-version }}-${{ hashFiles('**/python/pyproject.toml') }}
```

## Artifact Management

| Artifact | Retention | Format |
|----------|-----------|--------|
| Build Binaries | 5 days | Linux/macOS CLI + libs |
| Python Distribution | 5 days | wheel + sdist |

## Environment Details

### Action Versions
- actions/checkout: v4
- actions/setup-go: v5
- actions/setup-python: v4
- actions/cache: v3
- actions/upload-artifact: v3
- codecov/codecov-action: v3
- golangci/golangci-lint-action: v3

### Tool Versions
- Go: 1.25 (latest)
- Python: 3.8 - 3.12
- golangci-lint: latest
- pytest: latest
- black, flake8, isort: latest

## Local Development

To replicate CI environment locally:

### Go Testing
```bash
go test -v -race ./...
go test -v -coverprofile=coverage.out ./...
```

### Python Testing
```bash
cd python
CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api
pip install pytest pytest-cov
python -m pytest tests/ -v --cov=jira_client
```

### Linting
```bash
# Go
golangci-lint run
go fmt ./...

# Python
cd python
black --check .
flake8 .
isort --check-only .
```

## Troubleshooting

### Tests Failing Locally but Passing in CI
- Check Python version (CI tests 3.9, 3.11, 3.12)
- Check Go version (CI uses Go 1.25)
- Ensure C library is built: `cd python && CGO_ENABLED=1 go build -buildmode=c-shared -o libjira.so ../internal/api`

### Artifacts Not Found
- Check workflow logs for build failures
- Verify path matches artifact upload configuration
- Check artifact retention (5 days default)

### Coverage Not Uploading
- Ensure Codecov token is configured (usually auto-detected for public repos)
- Check coverage file paths match workflow configuration
- Verify `fail_ci_if_error: false` allows pipeline to continue

## Best Practices

1. **Commit Before Push**: Ensure local tests pass before pushing
2. **Monitor Workflows**: Check Actions tab regularly for failures
3. **Update Dependencies**: Keep action versions current
4. **Review Logs**: Check detailed logs for warnings or issues
5. **Test Matrix**: Ensure new code works on all tested versions

## Adding New Workflows

To add a new workflow:

1. Create `.github/workflows/new-workflow.yml`
2. Add appropriate triggers (`on:` section)
3. Define jobs and steps
4. Test locally first
5. Monitor first run in Actions tab

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Toolchain](https://go.dev/)
- [Python Packaging](https://packaging.python.org/)
- [Codecov Integration](https://codecov.io/)
- [golangci-lint](https://golangci-lint.run/)

---

**Last Updated**: 2024-01-29
**Status**: Active and Monitoring
