#!/bin/bash
# Generate BUILD files for all packages in the Gno project

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "Generating BUILD files for all packages..."

# Function to create a basic BUILD.bazel file
create_build_file() {
    local dir="$1"
    local build_file="$dir/BUILD.bazel"
    
    # Skip if BUILD file already exists
    if [[ -f "$build_file" ]]; then
        return
    fi
    
    # Get package name from directory
    local pkg_name=$(basename "$dir")
    local import_path="github.com/gnolang/gno/${dir#./}"
    
    # Check if directory contains Go files
    if ! ls "$dir"/*.go >/dev/null 2>&1; then
        return
    fi
    
    echo "Creating BUILD file for $dir"
    
    cat > "$build_file" << EOF
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "$pkg_name",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "$import_path",
    visibility = ["//visibility:public"],
)

go_test(
    name = "${pkg_name}_test",
    srcs = glob(["*_test.go"]),
    embed = [":$pkg_name"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
EOF
}

# Generate BUILD files for tm2 packages
find tm2/pkg -type d -name "*.go" -prune -o -type d -print | while read -r dir; do
    create_build_file "$dir"
done

# Generate BUILD files for gnovm packages
find gnovm/pkg -type d -name "*.go" -prune -o -type d -print | while read -r dir; do
    create_build_file "$dir"
done

# Generate BUILD files for gno.land packages
find gno.land/pkg -type d -name "*.go" -prune -o -type d -print | while read -r dir; do
    create_build_file "$dir"
done

echo "BUILD file generation complete!"
echo ""
echo "Note: Some packages may need manual dependency adjustments."
echo "Run 'bazelisk build //...' to identify missing dependencies."