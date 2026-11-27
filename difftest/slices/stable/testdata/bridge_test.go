package slices

func GoLeanSymMergeOrdered(data []int64, a, m, b int) { symMergeOrdered(data, a, m, b) }
func GoLeanSymMergeCmpFunc(data []int64, a, m, b int, compare func(int64, int64) int) {
	symMergeCmpFunc(data, a, m, b, compare)
}
func GoLeanStableOrdered(data []int64, n int) { stableOrdered(data, n) }
func GoLeanStableCmpFunc(data []int64, n int, compare func(int64, int64) int) {
	stableCmpFunc(data, n, compare)
}
