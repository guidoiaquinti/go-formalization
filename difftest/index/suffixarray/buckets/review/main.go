package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Profile, FSize, BSize   int
	Text, Frequency, Bucket []int64
	Want                    string
}

func view(size int, cells []int64) string {
	if cells == nil {
		return "nil"
	}
	var runs []string
	for i := 0; i < len(cells); {
		j := i + 1
		for j < len(cells) && cells[j] == cells[i] {
			j++
		}
		runs = append(runs, fmt.Sprintf("%d.%d", cells[i], j-i))
		i = j
	}
	return fmt.Sprintf("%d:%s", size, strings.Join(runs, ","))
}

func reports(freq, minStatus, maxStatus string, fsize, bsize int, originalFreq, originalBucket, newFreq, min, max []int64) string {
	f := view(fsize, newFreq)
	b := view(bsize, originalBucket)
	if originalFreq == nil && freq != "panic" {
		f = "nil"
		b = freq
	}
	return freq + "/" + f + "/" + b + "|" +
		minStatus + "/" + f + "/" + view(bsize, min) + "|" +
		maxStatus + "/" + f + "/" + view(bsize, max)
}

func main() {
	var cases []fixture
	for profile := range 2 {
		f, b, lo, hi := make([]int64, 257), make([]int64, 257), make([]int64, 257), make([]int64, 257)
		f[0], f[1], f[256] = 2, 1, 77
		for i := range b {
			b[i], lo[i], hi[i] = 8, 3, 3
		}
		lo[0], lo[1], lo[256] = 0, 2, 8
		hi[0], hi[256] = 2, 8
		cases = append(cases, fixture{profile, 256, 256, []int64{255}, f, b, reports(view(256, f), "ok", "ok", 256, 256, f, b, f, lo, hi)})
	}
	for _, profile := range []int{2, 3} {
		b := []int64{8, 8, 8, 99}
		cases = append(cases,
			fixture{profile, -1, 3, []int64{2, 0, 2}, nil, b,
				reports("3:1.1,0.1,2.1,99.1", "ok", "ok", -1, 3, nil, b, nil, []int64{0, 1, 1, 99}, []int64{1, 1, 3, 99})},
			fixture{profile, 3, 3, []int64{0, 3}, []int64{-1, 8, 8, 77}, b,
				reports("panic", "panic", "panic", 3, 3, []int64{-1}, b, []int64{1, 0, 0, 77}, b, b)},
		)
		maximum, minimum := int64(math.MaxInt32), int64(math.MinInt32)
		if profile == 3 {
			maximum, minimum = math.MaxInt64, math.MinInt64
		}
		f := []int64{maximum, 1, 7}
		cases = append(cases, fixture{profile, 2, 2, nil, f, b,
			reports(view(2, f), "ok", "ok", 2, 2, f, b, f, []int64{0, maximum, 8, 99}, []int64{maximum, minimum, 8, 99})})
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d %d", c.Profile, c.FSize, len(c.Frequency), c.BSize, len(c.Bucket), len(c.Text))
		for _, values := range [][]int64{c.Text, c.Frequency, c.Bucket} {
			for _, v := range values {
				fmt.Fprintf(&input, " %d", v)
			}
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-suffixarray-buckets", len(cases), &input)
	file, err := os.CreateTemp("", "golean-suffixarray-buckets-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SUFFIXARRAY_BUCKET_CASES", file.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayBuckets",
		"difftest/index/suffixarray/buckets/testdata/buckets_test.go"))
	if !strings.Contains(rows[2], "|ok/nil/3:0.1,1.2,99.1|ok/nil/3:1.2,3.1,99.1") ||
		!strings.HasPrefix(rows[3], "panic/3:1.1,0.2,77.1/") ||
		!strings.Contains(rows[4], "-2147483648.1") || !strings.Contains(rows[7], "-9223372036854775808.1") {
		log.Fatal("alias, partial clearing, or modular-width mutation witness missed")
	}
	fmt.Println("Eight independent native and Lean rows passed, distinguishing cache, alias, partial writes and fixed-width arithmetic.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
