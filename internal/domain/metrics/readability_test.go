package metrics

import (
	"context"
	"testing"

	"github.com/jedi-knights/kyber/internal/domain"
)

func TestReadability_SimpleScoresHigh(t *testing.T) {
	fn := findFunc(t, parseFixture(t, "simple"), "Add")
	score, err := NewReadability().Analyze(context.Background(), fn, domain.MetricOptions{})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	// Length signal = 1 (3 lines ≤ 40)
	// Nesting signal = 1 (depth 1 ≤ 4)
	// Ident signal: identifiers are a, b — median 1 → 0
	// Comment signal: 0 comment lines inside body → 0
	// Score = (1 + 1 + 0 + 0) / 4 = 0.5
	if score.Value < 0.4 || score.Value > 0.6 {
		t.Errorf("Value = %v, want ~0.5", score.Value)
	}
}

func TestReadability_UnreadableScoresLow(t *testing.T) {
	fn := findFunc(t, parseFixture(t, "unreadable"), "Tangled")
	score, err := NewReadability().Analyze(context.Background(), fn, domain.MetricOptions{})
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if score.Value >= 0.6 {
		t.Errorf("Value = %v, want < 0.6 (unreadable fixture)", score.Value)
	}
	if len(score.Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(score.Findings))
	}
}

func TestReadability_LengthSignalDominates(t *testing.T) {
	fn := findFunc(t, parseFixture(t, "unreadable"), "Tangled")
	// Give length signal all the weight; expect a very low score because the
	// function is > 40 lines long.
	opts := domain.MetricOptions{
		Threshold: 0.6,
		Params: map[string]any{
			"max_lines":       20.0,
			"weight_length":   1.0,
			"weight_nesting":  0.0,
			"weight_idents":   0.0,
			"weight_comments": 0.0,
		},
	}
	score, err := NewReadability().Analyze(context.Background(), fn, opts)
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if score.Value != 0 {
		t.Errorf("Value = %v, want 0 (function exceeds max_lines=20 by more than 2x)", score.Value)
	}
}

func TestReadability_NestingDepth(t *testing.T) {
	// The unreadable fixture nests if-blocks five deep.
	funcs := parseFixture(t, "unreadable")
	fn := findFunc(t, funcs, "Tangled")
	got := computeNestingDepth(fn.FuncDecl)
	if got < 5 {
		t.Errorf("nesting depth = %d, want at least 5", got)
	}
}
