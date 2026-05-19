// Package pkgctx is a kyber test fixture exercising package-level context:
// it declares an interface that another file in the same package consumes,
// plus a global variable. The parser must surface both via Package.Interfaces
// and Package.Globals.
package pkgctx

// Sender is a package-local interface; testability metric should resolve
// "Send" parameters against this.
type Sender interface {
	Send(msg string) error
}

// PackageCounter is a package-level global; the testability metric uses
// Globals to detect side-effects in functions that read/write it.
var PackageCounter int
