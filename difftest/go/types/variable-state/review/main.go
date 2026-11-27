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
	Mode, Kind, Replacement, Pos, Pkg, Typ int
	Embedded                               bool
	Name                                   []byte
	Want                                   string
}

func main() {
	cases := []fixture{
		{0, 0, 255, 0, 0, 0, false, nil, "0:false:true|1:false:true|4:false:true|6:false:true|panic|panic|panic|panic|panic|panic|VarKind(0)"},
		{2, 6, 2, -1, 1, 1, false, []byte{255, 0}, "6:false:true|1:false:true|4:false:true|6:false:true|6|false|false|true|32|2:false:32|FieldVar"},
		{1, 255, 6, 17, 1, 3, true, []byte("field"), "255:false:true|1:false:true|4:false:true|6:true:true|255|true|true|false|31|6:true:31|VarKind(255)"},
		{2, 4, 0, 89, 0, 2, true, []byte("param"), "4:false:true|1:false:true|4:false:true|6:true:true|4|true|true|false|32|0:true:32|ParamVar"},
	}
	f, err := os.CreateTemp("", "golean-type-var-control-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		embed := 0
		if c.Embedded {
			embed = 1
		}
		fmt.Fprint(&input, c.Mode, " ", c.Kind, " ", c.Replacement, " ", embed, " ", c.Pos, " ", c.Pkg, " ", c.Typ)
		for _, b := range c.Name {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [11]int
	for i, got := range check.LeanRows("godiff-go-types-variable-state", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
		fields := strings.Split(got, "|")
		wrong := [11]string{"1:false:true", "4:false:true", "1:false:true", "6:false:true", "0", "false", "false", fields[6], "31", fields[9], "VarKind(0)"}
		if c.Mode != 0 {
			suffix := strings.Split(fields[9], ":")
			wrong[9] = fmt.Sprintf("%d:%s:%s", c.Kind, suffix[1], suffix[2])
		}
		for j, w := range wrong {
			if w != fields[j] {
				rejected[j]++
			}
		}
	}
	if want := [11]int{4, 4, 4, 2, 4, 3, 3, 3, 3, 3, 3}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_TYPES_VAR_ROWS", f.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanVariableState", "difftest/go/types/variable-state/testdata/var_test.go"))
	fmt.Printf("go/types variable review: %d pinned Go and Lean rows, mutation counts %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
