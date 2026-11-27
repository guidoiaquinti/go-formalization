package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s            string
	limit, spare int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	add := func(s string) {
		for _, limit := range []int{math.MinInt, -2, -1, 0, 1, 2, 3, 4, len(s), len(s) + 1, math.MaxInt} {
			cases = append(cases, input{s, limit, 3})
		}
	}
	for _, s := range []string{"", "ascii", "é€🚀", "a\x00z", "\xc0\x80", "\xe2\x82", "\xf4\x90\x80\x80", "\xed\xa0\x80", "�\xff"} {
		add(s)
	}
	for b := range 256 {
		lead := byte(b)
		add(string([]byte{lead}))
		add(string([]byte{lead, 0x80, 0x80, 0x80, 'a'}))
		add(string([]byte{lead, 0xbf, 0xbf, 0xbf, 'a'}))
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}
	var leanInput, fixture bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		fmt.Fprintln(&leanInput, encodeInput(c))
		wants[i] = expected(c)
		fmt.Fprintf(&fixture, "%d %d x%x %s\n", c.limit, c.spare, c.s, wants[i])
	}
	for i, got := range check.LeanRows("godiff-bytes-explode", len(cases), &leanInput) {
		if got != wants[i] {
			return fmt.Errorf("%+v: Lean=%q want=%q", cases[i], got, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-explode-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_EXPLODE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/explode/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanExplode", data+"bridge_test.go", data+"explode_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.explode: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func sample(rng *rand.Rand, i int) input {
	chunks := []string{"a", "é", "€", "🚀", "�", "\xc2", "\xe2\x82", "\xff"}
	var s strings.Builder
	for range rng.Intn(65) {
		if rng.Intn(2) == 0 {
			s.WriteByte(byte(rng.Intn(256)))
		} else {
			s.WriteString(chunks[rng.Intn(len(chunks))])
		}
	}
	limits := []int{math.MinInt, -1, 0, 1, 2, rng.Intn(s.Len() + 3), math.MaxInt}
	return input{s.String(), limits[i%len(limits)], rng.Intn(8)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.limit), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.spare), 10)
	for i := range len(c.s) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.s[i]), 10)
	}
	return string(row)
}

func expected(c input) string {
	// String range supplies rune boundaries while retaining the original bytes.
	var offsets []int
	for offset := range c.s {
		offsets = append(offsets, offset)
	}
	limit := len(offsets)
	if c.limit > 0 {
		limit = min(c.limit, limit)
	}
	parts := make([]string, limit)
	for i := 0; i+1 < limit; i++ {
		parts[i] = c.s[offsets[i]:offsets[i+1]]
	}
	if limit > 0 {
		parts[limit-1] = c.s[offsets[limit-1]:]
	}
	slots := len(c.s)
	if c.limit > 0 {
		slots = min(slots, c.limit)
	}
	return describe(parts, slots, c.spare)
}

func describe(parts []string, slots, spare int) string {
	fields := make([]string, len(parts))
	for i, part := range parts {
		encoded := make([]string, len(part))
		for j := range len(part) {
			encoded[j] = strconv.Itoa(int(part[j]))
		}
		capacity := len(part)
		if i+1 == len(parts) && len(parts) == slots {
			capacity += spare
		}
		fields[i] = "f" + strconv.Itoa(capacity) + ":" + strings.Join(encoded, ",")
	}
	return "slice" + strings.Join(fields, ";")
}
