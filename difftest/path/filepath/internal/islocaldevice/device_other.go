//go:build !windows

package islocaldevice

// IsDOSName is unused by the native Unix and Plan 9 profiles.
func IsDOSName(string) bool { return false }
