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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Mode, Pos, Pkg, Type, Value, Import, Sig, Scope int
	Name                                            []byte
	Want                                            string
}

func main() {
	f, err := os.CreateTemp("", "golean-types-constructors-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	enc := json.NewEncoder(out)
	var fixed []fixture
	for mode := range 3 {
		for sig := range 4 {
			for scope := range 2 {
				for typ := range 3 {
					for tag := range 3 {
						for pkg := range 2 {
							for imp := range 2 {
								fixed = append(fixed, fixture{mode, []int{math.MinInt, 0, math.MaxInt}[typ], pkg, typ, tag, imp, sig, scope, []byte{255, 0, 'x'}, ""})
							}
						}
					}
				}
			}
		}
	}
	compare("fixed", fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				name := make([]byte, rng.Intn(24))
				rng.Read(name)
				rows[i] = fixture{rng.Intn(3), int(rng.Uint64()), rng.Intn(2), rng.Intn(3), rng.Intn(3), rng.Intn(2), rng.Intn(4), rng.Intn(2), name, ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, enc)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_TYPES_OBJECT_CONSTRUCT", f.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanObjectConstructors", "difftest/go/types/object-constructors/testdata/construct_test.go"))
}

func compare(label string, rows []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, r := range rows {
		fmt.Fprint(&input, r.Mode, " ", r.Pos, " ", r.Pkg, " ", r.Type, " ", r.Value, " ", r.Import, " ", r.Sig, " ", r.Scope)
		for _, b := range r.Name {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, want := range check.LeanRows("godiff-go-types-object-constructors", len(rows), &input) {
		r := rows[i]
		r.Want = want
		must(enc.Encode(r))
	}
	fmt.Printf("Go types objects %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
