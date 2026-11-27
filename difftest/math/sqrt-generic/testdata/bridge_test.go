package math

// GoLeanSqrtGeneric bypasses any substitution of the public Sqrt entry point.
func GoLeanSqrtGeneric(word uint64) float64 { return sqrt(Float64frombits(word)) }
