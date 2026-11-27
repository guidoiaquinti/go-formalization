package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Data    []byte
	Backing []uint16
	Length  int
	Nil     bool
}

func expected(c fixture) string {
	fields := make([]string, 1, 1+len(c.Backing))
	counts := [256]uint64{}
	fields[0] = "1"
	if len(c.Backing) >= 256 {
		fields[0] = "0"
		for _, b := range c.Data {
			counts[b]++
		}
	}
	for i, initial := range c.Backing {
		if i < len(counts) {
			initial += uint16(counts[i])
		}
		fields = append(fields, strconv.FormatUint(uint64(initial), 10))
	}
	return strings.Join(fields, " ")
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-histogram")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		line := make([]byte, 0, 2048)
		for _, c := range cases {
			line = strconv.AppendInt(line[:0], int64(len(c.Backing)), 10)
			for _, value := range c.Backing {
				line = strconv.AppendUint(append(line, ' '), uint64(value), 10)
			}
			for _, b := range c.Data {
				line = strconv.AppendUint(append(line, ' '), uint64(b), 10)
			}
			input.Write(append(line, '\n'))
			must(enc.Encode(c))
		}
		for i, got := range check.LeanRows("godiff-flate-histogram", len(cases), &input) {
			want := expected(cases[i])
			if got != want+" "+want {
				log.Fatalf("%s row %d disagreement (length %d, capacity %d, data length %d)", label, i, cases[i].Length, len(cases[i].Backing), len(cases[i].Data))
			}
		}
		fmt.Printf("%s: %d rows agree with independent oracle for both Lean entry points\n", label, len(cases))
	}
	rng := rand.New(rand.NewSource(4242))
	fixed := []fixture{{Nil: true}, {Data: []byte{0, 255}, Nil: true}}
	for _, size := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 8188, 8189, 8190, 8191, 8192, 8193, 8194, 8195, 65535, 65536, 65537, 131073} {
		for _, cap := range []int{0, 1, 255, 256, 257, 300} {
			for pattern := range 3 {
				c := sample(rng, size, cap)
				for i := range c.Data {
					switch pattern {
					case 0:
						c.Data[i] = 255
					case 1:
						c.Data[i] = byte(i)
					}
				}
				if cap >= 256 {
					c.Backing[0], c.Backing[255] = 65534, 65535
				}
				fixed = append(fixed, c)
				if cap == 256 {
					c.Length = 0
					fixed = append(fixed, c)
				}
			}
		}
	}
	for i := 0; i < len(fixed); i += 16 {
		compare("boundary batch", fixed[i:min(i+16, len(fixed))])
	}
	fmt.Printf("%d fixed rows\n", len(fixed))
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		caps := [...]int{1, 255, 256, 257, 300}
		for batch := range 20 {
			cases := make([]fixture, 1000)
			for i := range cases {
				cap := caps[rng.Intn(len(caps))]
				cases[i] = sample(rng, rng.Intn(97), cap)
			}
			compare(fmt.Sprintf("seed %d batch %d", seed, batch), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_HISTOGRAM_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHistogram", "difftest/compress/flate/histogram/testdata/histogram_test.go"))
}

func sample(rng *rand.Rand, size, cap int) fixture {
	c := fixture{Data: make([]byte, size), Backing: make([]uint16, cap), Length: rng.Intn(cap + 1)}
	_, err := rng.Read(c.Data)
	must(err)
	for i := range c.Backing {
		c.Backing[i] = uint16(rng.Uint32())
	}
	return c
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
