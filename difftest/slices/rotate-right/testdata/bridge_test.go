package slices

func GoLeanRotateRight[S ~[]E, E any](s S, r int) {
	rotateRight(s, r)
}
