package main

import (
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
	Nil            bool
	Length, Offset int
	Argument       byte
	Backing        []byte
	Want           string
}

func main() {
	cases := []fixture{{Nil: true}, {Backing: nil}, {Backing: []byte{}}}
	for v := range 256 {
		for _, length := range []int{0, 1, 2, 3} {
			for _, offset := range []int{math.MinInt, -1, 0, 1, 2, 3, math.MaxInt} {
				cases = append(cases, fixture{false, length, offset, byte(v), []byte{byte(v), 17, byte(255 - v)}, ""})
			}
		}
	}
	boundaries := len(cases)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			backing := make([]byte, rng.Intn(32))
			rng.Read(backing)
			length := rng.Intn(len(backing) + 1)
			offset := rng.Intn(len(backing)+4) - 2
			if i%8 == 0 {
				offset = int(rng.Uint64())
			}
			cases = append(cases, fixture{false, length, offset, byte(rng.Uint32()), backing, ""})
		}
	}
	var input bytes.Buffer
	for _, c := range cases {
		tag := 1
		if c.Nil {
			tag = 0
		}
		fmt.Fprintf(&input, "%d %d %d %d", tag, c.Length, c.Offset, c.Argument)
		for _, b := range c.Backing {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-binary-coder-byte", len(cases), &input)
	file, err := os.CreateTemp("", "golean-binary-coder-byte")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		c.Want = rows[i]
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_BYTE_CASES", file.Name()))
	const base = "difftest/encoding/binary/coder-byte/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderByte", base+"bridge_test.go", base+"coder_byte_test.go"))
	fmt.Printf("%d boundaries and 20,000 rows per seed (17, 2718, 65537) passed.\n", boundaries)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
