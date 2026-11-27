package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Present                             bool
	Type, Nested                        int
	Value, Struct, Field, Source, Inner []byte
	Want                                string
}

func main() {
	f, err := os.CreateTemp("", "golean-json-compat-errors-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	enc := json.NewEncoder(out)
	var fixed []fixture
	for _, field := range []string{"", "0", "00", "0x", "a.0", "a.", ".", "1.2.3", "a.01", "a.-1", "a.١", "\xff.1", "a.\x00", "a.4294967296123456789012"} {
		for typ := range 4 {
			for nested := range 3 {
				for _, source := range []string{"", "MarshalText", "\xff\x00"} {
					fixed = append(fixed, fixture{true, typ, nested, []byte("arg\xff"), []byte(""), []byte(field), []byte(source), []byte("cause\x00"), ""},
						fixture{true, typ, nested, []byte(""), []byte("Sample\x00"), []byte(field), []byte(source), []byte(""), ""})
				}
			}
		}
	}
	fixed = append(fixed, fixture{})
	compare("fixed", fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				field := raw(rng)
				if i%3 == 0 {
					field = fmt.Appendf(field, ".%d", rng.Uint64())
				}
				rows[i] = fixture{i%17 != 0, rng.Intn(4), rng.Intn(3), raw(rng), raw(rng), field, raw(rng), raw(rng), ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, enc)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_JSON_COMPAT_ERRORS", f.Name()))
	must(check.TestInPackage("encoding/json", "TestGoLeanCompatErrors", "difftest/encoding/json/compat-errors/testdata/errors_test.go"))
}

func raw(rng *rand.Rand) []byte {
	b := make([]byte, rng.Intn(18))
	rng.Read(b)
	return b
}

func compare(label string, rows []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, r := range rows {
		present := 0
		if r.Present {
			present = 1
		}
		kind := r.Type
		if kind == 3 {
			kind = 1
		}
		fmt.Fprint(&input, present, " ", kind, " ", r.Nested)
		for _, raw := range [][]byte{[]byte([]string{"", "int", "*int", "[]string"}[r.Type]), r.Inner, r.Value, r.Struct, r.Field, r.Source} {
			fmt.Fprint(&input, " ", len(raw))
			for _, b := range raw {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
	}
	for i, result := range check.LeanRows("godiff-json-compat-errors", len(rows), &input) {
		r := rows[i]
		r.Want = result
		must(enc.Encode(r))
	}
	fmt.Printf("JSON compatibility errors %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
