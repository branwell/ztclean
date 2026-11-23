# ztclean — Agent Guide

## Mission

`ztclean` is a small command-line tool that converts Zoom `.vtt` transcripts into compact Markdown.

Development is exploratory and incremental. Begin with minimal scaffolding and evolve the application through small, reviewable changes. Favor readable, maintainable code, clear separation of concerns, and minimal coupling. Avoid speculative abstractions and unnecessary complexity.

When tradeoffs arise, prioritize:

1. Approved requirements and user-visible behavior.
2. Security and safe operation (see **Operating constraints**).
3. Simplicity, readability, and maintainability.
4. Performance when required by an approved use case.

## Operating constraints

These constraints are not up for silent revision. Recommend relaxing one only if keeping it would force hacks or workarounds that violate simplicity and best practices. Never relax one unless explicitly instructed.

- **Avoid network.** Runtime and build must not depend on internet access or external services. The tool only reads and writes local files.
- **Avoid outside dependencies.** The Go code uses the standard library only. Ask before adding any library, module or tool; if one is approved, it is vendored into the repository.
- **On-prem.** The application may run in a locked-down environment. Assume no package manager, no telemetry, and no outbound connections.

## Source of truth

- `spec.md` is authoritative for approved behavior: command-line usage, cleanup rules, examples, and design decisions.
- `README.md` is orientation only, not a source of requirements.
- If authoritative documents conflict with each other or with approved behavior, report the discrepancy and wait for instructions. Never silently reconcile them.

## Stack and boundaries

- **Language:** Go (standard library only)
- **Interface:** command line; no server, no UI

Layout (verify the working tree before editing any listed path):

```text
ztclean/
├── AGENTS.md            # Agent workflow and boundaries
├── README.md            # Orientation: what it is
├── spec.md              # Authoritative requirements and design decisions
├── go.mod
├── main.go              # Entry point: flags, file read/write
├── clean.go             # Parse VTT cues, group into paragraphs, render Markdown
├── ascii.go             # Unicode-to-ASCII lookup table
└── clean_test.go        # Tests, including the spec.md examples as golden cases
```

## Conventions

`gofmt`; wrap errors with operation context (`fmt.Errorf("read transcript: %w", err)`); keep the cleanup logic in pure functions (string in, string out) separate from file I/O so it stays easy to test. When a spec example changes, update the matching golden test.

## Commands

```sh
gofmt -l .        # must print nothing
go vet ./...
go test ./...
go build .        # produces ./ztclean
```

Run the tests, `go vet` and `gofmt` before reporting Go work as done.

## Scope rules

**ALWAYS**
- Implement only approved scope; deferred ideas are not implied tasks.
- Read `spec.md` before any task.
- Record agreed requirements and design decisions as they are made, in `spec.md`: requirements in the section that owns them, decisions under **Design decisions**.
- Document any new command-line flag or behavior in `spec.md` (**Running the tool** or **Cleanup Rules**).

**ASK FIRST**
- Before adding libraries, modules, flags, build steps, or features.
- Before introducing any dependency on the network or on a specific operating system.

**NEVER**
- Silently resolve conflicts between active documents; report the conflict and wait for instructions.
- Commit, reset, or discard user work unless explicitly asked.
- Claim that checks were run when they were not.
- Hard-code absolute or machine-specific paths.
- Read, search, use, create, or modify anything in `_ignore/`. It holds the human's private working files, which are not a source of truth. Treat it as if it does not exist.
