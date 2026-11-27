package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tokens []int
	Want   string
}

func main() {
	fixed := [][]int{{0}, {7}, {8, 2, 8, 5, 7}, {8, 2, 8, 1, 7}, {9, 0, 1}, {10, 0, 1},
		{10, 4, 8, 2, 8, 5, 7}, {11, 1, 5}, {11, 0, 10, 3, 4}, {11, 0, 11, 0, 2}}
	for first := 1; first <= 7; first++ {
		fixed = append(fixed, []int{first})
		for length := range 12 {
			for second := 1; second <= 7; second++ {
				fixed = append(fixed, []int{9, length, 8, first, 8, second, 7},
					[]int{10, length, 8, first, 8, second, 7}, []int{11, length % 2, 9, length, first})
			}
		}
	}
	for tag := 12; tag <= 22; tag++ {
		fixed = append(fixed, []int{tag}, []int{10, 3, tag}, []int{11, 0, tag})
	}
	file, err := os.CreateTemp("", "golean-binary-reflect-size")
	must(err)
	defer os.Remove(file.Name())
	buf := bufio.NewWriter(file)
	enc := json.NewEncoder(buf)
	run(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			rows := make([][]int, 5000)
			for i := range rows {
				rows[i] = randomLayout(rng, 3)
			}
			run(rows, enc)
		}
	}
	must(buf.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_REFLECT_SIZE_CASES", file.Name()))
	const testdata = "difftest/encoding/binary/reflect-size/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanReflectSize", testdata+"bridge_test.go", testdata+"reflect_size_test.go"))
	fmt.Printf("%d boundaries and 20,000 per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func randomLayout(r *rand.Rand, depth int) []int {
	if depth == 0 || r.Intn(3) == 0 {
		atoms := []int{1, 2, 3, 4, 5, 6, 7, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}
		return []int{atoms[r.Intn(len(atoms))]}
	}
	switch r.Intn(4) {
	case 0:
		result := []int{}
		for range r.Intn(4) {
			result = append(result, 8)
			result = append(result, randomLayout(r, depth-1)...)
		}
		return append(result, 7)
	default:
		return append([]int{9 + r.Intn(3), r.Intn(6)}, randomLayout(r, depth-1)...)
	}
}

func run(cases [][]int, enc *json.Encoder) {
	var input bytes.Buffer
	for _, tokens := range cases {
		var separator string
		for _, v := range tokens {
			fmt.Fprintf(&input, "%s%d", separator, v)
			separator = " "
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-binary-reflect-size", len(cases), &input) {
		must(enc.Encode(fixture{cases[i], row}))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
