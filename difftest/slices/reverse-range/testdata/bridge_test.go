package slices

func GoLeanReverseRangeOrdered(data []int, a, b int) {
	reverseRangeOrdered(data, a, b)
}

func GoLeanReverseRangeCmpFunc(data []int, a, b int, compare func(int, int) int) {
	reverseRangeCmpFunc(data, a, b, compare)
}
