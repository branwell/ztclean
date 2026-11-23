# ztclean
Clean up Zoom transcripts: converts a Zoom `.vtt` transcript into a compact Markdown file.

## Build

Requires Go (see `go.mod` for the version). No other dependencies and no network access needed.

```sh
go build .
```

This produces the `ztclean` binary in the repository root. To install it into your Go bin directory instead (`$(go env GOPATH)/bin`), run `go install .`.

## Usage

```sh
./ztclean -f meeting.vtt              # writes meeting.md
./ztclean -f meeting.vtt -o notes.md  # writes notes.md
```

See `spec.md` for the full cleanup rules.

## Development

```sh
gofmt -l .     # should print nothing
go vet ./...
go test ./...
```
