package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	seconds     int64
	nanos       int
	offset      int
	date, clock uint16
}

const yearOne = -62135596800

func expected(c fixture) string {
	if c.seconds != yearOne || c.nanos != 0 {
		return fmt.Sprintf("%d:%d", c.seconds, c.nanos)
	}
	d, t := int(c.date), int(c.clock)
	day := time.Date(1980+d/512, time.Month(d/32%16), d%32, 0, 0, 0, 0, time.UTC)
	seconds := day.Unix() + int64(t/2048*3600+t/32%64*60+t%32*2)
	return fmt.Sprintf("%d:0", seconds)
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.nanos, uint64(c.seconds-yearOne), c.date, c.clock)
	}
	for i, row := range check.LeanRows("godiff-zip-header-modtime", len(cases), &input) {
		c := cases[i]
		modern := time.Unix(c.seconds, int64(c.nanos)).In(time.FixedZone("offset", c.offset))
		info := (&zip.FileHeader{Modified: modern, ModifiedDate: c.date, ModifiedTime: c.clock}).FileInfo()
		result := info.ModTime()
		got := fmt.Sprintf("%d:%d", result.Unix(), result.Nanosecond())
		want := expected(c)
		if row != want || got != want || result.Location() != time.UTC {
			log.Fatalf("%s row %d: Go %q Lean %q want %q", label, i, got, row, want)
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean rows agree\n", label, len(cases))
}

func main() {
	var fixed []fixture
	for _, sec := range []int64{math.MinInt64, yearOne - 1, yearOne, yearOne + 1, -1, 0, 1, math.MaxInt64} {
		for _, nanos := range []int{0, 1, 999999999} {
			for _, off := range []int{-43200, -1, 0, 1, 50400} {
				for _, word := range [][2]uint16{{0, 0}, {33, 0}, {0xffff, 0xffff}, {32, 31}, {63, 32}} {
					fixed = append(fixed, fixture{sec, nanos, off, word[0], word[1]})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		rows := make([]fixture, 20_000)
		for i := range rows {
			c := fixture{int64(r.Uint64()), r.Intn(1_000_000_000), r.Intn(100001) - 50000, uint16(r.Uint32()), uint16(r.Uint32())}
			if i%4 == 0 {
				c.seconds, c.nanos = yearOne, 0
			}
			rows[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), rows)
	}
	// A native monotonic observation checks the UTC conversion independently.
	now := time.Now()
	header := zip.FileHeader{Modified: now}
	if got := header.FileInfo().ModTime(); got != now.UTC() {
		log.Fatal("modern time did not match UTC with its monotonic reading stripped")
	}
}
