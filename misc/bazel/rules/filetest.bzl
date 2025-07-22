"""Rules for running Gno filetests."""

load("@bazel_skylib//lib:paths.bzl", "paths")

def _filetest_impl(ctx):
    """Implementation of filetest rule."""
    
    # Get the gno binary
    gno_bin = ctx.executable._gno
    
    # Prepare test runner
    runner = ctx.actions.declare_file(ctx.label.name + "_runner.sh")
    
    # Collect test files
    test_files = ctx.files.srcs
    
    # Build runner script
    runner_content = """#!/bin/bash
set -euo pipefail

# Set up environment
export GNOHOME="${TEST_TMPDIR}/gno"
mkdir -p "$GNOHOME"

# Track test results
FAILED_TESTS=()
TOTAL_TESTS=0
PASSED_TESTS=0

# Run each filetest
"""
    
    for test_file in test_files:
        runner_content += """
echo "Running filetest: %s"
TOTAL_TESTS=$((TOTAL_TESTS + 1))
if %s test -v --run-file "%s" 2>&1; then
    PASSED_TESTS=$((PASSED_TESTS + 1))
    echo "PASS: %s"
else
    FAILED_TESTS+=("%s")
    echo "FAIL: %s"
fi
echo ""
""" % (test_file.basename, gno_bin.path, test_file.path, 
       test_file.basename, test_file.basename, test_file.basename)
    
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
        files = [gno_bin] + test_files,
    )
    
    return [
        DefaultInfo(
            executable = runner,
            runfiles = runfiles,
        ),
    ]

filetest = rule(
    implementation = _filetest_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = ["_filetest.gno"],
            mandatory = True,
            doc = "Filetest source files (*_filetest.gno)",
        ),
        "_gno": attr.label(
            default = "//gnovm/cmd/gno",
            executable = True,
            cfg = "exec",
            doc = "The gno binary",
        ),
    },
    test = True,
    doc = "Run Gno filetests",
)

def filetest_suite(name, srcs, **kwargs):
    """Create a test suite for multiple filetests.
    
    Args:
        name: Name of the test suite
        srcs: List of filetest source files
        **kwargs: Additional arguments passed to the test rule
    """
    tests = []
    
    for src in srcs:
        # Create individual test name from source file
        test_name = "%s_%s" % (name, paths.basename(src).replace("_filetest.gno", ""))
        
        filetest(
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