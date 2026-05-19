package reporter

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
	"time"

	"github.com/jedi-knights/kyber/internal/domain"
)

// makeReport builds a small, deterministic Report for golden-style assertions
// without going through the full walker+parser pipeline.
func makeReport(t *testing.T) *domain.Report {
	t.Helper()
	fset := token.NewFileSet()
	const src = `package x
func Foo() {}
`
	file, err := parser.ParseFile(fset, "x/x.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	fn := file.Decls[0].(*ast.FuncDecl)
	pkg := &domain.Package{Name: "x", FileSet: fset}
	function := &domain.Function{
		Name:     "Foo",
		Package:  pkg,
		File:     "x/x.go",
		FuncDecl: fn,
		FileSet:  fset,
	}
	r := &domain.Report{
		StartTime:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC),
		FilesScanned: 1,
		Scores: []domain.Score{
			{MetricID: "cyclomatic", Function: function, Value: 12,
				Findings: []domain.Finding{{
					Severity: domain.SeverityWarning, Line: 2,
					Message: "cyclomatic complexity 12 exceeds threshold 7",
				}}},
			{MetricID: "readability", Function: function, Value: 0.42,
				Findings: []domain.Finding{{
					Severity: domain.SeverityWarning, Line: 2,
					Message: "readability 0.42 below threshold 0.60",
				}}},
		},
	}
	return r
}

func TestText_RendersFileAndFunction(t *testing.T) {
	var buf bytes.Buffer
	if err := NewText().Render(&buf, makeReport(t)); err != nil {
		t.Fatalf("Render: %v", err)
	}
	s := buf.String()
	if !strings.Contains(s, "x/x.go") {
		t.Errorf("output missing file header: %q", s)
	}
	if !strings.Contains(s, "Foo") {
		t.Errorf("output missing function name: %q", s)
	}
	if !strings.Contains(s, "cyclomatic=12") {
		t.Errorf("output missing cyclomatic value: %q", s)
	}
	if !strings.Contains(s, "readability=0.42") {
		t.Errorf("output missing readability value: %q", s)
	}
	if !strings.Contains(s, "Functions: 1") {
		t.Errorf("output missing summary line: %q", s)
	}
	// Findings increment the count: 2 findings expected from the fixture.
	if !strings.Contains(s, "Findings: 2") {
		t.Errorf("expected 'Findings: 2' in summary; got %q", s)
	}
}

func TestJSON_RoundtripsScores(t *testing.T) {
	var buf bytes.Buffer
	if err := NewJSON().Render(&buf, makeReport(t)); err != nil {
		t.Fatalf("Render: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	scores, ok := out["scores"].([]any)
	if !ok || len(scores) != 2 {
		t.Fatalf("scores should be a 2-element array, got %v", out["scores"])
	}
	first := scores[0].(map[string]any)
	if first["metric_id"] != "cyclomatic" {
		t.Errorf("first score metric_id = %v, want cyclomatic", first["metric_id"])
	}
	if first["value"].(float64) != 12 {
		t.Errorf("first score value = %v, want 12", first["value"])
	}
}

func TestSARIF_ProducesValidStructure(t *testing.T) {
	var buf bytes.Buffer
	if err := NewSARIF("0.1.0").Render(&buf, makeReport(t)); err != nil {
		t.Fatalf("Render: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if doc["version"] != "2.1.0" {
		t.Errorf("version = %v, want 2.1.0", doc["version"])
	}
	runs := doc["runs"].([]any)
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	run := runs[0].(map[string]any)
	results := run["results"].([]any)
	if len(results) != 2 {
		t.Errorf("expected 2 results (one per finding), got %d", len(results))
	}
	tool := run["tool"].(map[string]any)
	driver := tool["driver"].(map[string]any)
	if driver["name"] != "kyber" {
		t.Errorf("tool name = %v, want kyber", driver["name"])
	}
	if driver["version"] != "0.1.0" {
		t.Errorf("tool version = %v, want 0.1.0", driver["version"])
	}
}

func TestSARIF_EmptyReportStillValid(t *testing.T) {
	var buf bytes.Buffer
	report := &domain.Report{
		StartTime: time.Now(),
		EndTime:   time.Now(),
	}
	if err := NewSARIF("dev").Render(&buf, report); err != nil {
		t.Fatalf("Render: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if doc["version"] != "2.1.0" {
		t.Errorf("version = %v, want 2.1.0", doc["version"])
	}
}
