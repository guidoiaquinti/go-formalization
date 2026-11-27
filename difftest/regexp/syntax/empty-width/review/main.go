package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	cases := []fixture{
		{-9, -2, true, 8, 260, "false|false|47|false|7"},
		{-1, -2, true, 9, 260, "false|false|47|true|7"},
		{10, '_', true, 10, 16, "false|true|17|true|7"},
		{0x212a, 0xe9, true, 255, 3, "false|false|32|unknown|255"},
		{'A', -1, true, 2, 8, "true|false|26|true|2"},
		{'0', 10, false, 0, 0, "true|false|18|nil|nil"},
		{0, 0, true, 7, 32, "false|false|32|true|7"},
	}
	mutants := []string{
		"false|false|47|true|7", "false|false|47|unknown|7", "false|true|17|true|10",
		"false|false|32|true|255", "true|false|24|true|2", "true|false|18|false|nil",
		"false|false|32|false|7",
	}
	var in bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.Present {
			present = 1
		}
		fmt.Fprintf(&in, "%d %d %d %d %d\n", c.Before, c.After, present, c.Op, c.Arg)
	}
	rejects := 0
	for i, got := range check.LeanRows("godiff-regexp-syntax-empty-width", len(cases), &in) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].Want)
		}
		if got != mutants[i] {
			rejects++
		}
		fields := strings.Split(got, "|")
		if i == 2 && fields[1] != "false" {
			rejects++
		}
		if i == 3 && fields[0] != "true" {
			rejects++
		}
	}
	if rejects != len(cases)+2 {
		log.Fatalf("wrong rejection count %d", rejects)
	}
	f, err := os.CreateTemp("", "golean-regexp-empty-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_REGEXP_SYNTAX_EMPTY_CASES", f.Name()))
	must(check.TestInPackage("regexp/syntax", "TestGoLeanEmptyWidth", "difftest/regexp/syntax/empty-width/testdata/empty_width_test.go"))
	fmt.Printf("Regexp empty width: %d independent paired cases, %d altered outcomes rejected\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
