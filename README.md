<div align="center">

# kyber

**A function-level Go code-quality analyzer — cyclomatic complexity, testability, readability, and more.**

[![CI](https://github.com/jedi-knights/kyber/actions/workflows/ci.yml/badge.svg)](https://github.com/jedi-knights/kyber/actions/workflows/ci.yml)
[![Release](https://github.com/jedi-knights/kyber/actions/workflows/release.yml/badge.svg)](https://github.com/jedi-knights/kyber/actions/workflows/release.yml)
[![GoReleaser](https://github.com/jedi-knights/kyber/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/jedi-knights/kyber/actions/workflows/goreleaser.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedi-knights/kyber)](https://goreportcard.com/report/github.com/jedi-knights/kyber)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[Installation](#installation) · [Quickstart](#quickstart) · [Metrics](#metrics) · [Configuration](#configuration) · [GitHub Action](#github-action) · [Development](#development) · [Contributing](#contributing)

</div>

---

kyber walks a Go codebase, parses every function, and scores each one against a registered set of code-quality metrics. Output is human-readable text, machine-readable JSON, or SARIF v2.1.0 for GitHub code scanning.

```
$ kyber analyze ./...

internal/adapters/parser/goast.go
  parseFile                 cyclomatic=12 ⚠   readability=0.58 ⚠   testability=0.71
  extractInterfaces         cyclomatic=4      readability=0.81      testability=0.85

internal/domain/metrics/cyclomatic.go
  Analyze                   cyclomatic=6      readability=0.74      testability=0.92

Functions: 47   Findings: 3   Files: 12   Time: 28ms
```

The architecture is built around a single `Metric` interface. Adding a new metric (cognitive complexity, halstead, nesting depth, parameter count) is one new file in `internal/domain/metrics/` plus one registration line — nothing else needs to change.

## Why kyber?

Existing Go linters (`golangci-lint`, `gocyclo`, `revive`) enforce per-rule thresholds across an entire codebase. They tell you which lines violate a rule; they do not tell you which functions are healthy, mediocre, or risky overall. kyber works at the **function level**, producing per-function scores so you can:

- See at a glance which functions in a package carry the most quality risk
- Track readability and testability over time, not just complexity
- Plug a new metric into the same scoring pipeline without writing another linter

## Installation

### Go install (recommended for Go developers)

Requires Go 1.23 or later.

```bash
# Latest release
go install github.com/jedi-knights/kyber/cmd/kyber@latest

# Pinned to a specific version
go install github.com/jedi-knights/kyber/cmd/kyber@v0.1.0

# HEAD of main (unstable; pre-release work)
go install github.com/jedi-knights/kyber/cmd/kyber@main
```

`go install` places the binary in `$GOBIN` if set, otherwise in `$GOPATH/bin` (default: `$HOME/go/bin`). Make sure that directory is on your `PATH`:

```bash
# bash / zsh — append to ~/.bashrc or ~/.zshrc
export PATH="$(go env GOPATH)/bin:$PATH"
```

Verify the install:

```bash
$ kyber version
v0.1.0

$ kyber list-metrics
ID           NAME                   DEFAULT  DIRECTION        DESCRIPTION
cyclomatic   Cyclomatic Complexity  7        higher is worse  McCabe decision-point count.
readability  Readability Score      0.6      lower is worse   …
testability  Testability Score      0.6      lower is worse   …
```

When you install a tagged release this way, `kyber version` reports the tag automatically — kyber reads its module version from the embedded build info.

### Homebrew

```bash
brew install jedi-knights/tap/kyber
```

### Pre-built binaries

Download from the [releases page](https://github.com/jedi-knights/kyber/releases) — Linux, macOS, and Windows binaries for amd64 and arm64 (Windows arm64 is not built).

### Docker

```bash
docker run --rm -v "$PWD:/src" -w /src ghcr.io/jedi-knights/kyber:latest analyze ./...
```

## Quickstart

```bash
# Analyze the current module recursively
kyber analyze ./...

# Show every registered metric and its default threshold
kyber list-metrics

# Fail the build if any function exceeds its threshold
kyber analyze --fail-on-threshold ./...

# Emit SARIF for GitHub code scanning
kyber analyze --format=sarif -o kyber.sarif ./...

# Limit to specific metrics
kyber analyze --metric=cyclomatic --metric=readability ./...
```

## Metrics

| ID | What it measures | Default threshold |
|---|---|---|
| `cyclomatic` | McCabe decision-point count — number of linearly independent paths through a function | `> 7` flags |
| `readability` | Weighted 0–1 score from function length, nesting depth, identifier length distribution, and comment density | `< 0.6` flags |
| `testability` | Weighted 0–1 score from parameter count, observed side effects, interface-vs-concrete parameters, and length | `< 0.6` flags |

All three are configurable per-project via `kyber.toml`. See [Configuration](#configuration).

### Adding a new metric

1. Implement `ports.Metric` in `internal/domain/metrics/<name>.go`
2. Register it in `internal/domain/metrics/all.go`
3. Add a test in `<name>_test.go`

That's it.

## Configuration

kyber reads `kyber.toml` from the current directory by default. Precedence (highest first):

1. CLI flags (e.g. `--threshold-cyclomatic=10`)
2. Environment variables (`KYBER_FORMAT`, `KYBER_FAIL_ON_THRESHOLD`, `KYBER_VERBOSE`, `KYBER_PATHS`)
3. `kyber.toml`
4. Built-in defaults

See `kyber.toml.example` in the repo for an annotated reference.

## GitHub Action

```yaml
- uses: jedi-knights/kyber@v1
  with:
    paths: "./..."
    format: sarif
    fail-on-threshold: "true"
```

When `format: sarif`, the action uploads the report to the GitHub code-scanning view automatically.

## Development

```bash
git clone https://github.com/jedi-knights/kyber.git
cd kyber
go mod download
go test -race -count=1 ./...
git config core.hooksPath .githooks    # enable pre-push linting
```

Useful Makefile targets: `make build`, `make test`, `make lint`, `make install`, `make run`, `make clean`.

## Contributing

Bug reports and PRs are welcome. Run `make test` and `make lint` locally before pushing. Conventional Commits (`feat`, `fix`, `refactor`, ...) drive the release pipeline — non-conformant commits are rejected by CI.

## License

[MIT](LICENSE) © Jedi Knights
