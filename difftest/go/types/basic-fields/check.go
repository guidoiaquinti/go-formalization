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
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Info int
	Present    bool
	Name       []byte
	Want       string
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func expected(c fixture) string {
	if !c.Present {
		return "panic|panic|panic|typedNil|panic"
	}
	rendered := wire(c.Name)
	if string(c.Name) == "Pointer" {
		rendered = wire([]byte("unsafe.Pointer"))
	}
	return fmt.Sprintf("%d|%d|%s|same|%s", c.Kind, c.Info, wire(c.Name), rendered)
}

func main() {
	file, err := os.CreateTemp("", "golean-types-basic-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	compare := func(label string, cases []fixture) {
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
		}
		for i, row := range check.LeanRows("godiff-go-types-basic-fields", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q, want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	bounds := []int{math.MinInt, -1, 0, 1, 18, math.MaxInt}
	names := []string{"", "Pointer", "pointer", "Sizeof", "Int", "unsafe.Pointer", "é", "\xff\x00"}
	var fixed []fixture
	for _, kind := range bounds {
		for _, info := range bounds {
			for _, name := range names {
				for _, present := range []bool{false, true} {
					fixed = append(fixed, fixture{kind, info, present, []byte(name), ""})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			name := make([]byte, r.Intn(25))
			_, err := r.Read(name)
			must(err)
			if i%4 == 0 {
				name = []byte(names[r.Intn(len(names))])
			}
			cases[i] = fixture{int(r.Uint64()), int(r.Uint64()), i%7 != 0, name, ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BASIC_TYPE_ROWS", file.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanBasicFields", "difftest/go/types/basic-fields/testdata/basic_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
