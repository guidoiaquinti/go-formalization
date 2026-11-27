package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const step = int64(15 * time.Minute)

func expected(value int64) int64 {
	// Outside the furthest useful half-step, rounding or saturation
	// cannot bring an input into the accepted range.
	if value <= -int64(12*time.Hour)-step/2 || value >= int64(14*time.Hour)+step/2 {
		return 0
	}
	sign := int64(1)
	if value < 0 {
		sign, value = -1, -value
	}
	units := (value + step/2) / step
	return sign * units * 900
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-zone")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	out := bufio.NewWriter(f)
	compare := func(label string, offsets []int64) {
		var input bytes.Buffer
		for _, x := range offsets {
			fmt.Fprintln(&input, x)
		}
		for i, row := range check.LeanRows("godiff-zip-time-zone", len(offsets), &input) {
			x := offsets[i]
			offset := expected(x)
			want := fmt.Sprintf("%d %t", offset, offset%3600 == 0)
			if row != want {
				log.Fatalf("%s at %d: Lean %q want %q", label, x, row, want)
			}
			_, err := fmt.Fprintln(out, x, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(offsets))
	}
	fixed := []int64{math.MinInt64, math.MinInt64 + 1, math.MaxInt64 - 1, math.MaxInt64}
	for units := -60; units <= 64; units++ {
		for _, delta := range []int64{-step/2 - 1, -step / 2, -step/2 + 1, -1, 0, 1, step/2 - 1, step / 2, step/2 + 1} {
			fixed = append(fixed, int64(units)*step+delta)
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		values := make([]int64, 20_000)
		for i := range values {
			values[i] = int64(r.Uint64())
			if i%2 == 0 {
				values[i] = r.Int63n(int64(32*time.Hour)) - int64(16*time.Hour)
			}
		}
		compare(fmt.Sprintf("seed %d", seed), values)
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_ZONE_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipZone", "difftest/archive/zip/time-zone/testdata/zone_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
