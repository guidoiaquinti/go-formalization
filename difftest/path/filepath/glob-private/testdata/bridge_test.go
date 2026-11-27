package filepath

func GoLeanPrivateGlob(dir, pattern string, paths []string) ([]string, error) {
	return glob(dir, pattern, paths)
}
