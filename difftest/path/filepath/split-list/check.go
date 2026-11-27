package main

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", ":", ";", "a::b:", "a;;b;", "\"\"", "\"a;b\";c",
		"\"a;b", "a\";b\";c", "C:\\X;D:\\Y", "\x00:\xff;\"", "a\"\";b"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{'"', byte(b), ';', '"', ':', byte(b)}))
	}
	check.CompareSamples("godiff-path-filepath-split-list", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) string {
	raw := make([]byte, r.Intn(128))
	r.Read(raw)
	if i%2 == 0 {
		punctuation := []byte{'"', ':', ';'}
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
	actual := filepath.SplitList(path)
	switch filepath.ListSeparator {
	case ':':
		return wire(actual) + "|" + wire(oracle(path, ';')) + "|" + wire(oracle(path, 0))
	case ';':
		return wire(oracle(path, ':')) + "|" + wire(actual) + "|" + wire(oracle(path, 0))
	case 0:
		return wire(oracle(path, ':')) + "|" + wire(oracle(path, ';')) + "|" + wire(actual)
	default:
		panic("unsupported path list separator")
	}
}

func expected(path string) string {
	return wire(oracle(path, ':')) + "|" + wire(oracle(path, ';')) + "|" + wire(oracle(path, 0))
}

func oracle(path string, separator byte) []string {
	if path == "" {
		return []string{}
	}
	if separator != ';' {
		return strings.Split(path, string(separator))
	}
	var parts []string
	start := 0
	// A semicolon delimits exactly when an even number of quotes precede it.
	for i := range len(path) {
		if path[i] == ';' && strings.Count(path[:i], "\"")%2 == 0 {
			parts = append(parts, strings.ReplaceAll(path[start:i], "\"", ""))
			start = i + 1
		}
	}
	return append(parts, strings.ReplaceAll(path[start:], "\"", ""))
}

func wire(paths []string) string {
	if paths == nil {
		return "nil"
	}
	parts := make([]string, len(paths))
	for i, path := range paths {
		fields := make([]string, len(path))
		for j := range len(path) {
			fields[j] = strconv.Itoa(int(path[j]))
		}
		parts[i] = strings.Join(fields, ",")
	}
	return "slice:" + strconv.Itoa(len(paths)) + ":" + strings.Join(parts, ";")
}
