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
	Mode, Pos, Pkg, Type, Value, Import, Sig, Scope int
	Name                                            []byte
	Want                                            string
}

func main() {
	rows := []fixture{
		{1, -1, 1, 1, 1, 1, 0, 0, []byte{255, 0}, "typedNil1:typedNil2|9:22|typedNil1|9:false|nil|typedNil2|22|fresh|panic|35|21"},
		{2, 19, 0, 2, 2, 0, 1, 1, []byte("signed"), "11:12|9:nil|11|9:false|10|12|nil|10|50|36|nil"},
		{1, 0, 1, 0, 0, 1, 2, 1, nil, "nil:nil|9:22|nil|9:false|10|nil|22|typedNil10|panic|35|21"},
		{2, 77, 1, 2, 1, 0, 3, 0, []byte{0x80}, "11:typedNil2|9:nil|11|9:false|10|typedNil2|nil|panic|panic|36|21"},
		{0, 9, 0, 1, 2, 1, 1, 0, []byte("nil"), "typedNil1:12|9:22|typedNil1|9:false|10|panic|panic|panic|panic|panic|panic"},
		{1, 4, 0, 2, 0, 1, 1, 0, []byte("no scope"), "11:nil|9:22|11|9:false|10|nil|22|10|nil|35|nil"},
	}
	f, err := os.CreateTemp("", "golean-objects-review-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, r := range rows {
		must(enc.Encode(r))
		fmt.Fprint(&input, r.Mode, " ", r.Pos, " ", r.Pkg, " ", r.Type, " ", r.Value, " ", r.Import, " ", r.Sig, " ", r.Scope)
		for _, b := range r.Name {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	must(f.Close())
	wrong := [11]string{"nil:nil", "9:nil", "nil", "9:true", "nil", "nil", "nil", "typedNil10", "nil", "35", "nil"}
	var rejected [11]int
	for i, got := range check.LeanRows("godiff-go-types-object-constructors", len(rows), &input) {
		if got != rows[i].Want {
			log.Fatalf("row %d: Lean %q, expected %q", i, got, rows[i].Want)
		}
		for j, field := range strings.Split(got, "|") {
			if field != wrong[j] {
				rejected[j]++
			}
		}
	}
	if want := [11]int{5, 4, 5, 6, 5, 4, 4, 5, 5, 3, 4}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_TYPES_OBJECT_CONSTRUCT", f.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanObjectConstructors", "difftest/go/types/object-constructors/testdata/construct_test.go"))
	fmt.Printf("go/types object constructors: %d Go/Lean rows, mutation rejections %v\n", len(rows), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
