#!/bin/bash
# Clean script for Bazel - removes symlinks and cache

set -euo pipefail

# Get the repository root
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "Cleaning Bazel setup..."

# Shutdown bazel server if running
if command -v bazel >/dev/null 2>&1 || command -v bazelisk >/dev/null 2>&1; then
    echo "Shutting down Bazel server..."
    (bazel shutdown 2>/dev/null || bazelisk shutdown 2>/dev/null) || true
fi

cd "$REPO_ROOT"

# Remove symlinks
rm -f WORKSPACE BUILD.bazel .bazelrc .bazelignore .bazelversion deps.bzl gno_deps.bzl gno_test.bzl filetest.bzl txtar_test.bzl

# Remove Bazel output directories
rm -rf bazel-*

# Optional: Clean Bazel cache
if [ "${1:-}" = "--cache" ]; then
    echo "Cleaning Bazel cache..."
    rm -rf ~/.cache/bazel-gno
fi

echo "Bazel cleanup complete!"