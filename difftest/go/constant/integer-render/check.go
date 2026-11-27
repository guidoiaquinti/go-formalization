package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Word    int64
	Big     string
	Present bool
	Want    string
}

func decimal(magnitude []byte, negative bool) string {
	// Convert from base 256 to base 1e9 without calling an integer formatter.
	const base = 1_000_000_000
	var chunks []uint64
	for _, b := range magnitude {
		carry := uint64(b)
		for i, chunk := range chunks {
			acc := chunk*256 + carry
			chunks[i], carry = acc%base, acc/base
		}
		if carry != 0 {
			chunks = append(chunks, carry)
		}
	}
	if len(chunks) == 0 {
		return "0"
	}
	var out strings.Builder
	if negative {
		out.WriteByte('-')
	}
	fmt.Fprint(&out, chunks[len(chunks)-1])
	for i := len(chunks) - 2; i >= 0; i-- {
		fmt.Fprintf(&out, "%09d", chunks[i])
	}
	return out.String()
}

func newFixture(word int64, negative bool, magnitude []byte, present bool) fixture {
	value := uint64(word)
	if word < 0 {
		value = -value
	}
	var raw [8]byte
	binary.BigEndian.PutUint64(raw[:], value)
	short := decimal(raw[:], word < 0)
	large := decimal(magnitude, negative)
	rendered := large
	if !present {
		rendered = "<nil>"
	}
	magnitude = bytes.TrimLeft(magnitude, "\x00")
	n := 0
	if len(magnitude) != 0 {
		n = (len(magnitude)-1)*8 + bits.Len8(magnitude[0])
	}
	return fixture{word, large, present, strings.Join([]string{short, short, rendered, rendered,
		strconv.FormatBool(int64(int32(word)) == word),
		strconv.FormatBool((word<<1)>>1 == word), strconv.FormatBool(n < 4096)}, "|")}
}

func main() {
	file, err := os.CreateTemp("", "golean-constant-integer-render-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	writer := bufio.NewWriter(file)
	enc := json.NewEncoder(writer)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.Present {
				present = 1
			}
			fmt.Fprintln(&input, c.Word, c.Big, present)
		}
		for i, result := range check.LeanRows("godiff-go-constant-integer-render", len(cases), &input) {
			c := cases[i]
			if result != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, result, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	power := func(bit int) []byte {
		result := make([]byte, bit/8+1)
		result[0] = 1 << (bit % 8)
		return result
	}
	var fixed []fixture
	edges := []int64{math.MinInt64, -(1 << 62) - 1, -(1 << 62), -(1 << 62) + 1, math.MinInt32 - 1, math.MinInt32,
		-1, 0, 1, math.MaxInt32, math.MaxInt32 + 1, 1<<62 - 1, 1 << 62, math.MaxInt64}
	for _, word := range edges {
		for _, bits := range [][]byte{nil, {0, 0, 9}, {255}, power(4094), power(4095), power(4096)} {
			for _, neg := range []bool{false, true} {
				for _, present := range []bool{false, true} {
					fixed = append(fixed, newFixture(word, neg, bits, present))
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			raw := make([]byte, r.Intn(33))
			_, err := r.Read(raw)
			must(err)
			word := int64(r.Uint64())
			if i%3 == 0 {
				word = edges[r.Intn(len(edges))]
			}
			cases[i] = newFixture(word, r.Intn(2) == 0, raw, i%7 != 0)
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CONSTANT_INTEGER_RENDER", file.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanIntegerRender", "difftest/go/constant/integer-render/testdata/render_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
