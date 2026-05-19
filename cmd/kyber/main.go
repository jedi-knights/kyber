// Command kyber is a Go function-level code-quality analyzer.
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jedi-knights/kyber/cmd/kyber/commands"
)

// version is the build-time-injected release tag; see Makefile and
// .goreleaser.yml. Falls back to "dev" for unreleased builds.
// When users run `go install` against a tagged version, the module proxy
// records that version in debug.BuildInfo, and resolveVersion() picks it up
// — so binaries built without ldflags still report a meaningful version.
var version = "dev"

func resolveVersion() string {
	if version != "dev" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	return version
}

func main() {
	root := commands.NewRoot(resolveVersion())
	if err := root.Execute(); err != nil {
		// Cobra already prints user-facing errors; this branch covers cases
		// where Execute itself failed (e.g. flag parsing) before its handler
		// could emit an error message.
		fmt.Fprintln(os.Stderr, "kyber:", err)
		os.Exit(1)
	}
}
