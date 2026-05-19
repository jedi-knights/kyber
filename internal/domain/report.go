package domain

import "time"

// Report is the result of an analysis run — every Score from every metric
// against every function, plus run-time bookkeeping. Reporters read this
// structure to produce text, JSON, or SARIF output.
type Report struct {
	Scores       []Score
	StartTime    time.Time
	EndTime      time.Time
	FilesScanned int
	Errors       []error
}

// ByFunction groups scores by their function for per-function summaries.
// The iteration order of the result is not specified; callers that need a
// stable order should sort by Function.Position().
func (r *Report) ByFunction() map[*Function][]Score {
	out := make(map[*Function][]Score)
	for _, s := range r.Scores {
		out[s.Function] = append(out[s.Function], s)
	}
	return out
}

// ByMetric returns scores for one metric ID, in original order.
func (r *Report) ByMetric(id string) []Score {
	out := make([]Score, 0)
	for _, s := range r.Scores {
		if s.MetricID == id {
			out = append(out, s)
		}
	}
	return out
}

// ExceedingThresholds returns the subset of scores whose Findings list is
// non-empty (the metric already evaluated the breach).
func (r *Report) ExceedingThresholds() []Score {
	out := make([]Score, 0)
	for _, s := range r.Scores {
		if len(s.Findings) > 0 {
			out = append(out, s)
		}
	}
	return out
}
