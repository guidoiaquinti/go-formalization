package slices

func GoLeanRotateLeft[S ~[]E, E any](s S, r int) {
	rotateLeft(s, r)
}
