package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	payload  string
	mode     int
	capacity int
	request  int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for mode := range 4 {
		for _, payload := range []string{"", "abc", "é🚀", "\xff\x00"} {
			for _, reserve := range []int{-1, 0, 1, 7, 64, 32768} {
				capacity := len(payload) + reserve
				if reserve == -1 {
					if len(payload) != 0 {
						continue
					}
					capacity = -1
				}
				base := 2 * max(capacity, 0)
				for _, size := range []int{-1, 0, 1, len(payload) - 1, len(payload), len(payload) + 1, 8, 32768, 32769} {
					cases = append(cases, input{payload, mode, capacity, size - base})
				}
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			payload := make([]byte, rng.Intn(65))
			rng.Read(payload)
			capacity := len(payload) + rng.Intn(129)
			if i%7 == 0 {
				payload, capacity = nil, -1
			}
			sizes := []int{-1, 0, len(payload) - 1, len(payload), rng.Intn(4096)}
			request := sizes[i%len(sizes)] - 2*max(capacity, 0)
			cases = append(cases, input{string(payload), rng.Intn(4), capacity, request})
		}
	}

	var modelInput, fixture bytes.Buffer
	want := make([]string, len(cases))
	for i, c := range cases {
		size := 2*max(c.capacity, 0) + c.request
		rounded := 0
		if size >= 0 && c.mode != 3 {
			rounded = cap(slices.Grow([]byte(nil), size))
		}
		row := fmt.Appendf(nil, "%d %d %d %d %d", c.mode, c.capacity, c.request, size, rounded)
		for j := range len(c.payload) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.payload[j]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		want[i] = expected(c, size, rounded)
		fmt.Fprintf(&fixture, "%d %d %d x%x %s\n", c.mode, c.capacity, c.request, c.payload, want[i])
	}
	for i, got := range check.LeanRows("godiff-strings-builder-grow-private", len(cases), &modelInput) {
		if got != want[i] {
			return fmt.Errorf("%+v: Lean=%q want=%q", cases[i], got, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-builder-grow-private-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_STRINGS_GROW_PRIVATE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/strings/builder-grow-private/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanBuilderGrowPrivate", data+"bridge_test.go", data+"grow_test.go", "difftest/strings/builder-copy-check/testdata/bridge_test.go"); err != nil {
		return err
	}
	fmt.Printf("Builder.grow: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input, size, rounded int) string {
	if c.mode == 3 {
		return "nil"
	}
	status := "ok"
	identity := [...]string{"unset", "self", "foreign"}[c.mode]
	capacity := c.capacity
	switch {
	case size < 0:
		status = "allocation"
	case rounded < len(c.payload):
		status = "bounds"
	default:
		capacity = rounded
	}
	payload := make([]string, len(c.payload))
	for i := range len(c.payload) {
		payload[i] = strconv.Itoa(int(c.payload[i]))
	}
	return fmt.Sprintf("%s:%s:%d:%s", status, identity, capacity, strings.Join(payload, ","))
}
