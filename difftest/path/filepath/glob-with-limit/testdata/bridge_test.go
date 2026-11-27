package filepath

func GoLeanGlobWithLimit(pattern string, depth int) ([]string, error) {
	return globWithLimit(pattern, depth)
}
