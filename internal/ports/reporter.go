package ports

import (
	"io"

	"github.com/jedi-knights/kyber/internal/domain"
)

// Reporter renders a Report to a writer. Implementations exist for text,
// JSON, and SARIF v2.1.0.
type Reporter interface {
	Render(w io.Writer, r *domain.Report) error
}
