# ADR: gnopie — httpie-inspired CLI for gno.land (PR #5444)

## Context

Interacting with gno.land chains today requires knowing the right `gnokey` invocations for every operation, understanding which RPC endpoints to use, and manually handling gas estimation. There is no ergonomic developer tool for quickly querying realms, reading source, or sending transactions — comparable to what `httpie` provides for HTTP APIs.

## Decision

Add `gnopie` to `contribs/` as a standalone CLI tool that provides:

- **Verb-based API** (GET, EVAL, READ, INSPECT, CALL, RUN) inspired by httpie's HTTP verbs
- **Auto-discovery** of RPC endpoints via `<meta name="gnoconnect:rpc">` tags — no manual config needed
- **Smart dispatch** under the default GET verb: realm paths → Render, function calls → EVAL, symbols → READ
- **gnoweb URL support**: paste any `https://gno.land/...` URL directly
- **Auto-gas**: simulate → estimate → broadcast with configurable buffer (default 20%)
- **Cross-realm auto-injection**: `cross` prepended automatically for crossing functions
- **File-based caching**: discovery results cached 24h, source/function queries 1h
- **JSON output** (`--json`) for scripting and piping
- **`--print-gnokey-command`** to show the equivalent gnokey invocation without executing

## Alternatives Considered

- **Shell script wrapper around `gnokey`**: too brittle, can't provide smart dispatch or auto-discovery.
- **Extending `gnokey` itself**: gnokey is a key management and transaction tool; adding httpie-style verbs would mix concerns.
- **Using gnoweb**: web UI, not scriptable.
- **txtar/integration tests only**: gnopie is a developer UX tool, not just test infrastructure.

## Architecture

```
gnopie [flags] [VERB] <expression>
```

- `paths.go` — `ParsePath()` parses any expression into a typed `GnoPath` (domain, pkg, symbol, args, file, address)
- `get.go` — GET, EVAL, READ, INSPECT, plus sub-functions: `getRender`, `readSource`, `readFile`, `inspectPackage`, `inspectAddress`, etc.
- `call.go` — CALL with auto-gas simulation, `printGnokeyCmd` for dry-run display
- `run.go` — RUN: generates `main.gno` with the call, submits via `MsgRun`
- `discover.go` — HTTP fetch + meta-tag parsing, TOML cache per domain
- `querycache.go` — SHA-256-keyed file cache for `vm/qfile` and `vm/qfuncs` results
- `config.go` / `cmd_config.go` — TOML-based config (`key`, `gas-buffer`)
- `cmd_completion.go` — bash/zsh/fish shell completions
- `cmd_version.go` — version subcommand

## Consequences

- Developers can query any gno.land realm with a single command and zero configuration.
- CALL/RUN transactions require a configured key (`gnopie config set key=<name>`).
- Discovery cache means first call to a domain may take a second; subsequent calls are instant.
- Source code queries (qfile) are cached for 1h — stale source is possible during active development; use `--debug` to inspect cache hits.
- gnopie is intentionally opinionated: it makes decisions (e.g., auto-cross injection, Render as default for packages) that may surprise users expecting raw RPC behavior.
