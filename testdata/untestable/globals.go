// Package untestable is a kyber test fixture for the testability metric.
// The single function has many parameters, touches a package-level global,
// calls os.Exit (a recognized side-effect), and uses concrete types — all
// signals that should push testability low.
package untestable

import (
	"fmt"
	"os"
)

// counter is a package-level global; reading or writing it from a function
// is a side-effect signal.
var counter int

// Dispatch is intentionally hard to test: 7 parameters, side-effects on
// globals, and a call to os.Exit on the error branch.
func Dispatch(a, b, c, d, e, f, g int) int {
	counter++
	if a < 0 {
		fmt.Fprintln(os.Stderr, "negative a")
		os.Exit(1)
	}
	return a + b + c + d + e + f + g + counter
}
