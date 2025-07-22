#!/usr/bin/env bash

# Script to bump Go versions across the repository
# Usage: ./misc/bump-go.sh [--check-only]

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
VERSIONS_FILE="${ROOT_DIR}/.versions"
CHECK_ONLY=false

# Parse arguments
for arg in "$@"; do
    case $arg in
        --check-only)
            CHECK_ONLY=true
            shift
            ;;
        *)
            echo "Unknown option: $arg"
            echo "Usage: $0 [--check-only]"
            exit 1
            ;;
    esac
done

# Load version configuration
if [ ! -f "$VERSIONS_FILE" ]; then
    echo -e "${RED}Error: .versions file not found${NC}"
    exit 1
fi

source "$VERSIONS_FILE"

echo "Go Version Bump Tool"
echo "===================="
echo ""
echo "Target versions:"
echo "  Minimum (N-1): ${GO_MIN_VERSION}"
echo "  Current (N):   ${GO_CURRENT_VERSION}"
echo "  Docker:        ${GO_DOCKER_VERSION}"
echo "  CI/CD:         ${GO_CI_VERSION}"
echo ""

if [ "$CHECK_ONLY" = true ]; then
    echo -e "${BLUE}Running in check-only mode${NC}"
    echo ""
fi

# Track changes
CHANGES=0

# Function to update file content
update_file() {
    local file=$1
    local pattern=$2
    local replacement=$3
    local description=$4

    if [ -f "$file" ]; then
        if grep -q "$pattern" "$file"; then
            if [ "$CHECK_ONLY" = true ]; then
                echo -e "  ${YELLOW}Would update${NC} ${file#$ROOT_DIR/}: $description"
                ((CHANGES++))
            else
                sed -i "$replacement" "$file"
                echo -e "  ${GREEN}Updated${NC} ${file#$ROOT_DIR/}: $description"
                ((CHANGES++))
            fi
        fi
    fi
}

# Function to update YAML files (special handling to preserve syntax)
update_yaml_go_version() {
    local file=$1
    local new_version=$2
    local description=$3
    
    if [ -f "$file" ]; then
        # Create a temporary file for the updated content
        local tmp_file="${file}.tmp"
        local changed=false
        
        # Process the file line by line to preserve YAML structure
        while IFS= read -r line; do
            if [[ "$line" =~ ^[[:space:]]*go-version:[[:space:]]* ]]; then
                # Extract indentation and any existing structure
                local indent="${line%%go-version:*}"
                local rest="${line#*go-version:}"
                
                # Check if it's an array or single value
                if [[ "$rest" =~ ^[[:space:]]*\[ ]]; then
                    # It's an array format
                    echo "${indent}go-version: [\"${new_version}\"]" >> "$tmp_file"
                    changed=true
                elif [[ "$rest" =~ ^[[:space:]]*\" ]]; then
                    # It's a quoted single value
                    echo "${indent}go-version: \"${new_version}\"" >> "$tmp_file"
                    changed=true
                else
                    # It's an unquoted single value
                    echo "${indent}go-version: \"${new_version}\"" >> "$tmp_file"
                    changed=true
                fi
            else
                echo "$line" >> "$tmp_file"
            fi
        done < "$file"
        
        if [ "$changed" = true ]; then
            if [ "$CHECK_ONLY" = true ]; then
                echo -e "  ${YELLOW}Would update${NC} ${file#$ROOT_DIR/}: $description"
                ((CHANGES++))
                rm -f "$tmp_file"
            else
                mv "$tmp_file" "$file"
                echo -e "  ${GREEN}Updated${NC} ${file#$ROOT_DIR/}: $description"
                ((CHANGES++))
            fi
        else
            rm -f "$tmp_file"
        fi
    fi
}

# Update go.mod files
echo "Updating go.mod files..."
while IFS= read -r -d '' gomod; do
    if [ -f "$gomod" ]; then
        CURRENT=$(grep -E '^go [0-9]+\.[0-9]+(\.[0-9]+)?' "$gomod" | awk '{print $2}' || echo "")
        if [ -n "$CURRENT" ] && [ "$CURRENT" != "$GO_MIN_VERSION" ]; then
            # Extract major.minor version only
            CURRENT_BASE=$(echo "$CURRENT" | cut -d. -f1,2)
            if [ "$CURRENT_BASE" != "$GO_MIN_VERSION" ]; then
                update_file "$gomod" "^go [0-9]\\+\\.[0-9]\\+\\(\\.[0-9]\\+\\)\\?" "s/^go [0-9]\\+\\.[0-9]\\+\\(\\.[0-9]\\+\\)\\?/go ${GO_MIN_VERSION}/" "go ${CURRENT} → go ${GO_MIN_VERSION}"
            else
                # Version is 1.23.x, change to just 1.23
                update_file "$gomod" "^go ${GO_MIN_VERSION}\\.[0-9]\\+" "s/^go ${GO_MIN_VERSION}\\.[0-9]\\+/go ${GO_MIN_VERSION}/" "go ${CURRENT} → go ${GO_MIN_VERSION}"
            fi
        fi
    fi
done < <(find "$ROOT_DIR" -name "go.mod" -type f -print0)
echo ""

# Update CI workflows (using special YAML handling)
echo "Updating CI workflows..."
CI_FILES=(
    ".github/workflows/main_template.yml"
    ".github/workflows/gnofmt_template.yml"
    ".github/workflows/examples.yml"
    ".github/workflows/gnoland.yml"
)

for ci_file in "${CI_FILES[@]}"; do
    FULL_PATH="${ROOT_DIR}/${ci_file}"
    if [ -f "$FULL_PATH" ]; then
        update_yaml_go_version "$FULL_PATH" "$GO_CI_VERSION" "go-version → ${GO_CI_VERSION}"
    fi
done
echo ""

# Update Dockerfiles
echo "Updating Dockerfiles..."
DOCKER_FILES=$(find "$ROOT_DIR" -name "Dockerfile*" -type f)

for dockerfile in $DOCKER_FILES; do
    if [ -f "$dockerfile" ]; then
        # Update golang base image versions
        update_file "$dockerfile" \
            "FROM golang:[0-9]\\+\\.[0-9]\\+-alpine" \
            "s/FROM golang:[0-9]\\+\\.[0-9]\\+-alpine/FROM golang:${GO_DOCKER_VERSION}/" \
            "golang base image → ${GO_DOCKER_VERSION}"
        
        # Also update non-alpine golang images
        update_file "$dockerfile" \
            "FROM golang:[0-9]\\+\\.[0-9]\\+\\s*$" \
            "s/FROM golang:[0-9]\\+\\.[0-9]\\+/FROM golang:${GO_DOCKER_VERSION}/" \
            "golang base image → ${GO_DOCKER_VERSION}"
    fi
done
echo ""

# Update CONTRIBUTING.md
echo "Updating documentation..."
CONTRIB_FILE="${ROOT_DIR}/CONTRIBUTING.md"
if [ -f "$CONTRIB_FILE" ]; then
    # Update the version requirement
    update_file "$CONTRIB_FILE" \
        "- Go (version [0-9]\\+\\.[0-9]\\++)" \
        "s/- Go (version [0-9]\\+\\.[0-9]\\++)/- Go (version ${GO_MIN_VERSION}+)/" \
        "minimum Go version requirement"
    
    # Update example versions in the policy section
    update_file "$CONTRIB_FILE" \
        "- \\*\\*N\\*\\*: Current stable Go version (e.g., Go [0-9]\\+\\.[0-9]\\+)" \
        "s/- \\*\\*N\\*\\*: Current stable Go version (e.g., Go [0-9]\\+\\.[0-9]\\+)/- **N**: Current stable Go version (e.g., Go ${GO_CURRENT_VERSION})/" \
        "current version example"
    
    update_file "$CONTRIB_FILE" \
        "- \\*\\*N-1\\*\\*: Previous stable Go version (e.g., Go [0-9]\\+\\.[0-9]\\+)" \
        "s/- \\*\\*N-1\\*\\*: Previous stable Go version (e.g., Go [0-9]\\+\\.[0-9]\\+)/- **N-1**: Previous stable Go version (e.g., Go ${GO_MIN_VERSION})/" \
        "previous version example"
    
    update_file "$CONTRIB_FILE" \
        "- We can use features from the N-1 version (currently Go [0-9]\\+\\.[0-9]\\+)" \
        "s/- We can use features from the N-1 version (currently Go [0-9]\\+\\.[0-9]\\+)/- We can use features from the N-1 version (currently Go ${GO_MIN_VERSION})/" \
        "N-1 features note"
fi
echo ""

# Run go mod tidy if not in check-only mode
if [ "$CHECK_ONLY" = false ] && [ $CHANGES -gt 0 ]; then
    echo "Running go mod tidy..."
    if [ -f "${ROOT_DIR}/misc/mod_tidy.sh" ]; then
        echo -e "  ${BLUE}Running${NC} misc/mod_tidy.sh"
        cd "$ROOT_DIR" && ./misc/mod_tidy.sh
        echo -e "  ${GREEN}✓ Completed${NC}"
    else
        echo -e "  ${YELLOW}⚠ misc/mod_tidy.sh not found, skipping${NC}"
    fi
    echo ""
fi

# Summary
echo "Summary"
echo "-------"
if [ $CHANGES -eq 0 ]; then
    echo -e "${GREEN}✓ All files are already up to date!${NC}"
else
    if [ "$CHECK_ONLY" = true ]; then
        echo -e "${YELLOW}Found ${CHANGES} file(s) that need updating${NC}"
        echo ""
        echo "Run without --check-only to apply changes:"
        echo "  $0"
    else
        echo -e "${GREEN}✓ Updated ${CHANGES} file(s) successfully!${NC}"
        echo ""
        echo "Next steps:"
        echo "1. Review the changes: git diff"
        echo "2. Run tests: make test"
        echo "3. Commit changes: git commit -am \"chore: bump Go version to ${GO_MIN_VERSION}\""
    fi
fi

# Exit with error if check-only found issues
if [ "$CHECK_ONLY" = true ] && [ $CHANGES -gt 0 ]; then
    exit 1
fi