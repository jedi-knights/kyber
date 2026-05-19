package domain

import (
	"fmt"
	"sort"
	"sync"
)

// Registry holds the metrics available to a kyber run. Metrics are explicitly
// registered (no init() side effects) so tests can construct registries
// containing only the metrics they exercise.
//
// Registry is safe for concurrent reads after construction.
type Registry struct {
	mu      sync.RWMutex
	metrics map[string]Metric
}

// NewRegistry constructs an empty registry.
func NewRegistry() *Registry {
	return &Registry{metrics: make(map[string]Metric)}
}

// Register adds m to the registry. Duplicate IDs return an error so wiring
// bugs surface at startup, not silently.
func (r *Registry) Register(m Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.metrics[m.ID()]; dup {
		return fmt.Errorf("kyber: metric %q already registered", m.ID())
	}
	r.metrics[m.ID()] = m
	return nil
}

// MustRegister panics if Register returns an error. Useful in package-level
// helpers (e.g. metrics.DefaultRegistry) where a duplicate ID is a programmer
// error that should surface at startup.
func (r *Registry) MustRegister(m Metric) {
	if err := r.Register(m); err != nil {
		panic(err)
	}
}

// Get returns the metric for an ID, or nil if no such metric is registered.
func (r *Registry) Get(id string) Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.metrics[id]
}

// All returns every registered metric, sorted by ID for deterministic output.
func (r *Registry) All() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Metric, 0, len(r.metrics))
	for _, m := range r.metrics {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Enabled returns the subset of registered metrics whose IDs appear in ids.
// IDs that are not registered are silently skipped (callers can detect this
// by comparing len(ids) to len(result)). When ids is empty, Enabled returns
// All() — "no config means everything on."
func (r *Registry) Enabled(ids []string) []Metric {
	if len(ids) == 0 {
		return r.All()
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Metric, 0, len(ids))
	for _, id := range ids {
		if m, ok := r.metrics[id]; ok {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}
