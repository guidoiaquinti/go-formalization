package math

// GoLeanIsNegInt exposes the pinned helper to the external test fixture.
func GoLeanIsNegInt(word uint64) bool { return isNegInt(Float64frombits(word)) }
