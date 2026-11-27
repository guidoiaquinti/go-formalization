package math

// GoLeanShl exposes the pinned helper to the external test fixture.
func GoLeanShl(high, low uint64, shift uint) (uint64, uint64) { return shl(high, low, shift) }
