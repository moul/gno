# Handoff: misc/build-wasm.sh (PR #5359)

## Branch
`moul:gno-release-wasm` → PR https://github.com/gnolang/gno/pull/5359

## What this PR does
Adds `misc/build-wasm.sh`: a script to build `gno.wasm` + `root.zip` from `gnovm/cmd/gno`
and optionally publish them to the GitHub release for the given tag.

## Usage
```
./misc/build-wasm.sh                        # prints help
./misc/build-wasm.sh chain/gnoland1.0       # build only → gnovm/build/
./misc/build-wasm.sh chain/gnoland1.0 --push  # build + publish to GH release
./misc/build-wasm.sh HEAD --push            # build current HEAD + publish (must be on a tag)
./misc/build-wasm.sh chain/gnoland1.0 --push ./out/  # custom output dir
```

## What the script does
1. `git checkout <tag>` (restores original HEAD on exit via trap)
2. `GOARCH=wasm GOOS=js go build -ldflags "-X .../gnoenv._GNOROOT=<repo>/" -o gnovm/build/gno.wasm ./cmd/gno` (from gnovm/)
3. `zip root.zip gnovm/stdlibs gnovm/tests/stdlibs examples`
4. If `--push`: `gh release create <tag>` (if not exists) + `gh release upload --clobber`

## Published asset URLs (URL-encoded tag)
- `https://github.com/gnolang/gno/releases/download/chain%2Fgnoland1.0/gno.wasm`
- `https://github.com/gnolang/gno/releases/download/chain%2Fgnoland1.0/root.zip`

## Tag created
`chain/gnoland1.0` at `e3d37187d4747a5a371c9a0b668d314ce28b55ba`

## TODO
1. Run `./misc/build-wasm.sh chain/gnoland1.0 --push` to publish WASM for the betanet tag
2. Merge this PR (or land it however the gno team prefers — might go in misc/Makefile instead)
3. Future: can re-add the GitHub Actions workflow (was removed to keep PR minimal) for auto-publishing on chain/* tag push

## Context: why this exists
play.gno.land (gnostudio/studio PR #1612) switched from GCS to GitHub releases for WASM hosting.
Each `chain/*` tag in gnolang/gno corresponds to a GnoVM version in play.gno.land's runtime selector.
