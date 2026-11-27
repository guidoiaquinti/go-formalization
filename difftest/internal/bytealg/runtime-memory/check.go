package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Left, Right                           []byte
	Length, Size, LeftOffset, RightOffset int
	Alias                                 bool
	Want                                  string
}

func main() {
	f, err := os.CreateTemp("", "golean-bytealg-runtime-memory-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	var edges []fixture
	for n := 0; n <= 512; n++ {
		edges = append(edges, fixture{Length: n})
	}
	edges = append(edges, fixture{Length: -1}, fixture{Length: math.MinInt}, fixture{Length: -512})
	if strconv.IntSize == 64 {
		edges = append(edges, fixture{Length: math.MaxInt})
	}
	for offset := range 64 {
		for _, size := range []int{0, 1, 2, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129, 255, 256, 257, 511, 512} {
			left := bytes.Repeat([]byte{0xff}, size+1)
			right := bytes.Clone(left)
			right[size] = 0x80
			edges = append(edges, fixture{Left: left, Right: right, Length: size, Size: size, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: left, Right: right, Length: size, Size: size + 1, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: left, Right: left, Length: size, Size: size, LeftOffset: offset, RightOffset: offset, Alias: true})
		}
	}
	compare("boundaries", edges, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				a, b := make([]byte, rng.Intn(257)), make([]byte, rng.Intn(257))
				rng.Read(a)
				rng.Read(b)
				alias := false
				switch i % 5 {
				case 0:
					b = a
					alias = true
				case 1:
					b = bytes.Clone(a)
				case 2:
					b = bytes.Clone(a)
					if len(b) > 0 {
						b[rng.Intn(len(b))] ^= 255
					}
				case 3:
					b = append(bytes.Clone(a), 0)
				}
				size := rng.Intn(min(len(a), len(b)) + 1)
				if i%2 == 0 {
					size = min(len(a), len(b))
				}
				cases[i] = fixture{Left: a, Right: b, Length: rng.Intn(520) - 7, Size: size,
					LeftOffset: rng.Intn(64), RightOffset: rng.Intn(64), Alias: alias}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_BYTEALG_RUNTIME_MEMORY_CASES", f.Name()))
	const data = "difftest/internal/bytealg/runtime-memory/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgRuntimeMemory", data+"bridge_test.go", data+"runtime_memory_test.go"))
}

func expected(c fixture) string {
	allocation := "panic"
	if c.Length >= 0 && c.Length <= 512 {
		classes := []int{0, 8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
		i, _ := slices.BinarySearch(classes, c.Length)
		allocation = fmt.Sprintf("%d:%d:true", c.Length, classes[i])
	}
	return fmt.Sprintf("%s|%t", allocation, slices.Equal(c.Left[:c.Size], c.Right[:c.Size]))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d", c.Length, c.Size, len(c.Left))
		for _, raw := range [][]byte{c.Left, c.Right} {
			for _, b := range raw {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-runtime-memory", len(cases), &input) {
		c := cases[i]
		c.Want = expected(c)
		if got != c.Want {
			log.Fatalf("%s row %d %+v: Lean=%s want=%s", label, i, c, got, c.Want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("bytealg runtime memory %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
