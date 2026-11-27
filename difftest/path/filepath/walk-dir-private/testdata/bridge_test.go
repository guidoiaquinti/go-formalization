package filepath

import "io/fs"

func GoLeanPrivateWalkDir(path string, entry fs.DirEntry, callback fs.WalkDirFunc) error {
	return walkDir(path, entry, callback)
}
