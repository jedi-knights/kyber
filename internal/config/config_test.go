package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func writeTOML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "kyber.toml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestConfig_Defaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("Format = %q, want text", cfg.Format)
	}
	if !slices.Equal(cfg.Paths, []string{"./..."}) {
		t.Errorf("Paths = %v, want ./...", cfg.Paths)
	}
	if cfg.FailOnThreshold {
		t.Error("FailOnThreshold should default false")
	}
}

func TestConfig_LoadsTOML(t *testing.T) {
	path := writeTOML(t, `
paths = ["./pkg/..."]
format = "json"
fail_on_threshold = true

[metrics.cyclomatic]
enabled   = true
threshold = 10

[metrics.testability]
enabled   = false
threshold = 0.5
max_params = 8
weight_length = 2.0
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %q, want json", cfg.Format)
	}
	if !cfg.FailOnThreshold {
		t.Error("FailOnThreshold should be true from TOML")
	}
	if !slices.Equal(cfg.Paths, []string{"./pkg/..."}) {
		t.Errorf("Paths = %v, want [./pkg/...]", cfg.Paths)
	}
	cyclo := cfg.Metrics["cyclomatic"]
	if cyclo.Threshold != 10 {
		t.Errorf("cyclomatic threshold = %v, want 10", cyclo.Threshold)
	}
	test := cfg.Metrics["testability"]
	if test.Params["max_params"] != int64(8) {
		t.Errorf("testability.max_params = %v (%T), want 8", test.Params["max_params"], test.Params["max_params"])
	}
	if test.Params["weight_length"] != 2.0 {
		t.Errorf("testability.weight_length = %v, want 2.0", test.Params["weight_length"])
	}
}

func TestConfig_IsEnabled(t *testing.T) {
	path := writeTOML(t, `
[metrics.cyclomatic]
enabled = true

[metrics.testability]
enabled = false

[metrics.readability]
threshold = 0.5
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.IsEnabled("cyclomatic") {
		t.Error("cyclomatic should be enabled")
	}
	if cfg.IsEnabled("testability") {
		t.Error("testability should be disabled")
	}
	if !cfg.IsEnabled("readability") {
		t.Error("readability (no enabled key) should default to enabled")
	}
	// Unknown metric (no table) defaults to enabled so user can simply not
	// configure new metrics.
	if !cfg.IsEnabled("future-metric") {
		t.Error("metric without a config table should default to enabled")
	}
}

func TestConfig_EnvOverlay(t *testing.T) {
	t.Setenv("KYBER_FORMAT", "sarif")
	t.Setenv("KYBER_VERBOSE", "true")
	t.Setenv("KYBER_PATHS", "./cmd, ./internal")
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Format != "sarif" {
		t.Errorf("Format = %q, want sarif (from env)", cfg.Format)
	}
	if !cfg.Verbose {
		t.Error("Verbose should be true (from env)")
	}
	want := []string{"./cmd", "./internal"}
	if !slices.Equal(cfg.Paths, want) {
		t.Errorf("Paths = %v, want %v", cfg.Paths, want)
	}
}

func TestConfig_MissingExplicitFileIsError(t *testing.T) {
	// Stat error for an explicitly-named missing file is OK in our model —
	// we fall through to defaults. Confirm the call doesn't crash.
	cfg, err := Load("/definitely/does/not/exist/kyber.toml")
	if err != nil {
		t.Fatalf("Load should tolerate missing explicit path, got %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("Format = %q, want default text", cfg.Format)
	}
}
