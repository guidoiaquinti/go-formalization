package math

// GoLeanFmaSplit exposes the pinned helper to the external test fixture.
func GoLeanFmaSplit(word uint64) (uint32, int32, uint64) { return split(word) }
