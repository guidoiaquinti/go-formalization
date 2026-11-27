package main

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", "a", ":", "C:", "C:rel", "a:/abs", "//host/share/dir", "//host",
		"//host/../dir", "//../share/dir", "//./dev/path", "//?/UNC/host/share/path",
		"//?/unc", "//?/unc/", "//?/UNCx/host/share", "/??/C:/file", "/??",
		"//./", "///", "//", "/root", "//?/../x"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b), ':', '/', 'x'}),
			string([]byte{'/', '/', byte(b), '/', 's', '/', 't'}))
	}
	check.CompareSamples("godiff-path-filepath-volume-name", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(96))
	r.Read(raw)
	for j := range raw {
		if j%5 == 0 {
			raw[j] = "/\\:."[r.Intn(4)]
		}
	}
	prefixes := []string{"", "A:", "/", "\\\\", "//./", "\\\\?\\", "/??/", "//?/unc/"}
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
	got := filepath.VolumeName(path)
	switch filepath.Separator {
	case '/':
		return wire(got) + "|" + wire(windowsVolume(path))
	case '\\':
		return wire("") + "|" + wire(got)
	default:
		panic("unsupported path separator")
	}
}

func expected(path string) string {
	return wire("") + "|" + wire(windowsVolume(path))
}

// Normalize separators first, then select whole components in the normalized text.
func windowsVolume(input string) string {
	path := strings.ReplaceAll(input, "/", "\\")
	if len(path) >= 2 && path[1] == ':' {
		return path[:2]
	}
	if !strings.HasPrefix(path, "\\") {
		return ""
	}
	componentPrefix := func(s, prefix string) bool {
		return s == prefix || strings.HasPrefix(s, prefix+"\\")
	}
	device := componentPrefix(path, "\\\\.") || componentPrefix(path, "\\\\?") ||
		componentPrefix(path, "\\??")
	volume := ""
	switch {
	case device && len(path) == 3:
		return path
	case device:
		tail := path[4:]
		parts := strings.Split(tail, "\\")
		if strings.ToUpper(parts[0]) == "UNC" {
			volume = path[:4] + strings.Join(parts[:min(3, len(parts))], "\\")
		} else {
			volume = path[:4] + parts[0]
		}
	case strings.HasPrefix(path, "\\\\"):
		parts := strings.Split(path[2:], "\\")
		volume = path[:2] + strings.Join(parts[:min(2, len(parts))], "\\")
	}
	for _, part := range strings.Split(volume, "\\") {
		if part == ".." {
			return ""
		}
	}
	return volume
}

func wire(path string) string {
	fields := make([]string, len(path))
	for i := range len(path) {
		fields[i] = strconv.Itoa(int(path[i]))
	}
	return "bytes:" + strings.Join(fields, ",")
}
