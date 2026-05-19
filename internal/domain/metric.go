package domain

import "context"

// Metric is the single extension point of kyber. Implement this interface,
// register it in DefaultRegistry (internal/domain/metrics/all.go), and the
// new metric will appear in every report.
//
// Implementations must be stateless: the same instance is invoked across many
// functions, potentially concurrently. Confine state to local variables
// inside Analyze.
type Metric interface {
	// ID is a stable machine-readable identifier (e.g. "cyclomatic").
	// Used as the config key, JSON field name, and SARIF ruleId.
	ID() string

	// Name is the human-readable label shown by `kyber list-metrics` and the
	// text reporter.
	Name() string

	// Description is one or two sentences explaining what the metric measures
	// and why it matters.
	Description() string

	// DefaultThreshold is the value above (or below — see HigherIsWorse)
	// which Analyze emits a Finding. Overridable via TOML / CLI.
	DefaultThreshold() float64

	// HigherIsWorse reports the direction of "bad". True for counters like
	// cyclomatic complexity (more = worse); false for normalized scores like
	// readability and testability (less = worse).
	HigherIsWorse() bool

	// Analyze inspects fn and returns a Score plus any Findings. ctx is
	// honored for cancellation; metrics that walk large ASTs should check
	// ctx.Err() periodically.
	Analyze(ctx context.Context, fn *Function, opts MetricOptions) (Score, error)
}

// MetricOptions carries per-metric configuration resolved from TOML and CLI.
// Threshold is the effective threshold (defaulted to the metric's
// DefaultThreshold when the user did not override). Params holds metric-
// specific knobs like max_params or weight_length.
type MetricOptions struct {
	Threshold float64
	Params    map[string]any
}

// FloatParam returns the param value as a float64, or def if missing or of
// the wrong type. Useful for the many numeric knobs metrics expose.
func (o MetricOptions) FloatParam(key string, def float64) float64 {
	v, ok := o.Params[key]
	if !ok {
		return def
	}
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return def
	}
}

// IntParam returns the param value as an int, or def if missing/wrong type.
func (o MetricOptions) IntParam(key string, def int) int {
	v, ok := o.Params[key]
	if !ok {
		return def
	}
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return def
	}
}
