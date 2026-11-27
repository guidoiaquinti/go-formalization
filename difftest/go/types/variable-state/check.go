package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
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

func oracle(c fixture) string {
	names := []string{"VarKind(0)", "PackageVar", "LocalVar", "RecvVar", "ParamVar", "ResultVar", "FieldVar"}
	label := fmt.Sprintf("VarKind(%d)", c.Kind)
	if c.Kind < len(names) {
		label = names[c.Kind]
	}
	head := fmt.Sprintf("%d:false:true|1:false:true|4:false:true|6:%t:true", c.Kind, c.Embedded)
	if c.Mode == 0 {
		return head + strings.Repeat("|panic", 6) + "|" + label
	}
	id := 31
	if c.Mode == 2 {
		id = 32
	}
	return fmt.Sprintf("%s|%d|%t|%t|%t|%d|%d:%t:%d|%s",
		head, c.Kind, c.Embedded, c.Embedded, c.Kind == 6, id, c.Replacement, c.Embedded, id, label)
}

func main() {
	file, err := os.CreateTemp("", "golean-types-var-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	writer := bufio.NewWriter(file)
	enc := json.NewEncoder(writer)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			embedded := 0
			if c.Embedded {
				embedded = 1
			}
			fmt.Fprint(&input, c.Mode, " ", c.Kind, " ", c.Replacement, " ", embedded, " ", c.Pos, " ", c.Pkg, " ", c.Typ)
			for _, b := range c.Name {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-go-types-variable-state", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	var fixed []fixture
	positions := []int{math.MinInt, -1, 0, 1, math.MaxInt}
	for kind := range 256 {
		for mode := range 3 {
			for _, e := range []bool{false, true} {
				fixed = append(fixed, fixture{mode, kind, 255 - kind, positions[kind%5], kind % 2,
					kind % 4, e, []byte("raw\x00\xff"), ""})
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			name := make([]byte, r.Intn(17))
			_, err := r.Read(name)
			must(err)
			cases[i] = fixture{r.Intn(3), r.Intn(256), r.Intn(256), int(r.Uint64()), r.Intn(2),
				r.Intn(4), r.Intn(2) == 1, name, ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_TYPES_VAR_ROWS", file.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanVariableState", "difftest/go/types/variable-state/testdata/var_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
