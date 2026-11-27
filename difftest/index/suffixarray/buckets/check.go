package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Profile, FSize, BSize   int
	Text, Frequency, Bucket []int64
	Want                    string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-suffixarray-buckets")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	jsonRows := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&rows, c.Profile, " ", c.FSize, " ", len(c.Frequency), " ", c.BSize, " ", len(c.Bucket), " ", len(c.Text))
			for _, block := range [][]int64{c.Text, c.Frequency, c.Bucket} {
				for _, n := range block {
					fmt.Fprint(&rows, " ", n)
				}
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-suffixarray-buckets", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v: Lean %q want %q", c, got, c.Want)
			}
			if err := jsonRows.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("suffixarray frequencies and buckets %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	for profile := range 4 {
		var edges []fixture
		shapes := [][2]int{{-1, 0}, {0, 0}, {0, 3}, {1, 1}, {2, 3}, {10, 256}, {255, 255}, {255, 256}, {256, 256}, {256, 258}, {258, 258}}
		for _, f := range shapes {
			for _, b := range shapes {
				for _, head := range []int64{-1, 0, math.MaxInt32, math.MaxInt64} {
					for _, text := range [][]int64{nil, {0}, {1, 0, 255}, {-1, 256, 0}, {0, 0, 1, 2}, {0, math.MaxInt64, 2}} {
						edges = append(edges, newFixture(profile, f, b, head, text))
					}
				}
			}
		}
		edges = append(edges, newFixture(profile, [2]int{256, 256}, [2]int{256, 256}, 3, []int64{255}),
			newFixture(profile, [2]int{2, 2}, [2]int{2, 2}, 3, []int64{1}),
			newFixture(profile, [2]int{256, 256}, [2]int{256, 256}, 3, nil),
			newFixture(profile, [2]int{2, 2}, [2]int{2, 2}, 3, nil))
		if err := compare(fmt.Sprintf("profile=%d boundaries", profile), edges); err != nil {
			return err
		}
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			// Keep the wire batches bounded when tables have 256 cells.
			for block := range 4 {
				cases := make([]fixture, 5_000)
				for i := range cases {
					f, b := shapes[rng.Intn(len(shapes))], shapes[rng.Intn(len(shapes))]
					cells := min(256, max(f[0], b[0]))
					text := make([]int64, rng.Intn(24))
					for j := range text {
						text[j] = int64(rng.Intn(cells+2) - 1)
					}
					head := []int64{-1, 0, 1, math.MaxInt32, math.MaxInt64}[rng.Intn(5)]
					if i%4 == 0 {
						f, b = [2]int{cells, cells}, [2]int{cells, cells}
					}
					cases[i] = newFixture(profile, f, b, head, text)
				}
				if err := compare(fmt.Sprintf("profile=%d seed=%d block=%d", profile, seed, block), cases); err != nil {
					return err
				}
			}
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SUFFIXARRAY_BUCKET_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayBuckets", "difftest/index/suffixarray/buckets/testdata/buckets_test.go")
}

func newFixture(profile int, f, b [2]int, head int64, text []int64) fixture {
	c := fixture{Profile: profile, FSize: f[0], BSize: b[0], Text: text}
	if f[0] >= 0 {
		c.Frequency = make([]int64, f[1])
		for i := range c.Frequency {
			c.Frequency[i] = 3
		}
		if len(c.Frequency) > 0 {
			c.Frequency[0] = head
		}
	}
	if b[0] >= 0 {
		c.Bucket = make([]int64, b[1])
		for i := range c.Bucket {
			c.Bucket[i] = 9
		}
	}
	return c
}

func expected(c fixture) string {
	var out []string
	byteInput := c.Profile < 2
	narrow := func(x int64) int64 {
		if c.Profile%2 == 0 {
			return int64(int32(x))
		}
		return x
	}
	for mode := range 3 {
		f, b := slices.Clone(c.Frequency), slices.Clone(c.Bucket)
		for _, data := range [][]int64{f, b} {
			for i, x := range data {
				data[i] = narrow(x)
			}
		}
		render := func(status string) string { return status + "/" + wire(c.FSize, f) + "/" + wire(c.BSize, b) }
		counts, n := f, c.FSize
		if f == nil {
			counts, n = b, max(c.BSize, 0)
		}
		if f != nil && n == 0 {
			out = append(out, render("panic"))
			continue
		}
		if f == nil || counts[0] < 0 {
			if byteInput {
				n = 256
			}
			if n > len(counts) {
				out = append(out, render("panic"))
				continue
			}
			valid := true
			hist := make(map[int]int64)
			for _, symbol := range c.Text {
				i := narrow(symbol)
				if byteInput {
					i = int64(byte(symbol))
				}
				if i < 0 || i >= int64(n) {
					valid = false
					break
				}
				hist[int(i)]++
			}
			for i := range n {
				counts[i] = hist[i]
			}
			if !valid {
				out = append(out, render("panic"))
				continue
			}
		}
		if mode == 0 {
			out = append(out, render(wire(n, counts)))
			continue
		}
		target := max(c.BSize, 0)
		if byteInput {
			n, target = 256, 256
			if len(counts) < n || len(b) < target {
				out = append(out, render("panic"))
				continue
			}
		}
		cumulative := make([]uint64, n+1)
		for i, x := range counts[:n] {
			cumulative[i+1] = cumulative[i] + uint64(x)
		}
		for i := range min(n, target) {
			b[i] = narrow(int64(cumulative[i+mode-1]))
		}
		status := "ok"
		if n > target {
			status = "panic"
		}
		out = append(out, render(status))
	}
	return strings.Join(out, "|")
}

func wire(length int, data []int64) string {
	if data == nil {
		return "nil"
	}
	var runs []string
	for i := 0; i < len(data); {
		j := i + 1
		for j < len(data) && data[j] == data[i] {
			j++
		}
		runs = append(runs, fmt.Sprintf("%d.%d", data[i], j-i))
		i = j
	}
	return fmt.Sprintf("%d:%s", length, strings.Join(runs, ","))
}
