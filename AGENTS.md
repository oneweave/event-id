# Agent Instructions

This repository is a small Go library for generating and decoding prefixed UUIDv7 identifiers. Keep changes focused, minimal, and aligned with the existing public API documented in [README.md](README.md).

## Project Conventions

- Treat this project as a library package, not an application.
- Preserve the canonical format: 3 lowercase prefix characters plus 26 Crockford Base32 characters.
- Keep prefix validation strict: only lowercase ASCII letters, exactly 3 characters.
- Do not reintroduce regex-based parsing or look-alike Crockford character mappings unless explicitly requested.
- `Decode` accepts normalized input, including uppercase letters and separator whitespace/hyphens, but still rejects invalid prefixes and malformed payloads.

## Working Rules

- Prefer small edits over broad refactors.
- Use the existing tests as the primary safety net.
- Update or add tests when behavior changes.
- Do not duplicate README content here; link to it instead.

## Validation

- Run `go test ./...` after code changes.
- Use `go test -bench . ./...` only when changing performance-sensitive code.

## Reference Material

- API and examples: [README.md](README.md)