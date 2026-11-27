package main

import (
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", "/", "\\", "#d/one", "#", "a#d", "a", "C:", "C:rel", "C:/abs",
		"//host/share", "//host/../x", "//?/UNC/host/share", "\\??\\C:\\one", "\\??\\UNC\\host\\share",
		"\\??\\UNC\\host\\share\\one", "//./", "//.", "/??", "\xff:/\x00"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b), ':', '/', 'x'}), string([]byte{byte(b), '#', 0xff}))
	}
	check.CompareSamples("godiff-path-filepath-is-abs", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(80))
	r.Read(raw)
	for j := 0; j < len(raw); j += 4 {
		raw[j] = "/\\:?#."[r.Intn(6)]
	}
	prefixes := []string{"", "/", "A:", "a:/", "#", "\\", "//", "//?/UNC/", "\\??\\", "//./"}
	return prefixes[i%len(prefixes)] + string(raw)
}

func encode(path string) string {
	row := []byte("0")
	for i := range len(path) {
		row = strconv.AppendUint(append(row, ' '), uint64(path[i]), 10)
	}
	return string(row)
}

func native(path string) string {
	got := filepath.IsAbs(path)
	profiles := oracle(path)
	switch runtime.GOOS {
	case "windows":
		profiles[1] = got
	case "plan9":
		profiles[2] = got
	default:
		profiles[0] = got
	}
	return wire(profiles)
}

func expected(path string) string {
	return wire(oracle(path))
}

func oracle(path string) [3]bool {
	unix := strings.HasPrefix(path, "/")
	return [3]bool{unix, windowsAbs(path), unix || strings.HasPrefix(path, "#")}
}

// Classify normalized component lists rather than finding a byte volume offset.
func windowsAbs(path string) bool {
	normalized := strings.ReplaceAll(path, "/", "\\")
	if len(normalized) >= 2 && normalized[1] == ':' {
		return len(normalized) >= 3 && normalized[2] == '\\'
	}
	if !strings.HasPrefix(normalized, "\\") {
		return false
	}
	parts := strings.Split(normalized, "\\")
	doubled := strings.HasPrefix(normalized, "\\\\")
	start, count := 2, 2
	switch {
	case doubled && (parts[2] == "?" || parts[2] == "."):
		start, count = 3, 1
	case !doubled && parts[1] == "??":
		start, count = 2, 1
	case !doubled:
		return false
	}
	if count == 1 && len(parts) > start && strings.ToUpper(parts[start]) == "UNC" {
		count = 3
	}
	for _, part := range parts[start:min(len(parts), start+count)] {
		if part == ".." {
			return false
		}
	}
	return doubled || len(parts) > start+count
}

func wire(out [3]bool) string {
	return strconv.FormatBool(out[0]) + "|" + strconv.FormatBool(out[1]) + "|" + strconv.FormatBool(out[2])
}
