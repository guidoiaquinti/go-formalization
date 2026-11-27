package slices

func GoLeanSwapRangeOrdered(data []int, a, b, count int) {
	swapRangeOrdered(data, a, b, count)
}

func GoLeanSwapRangeCmpFunc(data []int, a, b, count int, compare func(int, int) int) {
	swapRangeCmpFunc(data, a, b, count, compare)
}
