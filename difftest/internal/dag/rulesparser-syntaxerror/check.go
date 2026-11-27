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
	Nil       bool
	Line      int
	Msg, Last []byte
	Want      []byte
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	values := [][]byte{nil, []byte("missing <"), []byte("%d %s"), []byte("a\x00\nb"), []byte("λ🐹"), {0xc0, 0xff}}
	var cases []input
	for _, line := range []int{math.MinInt, -10, -1, 0, 1, 9, 10, 99, 100, math.MaxInt} {
		for _, msg := range values {
			for _, last := range values {
				cases = append(cases, input{Line: line, Msg: msg, Last: last})
			}
		}
	}
	for _, msg := range values {
		cases = append(cases, input{Nil: true, Msg: msg})
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range 20_000 {
			msg, last := make([]byte, rng.Intn(33)), make([]byte, rng.Intn(33))
			rng.Read(msg)
			rng.Read(last)
			cases = append(cases, input{Nil: rng.Intn(16) == 0, Line: int(rng.Uint64()), Msg: msg, Last: last})
		}
	}

	var stdin bytes.Buffer
	for i := range cases {
		c := &cases[i]
		flag := 0
		if !c.Nil {
			flag = 1
			c.Want = strconv.AppendInt([]byte("parsing graph: line "), int64(c.Line), 10)
			c.Want = append(c.Want, ": syntax error: "...)
			c.Want = append(c.Want, c.Msg...)
			c.Want = append(c.Want, " near "...)
			c.Want = append(c.Want, c.Last...)
		}
		fmt.Fprint(&stdin, flag, " ", c.Line, " ", len(c.Msg))
		for _, field := range [][]byte{c.Msg, c.Last} {
			for _, b := range field {
				fmt.Fprint(&stdin, " ", b)
			}
		}
		fmt.Fprintln(&stdin)
	}
	rows := check.LeanRows("godiff-dag-rulesparser-syntaxerror", len(cases), &stdin)
	for i, c := range cases {
		want := []byte("nil")
		if !c.Nil {
			want = []byte("syntax")
			for _, b := range c.Want {
				want = strconv.AppendUint(append(want, ' '), uint64(b), 10)
			}
		}
		if rows[i] != string(want) {
			return fmt.Errorf("case %d %+v: Lean=%q want=%q", i, c, rows[i], want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-dag-syntaxerror-")
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
	if err := os.Setenv("GOLEAN_DAG_SYNTAXERROR_INPUT", path); err != nil {
		return err
	}
	const fixture = "difftest/internal/dag/rulesparser-syntaxerror/testdata/syntax_error_test.go"
	if err := check.TestInPackage("internal/dag", "TestGoLeanRulesSyntaxError", fixture); err != nil {
		return err
	}
	fmt.Printf("internal/dag.rulesParser.syntaxError: %d boundary and %d sampled cases agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}
