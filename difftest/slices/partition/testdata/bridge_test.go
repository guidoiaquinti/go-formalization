package slices

func GoLeanPartitionOrdered(data []int64, a, b, pivot int) (int, bool) {
	return partitionOrdered(data, a, b, pivot)
}
func GoLeanPartitionCmpFunc(data []int64, a, b, pivot int, cmp func(int64, int64) int) (int, bool) {
	return partitionCmpFunc(data, a, b, pivot, cmp)
}
func GoLeanPartitionEqualOrdered(data []int64, a, b, pivot int) int {
	return partitionEqualOrdered(data, a, b, pivot)
}
func GoLeanPartitionEqualCmpFunc(data []int64, a, b, pivot int, cmp func(int64, int64) int) int {
	return partitionEqualCmpFunc(data, a, b, pivot, cmp)
}
