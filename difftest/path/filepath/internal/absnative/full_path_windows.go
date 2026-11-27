package absnative

import "syscall"

func FullPath(path string) (string, error) { return syscall.FullPath(path) }
