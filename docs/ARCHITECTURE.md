# CLI Tools Architecture

## Goals
- Single Go binary for UX, concurrency, and cross-platform distribution.
- Python plugins for specialized scanning and OSINT workflows.
- Consistent JSON-friendly interface between layers.

## Layout
- cmd/ct: main entrypoint
- internal/cli: subcommand routing
- internal/tui: dashboard UI (Bubble Tea)
- internal/runner: process runners for plugins and external tools
- plugins/python: Python task plugins
- internal/runner/queue.go: concurrency-limited job queue

## Plugin contract (initial)
- Executable python script.
- Accepts arguments from Go.
- Writes human-readable or JSON output to stdout.
- Exit code indicates success (0) or failure (non-zero).

## Result schema
Each job returns a JSON-friendly result object:
- id: unique job id
- command: executed command
- args: command arguments
- started_at, finished_at: timestamps
- duration_ms: elapsed time
- exit_code: process exit code
- status: success | failed
- stdout, stderr: captured output
- error: error message if any

## Dependency checks
Commands verify required tools on the host. Missing dependencies prompt the user for consent, then attempt installation with the available package manager for the OS.

## Configuration
Default config lives in `configs/default.yaml`. Override with `CT_CONFIG` or `ct --config path`.

## Next steps
- Define per-command JSON output parsing for richer structured data.
- Add a persistent results store (sqlite or local files).
- Introduce config for tool paths (nmap, dns tools, etc.).


Signed-off-by: ronikoz
