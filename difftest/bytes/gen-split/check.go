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
	s, sep             string
	save, limit, spare int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	for _, s := range []string{"", "a", "aa", "abaaba", " a , b, ", "é\xff\xc0\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "aa", ",", "é", "\xa9", "\xff", "\x00", "absent"} {
			for _, save := range []int{math.MinInt, -2, -1, 0, 1, len(sep), len(sep) + 1, len(s), math.MaxInt} {
				for _, limit := range []int{math.MinInt, -1, 0, 1, 2, 3, math.MaxInt} {
					cases = append(cases, input{s, sep, save, limit, (len(s) + len(sep)) % 5})
				}
			}
		}
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
		fmt.Fprintf(&fixture, "%d %d %d x%x x%x %s\n", c.save, c.limit, c.spare, c.s, c.sep, wants[i])
	}
	for i, got := range check.LeanRows("godiff-bytes-gen-split", len(cases), &leanInput) {
		if got != wants[i] {
			return fmt.Errorf("%+v: Lean=%q want=%q", cases[i], got, wants[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-gen-split-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_GEN_SPLIT_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/gen-split/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanBytesGenSplit", data+"bridge_test.go", data+"gen_split_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.genSplit: %d boundaries and %d samples agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	start := rng.Intn(len(raw) + 1)
	end := start + rng.Intn(min(6, len(raw)-start)+1)
	sep := string(raw[start:end])
	if rng.Intn(3) == 0 {
		sep = strings.Repeat("a", rng.Intn(4))
	}
	if rng.Intn(3) == 0 {
		raw = []byte(sep + string(raw) + sep + sep)
	}
	saves := []int{0, len(sep), rng.Intn(len(raw)+4) - 2, math.MinInt, math.MaxInt}
	limits := []int{-1, 0, 1, 2, rng.Intn(len(raw) + 3), math.MaxInt}
	return input{string(raw), sep, saves[i%len(saves)], limits[rng.Intn(len(limits))], rng.Intn(8)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.limit), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.save), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.sep)), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.spare), 10)
	for _, s := range []string{c.sep, c.s} {
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}

type fragment struct {
	text string
	cap  int
}

func expected(c input) string {
	if c.limit == 0 {
		return "nil"
	}
	if c.sep == "" {
		parts := strings.Split(c.s, "")
		slots := len(c.s)
		if c.limit > 0 {
			slots = min(slots, c.limit)
		}
		if len(parts) > slots {
			parts[slots-1] = strings.Join(parts[slots-1:], "")
			parts = parts[:slots]
		}
		out := make([]fragment, len(parts))
		for i, p := range parts {
			out[i] = fragment{p, len(p)}
		}
		if len(out) != 0 && len(out) == slots {
			out[len(out)-1].cap += c.spare
		}
		return describe(out)
	}
	s := c.s
	padding := strings.Repeat("\xa5", c.spare)
	var parts []fragment
	for c.limit < 0 || len(parts) < c.limit-1 {
		i := strings.Index(s, c.sep)
		if i == -1 {
			break
		}
		if c.save < -i || c.save > len(s)+c.spare-i {
			return "panic"
		}
		end := i + c.save
		parts = append(parts, fragment{(s + padding)[:end], end})
		s = s[i+len(c.sep):]
	}
	return describe(append(parts, fragment{s, len(s) + c.spare}))
}

func describe(parts []fragment) string {
	fields := make([]string, len(parts))
	for i, part := range parts {
		encoded := make([]string, len(part.text))
		for j := range len(part.text) {
			encoded[j] = strconv.Itoa(int(part.text[j]))
		}
		fields[i] = "f" + strconv.Itoa(part.cap) + ":" + strings.Join(encoded, ",")
	}
	return "slice" + strings.Join(fields, ";")
}
