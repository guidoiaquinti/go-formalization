package math

// GoLeanFrexpGeneric exposes the pinned generic helper to the external test fixture.
func GoLeanFrexpGeneric(x float64) (float64, int) { return frexp(x) }
