package slices

func GoLeanOrder2Int(data []int64, a, b int, swaps *int) (int, int) {
	return order2Ordered(data, a, b, swaps)
}

func GoLeanOrder2Float(data []float64, a, b int, swaps *int) (int, int) {
	return order2Ordered(data, a, b, swaps)
}

func GoLeanOrder2Cmp(data []uint64, a, b int, swaps *int, cmp func(uint64, uint64) int) (int, int) {
	return order2CmpFunc(data, a, b, swaps, cmp)
}
