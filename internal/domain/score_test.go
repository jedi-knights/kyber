package domain

import "testing"

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityInfo, "info"},
		{SeverityWarning, "warning"},
		{SeverityError, "error"},
		{Severity(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", int(tt.sev), got, tt.want)
		}
	}
}

func TestScore_ExceedsThreshold(t *testing.T) {
	tests := []struct {
		name          string
		value         float64
		threshold     float64
		higherIsWorse bool
		want          bool
	}{
		{"complexity 8 over threshold 7 (higher worse)", 8, 7, true, true},
		{"complexity 7 at threshold 7 (higher worse, equal is ok)", 7, 7, true, false},
		{"complexity 5 under threshold 7 (higher worse)", 5, 7, true, false},
		{"readability 0.5 under threshold 0.6 (lower worse)", 0.5, 0.6, false, true},
		{"readability 0.6 at threshold 0.6 (lower worse, equal is ok)", 0.6, 0.6, false, false},
		{"readability 0.8 above threshold 0.6 (lower worse)", 0.8, 0.6, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Score{Value: tt.value}
			if got := s.ExceedsThreshold(tt.threshold, tt.higherIsWorse); got != tt.want {
				t.Errorf("ExceedsThreshold(%v, higherIsWorse=%v) = %v, want %v",
					tt.threshold, tt.higherIsWorse, got, tt.want)
			}
		})
	}
}
