# Bazel Setup Status

## What Was Accomplished

1. **Complete Bazel Infrastructure**:
   - WORKSPACE file with rules_go, gazelle, protobuf, and skylib
   - Bazel 7.4.1 configuration (to avoid Bazel 8's WORKSPACE deprecation)
   - Custom rules for gno test, filetest, and txtar tests
   - GitHub Actions integration with BuildBuddy caching

2. **Transparent Integration**:
   - All Bazel files in `misc/bazel/` directory
   - Symlinks created lazily from root (gitignored)
   - Simple Makefile integration: `make bazel-test`, `make bazel-build`, etc.
   - Automatic Bazel installation check with helpful error messages

3. **Working Components**:
   - Basic Go tests run successfully
   - Test caching works (see "(cached)" in output)
   - Simple packages without complex dependencies work
   - Example working test in `test_working/` directory

## Current Blockers

1. **Amino Package Initialization**:
   - The amino package uses `runtime.Caller` to get absolute paths
   - In Bazel's sandbox, this returns relative paths causing panics
   - Affects most tm2 and gnovm packages that use amino

2. **Dependency Issues**:
   - Private repository: `github.com/gnolang/libtm`
   - Module path mismatches: `sig-0/insertion-queue` vs `sig-0/iq`
   - Complex transitive dependencies requiring manual BUILD file creation

3. **Gazelle Limitations**:
   - Cannot automatically resolve all dependencies
   - Requires manual BUILD file creation for many packages
   - Excludes problematic paths like `tm2/pkg/libtm`

## Next Steps to Make It Fully Work

1. **Fix Amino Path Issue**:
   - Option 1: Patch amino to handle relative paths in tests
   - Option 2: Use Bazel's `--test_env=PWD=$(pwd)` to set absolute paths
   - Option 3: Create a test wrapper that sets up the environment

2. **Handle Private Dependencies**:
   - Vendor the private dependencies
   - Or use go_repository with authentication
   - Or replace with public alternatives

3. **Generate Working BUILD Files**:
   - Create a script that generates BUILD files with proper dependencies
   - Start with core packages and work outward
   - Manually specify dependencies for complex packages

4. **Implement Gno-specific Tests**:
   - Complete the gno_test rule implementation
   - Add support for filetest execution
   - Implement txtar test runner

## Commands That Work Now

```bash
# Run the simple test
make bazel-test

# Run specific working test
bazelisk test //test_working:simple_test

# Clean everything
make bazel-clean

# Generate dependencies
make bazel-generate
```

## Files Created

- `/misc/bazel/WORKSPACE` - Main Bazel workspace
- `/misc/bazel/BUILD.bazel` - Root BUILD file
- `/misc/bazel/.bazelrc` - Bazel configuration
- `/misc/bazel/.bazelversion` - Pinned Bazel version
- `/misc/bazel/deps.bzl` - Go dependencies
- `/misc/bazel/Makefile` - Integration commands
- `/misc/bazel/rules/*.bzl` - Custom test rules
- `/misc/bazel/scripts/*.sh` - Helper scripts
- `/misc/bazel/.github/workflows/bazel-ci.yml` - GitHub Actions
- `/test_working/` - Working example tests
- Various BUILD.bazel files for tm2/pkg packages