package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const yearOne = -62135596800
	zone := time.FixedZone("east", 19800)
	legacy := time.Unix(315532800, 0).UTC()
	cases := []struct {
		modern, want time.Time
		date, clock  uint16
	}{
		{time.Time{}, legacy, 33, 0},
		{time.Time{}.In(zone), time.Unix(315576002, 0).UTC(), 33, 12<<11 | 1},
		{time.Unix(yearOne, 1).In(zone), time.Unix(yearOne, 1).UTC(), 33, 0},
		{time.Unix(0, 987).In(zone), time.Unix(0, 987).UTC(), 33, 0},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.modern.Nanosecond(), uint64(c.modern.Unix()-yearOne), c.date, c.clock)
	}
	var rejected [4]int
	for i, got := range check.LeanRows("godiff-zip-header-modtime", len(cases), &input) {
		c := cases[i]
		h := zip.FileHeader{Modified: c.modern, ModifiedDate: c.date, ModifiedTime: c.clock}
		result := h.FileInfo().ModTime()
		want := fmt.Sprintf("%d:%d", c.want.Unix(), c.want.Nanosecond())
		if got != want || result != c.want || result.Location() != time.UTC {
			log.Fatalf("row %d: Go %v Lean %q want %v", i, result, got, c.want)
		}
		legacyOnly, secondsOnly, withoutUTC := h.ModTime(), result, result
		if c.modern.Unix() == yearOne {
			secondsOnly = legacyOnly
		}
		if !c.modern.IsZero() {
			withoutUTC = c.modern
		}
		for j, wrong := range []time.Time{legacyOnly, secondsOnly, withoutUTC} {
			if wrong != result {
				rejected[j]++
			}
		}
	}
	now := time.Now()
	h := zip.FileHeader{Modified: now}
	got := h.FileInfo().ModTime()
	if got != now.UTC() || got.Location() != time.UTC {
		log.Fatal("modern UTC result differs from now.UTC")
	}
	if want := [3]int{2, 1, 2}; [3]int(rejected[:3]) != want {
		log.Fatalf("mutation counts %v", rejected)
	}
	// Round(0) strips monotonic data without changing the original location.
	hasReading := now != now.Round(0)
	if hasReading {
		if got != now.Round(0).UTC() {
			log.Fatal("monotonic stripping changed the instant")
		}
		rejected[3] = 1
	}
	fmt.Printf("ZIP ModTime: %d pinned Go and Lean rows, mutation counts %v, native monotonic present=%t\n",
		len(cases), rejected, hasReading)
}
