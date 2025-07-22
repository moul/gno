#!/bin/bash
# Fix dependency issues for Bazel

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$REPO_ROOT"

echo "Fixing dependency issues for Bazel..."

# Step 1: Create go.work file to handle module issues
cat > go.work << 'EOF'
go 1.23

use .
EOF

# Step 2: Add replace directives for problematic dependencies
cat >> go.mod << 'EOF'

// Bazel dependency fixes
replace (
    // Fix sig-0/insertion-queue module path mismatch
    github.com/sig-0/insertion-queue => github.com/sig-0/iq v0.1.1
)
EOF

# Step 3: Download dependencies with fixes
echo "Downloading dependencies..."
go mod download

# Step 4: Update Gazelle to handle vendored dependencies
cat > misc/bazel/.bazelrc.vendor << 'EOF'
# Vendored dependencies configuration
build --repo_env=GOFLAGS=-mod=readonly
build --repo_env=GO111MODULE=on
EOF

# Step 5: Create a custom Gazelle config
cat > .gazelle.yaml << 'EOF'
# Gazelle configuration
go_prefix: github.com/gnolang/gno
go_generate_proto: false
build_file_name: BUILD.bazel
go_naming_convention: import_alias

# Exclude problematic directories
exclude:
  - misc/bazel
  - vendor
  - bazel-*
  - contribs/gnogenesis
  - misc/genstd
  - tm2/pkg/libtm  # Exclude libtm due to private dep

# Map external repos
go_repository_config:
  sig-0/insertion-queue:
    importpath: github.com/sig-0/iq
    sum: h1:YOUR_SUM_HERE
    version: v0.1.1
EOF

echo "Dependency fixes applied!"
echo ""
echo "Next steps:"
echo "1. Run 'go mod tidy' to clean up"
echo "2. Run 'make bazel-generate' to regenerate BUILD files"
echo "3. Run 'make bazel-test' to test"