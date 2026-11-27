package slices

func GoLeanMedianOrdered(data []int64, a, b, c int, count *int) int {
	return medianOrdered(data, a, b, c, count)
}
func GoLeanMedianCmpFunc(data []int64, a, b, c int, count *int, compare func(int64, int64) int) int {
	return medianCmpFunc(data, a, b, c, count, compare)
}
func GoLeanMedianAdjacentOrdered(data []int64, a int, count *int) int {
	return medianAdjacentOrdered(data, a, count)
}
func GoLeanMedianAdjacentCmpFunc(data []int64, a int, count *int, compare func(int64, int64) int) int {
	return medianAdjacentCmpFunc(data, a, count, compare)
}
