package math

// GoLeanShr exposes the pinned helper to the external test fixture.
func GoLeanShr(high, low uint64, shift uint) (uint64, uint64) { return shr(high, low, shift) }
