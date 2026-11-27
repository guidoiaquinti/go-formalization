package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Length  int
	Index, Value uint64
	Cells        []byte
	Want         string
}

var widths = [...]int{8, 8, 16, 16, 32, 32, 64, 64, strconv.IntSize, strconv.IntSize}

func normalize(word uint64, width, signed int) *big.Int {
	var value, modulus big.Int
	modulus.SetUint64(1)
	modulus.Lsh(&modulus, uint(width))
	value.SetUint64(word)
	value.Mod(&value, &modulus)
	if signed == 1 && value.Bit(width-1) == 1 {
		value.Sub(&value, &modulus)
	}
	return &value
}

func expected(c fixture, width, signed int, start *big.Int) string {
	results := []string{"bounds", "bounds", "bounds", "false", ""}
	if start.IsUint64() && start.Uint64() < uint64(len(c.Cells)) {
		at := int(start.Uint64())
		if at < c.Length {
			results[0] = strconv.Itoa(int(c.Cells[at]))
		}
		maxIndex := uint64(1)<<(width-signed) - 1
		if width == 64 && signed == 0 {
			maxIndex = ^uint64(0)
		}
		if len(c.Cells)-at >= 4 && uint64(at+4) <= maxIndex {
			results[1] = fmt.Sprint(binary.LittleEndian.Uint32(c.Cells[at:]))
		}
		if len(c.Cells)-at >= 8 && uint64(at+8) <= maxIndex {
			results[2] = fmt.Sprint(binary.LittleEndian.Uint64(c.Cells[at:]))
		}
	}
	cells := slices.Clone(c.Cells)
	if c.Length >= 8 {
		binary.LittleEndian.PutUint64(cells, c.Value)
		results[3] = "true"
	}
	output := make([]string, len(cells))
	for i, b := range cells {
		output[i] = strconv.Itoa(int(b))
	}
	results[4] = strings.Join(output, ",")
	return strings.Join(results, "|")
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-endian")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for i, c := range cases {
			width, signed := widths[c.Tag], 1-c.Tag%2
			start := normalize(c.Index, width, signed)
			fmt.Fprint(&input, width, " ", signed, " ", start, " ", c.Length, " ", c.Value)
			for _, b := range c.Cells {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
			cases[i].Want = expected(c, width, signed, start)
		}
		for i, row := range check.LeanRows("godiff-flate-little-endian", len(cases), &input) {
			c := cases[i]
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for tag, width := range widths {
		step := uint64(1) << (width - 1)
		for _, index := range []uint64{0, 1, 3, 7, 8, 15, step - 8, step - 4, step - 1, step, ^uint64(0) - 7, ^uint64(0) - 3, ^uint64(0)} {
			for _, size := range []int{0, 1, 3, 4, 7, 8, 9, 20, 270} {
				for _, length := range []int{0, size / 2, size} {
					c := fixture{Tag: tag, Length: length, Index: index, Value: 0xffeeddccbbaa9988, Cells: make([]byte, size)}
					for i := range c.Cells {
						c.Cells[i] = byte(i + 1)
					}
					fixed = append(fixed, c)
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Tag: r.Intn(10), Index: r.Uint64(), Value: r.Uint64(), Cells: make([]byte, r.Intn(64))}
			c.Length = r.Intn(len(c.Cells) + 1)
			if i%2 == 0 {
				c.Index = uint64(r.Intn(67))
			}
			r.Read(c.Cells)
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_ENDIAN_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanFlateEndian", "difftest/compress/flate/little-endian/testdata/little_endian_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
