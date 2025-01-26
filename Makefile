# Makefile for fastmap Package Testing

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean
GOVENDOR=$(GOCMD) mod vendor

# Project parameters
PACKAGE=github.com/billowdev/fastmap
TEST_PACKAGE=./...

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

# Default target
.PHONY: all
all: test coverage

# Initialize dependencies
.PHONY: init
init:
	@echo "$(YELLOW)Initializing dependencies...$(NC)"
	$(GOMOD) tidy
	$(GOMOD) download

# Run all tests
.PHONY: test
test:
	@echo "$(YELLOW)Running all tests...$(NC)"
	$(GOTEST) $(TEST_PACKAGE)

# Run tests with verbose output
.PHONY: test-verbose
test-verbose:
	@echo "$(YELLOW)Running tests with verbose output...$(NC)"
	$(GOTEST) -v $(TEST_PACKAGE)

# Generate coverage report
.PHONY: coverage
coverage:
	@echo "$(YELLOW)Generating test coverage report...$(NC)"
	$(GOTEST) -cover $(TEST_PACKAGE)
	$(GOTEST) -coverprofile=coverage.out $(TEST_PACKAGE)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated at coverage.html$(NC)"

# Run specific test
.PHONY: test-specific
test-specific:
	@read -p "Enter test name to run (e.g., TestHashMap_Put): " TEST_NAME; \
	$(GOTEST) -v -run $$TEST_NAME $(TEST_PACKAGE)

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. $(TEST_PACKAGE)

# Memory profiling
.PHONY: profile-memory
profile-memory:
	@echo "$(YELLOW)Running memory profiling...$(NC)"
	$(GOTEST) -bench=. -memprofile=mem.prof $(TEST_PACKAGE)
	go tool pprof mem.prof

# CPU profiling
.PHONY: profile-cpu
profile-cpu:
	@echo "$(YELLOW)Running CPU profiling...$(NC)"
	$(GOTEST) -bench=. -cpuprofile=cpu.prof $(TEST_PACKAGE)
	go tool pprof cpu.prof

# Race condition detection
.PHONY: race
race:
	@echo "$(YELLOW)Running race condition tests...$(NC)"
	$(GOTEST) -race $(TEST_PACKAGE)

# Clean test artifacts
.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning test artifacts...$(NC)"
	$(GOCLEAN)
	rm -f coverage.out coverage.html
	rm -f *.prof

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests and generate coverage report"
	@echo "  init          - Initialize dependencies"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  coverage      - Generate test coverage report"
	@echo "  test-specific - Run a specific test by name"
	@echo "  benchmark     - Run performance benchmarks"
	@echo "  profile-memory- Generate memory profiling"
	@echo "  profile-cpu   - Generate CPU profiling"
	@echo "  race          - Run race condition tests"
	@echo "  clean         - Clean test artifacts"
	@echo "  help          - Show this help message"