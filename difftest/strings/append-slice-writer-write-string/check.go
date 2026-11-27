package main

import (
	"bufio"
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

const nilPointer, nilSlice, independent = -1, -2, -1

type input struct {
	width, source int
	backing, text string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-append-slice-writer-write-string-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := os.Setenv("GOLEAN_STRINGS_APPEND_SLICE_WRITER_WRITE_STRING_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	texts := []string{"", "a", "\xff\x00", "\xe2\x82", "🚀", strings.Repeat("\x00\xff", 128)}
	edges := []input{
		{4, 1, "abcdefgh", "bcd"},
		{1, 7, "abcdefghij", "hij"},
		{6, 3, "abcdefgh", "defgh"},
	}
	for _, text := range texts {
		edges = append(edges, input{nilPointer, independent, "", text}, input{nilSlice, independent, "", text})
	}
	for capacity := range 21 {
		raw := make([]byte, capacity)
		for i := range raw {
			raw[i] = byte(255 - i)
		}
		for width := 0; width <= capacity; width++ {
			for size := range 21 {
				edges = append(edges, input{width, independent, string(raw), strings.Repeat("\x80", size)})
			}
			for start := 0; start <= capacity; start++ {
				candidate := input{width, start, string(raw), string(raw[start:])}
				if sourceStable(candidate) {
					edges = append(edges, candidate)
				}
			}
		}
	}
	for value := range 256 {
		raw := string([]byte{byte(value)})
		edges = append(edges, input{1, independent, raw + "\x00", "\xff" + raw})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fixture, fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const fixtures = "difftest/strings/append-slice-writer-write-string/testdata/"
	return check.TestInPackage("strings", "TestGoLeanAppendSliceWriterWriteString",
		fixtures+"bridge_test.go", fixtures+"append_slice_writer_write_string_test.go")
}

func randomBytes(rng *rand.Rand, n int) string {
	raw := make([]byte, n)
	rng.Read(raw)
	return string(raw)
}

func sample(rng *rand.Rand, i int) input {
	if i%19 == 0 {
		return input{nilSlice, independent, "", randomBytes(rng, rng.Intn(20))}
	}
	capacity := rng.Intn(64)
	raw := randomBytes(rng, capacity)
	width := rng.Intn(capacity + 1)
	switch i % 5 {
	case 0:
		return input{width, independent, raw, ""}
	case 1, 2:
		start := rng.Intn(capacity + 1)
		end := start + rng.Intn(capacity-start+1)
		candidate := input{width, start, raw, raw[start:end]}
		if sourceStable(candidate) {
			return candidate
		}
		return input{width, independent, raw, candidate.text}
	default:
		return input{width, independent, raw, randomBytes(rng, rng.Intn(40))}
	}
}

// A string view may share backing storage, but Go requires its bytes
// to remain immutable. Test only disjoint writes or reallocation.
func sourceStable(c input) bool {
	if len(c.text) == 0 {
		return true
	}
	end := c.width + len(c.text)
	return end > len(c.backing) || c.source+len(c.text) <= c.width || end <= c.source
}

func expected(c input) string {
	if c.width == nilPointer {
		return "panic"
	}
	width := max(c.width, 0)
	data := c.backing[:width] + c.text
	nilness := "slice"
	if c.width == nilSlice && len(c.text) == 0 {
		nilness = "nil"
	}
	capacity := strconv.Itoa(len(c.backing))
	if len(data) > len(c.backing) {
		capacity = "growth"
	}
	row := fmt.Appendf(nil, "%d:false:%s:%s:", len(c.text), nilness, capacity)
	for i := range len(data) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(data[i]), 10)
	}
	return string(row)
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := fmt.Appendf(nil, "%d %d", c.width, len(c.backing))
		for _, raw := range []string{c.backing[:max(c.width, 0)], c.text} {
			for j := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
			}
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = expected(c)
		if _, err := fmt.Fprintf(fixture, "%d\t%d\t%x\t%x\t%s\n", c.width, c.source, c.backing, c.text, wants[i]); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-append-slice-writer-write-string", len(cases), &modelInput) {
		if got != wants[i] {
			c := cases[i]
			return fmt.Errorf("%s appendSliceWriter.WriteString(%+v): Lean=%s want=%s", label, c, got, wants[i])
		}
	}
	fmt.Printf("strings.appendSliceWriter.WriteString %s: %d Lean results agree.\n", label, len(cases))
	return nil
}
