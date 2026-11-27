package filepath

func GoLeanPrivateEvalSymlinks(path string) (string, error) { return evalSymlinks(path) }
func GoLeanWalkSymlinks(path string) (string, error)        { return walkSymlinks(path) }
