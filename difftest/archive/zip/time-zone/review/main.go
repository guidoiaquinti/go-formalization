package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct {
		offset  time.Duration
		seconds int
		cached  bool
	}{
		{math.MinInt64, 0, true},
		{math.MaxInt64, 0, true},
		{-12*time.Hour - 450*time.Second, 0, true},
		{-12*time.Hour - 450*time.Second + 1, -43200, true},
		{14*time.Hour + 450*time.Second - 1, 50400, true},
		{14*time.Hour + 450*time.Second, 0, true},
		{-450*time.Second - 1, -900, false},
		{-450 * time.Second, -900, false},
		{-450*time.Second + 1, 0, true},
		{450*time.Second - 1, 0, true},
		{450 * time.Second, 900, false},
		{-time.Hour + 449*time.Second, -3600, true},
		{-time.Hour + 450*time.Second + 1, -2700, false},
		{5*time.Hour + 45*time.Minute, 20700, false},
	}
	f, err := os.CreateTemp("", "golean-zip-zone-review")
	must(err)
	defer os.Remove(f.Name())
	var input bytes.Buffer
	mutations := 0
	for _, c := range fixed {
		fmt.Fprintln(&input, int64(c.offset))
		_, err := fmt.Fprintln(f, int64(c.offset), c.seconds, c.cached)
		must(err)
		if (c.offset < -12*time.Hour || c.offset > 14*time.Hour) && c.seconds != 0 {
			mutations++
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-zip-time-zone", len(fixed), &input) {
		want := fmt.Sprintf("%d %t", fixed[i].seconds, fixed[i].cached)
		if got != want {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, want)
		}
	}
	must(os.Setenv("GOLEAN_ZIP_ZONE_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipZone", "difftest/archive/zip/time-zone/testdata/zone_test.go"))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipZoneReview", "difftest/archive/zip/time-zone/review/testdata/zone_review_test.go"))
	if mutations != 2 {
		log.Fatalf("clamp-before-rounding mutation disagrees on %d rows, want 2", mutations)
	}
	fmt.Printf("%d literal Go/Lean controls; clamp-before-rounding disagrees on %d rows; native cache identity passed\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
