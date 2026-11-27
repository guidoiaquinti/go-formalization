package math

// GoLeanShrCompress exposes the pinned helper to the external test fixture.
func GoLeanShrCompress(high, low uint64, shift uint) (uint64, uint64) {
	return shrcompress(high, low, shift)
}
