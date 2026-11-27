package main

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", "/", "\\", "///", "a\\\\b//c", "C:\\a\\b", "\\\\server\\share", "a\x00\\\xff/"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{'\\', byte(b), '/', byte(b)}))
	}
	check.CompareSamples("godiff-path-filepath-to-slash", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(128))
	r.Read(raw)
	if i%2 == 0 {
		for j := 0; j < len(raw); j += 4 {
			raw[j] = '\\'
		}
	}
	return string(raw)
}

func encode(path string) string {
	result := []byte("0")
	for i := range len(path) {
		result = strconv.AppendUint(append(result, ' '), uint64(path[i]), 10)
	}
	return string(result)
}

func native(path string) string {
	actual := filepath.ToSlash(path)
	switch filepath.Separator {
	case '/':
		return wire(actual) + "|" + wire(strings.ReplaceAll(path, "\\", "/"))
	case '\\':
		return wire(path) + "|" + wire(actual)
	default:
		panic("unsupported path separator")
	}
}

func expected(path string) string {
	return wire(path) + "|" + wire(strings.ReplaceAll(path, "\\", "/"))
}

func wire(path string) string {
	rows := make([]string, len(path))
	for i := range len(path) {
		rows[i] = strconv.Itoa(int(path[i]))
	}
	return strings.Join(rows, ",")
}
