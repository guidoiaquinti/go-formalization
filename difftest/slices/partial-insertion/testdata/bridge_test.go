package slices

func GoLeanPartialInsertionOrdered(data []int64, a, b int) bool {
	return partialInsertionSortOrdered(data, a, b)
}
func GoLeanPartialInsertionCmpFunc(data []int64, a, b int, cmp func(int64, int64) int) bool {
	return partialInsertionSortCmpFunc(data, a, b, cmp)
}
