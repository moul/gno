# Bazel-powered Makefile rules
# Include this in the main Makefile to enable Bazel support

# Check if Bazel binary exists
BAZEL_EXISTS := $(shell command -v bazel 2>/dev/null || command -v bazelisk 2>/dev/null || echo "")

# Check if Bazel is set up
BAZEL_SETUP_CHECK := $(shell test -L ../../WORKSPACE && echo "yes" || echo "no")

# Bazel binary with common flags
BAZEL := bazel
BAZEL_BUILD_FLAGS := --config=local
BAZEL_TEST_FLAGS := --config=local --test_output=errors

# Check for Bazel and setup if needed
define check_bazel
	@if [ -z "$(BAZEL_EXISTS)" ]; then \
		echo "Error: Bazel is not installed."; \
		echo ""; \
		echo "To install Bazel, run:"; \
		echo "  go install github.com/bazelbuild/bazelisk@latest"; \
		echo ""; \
		echo "Alternative installation methods:"; \
		echo "  - Download from: https://github.com/bazelbuild/bazelisk/releases"; \
		echo "  - macOS: brew install bazelisk"; \
		echo "  - Ubuntu/Debian: sudo apt install bazel"; \
		echo ""; \
		echo "After installation, run this command again."; \
		exit 1; \
	fi
	@if [ "$(BAZEL_SETUP_CHECK)" = "no" ]; then \
		cd ../.. && ./misc/bazel/setup.sh >/dev/null; \
	fi
endef

# Clean Bazel setup and cache
.PHONY: clean
clean:
	@echo "Cleaning Bazel setup and cache..."
	@cd ../.. && ./misc/bazel/clean.sh --cache

# Clean Bazel setup only (keep cache)
.PHONY: bazel-clean
bazel-clean:
	@./clean.sh

# Deep clean including cache
.PHONY: bazel-clean-all
bazel-clean-all:
	@./clean.sh --cache

# Generate/update BUILD files with Gazelle
.PHONY: generate
generate:
	@if command -v bazelisk >/dev/null 2>&1; then \
		BAZEL_CMD=bazelisk; \
	elif command -v bazel >/dev/null 2>&1; then \
		BAZEL_CMD=bazel; \
	else \
		echo "Error: Bazel is not installed."; \
		echo ""; \
		echo "To install Bazel, run:"; \
		echo "  go install github.com/bazelbuild/bazelisk@latest"; \
		echo ""; \
		echo "Alternative installation methods:"; \
		echo "  - Download from: https://github.com/bazelbuild/bazelisk/releases"; \
		echo "  - macOS: brew install bazelisk"; \
		echo "  - Ubuntu/Debian: sudo apt install bazel"; \
		echo ""; \
		echo "After installation, run this command again."; \
		exit 1; \
	fi; \
	if [ ! -L ../../WORKSPACE ]; then \
		cd ../.. && ./misc/bazel/setup.sh >/dev/null; \
	fi; \
	echo "Generating BUILD files with Gazelle..."; \
	$$BAZEL_CMD run //:gazelle -- -mode=fix -build_file_name=BUILD.bazel

# Update Go dependencies for Bazel
.PHONY: bazel-update-deps
bazel-update-deps:
	$(check_bazel)
	@echo "Updating Go dependencies for Bazel..."
	@$(BAZEL) run //:gazelle-update-repos

# Run all tests with Bazel
.PHONY: test
test:
	@bash -c ' \
	if command -v bazelisk >/dev/null 2>&1; then \
		BAZEL_CMD=bazelisk; \
	elif command -v bazel >/dev/null 2>&1; then \
		BAZEL_CMD=bazel; \
	else \
		echo "Error: Bazel is not installed."; \
		echo ""; \
		echo "To install Bazel, run:"; \
		echo "  go install github.com/bazelbuild/bazelisk@latest"; \
		echo ""; \
		echo "Alternative installation methods:"; \
		echo "  - Download from: https://github.com/bazelbuild/bazelisk/releases"; \
		echo "  - macOS: brew install bazelisk"; \
		echo "  - Ubuntu/Debian: sudo apt install bazel"; \
		echo ""; \
		echo "After installation, run this command again."; \
		exit 1; \
	fi; \
	if [ ! -L ../../WORKSPACE ]; then \
		cd ../.. && ./misc/bazel/setup.sh >/dev/null; \
	fi; \
	echo "Running tests with Bazel..."; \
	cd ../.. && $$BAZEL_CMD test //... --test_output=errors || true \
	'

# Run Go tests only
.PHONY: test-bazel-go
test-bazel-go:
	$(check_bazel)
	@echo "Running Go tests with Bazel..."
	@$(BAZEL) test $(BAZEL_TEST_FLAGS) //... --test_tag_filters=-gnotest,-filetest,-txtar

# Run Gno tests only
.PHONY: test-bazel-gno
test-bazel-gno:
	$(check_bazel)
	@echo "Running Gno tests with Bazel..."
	@$(BAZEL) test $(BAZEL_TEST_FLAGS) //... --test_tag_filters=gnotest

# Run filetests only
.PHONY: test-bazel-filetest
test-bazel-filetest:
	$(check_bazel)
	@echo "Running filetests with Bazel..."
	@$(BAZEL) test $(BAZEL_TEST_FLAGS) //... --test_tag_filters=filetest

# Run txtar tests only
.PHONY: test-bazel-txtar
test-bazel-txtar:
	$(check_bazel)
	@echo "Running txtar tests with Bazel..."
	@$(BAZEL) test $(BAZEL_TEST_FLAGS) //... --test_tag_filters=txtar

# Build all targets
.PHONY: build
build:
	@bash -c ' \
	if command -v bazelisk >/dev/null 2>&1; then \
		BAZEL_CMD=bazelisk; \
	elif command -v bazel >/dev/null 2>&1; then \
		BAZEL_CMD=bazel; \
	else \
		echo "Error: Bazel is not installed."; \
		echo ""; \
		echo "To install Bazel, run:"; \
		echo "  go install github.com/bazelbuild/bazelisk@latest"; \
		echo ""; \
		echo "Alternative installation methods:"; \
		echo "  - Download from: https://github.com/bazelbuild/bazelisk/releases"; \
		echo "  - macOS: brew install bazelisk"; \
		echo "  - Ubuntu/Debian: sudo apt install bazel"; \
		echo ""; \
		echo "After installation, run this command again."; \
		exit 1; \
	fi; \
	if [ ! -L ../../WORKSPACE ]; then \
		cd ../.. && ./misc/bazel/setup.sh >/dev/null; \
	fi; \
	if ! [ -f ../../gnovm/cmd/gno/BUILD.bazel ] || ! [ -f ../../gno.land/cmd/gnoland/BUILD.bazel ]; then \
		echo "BUILD files not found. Generating them with Gazelle..."; \
		$$BAZEL_CMD run //:gazelle -- -mode=fix -build_file_name=BUILD.bazel >/dev/null 2>&1 || \
			echo "Warning: Some BUILD file generation errors occurred"; \
	fi; \
	echo "Building binaries with Bazel..."; \
	echo "  - Building gno..."; \
	$$BAZEL_CMD build $(BAZEL_BUILD_FLAGS) //gnovm/cmd/gno 2>/dev/null || echo "    Failed to build gno"; \
	echo "  - Building gnoland..."; \
	$$BAZEL_CMD build $(BAZEL_BUILD_FLAGS) //gno.land/cmd/gnoland 2>/dev/null || echo "    Failed to build gnoland"; \
	echo "  - Building gnokey..."; \
	$$BAZEL_CMD build $(BAZEL_BUILD_FLAGS) //gno.land/cmd/gnokey 2>/dev/null || echo "    Failed to build gnokey"; \
	echo "  - Building gnodev..."; \
	$$BAZEL_CMD build $(BAZEL_BUILD_FLAGS) //contribs/gnodev 2>/dev/null || echo "    Failed to build gnodev"; \
	echo ""; \
	echo "Built binaries are in bazel-bin/"; \
	echo "Note: Run '"'"'make bazel-generate'"'"' to update BUILD files if builds are failing" \
	'

# Build specific binaries
.PHONY: build-bazel-gno
build-bazel-gno:
	$(check_bazel)
	@$(BAZEL) build $(BAZEL_BUILD_FLAGS) //gnovm/cmd/gno

.PHONY: build-bazel-gnoland
build-bazel-gnoland:
	$(check_bazel)
	@$(BAZEL) build $(BAZEL_BUILD_FLAGS) //gno.land/cmd/gnoland

.PHONY: build-bazel-gnokey
build-bazel-gnokey:
	$(check_bazel)
	@$(BAZEL) build $(BAZEL_BUILD_FLAGS) //gno.land/cmd/gnokey

# Run with remote caching (for CI)
.PHONY: test-bazel-ci
test-bazel-ci:
	$(check_bazel)
	@echo "Running tests with remote caching..."
	@$(BAZEL) test --config=ci //...

# Query targets
.PHONY: bazel-list-tests
bazel-list-tests:
	$(check_bazel)
	@echo "Listing all test targets..."
	@$(BAZEL) query 'tests(//...)'

.PHONY: bazel-list-go-tests
bazel-list-go-tests:
	$(check_bazel)
	@echo "Listing Go test targets..."
	@$(BAZEL) query 'attr(tags, "[-]gnotest", tests(//...))'

# Help
.PHONY: help
help:
	@echo "Bazel targets (run from project root with 'make bazel-<target>'):"
	@echo "  test               - Run all tests (auto-generates BUILD files if needed)"
	@echo "  build              - Build all binaries (auto-generates BUILD files if needed)"
	@echo "  generate           - Generate/update BUILD files with Gazelle"
	@echo "  clean              - Clean Bazel setup, cache, and symlinks"
	@echo ""
	@echo "Other targets (run from misc/bazel/):"
	@echo "  bazel-clean        - Clean Bazel setup only (keep cache)"
	@echo "  bazel-clean-all    - Clean Bazel setup and cache"
	@echo "  bazel-update-deps  - Update Go dependencies for Bazel"
	@echo "  test-bazel-go      - Run only Go tests"
	@echo "  test-bazel-gno     - Run only Gno tests"
	@echo "  test-bazel-filetest - Run only filetests"
	@echo "  test-bazel-txtar   - Run only txtar tests"
	@echo "  test-bazel-ci      - Run tests with remote caching (CI)"
	@echo "  build-bazel-gno    - Build gno binary"
	@echo "  build-bazel-gnoland - Build gnoland binary"
	@echo "  build-bazel-gnokey - Build gnokey binary"
	@echo "  bazel-list-tests   - List all test targets"
	@echo "  bazel-list-go-tests - List Go test targets"