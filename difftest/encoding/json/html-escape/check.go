package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"

	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Present      bool
	Backing, Src []byte
	Length       int
	Want         string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-encoding-json-html-escape")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	js := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			present, dst := 0, 0
			if c.Present {
				present = 1
			}
			if c.Backing != nil {
				dst = 1
			}
			fmt.Fprint(&rows, present, " ", dst, " ", c.Length, " ", len(c.Backing))
			for _, block := range [][]byte{c.Backing, c.Src} {
				for _, b := range block {
					fmt.Fprint(&rows, " ", b)
				}
			}
			fmt.Fprintln(&rows)
		}
		for i, got := range check.LeanRows("godiff-encoding-json-html-escape", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("row %+v Lean=%q want=%q", c, got, c.Want)
			}
			if err := js.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("JSON HTML escape %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	var edges []fixture
	tokens := [][]byte{nil, {}, []byte("notes"), {60}, {62}, {38}, {226, 128, 168}, {226, 128, 169}, {226, 128}, {128, 169}, {226, 128, 170}, {255}}
	for _, a := range tokens {
		for _, b := range tokens {
			src := append(slices.Clone(a), b...)
			for length := range 3 {
				for _, spare := range []int{0, 1, 2, 5, 6, 7, 12, 24} {
					for _, live := range []bool{false, true} {
						edges = append(edges, fixture{Present: live, Backing: bytes.Repeat([]byte{99}, length+spare), Src: src, Length: length})
					}
				}
			}
			edges = append(edges, fixture{Present: true, Src: src})
		}
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			backing := make([]byte, rng.Intn(40))
			rng.Read(backing)
			var src []byte
			for range rng.Intn(20) {
				if rng.Intn(3) == 0 {
					src = append(src, byte(rng.Uint32()))
				} else {
					src = append(src, tokens[rng.Intn(len(tokens))]...)
				}
			}
			cases[i] = fixture{Present: i%7 != 0, Backing: backing, Src: src, Length: rng.Intn(len(backing) + 1)}
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_JSON_HTML_CASES", file.Name()); err != nil {
		return err
	}
	const overlay = "difftest/encoding/json/html-escape/testdata/"
	return check.TestInPackage("encoding/json", "TestGoLeanHTMLEscape", overlay+"bridge_test.go", overlay+"html_test.go")
}

func expected(c fixture) string {
	type span struct {
		start, end int
		escaped    string
	}
	var spans []span
	for i := 0; i < len(c.Src); i++ {
		switch c.Src[i] {
		case '<', '>', '&':
			spans = append(spans, span{i, i + 1, fmt.Sprintf("\\u%04x", c.Src[i])})
		case 0xe2:
			if i+2 < len(c.Src) && c.Src[i+1] == 0x80 && (c.Src[i+2] == 0xa8 || c.Src[i+2] == 0xa9) {
				spans = append(spans, span{i, i + 3, fmt.Sprintf("\\u%04x", 0x2028+int(c.Src[i+2]-0xa8))})
				i += 2
			}
		}
	}
	var encoded []byte
	old := slices.Clone(c.Backing)
	at, stored := 0, c.Length
	appendPart := func(part []byte) {
		encoded = append(encoded, part...)
		if stored+len(part) <= len(old) {
			copy(old[stored:], part)
			stored += len(part)
		} else {
			stored = len(old) + 1
		}
	}
	for _, span := range spans {
		appendPart(c.Src[at:span.start])
		appendPart([]byte(span.escaped))
		at = span.end
	}
	appendPart(c.Src[at:])
	output := append(slices.Clone(c.Backing[:c.Length]), encoded...)
	public := "panic"
	if c.Present {
		if output == nil {
			public = "slice:"
		} else {
			public = wire(output)
		}
	}
	return fmt.Sprintf("%s/%s/%t|%s", wire(output), wire(old), len(old) > 0 && stored <= len(old), public)
}

func wire(b []byte) string {
	if b == nil {
		return "nil"
	}
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = strconv.Itoa(int(v))
	}
	return "slice:" + strings.Join(parts, ",")
}
