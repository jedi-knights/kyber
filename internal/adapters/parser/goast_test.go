package parser

import (
	"context"
	"path/filepath"
	"testing"
)

func TestParser_ParsesSimpleFunction(t *testing.T) {
	p := New()
	path := absFixture(t, "simple/simple.go")
	funcs, err := p.ParseFiles(context.Background(), []string{path})
	if err != nil {
		t.Fatalf("ParseFiles: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("got %d functions, want 1", len(funcs))
	}
	fn := funcs[0]
	if fn.Name != "Add" {
		t.Errorf("Name = %q, want Add", fn.Name)
	}
	if fn.Receiver != "" {
		t.Errorf("Receiver = %q, want empty", fn.Receiver)
	}
	if got := fn.ParamCount(); got != 2 {
		t.Errorf("ParamCount = %d, want 2", got)
	}
	if got := fn.LineCount(); got != 3 {
		t.Errorf("LineCount = %d, want 3", got)
	}
	if fn.Package == nil || fn.Package.Name != "simple" {
		t.Errorf("Package.Name = %v, want simple", fn.Package)
	}
}

func TestParser_PopulatesPackageContext(t *testing.T) {
	p := New()
	paths := []string{
		absFixture(t, "package_context/interfaces.go"),
		absFixture(t, "package_context/concrete.go"),
	}
	funcs, err := p.ParseFiles(context.Background(), paths)
	if err != nil {
		t.Fatalf("ParseFiles: %v", err)
	}

	var pkg = funcs[0].Package
	if pkg.Name != "pkgctx" {
		t.Errorf("Package.Name = %q, want pkgctx", pkg.Name)
	}
	if _, ok := pkg.Interfaces["Sender"]; !ok {
		t.Errorf("Package.Interfaces missing Sender; got %v", keys(pkg.Interfaces))
	}
	if _, ok := pkg.Globals["PackageCounter"]; !ok {
		t.Errorf("Package.Globals missing PackageCounter; got %v", keys(pkg.Globals))
	}
	// All four functions (Notifier.Send, SendViaInterface, SendViaConcrete,
	// BumpCounter) must share the same *Package.
	for _, fn := range funcs {
		if fn.Package != pkg {
			t.Errorf("function %s has a different *Package than its siblings", fn.Name)
		}
	}
}

func TestParser_MethodReceiverIncludedInName(t *testing.T) {
	p := New()
	paths := []string{absFixture(t, "package_context/concrete.go"), absFixture(t, "package_context/interfaces.go")}
	funcs, err := p.ParseFiles(context.Background(), paths)
	if err != nil {
		t.Fatalf("ParseFiles: %v", err)
	}
	var found bool
	for _, fn := range funcs {
		if fn.Name == "Notifier.Send" {
			found = true
			if fn.Receiver != "Notifier" {
				t.Errorf("Receiver = %q, want Notifier", fn.Receiver)
			}
		}
	}
	if !found {
		t.Errorf("did not find Notifier.Send among parsed functions")
	}
}

func absFixture(t *testing.T, rel string) string {
	t.Helper()
	// Tests run from the package dir; testdata lives at the module root.
	path, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", rel))
	if err != nil {
		t.Fatalf("abs(%s): %v", rel, err)
	}
	return path
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
