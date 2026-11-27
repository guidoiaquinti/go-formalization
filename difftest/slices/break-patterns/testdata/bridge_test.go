package slices

func GoLeanBreakPatternsOrdered(data []int, a, b int) {
	breakPatternsOrdered(data, a, b)
}

func GoLeanBreakPatternsCmpFunc(data []int, a, b int, compare func(int, int) int) {
	breakPatternsCmpFunc(data, a, b, compare)
}
