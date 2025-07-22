"""Rules for running gno test."""

load("@io_bazel_rules_go//go:def.bzl", "go_context")
load("@bazel_skylib//lib:paths.bzl", "paths")

def _gno_test_impl(ctx):
    """Implementation of gno_test rule."""
    
    # Get the gno binary
    gno_bin = ctx.executable._gno
    
    # Prepare test runner script
    runner = ctx.actions.declare_file(ctx.label.name + "_runner.sh")
    
    # Collect test files and dependencies
    test_files = ctx.files.srcs
    data_files = ctx.files.data
    
    # Build the test command
    test_args = ["test"]
    if ctx.attr.verbose:
        test_args.append("-v")
    if ctx.attr.run:
        test_args.extend(["-run", ctx.attr.run])
    test_args.append(ctx.attr.package_path)
    
    # Create runner script
    runner_content = """#!/bin/bash
set -euo pipefail

# Set up environment
export GNOHOME="${TEST_TMPDIR}/gno"
mkdir -p "$GNOHOME"

# Copy test files to runfiles directory
TEST_DIR="${TEST_TMPDIR}/test_files"
mkdir -p "$TEST_DIR"
"""
    
    # Add file copy commands
    for f in test_files + data_files:
        runner_content += 'cp -r "%s" "$TEST_DIR/"\n' % f.path
    
    # Add the gno test command
    runner_content += """
# Run gno test
cd "$TEST_DIR"
exec %s %s
""" % (gno_bin.path, " ".join(test_args))
    
    ctx.actions.write(
        output = runner,
        content = runner_content,
        is_executable = True,
    )
    
    # Return test executable info
    runfiles = ctx.runfiles(
        files = [gno_bin] + test_files + data_files,
    )
    
    return [
        DefaultInfo(
            executable = runner,
            runfiles = runfiles,
        ),
    ]

gno_test = rule(
    implementation = _gno_test_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".gno"],
            doc = "Gno test source files",
        ),
        "data": attr.label_list(
            allow_files = True,
            doc = "Additional data files needed for tests",
        ),
        "package_path": attr.string(
            mandatory = True,
            doc = "Package path to test (e.g., ./examples/gno.land/p/demo/...)",
        ),
        "verbose": attr.bool(
            default = False,
            doc = "Run tests in verbose mode",
        ),
        "run": attr.string(
            doc = "Run only tests matching regex",
        ),
        "_gno": attr.label(
            default = "//gnovm/cmd/gno",
            executable = True,
            cfg = "exec",
            doc = "The gno binary",
        ),
    },
    test = True,
    doc = "Run gno tests on Gno packages",
)