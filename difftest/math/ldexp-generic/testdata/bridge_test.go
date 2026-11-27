package math

// GoLeanLdexpGeneric exposes the pinned generic helper to the external test fixture.
func GoLeanLdexpGeneric(x float64, exp int) float64 { return ldexp(x, exp) }
