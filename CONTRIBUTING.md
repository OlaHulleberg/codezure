# Contributing to Codzure

Thanks for your interest in contributing! This project mirrors conventions from `clauderock` and aims for a clean, lightweight CLI.

## Setup
- Requires Go 1.23+
- `az` CLI installed and logged in (`az login`)
- Run `make build` to build locally

## Development
- Commands implemented via Cobra under `cmd/`
- Internal packages under `internal/`
- Profiles stored as JSON at `~/.codzure/profiles/<name>.json`
- Current profile tracked in `~/.codzure/current-profile.txt`

## Release
- Tag as `vX.Y.Z` to trigger GoReleaser workflow
- Binaries published to GitHub Releases

## Code Style
- Keep changes minimal and focused
- Prefer small, composable functions
- Avoid storing secrets in files; fetch via `az` just-in-time

## Issues and PRs
- Please include clear reproduction steps and expected behavior
- Reference related issues where possible
