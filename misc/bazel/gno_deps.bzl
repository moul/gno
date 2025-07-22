"""Dependencies for Gno Bazel rules."""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")

def gno_dependencies():
    """Load dependencies required for Gno rules."""
    
    # You can add any Gno-specific external dependencies here
    # For now, we'll use the local gno binary built by the project
    pass