package main

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", ".", "..", ".git", "a.b", "a.", "a.b/", "a.b\\c",
		"C:\\d.ext", "//x/.a.b", "a\x00.\xff", "\\\\srv\\share"}
	for b := range 256 {
		edges = append(edges, string([]byte{'.', byte(b)}), string([]byte{'.', byte(b), '.', byte(b)}))
	}
	check.CompareSamples("godiff-path-filepath-ext", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(128))
	r.Read(raw)
	if i%2 == 0 {
		punctuation := []byte{'.', '/', '\\'}
		for j := 0; j < len(raw); j += 4 {
			raw[j] = punctuation[r.Intn(len(punctuation))]
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
	actual := filepath.Ext(path)
	switch filepath.Separator {
	case '/':
		return wire(actual) + "|" + wire(oracle(path, true))
	case '\\':
		return wire(oracle(path, false)) + "|" + wire(actual)
	default:
		panic("unsupported path separator")
	}
}

func expected(path string) string {
	return wire(oracle(path, false)) + "|" + wire(oracle(path, true))
}

func oracle(path string, windows bool) string {
	slash := strings.LastIndexByte(path, '/')
	if windows {
		slash = max(slash, strings.LastIndexByte(path, '\\'))
	}
	dot := strings.LastIndexByte(path, '.')
	if dot <= slash {
		return ""
	}
	return path[dot:]
}

func wire(path string) string {
	rows := make([]string, len(path))
	for i := range len(path) {
		rows[i] = strconv.Itoa(int(path[i]))
	}
	return strings.Join(rows, ",")
}
