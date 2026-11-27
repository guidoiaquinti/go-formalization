package slices

func GoLeanRotateOrdered(data []int, a, middle, b int) {
	rotateOrdered(data, a, middle, b)
}

func GoLeanRotateCmpFunc(data []int, a, middle, b int, compare func(int, int) int) {
	rotateCmpFunc(data, a, middle, b, compare)
}
