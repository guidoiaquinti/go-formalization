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
	Present    bool
	Kind, Info int
	Name       []byte
	Want       string
}

func main() {
	cases := []fixture{
		{false, 0, 0, nil, "panic|panic|panic|typedNil|panic"},
		{true, 0, 0, nil, "0|0||same|"},
		{true, -7, 91, []byte("Pointer"), "-7|91|80,111,105,110,116,101,114|same|117,110,115,97,102,101,46,80,111,105,110,116,101,114"},
		{true, 11, -3, []byte("Sizeof"), "11|-3|83,105,122,101,111,102|same|83,105,122,101,111,102"},
		{true, 22, 0, []byte{255, 0, 'A'}, "22|0|255,0,65|same|255,0,65"},
	}
	f, err := os.CreateTemp("", "golean-basic-fields-control-")
	must(err)
	defer os.Remove(f.Name())
	jsonRows := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.Present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", c.Kind, " ", c.Info)
		for _, b := range c.Name {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(jsonRows.Encode(c))
	}
	must(f.Close())
	var rejected [5]int
	for i, got := range check.LeanRows("godiff-go-types-basic-fields", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
		fields := strings.Split(got, "|")
		wrong := [5]string{"0", "0", "", fields[3], fields[2]}
		if !c.Present {
			wrong[3] = "nil"
		}
		for j, value := range wrong {
			if value != fields[j] {
				rejected[j]++
			}
		}
	}
	if want := [5]int{4, 3, 4, 1, 1}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_BASIC_TYPE_ROWS", f.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanBasicFields", "difftest/go/types/basic-fields/testdata/basic_test.go"))
	fmt.Printf("Basic review: %d pinned Go and Lean cases, mutation counts %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
