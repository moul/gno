#!/bin/bash
# Simple BUILD file generator that bypasses Gazelle issues

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "Generating simple BUILD.bazel files..."

# Function to create a basic go_test BUILD file
create_test_build() {
    local dir="$1"
    local pkg_name=$(basename "$dir")
    
    # Skip if no test files
    if ! ls "$dir"/*_test.go >/dev/null 2>&1; then
        return
    fi
    
    # Skip if BUILD.bazel already exists
    if [ -f "$dir/BUILD.bazel" ]; then
        return
    fi
    
    cat > "$dir/BUILD.bazel" << EOF
load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "${pkg_name}_test",
    srcs = glob(["*_test.go"]),
    size = "small",
    deps = [
        # Add dependencies as needed
    ],
)
EOF
    echo "Created $dir/BUILD.bazel"
}

# Create BUILD files for key test directories
create_test_build "tm2/pkg/crypto"
create_test_build "tm2/pkg/std" 
create_test_build "gnovm/pkg/gnolang"
create_test_build "gnovm/tests"

# Create a simple root test that aggregates others
cat > BUILD.bazel << 'EOF'
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle configuration
# gazelle:prefix github.com/gnolang/gno
# gazelle:go_generate_proto false
# gazelle:exclude misc/bazel
# gazelle:exclude vendor
# gazelle:exclude bazel-*
# gazelle:build_file_name BUILD.bazel

gazelle(
    name = "gazelle",
    command = "update",
    extra_args = [
        "-build_file_name=BUILD.bazel",
    ],
)

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

# Simple test suite that includes working tests
test_suite(
    name = "simple_tests",
    tests = [
        "//test_working:simple_test",
        # Add more as they work
    ],
)

# Root package exports
exports_files([
    "WORKSPACE",
    ".bazelrc", 
    ".bazelignore",
    "deps.bzl",
    "gno_deps.bzl",
    "gno_test.bzl",
    "filetest.bzl",
    "txtar_test.bzl",
])
EOF

echo "Done! You can now run: bazel test //:simple_tests"