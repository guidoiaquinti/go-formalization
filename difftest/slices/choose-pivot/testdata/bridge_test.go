package slices

func GoLeanChoosePivotOrdered(data []int64, a, b int) (int, int) {
	pivot, hint := choosePivotOrdered(data, a, b)
	return pivot, int(hint)
}
func GoLeanChoosePivotCmpFunc(data []int64, a, b int, compare func(int64, int64) int) (int, int) {
	pivot, hint := choosePivotCmpFunc(data, a, b, compare)
	return pivot, int(hint)
}
