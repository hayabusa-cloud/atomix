# Makefile for atomix with intrinsics-customized Go compiler
#
# The intrinsics compiler provides single-instruction atomic operations
# by replacing function calls with inline CPU instructions.
#
# Prerequisites:
#   - Git
#   - Go 1.25+ (for bootstrapping)
#   - Linux (or WSL2 on Windows)
#
# Usage:
#   make install-compiler  # Install or update the intrinsics compiler
#   make build             # Build atomix with intrinsics compiler
#   make test              # Test atomix with intrinsics compiler
#   make verify            # Verify intrinsics are applied (check assembly)

# Configuration
COMPILER_REPO    := https://github.com/hayabusa-cloud/go.git
COMPILER_BRANCH  := atomix
COMPILER_DIR     := $(HOME)/github.com/go
GOROOT_INTRINSIC := $(COMPILER_DIR)
GO_INTRINSIC     := $(GOROOT_INTRINSIC)/bin/go

# Default target
.DEFAULT_GOAL := build

# Compiler check macro
define require-compiler
	@if [ ! -x "$(GO_INTRINSIC)" ]; then \
		echo "Error: Intrinsics compiler not found at $(GO_INTRINSIC)"; \
		echo "Run 'make install-compiler' first."; \
		exit 1; \
	fi
endef

# ============================================================================
# Compiler Installation
# ============================================================================

.PHONY: install-compiler
install-compiler:
	@if [ -d "$(COMPILER_DIR)" ]; then \
		echo "Updating intrinsics compiler..."; \
		cd $(COMPILER_DIR) && git fetch origin && git checkout $(COMPILER_BRANCH) && git pull origin $(COMPILER_BRANCH); \
	else \
		echo "Cloning intrinsics compiler..."; \
		mkdir -p $(dir $(COMPILER_DIR)); \
		git clone --branch $(COMPILER_BRANCH) $(COMPILER_REPO) $(COMPILER_DIR); \
	fi
	@echo "Building compiler (this may take several minutes)..."
	cd $(COMPILER_DIR)/src && ./make.bash
	@echo "✓ Intrinsics compiler ready at $(GOROOT_INTRINSIC)"

# ============================================================================
# Build & Test with Intrinsics Compiler
# ============================================================================

.PHONY: build
build:
	$(require-compiler)
	@echo "Building atomix with intrinsics compiler..."
	GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) vet ./...
	GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) build ./...
	@echo "✓ Build successful"

.PHONY: test
test:
	$(require-compiler)
	@echo "Testing atomix with intrinsics compiler..."
	GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo "✓ Tests passed"

.PHONY: bench
bench:
	$(require-compiler)
	@echo "Running benchmarks with intrinsics compiler..."
	GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) test -bench=. -benchmem ./...

# ============================================================================
# Verification
# ============================================================================

.PHONY: verify
verify:
	$(require-compiler)
	@echo "Verifying intrinsics are applied..."
	@echo ""
	@echo "=== Checking for inline atomic instructions ==="
	@GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) build -gcflags='-S' ./... 2>&1 | \
		grep -E 'LDADDA|LDADDAL|STLR|STLRW|LDAR|LDARW|XCHG|CMPXCHG|MOVLstore|MOVQstore' | head -20 || true
	@echo ""
	@echo "=== Checking for function calls (should be empty if intrinsics work) ==="
	@GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) build -gcflags='-S' ./... 2>&1 | \
		grep -E 'CALL.*internal/arch\.' | head -10 || echo "✓ No function calls to internal/arch (intrinsics applied)"

# ============================================================================
# Utilities
# ============================================================================

.PHONY: clean
clean:
	rm -f coverage.out
	rm -f test_*
	GOROOT=$(GOROOT_INTRINSIC) $(GO_INTRINSIC) clean -cache 2>/dev/null || true

.PHONY: help
help:
	@echo "atomix Makefile - Build with intrinsics-customized Go compiler"
	@echo ""
	@echo "Compiler:"
	@echo "  install-compiler  Install or update the intrinsics compiler"
	@echo ""
	@echo "Build:"
	@echo "  build             Build with intrinsics compiler (includes vet)"
	@echo "  test              Run tests with race detection and coverage"
	@echo "  bench             Run benchmarks"
	@echo ""
	@echo "Verification:"
	@echo "  verify            Verify intrinsics are applied"
	@echo ""
	@echo "Other:"
	@echo "  clean             Remove build artifacts"
	@echo "  help              Show this help"
	@echo ""
	@echo "Configuration:"
	@echo "  COMPILER_DIR     = $(COMPILER_DIR)"
	@echo "  GO_INTRINSIC     = $(GO_INTRINSIC)"
