# Makefile for atomix with intrinsics-customized Go compiler
#
# The intrinsics compiler provides single-instruction atomic operations
# by replacing function calls with inline CPU instructions.
#
# Prerequisites:
#   - curl, tar and shasum
#   - Go 1.27+ (for bootstrapping the tagged compiler)
#   - Linux (or WSL2 on Windows)
#
# Usage:
#   make install-compiler  # Install the latest Atomix compiler release
#   make build             # Build atomix with intrinsics compiler
#   make test              # Test atomix with intrinsics compiler
#   make verify            # Verify intrinsics are applied (check assembly)

# Configuration
export GOTOOLCHAIN := local
COMPILER_REPO    := https://github.com/hayabusa-cloud/go
COMPILER_TAG     ?=
COMPILER_DIR     ?= $(HOME)/sdk/$(or $(COMPILER_TAG),go-atomix)
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

# The default SDK path points to a versioned release directory.
.PHONY: install-compiler
install-compiler:
	@set -eu; dir="$(COMPILER_DIR)"; tag="$(COMPILER_TAG)"; \
	parent=$$(dirname "$$dir"); dir="$$parent/$$(basename "$$dir")"; dest="$$dir"; \
	if test -z "$$tag"; then \
		release=$$(curl -fsSIL -o /dev/null -w '%{url_effective}' "$(COMPILER_REPO)/releases/latest"); \
		tag=$${release##*/}; dest="$$parent/$$tag"; \
	fi; \
	matches_tag() { test -x "$$1/bin/go" && test -f "$$1/VERSION" && test "$$(head -n 1 "$$1/VERSION")" = "$$tag"; }; \
	if matches_tag "$$dir"; then exit 0; fi; \
	if test -e "$$dir" && test ! -L "$$dir"; then echo "Compiler destination already exists: $$dir" >&2; exit 1; fi; \
	mkdir -p "$$parent"; \
	tmp="$$dest.installing"; mkdir -m 700 "$$tmp"; \
	trap 'rm -rf "$$tmp"' 0; \
	if test -e "$$dest" || test -L "$$dest"; then \
		matches_tag "$$dest" || { echo "Invalid compiler at $$dest" >&2; exit 1; }; \
	else \
		url="$(COMPILER_REPO)/releases/download/$$tag"; archive="$$tag.src.tar.gz"; \
		curl -fsSL "$$url/$$archive" -o "$$tmp/$$archive"; \
		curl -fsSL "$$url/SHA256SUMS" -o "$$tmp/SHA256SUMS"; \
		(cd "$$tmp" && awk -v file="$$archive" '$$2 == file { print; found=1 } END { if (!found) exit 1 }' SHA256SUMS > selected.sha256 && shasum -a 256 -c selected.sha256); \
		tar -xzf "$$tmp/$$archive" -C "$$tmp"; \
		test "$$(head -n 1 "$$tmp/go/VERSION")" = "$$tag"; \
		bootstrap=$$(GOWORK=off go env GOROOT); (cd "$$tmp/go/src" && GOWORK=off GOROOT_BOOTSTRAP="$$bootstrap" ./make.bash); \
		test -x "$$tmp/go/bin/go"; mv "$$tmp/go" "$$dest"; \
	fi; \
	if test "$$dest" != "$$dir"; then ln -sfn "$$tag" "$$dir"; fi

# ============================================================================
# Build & Test with Intrinsics Compiler
# ============================================================================

.PHONY: build
build:
	$(require-compiler)
	@echo "Building atomix with intrinsics compiler..."
	GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" vet ./...
	GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" build ./...
	@echo "✓ Build successful"

.PHONY: test
test:
	$(require-compiler)
	@echo "Testing atomix with intrinsics compiler..."
	GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" test -timeout 120s -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo "✓ Tests passed"

.PHONY: bench
bench:
	$(require-compiler)
	@echo "Running benchmarks with intrinsics compiler..."
	GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" test -timeout 30m -bench=. -benchmem ./...

# ============================================================================
# Verification
# ============================================================================

.PHONY: verify
verify:
	$(require-compiler)
	@echo "Verifying intrinsics are applied..."
	@echo ""
	@echo "=== Checking for inline atomic instructions ==="
	@GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" build -gcflags='-S' ./... 2>&1 | \
		grep -E 'LDADDA|LDADDAL|STLR|STLRW|LDAR|LDARW|XCHG|CMPXCHG|MOVLstore|MOVQstore' | head -20 || true
	@echo ""
	@echo "=== Checking for function calls (should be empty if intrinsics work) ==="
	@GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" build -gcflags='-S' ./... 2>&1 | \
		grep -E 'CALL.*internal/arch\.' | head -10 || echo "✓ No function calls to internal/arch (intrinsics applied)"

# ============================================================================
# Utilities
# ============================================================================

.PHONY: clean
clean:
	rm -f coverage.out
	rm -f test_*
	GOROOT="$(GOROOT_INTRINSIC)" "$(GO_INTRINSIC)" clean -cache 2>/dev/null || true

.PHONY: help
help:
	@echo "atomix Makefile - Build with intrinsics-customized Go compiler"
	@echo ""
	@echo "Compiler:"
	@echo "  install-compiler  Install the latest Atomix compiler or the release selected by COMPILER_TAG"
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
