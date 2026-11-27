package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Mode                             int
	Complete                         bool
	Path, Name, Version, Replacement []byte
	Want                             string
}

func wire(b []byte) string {
	parts := make([]string, len(b))
	for i, n := range b {
		parts[i] = strconv.Itoa(int(n))
	}
	return strings.Join(parts, ",")
}

func oracle(c fixture) string {
	state := func(name []byte, scope string, complete bool, version []byte) string {
		return fmt.Sprintf("%s|%s|%s|%t|%s", wire(c.Path), wire(name), scope, complete, wire(version))
	}
	created := state(c.Name, "31", false, nil)
	if c.Mode == 0 {
		return created + "#panic|panic|0|panic|panic#panic#panic"
	}
	scope := "31"
	if c.Mode == 2 {
		scope = "nil"
	}
	return strings.Join([]string{created, state(c.Name, scope, c.Complete, c.Version),
		state(c.Replacement, scope, c.Complete, c.Version), state(c.Name, scope, true, c.Version)}, "#")
}

func main() {
	file, err := os.CreateTemp("", "golean-types-package-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	compare := func(label string, cases []fixture) {
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
		}
		for i, row := range check.LeanRows("godiff-go-types-package-fields", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if c.Want != row {
				log.Fatalf("%s row %d got %q, want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	words := []string{"", "word", "go1.27", "x/y", "\xff\x00"}
	var fixed []fixture
	for mode := range 3 {
		for _, complete := range []bool{false, true} {
			for _, name := range words {
				for _, other := range words {
					fixed = append(fixed, fixture{mode, complete, []byte(name), []byte(other), []byte(name), []byte(other + "/new"), ""})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		sample := func() []byte {
			p := make([]byte, r.Intn(17))
			_, err := r.Read(p)
			must(err)
			return p
		}
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{r.Intn(3), r.Intn(2) == 0, sample(), sample(), sample(), sample(), ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_TYPES_PACKAGE_ROWS", file.Name()))
	must(check.TestInPackage("go/types", "TestGoLeanPackageFields", "difftest/go/types/package-fields/testdata/package_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
