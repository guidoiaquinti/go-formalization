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
	Form, Kind, Length int
	Bytes              []byte
	Want               string
}

func main() {
	var fixed []fixture
	for form := range 9 {
		for kind := range 20 {
			for _, backing := range [][]byte{nil, {}, {11, 22, 33, 44, 55}} {
				for length := 0; length <= len(backing); length++ {
					fixed = append(fixed, fixture{form, kind, length, backing, ""})
				}
			}
		}
	}
	file, err := os.CreateTemp("", "golean-binary-int-data-size")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	enc := json.NewEncoder(buffer)
	compare(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				cells := make([]byte, rng.Intn(36))
				rng.Read(cells)
				rows[i] = fixture{rng.Intn(9), rng.Intn(20), rng.Intn(len(cells) + 1), cells, ""}
			}
			compare(rows, enc)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_INT_DATA_SIZE_CASES", file.Name()))
	const data = "difftest/encoding/binary/int-data-size/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanIntDataSize", data+"bridge_test.go", data+"int_data_size_test.go"))
	fmt.Printf("%d fixed cases and 20,000 cases per seed (17, 2718, 65537) agreed.\n", len(fixed))
}

func compare(cases []fixture, enc *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.Bytes == nil {
			tag = 1
		}
		fmt.Fprintf(&input, "%d %d %d %d", c.Form, c.Kind, c.Length, tag)
		for _, cell := range c.Bytes {
			fmt.Fprintf(&input, " %d", cell)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-binary-int-data-size", len(cases), &input)
	for i, row := range rows {
		c := cases[i]
		c.Want = row
		must(enc.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
