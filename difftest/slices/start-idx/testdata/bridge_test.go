package slices

func GoLeanStartIdx[E any](haystack, needle []E) int {
	return startIdx(haystack, needle)
}
