package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Profile, Limit int
	Text, Needle   []byte
	Offsets        []int
	Want           string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-suffixarray-lookup")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	writer := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for i := range cases {
			c := &cases[i]
			c.Offsets = suffixOffsets(*c)
			fmt.Fprint(&rows, strconv.IntSize, " ", c.Profile, " ", len(c.Text), " ", len(c.Needle), " ", c.Limit)
			for _, data := range [][]byte{c.Text, c.Needle} {
				for _, value := range data {
					fmt.Fprint(&rows, " ", value)
				}
			}
			for _, offset := range c.Offsets {
				fmt.Fprint(&rows, " ", offset)
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-suffixarray-lookup", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v: Lean %q", c, got)
			}
			if err := writer.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("suffixarray lookup %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var edges []fixture
	for _, limit := range []int{math.MinInt, -1, 0, 1, 2, math.MaxInt} {
		for needle := range 15 {
			key := binaryText(needle)
			edges = append(edges, fixture{Profile: 0, Limit: limit, Needle: key},
				fixture{Profile: 3, Limit: limit, Needle: key})
			for text := range 31 {
				for profile := 1; profile <= 2; profile++ {
					edges = append(edges, fixture{Profile: profile, Limit: limit, Text: binaryText(text), Needle: key})
				}
			}
		}
		for profile := 1; profile <= 2; profile++ {
			edges = append(edges, fixture{Profile: profile, Limit: limit, Text: []byte("banana"), Needle: []byte("ana")})
		}
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			text := make([]byte, rng.Intn(45))
			rng.Read(text)
			if i%3 == 0 {
				for j := range text {
					text[j] %= 4
				}
			}
			key := make([]byte, rng.Intn(8))
			rng.Read(key)
			if i%2 == 0 && len(text) > 0 {
				start := rng.Intn(len(text))
				key = text[start:min(len(text), start+len(key))]
			}
			profile := 1 + rng.Intn(2)
			if i%11 == 0 {
				profile = 0
				text = nil
			}
			cases[i] = fixture{Profile: profile, Limit: rng.Intn(24) - 3, Text: text, Needle: key}
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SUFFIXARRAY_LOOKUP_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayLookup", "difftest/index/suffixarray/lookup/testdata/lookup_test.go")
}

func binaryText(index int) []byte {
	var out []byte
	for bits := index + 1; bits > 1; bits >>= 1 {
		out = append(out, byte(bits&1)*255)
	}
	return out
}

func suffixOffsets(c fixture) []int {
	if c.Profile != 1 && c.Profile != 2 {
		return nil
	}
	offsets := make([]int, len(c.Text))
	for i := range offsets {
		offsets[i] = i
	}
	slices.SortFunc(offsets, func(a, b int) int {
		return cmp.Compare(string(c.Text[a:]), string(c.Text[b:]))
	})
	return append(offsets, 99)
}

func expected(c fixture) string {
	public := "nil"
	if c.Profile == 0 {
		if len(c.Needle) > 0 && c.Limit != 0 {
			public = "panic"
		}
		return "panic|" + public
	}
	if c.Profile == 3 {
		return "nil/nil|nil"
	}
	var lo, count int
	for _, offset := range c.Offsets[:len(c.Text)] {
		suffix, key := string(c.Text[offset:]), string(c.Needle)
		if suffix < key {
			lo++
		}
		if len(suffix) >= len(key) && suffix[:len(key)] == key {
			count++
		}
	}
	backing := wire(count, c.Offsets[lo:])
	all := backing + "/nil"
	if c.Profile == 2 {
		all = "nil/" + backing
	}
	if len(c.Needle) > 0 && c.Limit != 0 && count > 0 {
		if c.Limit > 0 {
			count = min(count, c.Limit)
		}
		public = wire(count, c.Offsets[lo:lo+count])
	}
	return all + "|" + public
}

func wire(length int, values []int) string {
	fields := make([]string, len(values))
	for i, n := range values {
		fields[i] = strconv.Itoa(n)
	}
	return fmt.Sprintf("%d:%s", length, strings.Join(fields, ","))
}
