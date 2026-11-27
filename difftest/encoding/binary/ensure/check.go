package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing       []byte
	Length, Count int
	Want          string
}

func main() {
	cases := []fixture{{nil, 0, -1, ""}, {nil, 0, 0, ""}, {nil, 0, 1, ""}}
	for capacity := range 24 {
		backing := make([]byte, capacity)
		for i := range backing {
			backing[i] = byte(42 + i)
		}
		for length := 0; length <= capacity; length++ {
			for count := -2; count <= capacity+2; count++ {
				cases = append(cases, fixture{backing, length, count, ""})
			}
		}
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 20000 {
			capacity := rng.Intn(48)
			backing := make([]byte, capacity)
			rng.Read(backing)
			length := rng.Intn(capacity + 1)
			count := rng.Intn(70) - 8
			cases = append(cases, fixture{backing, length, count, ""})
		}
	}
	var input bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.Backing == nil {
			tag = 1
		}
		row := fmt.Appendf(nil, "%d %d %d", tag, c.Length, c.Count)
		for _, b := range c.Backing {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
		fmt.Fprintln(&input, string(row))
	}
	rows := check.LeanRows("godiff-binary-ensure", len(cases), &input)
	file, err := os.CreateTemp("", "golean-binary-ensure")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		c.Want = rows[i]
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_ENSURE_CASES", file.Name()))
	must(check.TestInPackage("encoding/binary", "TestGoLeanEnsure", "difftest/encoding/binary/ensure/testdata/bridge_test.go", "difftest/encoding/binary/ensure/testdata/ensure_test.go"))
	fmt.Printf("%d boundaries and 20,000 per seed (17, 2718, 65537) agreed.\n", boundaries)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
