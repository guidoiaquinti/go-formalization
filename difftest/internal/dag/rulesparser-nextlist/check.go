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
	"regexp"
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
	edges := []string{"", "a", "a,b", "a,b,c<d", "a,b!<c", "a,b;", "a,b c,d", "a!,b",
		"!", "!x", "!<", ";", ",", "<", "a,", "a,,", "a,<", "a,!<", "a,;",
		"a,b,\n# comment", "a, # note\n b ; rest", " # hi\n\t!x",
		"#eof", "#x\n\nfoo", "\tfoo\x00\xff#x", "a\r\nb", "a,\x00\xff\n!"}
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
			switch i % 6 {
			case 1:
				text = append([]byte(" \t\n#comment\n"), text...)
			case 2:
				text = append([]byte("#comment\n"), []byte(edges[rng.Intn(len(edges))])...)
			case 3:
				text = append([]byte(" \t\n"), []byte(edges[rng.Intn(len(edges))])...)
			case 4, 5:
				var names bytes.Buffer
				for n := 1 + rng.Intn(12); n > 0; n-- {
					fmt.Fprintf(&names, "name%d", rng.Uint32())
					if n > 1 || i%6 == 5 {
						names.WriteString(", # between\n")
					}
				}
				text = append(names.Bytes(), text...)
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
	rows := check.LeanRows("godiff-dag-nextlist", len(cases), &stdin)
	for i, row := range rows {
		if c := cases[i]; row != c.Expected {
			return fmt.Errorf("case %d line %d last %q text %q: Lean=%q expected=%q", i, c.Line, c.Last, c.Text, row, c.Expected)
		}
	}
	dir, err := os.MkdirTemp("", "golean-dag-nextlist-")
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
	if err := os.Setenv("GOLEAN_DAG_NEXTLIST_INPUT", path); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/rulesparser-nextlist/testdata/nextlist_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanNextList", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.rulesParser.nextList: %d boundary and %d sampled cases agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

// Anchored expressions tokenize only the prefix reached by the list grammar.
// Byte offsets preserve the original data even for malformed UTF-8.
var trivia = regexp.MustCompile(`^(?:[ \t\n]|#[^\n]*)*`)
var token = regexp.MustCompile(`^(?:!<|[,;<]|[^!;,<#\n \t]+)`)

func expected(c input) string {
	if !c.Present {
		return "nil"
	}
	line, last, rest := c.Line, c.Last, c.Text
	syntax := func(msg string) string {
		payload := fmt.Appendf(nil, "parsing graph: line %d: syntax error: %s near %s", line, msg, last)
		return fmt.Sprintf("syntax %s %d %s %s", encode(payload), line, encode(last), encode(rest))
	}
	var names [][]byte
	for index := 0; ; index++ {
		prefix := trivia.Find(rest)
		line += bytes.Count(prefix, []byte{'\n'})
		rest = rest[len(prefix):]
		tok := token.Find(rest)
		if len(tok) == 0 && len(rest) != 0 {
			return syntax("unexpected token !")
		}
		rest = rest[len(tok):]
		separator := len(tok) == 0 || (len(tok) == 1 && bytes.IndexByte([]byte(",;<"), tok[0]) >= 0) || string(tok) == "!<"
		if !separator {
			last = tok
		}
		if index%2 == 0 {
			if len(tok) == 0 && index == 0 {
				return fmt.Sprintf("return 0 0 %d %s %s", line, encode(last), encode(rest))
			}
			if separator {
				return syntax("bad list syntax")
			}
			names = append(names, tok)
		} else if string(tok) != "," {
			encoded := strconv.AppendInt(nil, int64(len(names)), 10)
			for _, name := range names {
				encoded = append(append(encoded, ' '), encode(name)...)
			}
			return fmt.Sprintf("return %s %s %d %s %s", encoded, encode(tok), line, encode(last), encode(rest))
		}
	}
}

func encode(b []byte) string { return strconv.Itoa(len(b)) + byteFields(b) }

func byteFields(b []byte) string {
	var out []byte
	for _, x := range b {
		out = strconv.AppendUint(append(out, ' '), uint64(x), 10)
	}
	return string(out)
}
