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
	Shape, Variant int
	Bytes          [4][]byte
	Want           string
}

func main() {
	cases := []fixture{
		{0, 0, [4][]byte{}, "true|2|0|:leaf:|0|:leaf:"},
		{1, 0, [4][]byte{{255, 0}, {'x'}, {195, 169}, {'Z'}}, "false|2|6|255,0,120,195,169,90:leaf:255,0,120,195,169,90|6|255,0,120,195,169,90:leaf:255,0,120,195,169,90"},
		{3, 0, [4][]byte{{'a'}, {'b'}}, "false|2|4|97,98,97,98:leaf:97,98,97,98|4|97,98,97,98:leaf:97,98,97,98"},
		{2, 1, [4][]byte{{'a'}, nil, {'b'}}, "false|2|2|97,98:leaf:97,98|0|:unknown"},
		{0, 2, [4][]byte{{'a'}}, "false|2|1|97:leaf:97|panic:true not a String|panic:true not a String"},
		{0, 3, [4][]byte{{'a'}}, "false|2|1|97:leaf:97|panic:<nil> not a String|panic:<nil> not a String"},
		{0, 4, [4][]byte{}, "true|2|nilpanic|nilpanic|nilpanic|nilpanic"},
	}
	f, err := os.CreateTemp("", "golean-constant-string-control-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.Shape, " ", c.Variant)
		for _, val := range c.Bytes {
			fmt.Fprint(&input, " ", len(val))
			for _, b := range val {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [6]int
	for i, got := range check.LeanRows("godiff-go-constant-string-values", len(cases), &input) {
		c := cases[i]
		if c.Want != got {
			log.Fatalf("row %d: Lean %q want %q", i, got, c.Want)
		}
		fields := strings.Split(got, "|")
		wrong := [6]string{"false", fields[1], fields[2], fields[3], fields[4], fields[5]}
		if c.Variant == 4 {
			wrong[1] = "nilpanic"
		}
		if c.Shape == 3 {
			wrong[2] = "2"
		}
		if c.Shape != 0 {
			wrong[3] = "uncached"
		}
		if c.Variant == 1 {
			wrong[4] = "panic:unknown not a String"
		}
		if c.Variant == 0 && c.Shape != 0 {
			wrong[5] = "uncached"
		}
		for j, w := range wrong {
			if w != fields[j] {
				rejected[j]++
			}
		}
	}
	if want := [6]int{2, 1, 1, 3, 1, 2}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_CONSTANT_STRING_ROWS", f.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanStringValues", "difftest/go/constant/string-values/testdata/string_test.go"))
	fmt.Printf("constant strings: %d pinned Go and Lean rows, mutations %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
