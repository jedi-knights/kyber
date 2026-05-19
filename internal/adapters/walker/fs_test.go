package walker

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/jedi-knights/kyber/internal/ports"
)

// makeTree builds a temp tree with the given relative paths (and a "package x"
// stub for any .go file) so the walker can be tested without filesystem
// pollution.
func makeTree(t *testing.T, paths []string) string {
	t.Helper()
	root := t.TempDir()
	for _, p := range paths {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte("package x\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", full, err)
		}
	}
	return root
}

func relAll(t *testing.T, base string, paths []string) []string {
	t.Helper()
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		r, err := filepath.Rel(base, p)
		if err != nil {
			t.Fatalf("rel: %v", err)
		}
		out = append(out, filepath.ToSlash(r))
	}
	sort.Strings(out)
	return out
}

func TestWalker_SkipsDefaultExcludes(t *testing.T) {
	root := makeTree(t, []string{
		"main.go",
		"vendor/dep/lib.go",
		"testdata/fixture.go",
		"pkg/util.go",
		"pkg/util_test.go",
	})

	files, err := New().Walk(context.Background(), []string{root + "/..."}, ports.WalkOptions{})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	got := relAll(t, root, files)
	want := []string{"main.go", "pkg/util.go"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWalker_IncludeTests(t *testing.T) {
	root := makeTree(t, []string{
		"main.go",
		"main_test.go",
	})
	files, err := New().Walk(context.Background(), []string{root + "/..."}, ports.WalkOptions{IncludeTests: true})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	got := relAll(t, root, files)
	want := []string{"main.go", "main_test.go"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWalker_CustomExcludes(t *testing.T) {
	root := makeTree(t, []string{
		"main.go",
		"internal/secret/lib.go",
		"pkg/util.go",
	})
	files, err := New().Walk(context.Background(), []string{root + "/..."}, ports.WalkOptions{
		ExcludeGlobs: []string{"internal/secret/**"},
	})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	got := relAll(t, root, files)
	want := []string{"main.go", "pkg/util.go"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWalker_RecursiveDotsSyntax(t *testing.T) {
	root := makeTree(t, []string{
		"a.go",
		"pkg/b.go",
	})
	// Use root/... explicitly to exercise the dots-suffix branch.
	files, err := New().Walk(context.Background(), []string{root + "/..."}, ports.WalkOptions{})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("got %d files, want 2: %v", len(files), files)
	}
}

func TestWalker_SinglePath(t *testing.T) {
	root := makeTree(t, []string{"main.go", "other.go"})
	target := filepath.Join(root, "main.go")
	files, err := New().Walk(context.Background(), []string{target}, ports.WalkOptions{})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(files) != 1 || files[0] != target {
		t.Errorf("got %v, want [%s]", files, target)
	}
}

func TestWalker_Deduplicates(t *testing.T) {
	root := makeTree(t, []string{"a.go", "pkg/b.go"})
	files, err := New().Walk(context.Background(), []string{root + "/...", root + "/..."}, ports.WalkOptions{})
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 unique files, got %d (%v)", len(files), files)
	}
}
