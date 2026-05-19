package ports

import "context"

// FileWalker discovers .go source files under one or more root paths,
// applying include/exclude filters. It is the only seam through which kyber
// touches the filesystem to discover sources.
type FileWalker interface {
	Walk(ctx context.Context, roots []string, opts WalkOptions) ([]string, error)
}

// WalkOptions configures a walk. ExcludeGlobs are matched against each file's
// path relative to the walk root using path/filepath.Match-style globs (with
// "**" expanded as "any number of path components"). Test files (*_test.go)
// are excluded unless IncludeTests is true.
type WalkOptions struct {
	ExcludeGlobs []string
	IncludeTests bool
}

// DefaultExcludes returns the conventional exclude list applied when the
// caller does not specify any.
func DefaultExcludes() []string {
	return []string{"vendor/**", "testdata/**"}
}
