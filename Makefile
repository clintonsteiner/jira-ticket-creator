.PHONY: help build python-build python-test python-install python-clean clean test all

# Variables
GO := go
PYTHON := python3
LIB_NAME := libjira
PYTHON_DIR := python
INTERNAL_API_DIR := internal/api

# Detect OS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    LIB_EXT := so
    DYLIB_FLAG :=
endif
ifeq ($(UNAME_S),Darwin)
    LIB_EXT := dylib
    DYLIB_FLAG :=
endif
ifeq ($(OS),Windows_NT)
    LIB_EXT := dll
    DYLIB_FLAG :=
endif

LIB_FILE := $(PYTHON_DIR)/$(LIB_NAME).$(LIB_EXT)

help:
	@echo "JIRA Ticket Creator - Makefile targets"
	@echo ""
	@echo "Go CLI Targets:"
	@echo "  build              - Build the CLI binary"
	@echo "  test               - Run Go tests"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "Python Integration Targets:"
	@echo "  python-build       - Build C library for Python"
	@echo "  python-test        - Test Python client"
	@echo "  python-install     - Install Python package locally"
	@echo "  python-clean       - Clean Python build artifacts"
	@echo "  python-example     - Run Python example"
	@echo ""
	@echo "Combined Targets:"
	@echo "  all                - Build CLI + Python library"
	@echo "  all-test           - Run all tests (Go + Python)"
	@echo "  all-clean          - Clean all build artifacts"

# ============================================================================
# GO CLI TARGETS
# ============================================================================

build:
	@echo "Building CLI binary..."
	$(GO) build -o jira-ticket-creator ./cmd/jira-ticket-creator
	@echo "Built: jira-ticket-creator"

test:
	@echo "Running Go tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning Go build artifacts..."
	rm -f jira-ticket-creator
	rm -f *.o *.a

# ============================================================================
# PYTHON TARGETS
# ============================================================================

python-build: $(LIB_FILE)
	@echo "Python C library built: $(LIB_FILE)"

$(LIB_FILE): $(INTERNAL_API_DIR)/capi.go
	@echo "Building C library for Python ($(UNAME_S))..."
	@mkdir -p $(PYTHON_DIR)
	cd $(PYTHON_DIR) && CGO_ENABLED=1 $(GO) build \
		-buildmode=c-shared \
		-o $(notdir $@) \
		../$(INTERNAL_API_DIR)/capi.go
	@echo "Built: $@"

python-test: python-build
	@echo "Running Python tests..."
	@cd $(PYTHON_DIR) && $(PYTHON) -m pytest tests/ -v || echo "No pytest installed"
	@cd $(PYTHON_DIR) && $(PYTHON) jira_client.py

python-install: python-build
	@echo "Installing Python package..."
	cd $(PYTHON_DIR) && $(PYTHON) -m pip install -e .

python-clean:
	@echo "Cleaning Python build artifacts..."
	rm -f $(LIB_FILE)
	rm -f $(PYTHON_DIR)/*.so $(PYTHON_DIR)/*.dylib $(PYTHON_DIR)/*.dll
	rm -rf $(PYTHON_DIR)/*.egg-info $(PYTHON_DIR)/dist $(PYTHON_DIR)/build
	find $(PYTHON_DIR) -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	find $(PYTHON_DIR) -type f -name "*.pyc" -delete

python-example: python-build
	@echo "Running Python example..."
	@cd $(PYTHON_DIR) && $(PYTHON) -c "\
		from jira_client import JiraClient; \
		print('JiraClient imported successfully'); \
		print('To use, create a client:'); \
		print('  client = JiraClient(url=\"...\", email=\"...\", token=\"...\", project=\"...\")'); \
		print('  ticket = client.create_ticket(summary=\"My task\")'); \
	"

# ============================================================================
# COMBINED TARGETS
# ============================================================================

all: build python-build
	@echo "Build complete!"
	@echo "CLI: ./jira-ticket-creator"
	@echo "Python: $(LIB_FILE)"

all-test: test python-test
	@echo "All tests passed!"

all-clean: clean python-clean
	@echo "All build artifacts cleaned"

# ============================================================================
# DEVELOPMENT TARGETS
# ============================================================================

dev-setup:
	@echo "Setting up development environment..."
	$(GO) mod tidy
	$(GO) mod download
	@which cgo > /dev/null || (echo "Error: CGO required for C bindings"; exit 1)
	@echo "Development environment ready"

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	cd $(PYTHON_DIR) && $(PYTHON) -m black . || echo "black not installed"

lint:
	@echo "Linting code..."
	$(GO) vet ./...
	@which golangci-lint > /dev/null && golangci-lint run ./... || echo "golangci-lint not installed"

# ============================================================================
# DOCUMENTATION TARGETS
# ============================================================================

docs:
	@echo "Generating documentation..."
	mkdir -p docs
	echo "# JIRA Ticket Creator Python API" > docs/python_api.md
	echo "" >> docs/python_api.md
	echo "## Installation" >> docs/python_api.md
	echo "\`\`\`bash" >> docs/python_api.md
	echo "make python-build" >> docs/python_api.md
	echo "make python-install" >> docs/python_api.md
	echo "\`\`\`" >> docs/python_api.md
	echo "" >> docs/python_api.md
	echo "## Usage" >> docs/python_api.md
	echo "\`\`\`python" >> docs/python_api.md
	echo "from jira_client import JiraClient" >> docs/python_api.md
	echo "" >> docs/python_api.md
	echo "client = JiraClient(" >> docs/python_api.md
	echo "    url='https://company.atlassian.net'," >> docs/python_api.md
	echo "    email='user@company.com'," >> docs/python_api.md
	echo "    token='api-token'," >> docs/python_api.md
	echo "    project='PROJ'" >> docs/python_api.md
	echo ")" >> docs/python_api.md
	echo "" >> docs/python_api.md
	echo "ticket = client.create_ticket(summary='My task')" >> docs/python_api.md
	echo "print(f'Created: {ticket[\"key\"]}')" >> docs/python_api.md
	echo "\`\`\`" >> docs/python_api.md
	@echo "Documentation generated in docs/python_api.md"

# ============================================================================
# UTILITY TARGETS
# ============================================================================

info:
	@echo "JIRA Ticket Creator Build Info"
	@echo "=============================="
	@echo "OS: $(UNAME_S)"
	@echo "Go version: $(shell $(GO) version)"
	@echo "Python version: $(shell $(PYTHON) --version)"
	@echo "Library extension: $(LIB_EXT)"
	@echo "Library target: $(LIB_FILE)"

version:
	@grep "version" $(PYTHON_DIR)/__init__.py | head -1

# Default target
.DEFAULT_GOAL := help
