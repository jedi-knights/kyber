package reporter

import (
	"encoding/json"
	"io"

	"github.com/jedi-knights/kyber/internal/domain"
)

// JSON renders a report as a structured JSON document. The shape is stable
// and intended for downstream tooling (CI dashboards, code-review bots).
type JSON struct{}

// NewJSON constructs a JSON reporter.
func NewJSON() *JSON { return &JSON{} }

type jsonReport struct {
	Scores       []jsonScore `json:"scores"`
	StartTime    string      `json:"start_time"`
	EndTime      string      `json:"end_time"`
	FilesScanned int         `json:"files_scanned"`
	Errors       []string    `json:"errors,omitempty"`
}

type jsonScore struct {
	MetricID string        `json:"metric_id"`
	Function jsonFunction  `json:"function"`
	Value    float64       `json:"value"`
	Findings []jsonFinding `json:"findings,omitempty"`
}

type jsonFunction struct {
	Name     string `json:"name"`
	Receiver string `json:"receiver,omitempty"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type jsonFinding struct {
	Severity string `json:"severity"`
	Line     int    `json:"line"`
	Column   int    `json:"column,omitempty"`
	Message  string `json:"message"`
}

// Render writes the report to w as indented JSON.
func (JSON) Render(w io.Writer, r *domain.Report) error {
	out := jsonReport{
		StartTime:    r.StartTime.UTC().Format("2006-01-02T15:04:05Z"),
		EndTime:      r.EndTime.UTC().Format("2006-01-02T15:04:05Z"),
		FilesScanned: r.FilesScanned,
	}
	for _, s := range r.Scores {
		js := jsonScore{
			MetricID: s.MetricID,
			Function: jsonFunction{
				Name:     s.Function.Name,
				Receiver: s.Function.Receiver,
				File:     s.Function.File,
				Line:     s.Function.Position().Line,
			},
			Value: s.Value,
		}
		for _, f := range s.Findings {
			js.Findings = append(js.Findings, jsonFinding{
				Severity: f.Severity.String(),
				Line:     f.Line,
				Column:   f.Column,
				Message:  f.Message,
			})
		}
		out.Scores = append(out.Scores, js)
	}
	for _, e := range r.Errors {
		out.Errors = append(out.Errors, e.Error())
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
