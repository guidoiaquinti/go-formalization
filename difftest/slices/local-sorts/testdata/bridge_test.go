package slices

func GoLeanInsertionSortOrdered(data []int64, a, b int)     { insertionSortOrdered(data, a, b) }
func GoLeanSiftDownOrdered(data []int64, lo, hi, first int) { siftDownOrdered(data, lo, hi, first) }
func GoLeanHeapSortOrdered(data []int64, a, b int)          { heapSortOrdered(data, a, b) }

func GoLeanInsertionSortCmpFunc(data []int64, a, b int, cmp func(int64, int64) int) {
	insertionSortCmpFunc(data, a, b, cmp)
}
func GoLeanSiftDownCmpFunc(data []int64, lo, hi, first int, cmp func(int64, int64) int) {
	siftDownCmpFunc(data, lo, hi, first, cmp)
}
func GoLeanHeapSortCmpFunc(data []int64, a, b int, cmp func(int64, int64) int) {
	heapSortCmpFunc(data, a, b, cmp)
}
