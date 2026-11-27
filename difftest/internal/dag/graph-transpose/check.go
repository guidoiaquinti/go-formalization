package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type transposeCase struct {
	size  int
	codes []int
}

func (c transposeCase) want() string {
	var want strings.Builder
	want.WriteString("ok")
	for from := range c.size {
		for to := range c.size {
			if c.codes[to*c.size+from] != 0 {
				want.WriteString(" 1")
			} else {
				want.WriteString(" 0")
			}
		}
	}
	return want.String()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	seeds := []int64{17, 2718, 65537}
	cases := []transposeCase{{}}
	for _, size := range []int{1, 2} {
		count := 1
		for range size * size {
			count *= 3
		}
		for mask := range count {
			codes := make([]int, size*size)
			for i := range codes {
				codes[i], mask = mask%3, mask/3
			}
			cases = append(cases, transposeCase{size, codes})
		}
	}
	boundaries := len(cases)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			size := rng.Intn(9)
			codes := make([]int, size*size)
			for i := range codes {
				codes[i] = rng.Intn(3)
			}
			cases = append(cases, transposeCase{size, codes})
		}
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.size)
		for _, code := range c.codes {
			fmt.Fprint(&input, " ", code)
		}
		fmt.Fprintln(&input)
	}
	dir, err := os.MkdirTemp("", "golean-transpose-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "inputs"), input.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_TRANSPOSE_DIR", dir); err != nil {
		return err
	}
	lean := check.LeanRows("godiff-internal-dag-graph-transpose", len(cases), &input)
	const fixture = "difftest/internal/dag/graph-transpose/testdata/transpose_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanTranspose", fixture); err != nil {
		return err
	}
	output, err := os.ReadFile(filepath.Join(dir, "outputs"))
	if err != nil {
		return err
	}
	goRows := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(goRows) != len(cases) {
		return fmt.Errorf("Go produced %d rows, want %d", len(goRows), len(cases))
	}
	for i, c := range cases {
		if want := c.want(); goRows[i] != want || lean[i] != want {
			return fmt.Errorf("case %d %+v, seeds %v: Go=%q Lean=%q want=%q", i, c, seeds, goRows[i], lean[i], want)
		}
	}
	fmt.Printf("internal/dag.Graph.Transpose: %d boundary and %d sampled cases agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}
