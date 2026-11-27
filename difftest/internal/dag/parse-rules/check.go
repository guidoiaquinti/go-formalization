package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Text     []byte
	Expected string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	edges := []string{"", "a", "a;", "a<", "a<b", "a<b;", "a<b!<c;",
		"a,b<c,d;", "a<b<c,d;", "a<b; c!<d;", "a<b; c,d", "a<b; c<d e",
		"a<b; c<d,", "a b", "a<b c", "a<;", "a<<b;", ";", "!<",
		"a<\n# note\nb;", "a<\n# note\n!", "# tail", "a<b;\n# tail",
		"\xff,\x00<\xc0\x80;", "a\r\nb;", "a<b;c<d!"}
	var cases []input
	for _, s := range edges {
		cases = append(cases, input{Text: []byte(s)})
	}
	for b := range 256 {
		cases = append(cases, input{Text: []byte{byte(b)}},
			input{Text: append([]byte("a<b; c<"), byte(b))})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			var text []byte
			switch i % 5 {
			case 0:
				text = make([]byte, rng.Intn(65))
				rng.Read(text)
			case 1:
				text = []byte(edges[rng.Intn(len(edges))] + edges[rng.Intn(len(edges))])
			default:
				for n := 1 + rng.Intn(8); n > 0; n-- {
					for j := 1 + rng.Intn(4); j > 0; j-- {
						fmtToken := "n%d"
						if rng.Intn(3) == 0 {
							fmtToken = "n%d, z"
						}
						text = fmt.Appendf(text, fmtToken, rng.Uint32())
						if j > 1 {
							text = append(text, []string{"<", "!<", " # trivia\n<"}[rng.Intn(3)]...)
						}
					}
					text = append(text, ';')
				}
				if i%5 == 3 {
					text = text[:rng.Intn(len(text))]
				} else if i%5 == 4 {
					text = append(text, edges[rng.Intn(len(edges))]...)
				}
			}
			cases = append(cases, input{Text: text})
		}
	}

	var stdin bytes.Buffer
	for i := range cases {
		cases[i].Expected = expected(cases[i].Text)
		fmt.Fprintln(&stdin, strconv.Itoa(strconv.IntSize)+byteFields(cases[i].Text))
	}
	for i, got := range check.LeanRows("godiff-dag-parse-rules", len(cases), &stdin) {
		if got != cases[i].Expected {
			return fmt.Errorf("case %d text %q: Lean=%q expected=%q", i, cases[i].Text, got, cases[i].Expected)
		}
	}
	dir, err := os.MkdirTemp("", "golean-dag-parse-rules-")
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
	if err := os.Setenv("GOLEAN_DAG_PARSE_RULES_INPUT", path); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/parse-rules/testdata/parse_rules_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanParseRules", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.parseRules: %d boundary and %d sampled strings agree, seeds %v.\n", boundaries, len(cases)-boundaries, seeds)
	return nil
}

var trivia = regexp.MustCompile(`^(?:[ \t\n]|#[^\n]*)*`)
var lexeme = regexp.MustCompile(`^(?:!<|[,;<]|[^!;,<#\n \t]+)`)

func expected(text []byte) string {
	var emitted []string
	line := 1
	var last []byte
	scan := func() (string, bool, bool) {
		prefix := trivia.Find(text)
		line += bytes.Count(prefix, []byte{'\n'})
		text = text[len(prefix):]
		token := lexeme.Find(text)
		if len(token) == 0 {
			return "", false, len(text) != 0
		}
		text = text[len(token):]
		name := !strings.ContainsAny(string(token[:1]), ",;<!")
		if name {
			last = token
		}
		return string(token), name, false
	}
	finish := func(message string) string {
		kind := "ok"
		if message != "" {
			kind = "error"
		}
		row := kind + " " + strconv.Itoa(len(emitted))
		for _, r := range emitted {
			row += " " + r
		}
		if message != "" {
			payload := fmt.Appendf(nil, "parsing graph: line %d: syntax error: %s near %s", line, message, last)
			row += " " + encode(payload)
		}
		return row
	}
	var prior, pending string
	for {
		name, valid, bad := scan()
		switch {
		case bad:
			return finish("unexpected token !")
		case name == "":
			if prior != "" {
				return finish("unexpected EOF")
			}
			return finish("")
		case !valid:
			return finish("bad list syntax")
		}
		names := []string{name}
		tok, _, bad := scan()
		for tok == "," && !bad {
			name, valid, bad = scan()
			if bad {
				break
			}
			if !valid {
				return finish("bad list syntax")
			}
			names = append(names, name)
			tok, _, bad = scan()
		}
		if bad {
			return finish("unexpected token !")
		}
		if tok == "" {
			if prior == "" {
				return finish("")
			}
			return finish("unexpected EOF")
		}
		list := strconv.Itoa(len(names))
		for _, name := range names {
			list += " " + encode([]byte(name))
		}
		if prior != "" {
			emitted = append(emitted, prior+" "+encode([]byte(pending))+" "+list)
		}
		switch tok {
		case ";":
			prior, pending = "", ""
		case "<", "!<":
			prior, pending = list, tok
		default:
			return finish("missing <")
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
