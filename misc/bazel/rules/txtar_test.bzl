"""Rules for running txtar integration tests."""

load("@io_bazel_rules_go//go:def.bzl", "go_test")
load("@bazel_skylib//lib:paths.bzl", "paths")

def _txtar_test_impl(ctx):
    """Implementation of txtar_test rule."""
    
    # Get required binaries
    gnoland_bin = ctx.executable._gnoland
    gnokey_bin = ctx.executable._gnokey
    
    # Prepare test runner
    runner = ctx.actions.declare_file(ctx.label.name + "_runner.sh")
    
    # Collect test files
    test_files = ctx.files.srcs
    
    # Build runner script
    runner_content = """#!/bin/bash
set -euo pipefail

# Set up environment
export TEST_TMPDIR="${TEST_TMPDIR:-/tmp}"
export TESTWORK="${TEST_TMPDIR}/txtar_work"
export PATH="%s:%s:${PATH}"

# Ensure binaries are accessible
ln -sf "%s" gnoland
ln -sf "%s" gnokey

# Create work directory
mkdir -p "$TESTWORK"

# Track test results
FAILED_TESTS=()
TOTAL_TESTS=0
PASSED_TESTS=0

# Function to run a single txtar test
run_txtar_test() {
    local test_file="$1"
    local test_name=$(basename "$test_file" .txtar)
    
    echo "Running txtar test: $test_name"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Create isolated test directory
    TEST_DIR="$TESTWORK/$test_name"
    mkdir -p "$TEST_DIR"
    cd "$TEST_DIR"
    
    # Run the test using testscript
    if go run github.com/rogpeppe/go-internal/cmd/testscript -v "$test_file" 2>&1; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "PASS: $test_name"
    else
        FAILED_TESTS+=("$test_name")
        echo "FAIL: $test_name"
    fi
    
    cd - > /dev/null
    echo ""
}

# Run each txtar test
""" % (paths.dirname(gnoland_bin.path), paths.dirname(gnokey_bin.path),
       gnoland_bin.path, gnokey_bin.path)
    
    for test_file in test_files:
        runner_content += 'run_txtar_test "%s"\n' % test_file.path
    
    runner_content += """
# Report results
echo "=============================="
echo "Test Summary:"
echo "Total tests: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS" 
echo "Failed: ${#FAILED_TESTS[@]}"

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo ""
    echo "Failed tests:"
    for test in "${FAILED_TESTS[@]}"; do
        echo "  - $test"
    done
    exit 1
fi

exit 0
"""
    
    ctx.actions.write(
        output = runner,
        content = runner_content,
        is_executable = True,
    )
    
    runfiles = ctx.runfiles(
        files = [gnoland_bin, gnokey_bin] + test_files,
    )
    
    return [
        DefaultInfo(
            executable = runner,
            runfiles = runfiles,
        ),
    ]

txtar_test = rule(
    implementation = _txtar_test_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".txtar"],
            mandatory = True,
            doc = "Txtar test files",
        ),
        "_gnoland": attr.label(
            default = "//gno.land/cmd/gnoland",
            executable = True,
            cfg = "exec",
            doc = "The gnoland binary",
        ),
        "_gnokey": attr.label(
            default = "//gno.land/cmd/gnokey",
            executable = True,
            cfg = "exec",
            doc = "The gnokey binary",
        ),
    },
    test = True,
    doc = "Run txtar integration tests",
)

def txtar_test_suite(name, srcs, **kwargs):
    """Create a test suite for multiple txtar tests.
    
    Args:
        name: Name of the test suite
        srcs: List of txtar test files
        **kwargs: Additional arguments passed to the test rule
    """
    tests = []
    
    for src in srcs:
        # Create individual test name from source file
        test_name = "%s_%s" % (name, paths.basename(src).replace(".txtar", ""))
        
        txtar_test(
            name = test_name,
            srcs = [src],
            **kwargs
        )
        tests.append(":" + test_name)
    
    # Create test suite
    native.test_suite(
        name = name,
        tests = tests,
    )