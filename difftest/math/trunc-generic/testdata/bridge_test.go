package math

// GoLeanTruncGeneric calls the pinned pure-Go helper directly.
func GoLeanTruncGeneric(word uint64) uint64 {
	return Float64bits(trunc(Float64frombits(word)))
}
