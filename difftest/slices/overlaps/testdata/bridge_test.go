package slices

func GoLeanOverlaps[E any](a, b []E) bool {
	return overlaps(a, b)
}
