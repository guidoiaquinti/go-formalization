package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Present    bool
	Line       int
	Last, Text []byte
	Expected   string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []input{{}}
	edges := []string{"", "a", "a!<b", "!", "!x", "!<", ";", ",", "<", " # hi\n\t!x", "#eof",
		"#x\n\nfoo", "\tfoo\x00\xff#x", "a\r\nb", "\n!<", "\x00\xff\n!"}
	for i := range 256 {
		edges = append(edges, string([]byte{byte(i)}))
	}
	for _, line := range []int{math.MinInt, 0, 1, math.MaxInt} {
		for _, s := range edges {
			cases = append(cases, input{true, line, []byte("old\x00\xff"), []byte(s), ""})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			line := int(rng.Uint64())
			last := make([]byte, rng.Intn(17))
			rng.Read(last)
			text := make([]byte, rng.Intn(65))
			rng.Read(text)
			switch i % 4 {
			case 1:
				text = append([]byte(" \t\n#comment\n"), text...)
			case 2:
				text = append([]byte("#comment\n"), []byte(edges[rng.Intn(len(edges))])...)
			case 3:
				text = append([]byte(" \t\n"), []byte(edges[rng.Intn(len(edges))])...)
			}
			cases = append(cases, input{true, line, last, text, ""})
		}
	}
	var stdin bytes.Buffer
	for i, c := range cases {
		cases[i].Expected = expected(c)
		flag := 0
		if c.Present {
			flag = 1
		}
		fmt.Fprintln(&stdin, strconv.IntSize, flag, c.Line, encode(c.Last)+byteFields(c.Text))
	}
	rows := check.LeanRows("godiff-dag-nexttoken", len(cases), &stdin)
	for i, row := range rows {
		if c := cases[i]; row != c.Expected {
			return fmt.Errorf("case %d line %d last %q text %q: Lean=%q expected=%q", i, c.Line, c.Last, c.Text, row, c.Expected)
		}
	}
	dir, err := os.MkdirTemp("", "golean-dag-nexttoken-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	data, err := json.Marshal(cases)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "input.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_DAG_NEXTTOKEN_INPUT", path); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/rulesparser-nexttoken/testdata/nexttoken_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanNextToken", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.rulesParser.nextToken: %d boundary and %d sampled cases agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

func expected(c input) string {
	if !c.Present {
		return "nil"
	}
	i, line, comment := 0, c.Line, false
scan:
	for i < len(c.Text) {
		b := c.Text[i]
		switch {
		case b == '\n':
			line++
			comment = false
		case comment:
		case b == '#':
			comment = true
		case b == ' ' || b == '\t':
		default:
			break scan
		}
		i++
	}
	text, last := c.Text[i:], c.Last
	format := func(kind string, tok []byte, end int) string {
		return fmt.Sprintf("%s %s %d %s %s", kind, encode(tok), line, encode(last), encode(text[end:]))
	}
	if len(text) == 0 {
		return format("return", nil, 0)
	}
	switch text[0] {
	case ';', ',', '<':
		return format("return", text[:1], 1)
	case '!':
		if len(text) > 1 && text[1] == '<' {
			return format("return", text[:2], 2)
		}
		msg := fmt.Appendf(nil, "parsing graph: line %d: syntax error: unexpected token ! near %s", line, last)
		return format("syntax", msg, 0)
	}
	n := 0
word:
	for n < len(text) {
		switch text[n] {
		case '!', ';', ',', '<', '#', '\n', ' ', '\t':
			break word
		}
		n++
	}
	last = text[:n]
	return format("return", last, n)
}

func encode(b []byte) string { return strconv.Itoa(len(b)) + byteFields(b) }

func byteFields(b []byte) string {
	var out []byte
	for _, x := range b {
		out = strconv.AppendUint(append(out, ' '), uint64(x), 10)
	}
	return string(out)
}
