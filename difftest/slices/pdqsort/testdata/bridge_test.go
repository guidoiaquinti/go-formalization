package slices

func GoLeanPdqsortOrdered(data []int64, a, b, limit int) {
	pdqsortOrdered(data, a, b, limit)
}
func GoLeanPdqsortCmpFunc(data []int64, a, b, limit int, cmp func(int64, int64) int) {
	pdqsortCmpFunc(data, a, b, limit, cmp)
}
