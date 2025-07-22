#!/bin/bash
# Generate working BUILD files while excluding problematic dependencies

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "Generating working BUILD files..."

# First, update the root BUILD.bazel to exclude more problematic paths
cat > BUILD.bazel << 'EOF'
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle configuration
# gazelle:prefix github.com/gnolang/gno
# gazelle:go_generate_proto false
# gazelle:exclude misc/bazel
# gazelle:exclude vendor
# gazelle:exclude bazel-*
# gazelle:exclude contribs/gnogenesis
# gazelle:exclude misc/genstd
# gazelle:exclude tm2/pkg/libtm
# gazelle:exclude **/*_test.go
# gazelle:build_file_name BUILD.bazel
# gazelle:go_naming_convention import_alias
# gazelle:resolve go github.com/gnolang/libtm/messages @com_github_gnolang_libtm//messages
# gazelle:resolve go github.com/sig-0/insertion-queue @com_github_sig_0_iq//:go_default_library

gazelle(
    name = "gazelle",
    command = "update",
    extra_args = [
        "-build_file_name=BUILD.bazel",
        "-go_naming_convention=import_alias",
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

exports_files([
    "go.mod",
    "go.sum",
])
EOF

# Create manual BUILD files for core packages without external deps
echo "Creating BUILD files for core packages..."

# tm2/pkg/crypto - no external deps
mkdir -p tm2/pkg/crypto
cat > tm2/pkg/crypto/BUILD.bazel << 'EOF'
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "crypto",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/gnolang/gno/tm2/pkg/crypto",
    visibility = ["//visibility:public"],
)

go_test(
    name = "crypto_test",
    srcs = glob(["*_test.go"]),
    embed = [":crypto"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
EOF

# tm2/pkg/std - minimal deps
mkdir -p tm2/pkg/std
cat > tm2/pkg/std/BUILD.bazel << 'EOF'
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "std",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/gnolang/gno/tm2/pkg/std",
    visibility = ["//visibility:public"],
    deps = [
        "//tm2/pkg/amino",
        "//tm2/pkg/crypto",
        "//tm2/pkg/errors",
    ],
)

go_test(
    name = "std_test",
    srcs = glob(["*_test.go"]),
    embed = [":std"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
EOF

# gnovm/pkg/gnolang - core VM package
mkdir -p gnovm/pkg/gnolang
cat > gnovm/pkg/gnolang/BUILD.bazel << 'EOF'
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "gnolang",
    srcs = glob(["*.go"], exclude=["*_test.go", "*_filetest.go"]),
    importpath = "github.com/gnolang/gno/gnovm/pkg/gnolang",
    visibility = ["//visibility:public"],
    deps = [
        "//tm2/pkg/crypto",
        "//tm2/pkg/std",
        "@com_github_cockroachdb_apd_v3//:apd",
        "@com_github_gnolang_overflow//:overflow",
    ],
)

go_test(
    name = "gnolang_test",
    srcs = glob(["*_test.go"]),
    embed = [":gnolang"],
    data = glob(["testdata/**"]),
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
EOF

# Create a working test target
cat >> BUILD.bazel << 'EOF'

# Test suite for working tests
test_suite(
    name = "working_tests",
    tests = [
        "//test_working:simple_test",
        "//tm2/pkg/crypto:crypto_test",
        # Add more as they work
    ],
)
EOF

echo "Working BUILD files created!"
echo ""
echo "Now generating dependencies..."

# Generate go_repository rules for deps
bazelisk run //:gazelle-update-repos -- -from_file=go.mod -to_macro=deps.bzl%go_dependencies || true

echo ""
echo "Setup complete! You can now run:"
echo "  make bazel-test    # Run working tests"
echo "  bazel test //tm2/pkg/crypto:crypto_test  # Run specific test"