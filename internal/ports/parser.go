// Package ports defines the interfaces by which the domain interacts with the
// outside world: parsing Go source, walking the filesystem, emitting reports.
// Concrete implementations live in internal/adapters/.
package ports

import (
	"context"

	"github.com/jedi-knights/kyber/internal/domain"
)

// SourceParser parses Go source files into the domain's Function model. It is
// the only seam through which raw source enters kyber; everything downstream
// (metrics, reporters) operates on domain.Function values.
type SourceParser interface {
	// ParseFiles parses the given .go file paths and returns every top-level
	// function and method found, with full Package context populated. Files
	// belonging to the same package directory are grouped into a single
	// Package; functions in that group share that Package value.
	ParseFiles(ctx context.Context, paths []string) ([]*domain.Function, error)
}
