package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Before, After rune
	Present       bool
	Op            uint8
	Arg           uint32
	Want          string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-regexp-syntax-empty-width")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	js := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.Present {
				present = 1
			}
			fmt.Fprintln(&rows, c.Before, c.After, present, c.Op, c.Arg)
		}
		for i, got := range check.LeanRows("godiff-regexp-syntax-empty-width", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v Lean=%q want=%q", c, got, c.Want)
			}
			if err := js.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("regexp syntax empty width %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	edges := []fixture{
		{Before: -9, After: -2, Present: true, Op: 3, Arg: 3},
		{Before: -2, After: -1, Present: true, Op: 3, Arg: 260},
		{Before: -1, After: -2, Present: true, Op: 3, Arg: 260},
	}
	runes := []rune{math.MinInt32, -9, -1, 0, 9, 10, 11, '/', '0', '9', ':', '@', 'A', 'Z', '[', '^', '_', '`', 'a', 'z', '{', 0x80, 0x212a, 0x10ffff, math.MaxInt32}
	for _, before := range runes {
		for _, after := range runes {
			for _, arg := range []uint32{0, 1, 2, 3, 4, 8, 16, 32, 64, 260, math.MaxUint32} {
				for _, op := range []uint8{3, 7, 8, 9, 10, 255} {
					edges = append(edges, fixture{Before: before, After: after, Present: true, Op: op, Arg: arg})
				}
			}
		}
	}
	for op := range 256 {
		edges = append(edges, fixture{Before: -1, After: -1, Present: true, Op: uint8(op), Arg: uint32(op)},
			fixture{Before: -1, After: -1, Op: uint8(op), Arg: uint32(op)})
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			before, after := rune(rng.Uint32()), rune(rng.Uint32())
			if i%2 == 0 {
				before = runes[rng.Intn(len(runes))]
			}
			if i%3 == 0 {
				after = runes[rng.Intn(len(runes))]
			}
			arg := rng.Uint32()
			if i%3 == 0 {
				arg = arg&^255 | (1 << rng.Intn(6))
			}
			cases[i] = fixture{Before: before, After: after, Present: i%11 != 0, Op: uint8(rng.Uint32()), Arg: arg}
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_REGEXP_SYNTAX_EMPTY_CASES", file.Name()); err != nil {
		return err
	}
	return check.TestInPackage("regexp/syntax", "TestGoLeanEmptyWidth", "difftest/regexp/syntax/empty-width/testdata/empty_width_test.go")
}

func expected(c fixture) string {
	word := func(r rune) bool {
		return strings.ContainsRune("0123456789_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", r)
	}
	l, r := word(c.Before), word(c.After)
	context := uint8(32)
	if l != r {
		context = 16
	}
	if c.Before <= -1 {
		context |= 5
	} else if c.Before == 10 {
		context |= 1
	}
	if c.After <= -1 {
		context |= 10
	} else if c.After == 10 {
		context |= 2
	}
	width, op := "nil", "nil"
	if c.Present {
		known := map[byte]bool{1: c.Before == 10 || c.Before == -1, 2: c.After == 10 || c.After == -1,
			4: c.Before == -1, 8: c.After == -1, 16: l != r, 32: l == r}
		width = "unknown"
		if matched, ok := known[byte(c.Arg)]; ok {
			width = fmt.Sprint(matched)
		}
		code := c.Op
		if 8 <= code && code <= 10 {
			code = 7
		}
		op = fmt.Sprint(code)
	}
	return fmt.Sprintf("%t|%t|%d|%s|%s", l, r, context, width, op)
}
