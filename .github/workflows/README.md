# Simplified CI Structure

This directory contains the simplified GitHub Actions workflows for the Gno project.

## Workflows Overview

### Core Workflows

1. **`ci.yml`** - Main CI pipeline
   - Runs on every PR and push to master/staging
   - Quick checks (go mod tidy, generated files)
   - Linting with golangci-lint
   - Gno formatting checks
   - Component tests (gnovm, gno.land, tm2)
   - Examples and contribs testing
   - Smart path-based triggering

2. **`release.yml`** - Release automation
   - Triggered by version tags (v*)
   - Builds release binaries
   - Creates GitHub releases with changelog
   - Marks pre-releases automatically

3. **`nightly.yml`** - Nightly builds
   - Runs daily at midnight UTC
   - Builds Docker images
   - Runs integration tests
   - Can be manually triggered

### Automation Workflows

4. **`bot.yml`** - Bot actions
   - Auto-labels PRs based on changed files
   - Handles bot commands (/bot run tests, /bot label)
   - Cleans up stale issues/PRs
   - Checks for merge conflicts

5. **`deploy.yml`** - Deployment
   - Deploys to staging on push
   - Manual deployment to test environments
   - Uses AWS for deployment (example)

## Key Improvements

1. **Consolidated workflows** - Reduced from 32 to 5 main workflows
2. **Clear naming** - Each workflow has a specific purpose
3. **No complex templates** - Direct, readable job definitions
4. **Smart triggers** - Path-based and label-based triggering
5. **Consistent patterns** - Similar structure across workflows

## Running Tests Locally

```bash
# Run all tests
make test

# Run specific component tests
cd gnovm && make test
cd gno.land && make test
cd tm2 && make test

# Run linting
make lint

# Check formatting
make fmt
```

## Workflow Triggers

- **CI**: PRs, pushes to master/staging
- **Release**: Version tags (v*)
- **Nightly**: Daily schedule or manual
- **Bot**: PR events, comments, schedule
- **Deploy**: Push to staging or manual

## Matrix Strategy

The CI uses matrices for:
- Contribs (gnodev, gnofaucet, etc.)
- Misc tools (autocounterd, genproto, etc.)

This allows parallel testing while keeping the workflow simple.

## Coverage

All tests generate coverage reports that are uploaded to Codecov with appropriate flags for each component.