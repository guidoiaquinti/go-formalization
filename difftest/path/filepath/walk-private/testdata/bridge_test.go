package filepath

import "io/fs"

func GoLeanPrivateWalk(path string, info fs.FileInfo, visit WalkFunc) error {
	return walk(path, info, visit)
}
