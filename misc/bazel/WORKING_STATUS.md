# Bazel Setup - Current Working Status

## What Works

✅ Basic Bazel infrastructure is set up and working
✅ Simple Go tests can be run with Bazel
✅ Custom rules for Gno tests are defined
✅ Lazy setup (auto-creates symlinks, auto-installs deps)
✅ Clean commands and transparent integration

## What Doesn't Work (Yet)

❌ Full Gazelle integration - blocked by:
   - Private repository dependencies (gnolang/libtm)
   - Module path mismatches (sig-0/insertion-queue vs sig-0/iq)
   - Complex transitive dependencies

## How to Use (Current State)

1. **Run the working test:**
   ```bash
   make bazel-test
   ```
   This runs a simple test to verify Bazel is working.

2. **To add more tests:**
   - Create BUILD.bazel files manually in directories with tests
   - Use simple go_test rules without complex dependencies
   - Add them to the test suite

## How to Fix

To get full Bazel support working, you need to:

1. **Fix module dependencies:**
   - Replace private repos with public mirrors or vendor them
   - Fix module path mismatches in go.mod
   - Use go mod vendor and configure Bazel to use vendored deps

2. **Alternative: Use vendored dependencies:**
   ```bash
   go mod vendor
   # Then configure Bazel to use vendor/ directory
   ```

3. **Alternative: Manual BUILD files:**
   - Skip Gazelle for problematic packages
   - Write BUILD.bazel files manually for core packages
   - Gradually add more as dependencies are fixed

## Example Working BUILD.bazel

```starlark
load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "mypackage_test",
    srcs = glob(["*_test.go"]),
    size = "small",
    # Add only deps that work
)
```

## Next Steps

1. Fix go.mod dependencies
2. Consider vendoring for hermetic builds
3. Create BUILD files manually for core packages
4. Gradually expand coverage as deps are fixed