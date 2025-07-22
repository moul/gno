# Bazel Setup for Gno

This directory contains the Bazel build configuration for the Gno project.

## Quick Start

```bash
# Run all Bazel tests
make bazel-test

# Build all targets
make bazel-build

# Generate/update BUILD files
make bazel-generate

# Clean Bazel cache and symlinks
make bazel-clean
```

## Current Status

### ✅ Working
- Basic Go test execution for simple packages
- Test caching and incremental builds
- GitHub Actions integration for remote caching
- Automatic symlink setup (lazy initialization)
- Integration with main Makefile

### ⚠️ Known Issues
1. **Amino Package Initialization**: Tests that use amino package fail due to runtime path resolution issues in Bazel's sandbox environment
2. **Private Dependencies**: `github.com/gnolang/libtm` is a private repository that causes dependency resolution failures
3. **Module Path Mismatches**: Some dependencies have import path conflicts (e.g., sig-0/insertion-queue vs sig-0/iq)
4. **Complex Dependencies**: Many packages have complex transitive dependencies that require manual BUILD file creation

### 🔧 Current Workarounds
- Use `test_working/` directory for simple tests that work with Bazel
- Manually create BUILD files for packages with complex dependencies
- Skip packages that depend on problematic external dependencies

## Directory Structure

```
misc/bazel/
├── WORKSPACE          # Bazel workspace configuration
├── BUILD.bazel        # Root BUILD file
├── .bazelrc           # Bazel configuration
├── .bazelversion      # Pinned Bazel version (7.4.1)
├── deps.bzl           # Go dependencies
├── Makefile           # Integration with main Makefile
├── rules/             # Custom Bazel rules
│   ├── gno_test.bzl   # Rule for gno test
│   ├── filetest.bzl   # Rule for filetest
│   └── txtar_test.bzl # Rule for txtar tests
├── scripts/           # Helper scripts
│   ├── setup.sh       # Setup symlinks
│   ├── clean.sh       # Clean Bazel artifacts
│   └── *.sh           # Other utility scripts
└── README.md          # This file
```

## Working Example

The `test_working/` directory contains examples of tests that work with Bazel:

```bash
# Run the working tests
bazelisk test //test_working:all

# Run specific test
bazelisk test //test_working:simple_test
bazelisk test //test_working:simple_crypto_test
```

## Adding New Tests

1. Create a BUILD.bazel file in your package directory:
```starlark
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "mypackage",
    srcs = glob(["*.go"], exclude=["*_test.go"]),
    importpath = "github.com/gnolang/gno/path/to/mypackage",
    visibility = ["//visibility:public"],
    deps = [
        # Add dependencies here
    ],
)

go_test(
    name = "mypackage_test",
    srcs = glob(["*_test.go"]),
    embed = [":mypackage"],
    deps = [
        # Add test dependencies here
    ],
)
```

2. Run the test:
```bash
bazelisk test //path/to/mypackage:mypackage_test
```

## Remote Caching

The GitHub Actions workflow enables remote caching via BuildBuddy:
- Builds are cached between CI runs
- Cache is shared across all branches
- Significantly reduces CI build times

To enable, set the `BUILDBUDDY_API_KEY` secret in your GitHub repository.


## Future Improvements

1. Fix amino package initialization for full test compatibility
2. Vendor or replace problematic external dependencies
3. Add support for all Gno-specific test types (gno test, filetest, txtar)
4. Enable Bazel's module system (bzlmod) when stable
5. Add test sharding for parallel execution
6. Implement proper dependency resolution for all packages

## Troubleshooting

### "Command not found: bazel"
The Makefile will provide installation instructions when Bazel is not found:
```bash
go install github.com/bazelbuild/bazelisk@latest
```

### "Missing strict dependencies"
This error occurs when BUILD files don't list all required dependencies. Add the missing dependencies to the `deps` attribute.

### Tests fail with "dirName if present should be absolute"
This is the amino package initialization issue. The package expects absolute paths but gets relative paths in Bazel's sandbox.

### "BUILD file not found"
Run `make bazel-generate` to create BUILD files using Gazelle.