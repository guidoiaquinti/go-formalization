package math

// GoLeanIsOddInt exposes the pinned helper to the external test fixture.
func GoLeanIsOddInt(word uint64) bool { return isOddInt(Float64frombits(word)) }
