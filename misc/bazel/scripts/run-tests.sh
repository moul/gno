#!/bin/bash
# Run tests with proper working directory setup

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "Running Bazel tests with compatibility fixes..."

# Run tests with the proper working directory
bazelisk test \
    --test_env="GOPATH=$GOPATH" \
    --test_env="HOME=$HOME" \
    --test_env="PWD=$REPO_ROOT" \
    --test_output=errors \
    --test_arg="-test.v" \
    "$@"