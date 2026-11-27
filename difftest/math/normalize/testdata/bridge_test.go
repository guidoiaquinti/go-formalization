package math

// GoLeanNormalize exposes the pinned helper to the external test fixture.
func GoLeanNormalize(x float64) (float64, int) { return normalize(x) }

// GoLeanIlogb exposes the pinned finite-only exponent helper to the external test fixture.
func GoLeanIlogb(x float64) int { return ilogb(x) }
