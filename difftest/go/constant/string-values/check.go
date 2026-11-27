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
	Shape, Variant int
	Bytes          [4][]byte
	Want           string
}

func wire(s string) string {
	parts := make([]string, len(s))
	for i, b := range []byte(s) {
		parts[i] = strconv.Itoa(int(b))
	}
	return strings.Join(parts, ",")
}

func oracle(c fixture) string {
	raw := string(c.Bytes[0])
	if c.Shape == 1 || c.Shape == 2 {
		raw += string(c.Bytes[1]) + string(c.Bytes[2]) + string(c.Bytes[3])
	} else if c.Shape == 3 {
		raw = strings.Repeat(raw+string(c.Bytes[1]), 2)
	}
	length := strconv.Itoa(len(raw))
	value := wire(raw) + ":leaf:" + wire(raw)
	if c.Variant == 4 {
		length, value = "nilpanic", "nilpanic"
	}
	publicLength, publicValue := length, value
	switch c.Variant {
	case 1:
		publicLength, publicValue = "0", ":unknown"
	case 2:
		publicLength, publicValue = "panic:true not a String", "panic:true not a String"
	case 3:
		publicLength, publicValue = "panic:<nil> not a String", "panic:<nil> not a String"
	}
	return fmt.Sprintf("%t|2|%s|%s|%s|%s", len(c.Bytes[0]) == 0, length, value, publicLength, publicValue)
}

func main() {
	file, err := os.CreateTemp("", "golean-constant-string-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Shape, " ", c.Variant)
			for _, raw := range c.Bytes {
				fmt.Fprint(&input, " ", len(raw))
				for _, b := range raw {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-go-constant-string-values", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d cases agree\n", label, len(cases))
	}
	words := []string{"", "left", "é", "\x00\xff", "字"}
	var fixed []fixture
	for shape := range 4 {
		for variant := range 5 {
			for _, a := range words {
				for _, b := range words {
					fixed = append(fixed, fixture{shape, variant, [4][]byte{[]byte(a), []byte(b), []byte("C"), []byte("D")}, ""})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			var payload [4][]byte
			for j := range payload {
				payload[j] = make([]byte, r.Intn(20))
				_, err := r.Read(payload[j])
				must(err)
			}
			cases[i] = fixture{r.Intn(4), r.Intn(5), payload, ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CONSTANT_STRING_ROWS", file.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanStringValues", "difftest/go/constant/string-values/testdata/string_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
