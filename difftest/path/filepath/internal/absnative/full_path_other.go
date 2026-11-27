//go:build !windows

package absnative

func FullPath(string) (string, error) { panic("FullPath is Windows-only") }
