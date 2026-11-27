package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	check.CompareSamples("godiff-path-filepath-split-list-private", edges, sample, encode, expected, expected)

	file, err := os.CreateTemp("", "golean-split-list-private-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	write := func(path string) {
		fmt.Fprintf(file, "x%s", hex.EncodeToString([]byte(path)))
		for _, part := range oracle(path, filepath.ListSeparator) {
			fmt.Fprintf(file, " x%s", hex.EncodeToString([]byte(part)))
		}
		fmt.Fprintln(file)
	}
	for _, path := range edges {
		write(path)
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			write(sample(r, i))
		}
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_PRIVATE_SPLIT_LIST_CASES", file.Name()); err != nil {
		log.Fatal(err)
	}
	const dir = "difftest/path/filepath/split-list-private/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanPrivateSplitList",
		dir+"bridge_test.go", dir+"split_list_test.go"); err != nil {
		log.Fatal(err)
	}
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
