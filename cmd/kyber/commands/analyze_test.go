package commands

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func runRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := NewRoot("test")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func testdataRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "testdata"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	return root
}

func TestListMetricsCommand(t *testing.T) {
	out, err := runRoot(t, "list-metrics")
	if err != nil {
		t.Fatalf("list-metrics: %v", err)
	}
	for _, want := range []string{"cyclomatic", "readability", "testability"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to mention %q; got %q", want, out)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	out, err := runRoot(t, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	if strings.TrimSpace(out) != "test" {
		t.Errorf("version output = %q, want test", out)
	}
}

func TestAnalyzeCommand_TextOutput(t *testing.T) {
	out, err := runRoot(t, "analyze", "--config", "/dev/null", testdataRoot(t)+"/...")
	if err != nil {
		t.Fatalf("analyze: %v\nout: %s", err, out)
	}
	if !strings.Contains(out, "Functions:") {
		t.Errorf("expected text summary; got %q", out)
	}
}

func TestAnalyzeCommand_JSONOutput(t *testing.T) {
	out, err := runRoot(t, "analyze", "--config", "/dev/null", "--format", "json", testdataRoot(t)+"/...")
	if err != nil {
		t.Fatalf("analyze: %v\nout: %s", err, out)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if _, ok := doc["scores"]; !ok {
		t.Errorf("JSON output missing 'scores' field: %v", doc)
	}
}

func TestAnalyzeCommand_FailOnThreshold(t *testing.T) {
	// The complex fixture's Branchy has cyclomatic 12 > default threshold 7,
	// so --fail-on-threshold must surface a non-nil error.
	_, err := runRoot(t,
		"analyze",
		"--config", "/dev/null",
		"--fail-on-threshold",
		testdataRoot(t)+"/...",
	)
	if err == nil {
		t.Error("expected non-nil error from --fail-on-threshold (threshold breaches present)")
	}
}

func TestAnalyzeCommand_MetricFilter(t *testing.T) {
	out, err := runRoot(t,
		"analyze",
		"--config", "/dev/null",
		"--format", "json",
		"--metric", "cyclomatic",
		testdataRoot(t)+"/...",
	)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, s := range doc["scores"].([]any) {
		sm := s.(map[string]any)
		if sm["metric_id"] != "cyclomatic" {
			t.Errorf("unexpected metric in filtered output: %v", sm["metric_id"])
		}
	}
}
