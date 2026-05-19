package metrics

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jedi-knights/kyber/internal/adapters/parser"
	"github.com/jedi-knights/kyber/internal/domain"
)

// parseFixture is a test helper that returns every function in the named
// testdata fixture directory.
func parseFixture(t *testing.T, fixture string) []*domain.Function {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata", fixture))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	funcs, err := parser.New().ParseFiles(context.Background(), matches)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return funcs
}

func findFunc(t *testing.T, funcs []*domain.Function, name string) *domain.Function {
	t.Helper()
	for _, fn := range funcs {
		if fn.Name == name {
			return fn
		}
	}
	t.Fatalf("function %q not found", name)
	return nil
}

func TestCyclomatic_Simple(t *testing.T) {
	funcs := parseFixture(t, "simple")
	fn := findFunc(t, funcs, "Add")
	score, err := NewCyclomatic().Analyze(context.Background(), fn, domain.MetricOptions{})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if score.Value != 1 {
		t.Errorf("Value = %v, want 1", score.Value)
	}
	if len(score.Findings) != 0 {
		t.Errorf("Findings = %v, want none", score.Findings)
	}
}

func TestCyclomatic_Complex(t *testing.T) {
	funcs := parseFixture(t, "complex")
	fn := findFunc(t, funcs, "Branchy")
	score, err := NewCyclomatic().Analyze(context.Background(), fn, domain.MetricOptions{})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	// Hand-computed in the fixture's package comment.
	if score.Value != 12 {
		t.Errorf("Value = %v, want 12", score.Value)
	}
	if len(score.Findings) != 1 {
		t.Fatalf("Findings = %v, want 1", score.Findings)
	}
	if score.Findings[0].Severity != domain.SeverityWarning {
		t.Errorf("Severity = %v, want Warning (12 > 7 but < 2×7)", score.Findings[0].Severity)
	}
}

func TestCyclomatic_EscalatesToError(t *testing.T) {
	funcs := parseFixture(t, "complex")
	fn := findFunc(t, funcs, "Branchy")
	// Threshold of 5 means 12 ≥ 2×5 → Error severity.
	score, err := NewCyclomatic().Analyze(context.Background(), fn, domain.MetricOptions{Threshold: 5})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(score.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(score.Findings))
	}
	if score.Findings[0].Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", score.Findings[0].Severity)
	}
}

func TestCyclomatic_ContextCancellation(t *testing.T) {
	funcs := parseFixture(t, "simple")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := NewCyclomatic().Analyze(ctx, funcs[0], domain.MetricOptions{}); err == nil {
		t.Error("expected error from cancelled context, got nil")
	}
}
