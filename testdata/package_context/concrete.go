package pkgctx

// Notifier is a concrete struct in the same package as Sender.
type Notifier struct{ name string }

// Send satisfies the Sender interface.
func (n *Notifier) Send(msg string) error { return nil }

// SendViaInterface accepts the Sender interface (more testable).
func SendViaInterface(s Sender, msg string) error {
	return s.Send(msg)
}

// SendViaConcrete accepts the concrete *Notifier (less testable — cannot
// mock without an interface).
func SendViaConcrete(n *Notifier, msg string) error {
	return n.Send(msg)
}

// BumpCounter touches a package-level global (side-effect signal).
func BumpCounter() int {
	PackageCounter++
	return PackageCounter
}
