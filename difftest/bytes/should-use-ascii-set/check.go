package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const casesEnv = "GOLEAN_BYTES_SHOULD_USE_ASCII_SET_CASES"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-bytes-should-use-ascii-set-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv(casesEnv, path); err != nil {
		return err
	}
	edges := make([]int, 0, 545)
	for length := -256; length <= 256; length++ {
		edges = append(edges, length)
	}
	for _, bound := range []int{math.MinInt, math.MinInt32, math.MaxInt32 - 7, math.MaxInt - 7} {
		for offset := range 8 {
			edges = append(edges, bound+offset)
		}
	}
	if err := compare(path, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		lengths := make([]int, 20_000)
		for i := range lengths {
			lengths[i] = int(rng.Uint64())
		}
		if err := compare(path, fmt.Sprintf("seed=%d", seed), lengths); err != nil {
			return err
		}
	}
	return nil
}

func compare(path, label string, lengths []int) error {
	var modelInput, fixture bytes.Buffer
	for _, length := range lengths {
		fmt.Fprintln(&modelInput, length)
		fmt.Fprintf(&fixture, "%d\t%t\n", length, length >= 9)
	}
	for i, got := range check.LeanRows("godiff-bytes-should-use-ascii-set", len(lengths), &modelInput) {
		length := lengths[i]
		want := strconv.FormatBool(length >= 9)
		if got != want {
			return fmt.Errorf("%s shouldUseASCIISet(%d): Lean=%s want=%s", label, length, got, want)
		}
	}
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/should-use-ascii-set/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanShouldUseASCIISet",
		fixtures+"bridge_test.go", fixtures+"should_use_ascii_set_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.shouldUseASCIISet %s: %d inputs agree.\n", label, len(lengths))
	return nil
}
