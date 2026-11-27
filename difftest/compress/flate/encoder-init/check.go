package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func expected(size int) string {
	if size < 0 {
		return "bounds|65535 65535"
	}
	capacity := 1
	for capacity < size {
		capacity *= 2
	}
	return fmt.Sprintf("%d %d true|65535 65535", size, capacity)
}

func main() {
	file, err := os.CreateTemp("", "golean-flate-init")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	native := bufio.NewWriter(file)
	compare := func(label string, sizes []int) {
		var input bytes.Buffer
		for _, size := range sizes {
			fmt.Fprintln(&input, strconv.IntSize, size)
		}
		for i, row := range check.LeanRows("godiff-flate-encoder-init", len(sizes), &input) {
			size := sizes[i]
			want := expected(size)
			if row != want {
				log.Fatalf("%s size %d: Lean=%q want=%q", label, size, row, want)
			}
			_, err := fmt.Fprintln(native, size, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean inputs agree\n", label, len(sizes))
	}
	edges := []int{math.MinInt, math.MinInt + 1, -4097, -1025, -513, -257, -2, -1}
	for i := 1; i <= 4096; i++ {
		edges = append(edges, i)
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		sizes := make([]int, 20_000)
		for i := range sizes {
			if i%4 == 0 {
				sizes[i] = -1 - rng.Intn(math.MaxInt)
			} else {
				sizes[i] = 1 + rng.Intn(1024)
			}
		}
		compare(fmt.Sprintf("seed %d", seed), sizes)
	}
	must(native.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_FLATE_INIT_ROWS", file.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanEncoderInit", "difftest/compress/flate/encoder-init/testdata/init_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
