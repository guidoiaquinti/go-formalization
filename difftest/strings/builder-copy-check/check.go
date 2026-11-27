package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	payload  string
	mode     int
	capacity int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for mode := range 4 {
		cases = append(cases, input{"", mode, -1})
		for _, s := range []string{"", "abc", "é🚀", "\xff\x00"} {
			for _, extra := range []int{0, 1, 7, 64} {
				cases = append(cases, input{s, mode, len(s) + extra})
			}
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			b := make([]byte, rng.Intn(65))
			rng.Read(b)
			mode := rng.Intn(4)
			if i%5 == 0 {
				cases = append(cases, input{"", mode, -1})
			} else {
				cases = append(cases, input{string(b), mode, len(b) + rng.Intn(65)})
			}
		}
	}

	var modelInput, fixture bytes.Buffer
	want := make([]string, len(cases))
	for i, c := range cases {
		row := fmt.Appendf(nil, "%d %d", c.mode, c.capacity)
		for j := range len(c.payload) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.payload[j]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		want[i] = expected(c)
		fmt.Fprintf(&fixture, "%d %d x%x %s\n", c.mode, c.capacity, c.payload, want[i])
	}
	for i, got := range check.LeanRows("godiff-strings-builder-copy-check", len(cases), &modelInput) {
		if got != want[i] {
			return fmt.Errorf("%+v: Lean=%q want=%q", cases[i], got, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-builder-copy-check-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_STRINGS_COPY_CHECK_CASES", path); err != nil {
		return err
	}
	const data = "difftest/strings/builder-copy-check/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanBuilderCopyCheck", data+"bridge_test.go", data+"copy_check_test.go"); err != nil {
		return err
	}
	fmt.Printf("Builder.copyCheck: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input) string {
	if c.mode == 3 {
		return "nil"
	}
	status, identity := "ok", "self"
	if c.mode == 2 {
		status, identity = "copied", "foreign"
	}
	payload := make([]string, len(c.payload))
	for i := range len(c.payload) {
		payload[i] = strconv.Itoa(int(c.payload[i]))
	}
	return fmt.Sprintf("%s:%s:%d:%s", status, identity, c.capacity, strings.Join(payload, ","))
}
