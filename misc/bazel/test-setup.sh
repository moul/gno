#!/bin/bash
# Test script to validate Bazel setup

set -euo pipefail

echo "=== Bazel Setup Test ==="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if bazel/bazelisk is installed
if command -v bazel >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Bazel found: $(bazel version | grep "Build label" | cut -d: -f2)"
elif command -v bazelisk >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Bazelisk found"
else
    echo -e "${RED}✗${NC} Bazel not found. Install with: go install github.com/bazelbuild/bazelisk@latest"
    exit 1
fi

# Check if setup has been run
if [ -L "WORKSPACE" ]; then
    echo -e "${GREEN}✓${NC} Bazel symlinks are set up"
else
    echo -e "${YELLOW}!${NC} Bazel not set up. Run: ./misc/bazel/setup.sh"
    exit 1
fi

# Try to run gazelle
echo
echo "Generating BUILD files with Gazelle..."
if bazel run //:gazelle >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Gazelle executed successfully"
else
    echo -e "${RED}✗${NC} Gazelle failed"
    exit 1
fi

# Try a simple query
echo
echo "Testing Bazel query..."
if bazel query '//gnovm/cmd/gno' >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Bazel can query targets"
else
    echo -e "${RED}✗${NC} Bazel query failed"
    exit 1
fi

# Try building gno binary
echo
echo "Testing build of gno binary..."
if bazel build //gnovm/cmd/gno >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Successfully built gno binary"
    echo "    Binary at: bazel-bin/gnovm/cmd/gno/gno_/gno"
else
    echo -e "${RED}✗${NC} Failed to build gno binary"
    exit 1
fi

# Check if remote cache is configured
echo
echo "Checking remote cache configuration..."
if [ -n "${BUILDBUDDY_API_KEY:-}" ]; then
    echo -e "${GREEN}✓${NC} BuildBuddy API key found"
else
    echo -e "${YELLOW}!${NC} No BuildBuddy API key set (optional for local development)"
fi

# Summary
echo
echo "=== Setup Summary ==="
echo -e "${GREEN}✓${NC} Bazel is properly configured!"
echo
echo "Next steps:"
echo "1. Run tests: make test-bazel"
echo "2. Run specific tests: bazel test //gnovm/pkg/gnolang:all"
echo "3. Build binaries: make build-bazel"
echo
echo "For faster builds with caching:"
echo "1. Sign up at https://buildbuddy.io"
echo "2. Export BUILDBUDDY_API_KEY=\"your-key\""
echo "3. Use --config=buildbuddy-auth flag"