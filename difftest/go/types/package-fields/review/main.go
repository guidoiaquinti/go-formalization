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
	Mode                             int
	Complete                         bool
	Path, Name, Version, Replacement []byte
	Want                             string
}

func main() {
	cases := []fixture{
		{0, false, []byte("x"), []byte("n"), nil, nil, "120|110|31|false|#panic|panic|0|panic|panic#panic#panic"},
		{1, false, []byte("x"), []byte("n"), []byte("v"), []byte("r"), "120|110|31|false|#120|110|31|false|118#120|114|31|false|118#120|110|31|true|118"},
		{2, true, []byte{255, 0}, []byte("m"), []byte{128}, nil, "255,0|109|31|false|#255,0|109|nil|true|128#255,0||nil|true|128#255,0|109|nil|true|128"},
		{1, true, nil, nil, nil, []byte("a"), "||31|false|#||31|true|#|97|31|true|#||31|true|"},
	}
	f, err := os.CreateTemp("", "golean-package-control-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		complete := 0
		if c.Complete {
			complete = 1
		}
		fmt.Fprint(&input, c.Mode, " ", complete, " ", len(c.Path), " ", len(c.Name), " ", len(c.Version))
		for _, field := range [][]byte{c.Path, c.Name, c.Version, c.Replacement} {
			for _, b := range field {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [8]int
	for i, got := range check.LeanRows("godiff-go-types-package-fields", len(cases), &input) {
		c := cases[i]
		if c.Want != got {
			log.Fatalf("row %d: Lean %q want %q", i, got, c.Want)
		}
		for j, wrong := range mutants(got, c.Mode) {
			if wrong != got {
				rejected[j]++
			}
		}
	}
	if want := [8]int{4, 3, 3, 4, 3, 3, 4, 2}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_TYPES_PACKAGE_ROWS", f.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanPackageFields", "difftest/go/types/package-fields/testdata/package_test.go"))
	fmt.Printf("go/types Package review: %d pinned Go and Lean rows, mutations %v\n", len(cases), rejected)
}

func mutants(s string, mode int) [8]string {
	replace := func(group, field int, value string) string {
		groups := strings.Split(s, "#")
		fields := strings.Split(groups[group], "|")
		fields[field] = value
		groups[group] = strings.Join(fields, "|")
		return strings.Join(groups, "#")
	}
	groups := strings.Split(s, "#")
	scope := "0"
	if mode == 0 {
		scope = "nil"
	}
	omitName, omitCompletion := s, s
	if mode == 0 {
		omitName = strings.Join([]string{groups[0], groups[1], "ignored", groups[3]}, "#")
		omitCompletion = strings.Join([]string{groups[0], groups[1], groups[2], "ignored"}, "#")
	} else {
		before := strings.Split(groups[1], "|")
		omitName = replace(2, 1, before[1])
		omitCompletion = replace(3, 3, before[3])
	}
	return [8]string{replace(0, 3, "true"), replace(1, 0, ""), replace(1, 1, ""), replace(1, 2, scope),
		replace(1, 3, "false"), replace(1, 4, ""), omitName, omitCompletion}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
