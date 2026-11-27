package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Hi, Lo, SeedHi, SeedLo uint64
	Length                 int
	Backing, Input         []byte
	Want                   string
}

func main() {
	file, err := os.CreateTemp("", "golean-rand-v2-pcg")
	must(err)
	defer os.Remove(file.Name())
	buf := bufio.NewWriter(file)
	enc := json.NewEncoder(buf)
	rng := rand.New(rand.NewSource(17))
	var fixed []fixture
	edges := []uint64{0, 1, 2, 1 << 31, 1 << 32, 1 << 63, ^uint64(0) - 1, ^uint64(0)}
	for _, hi := range edges {
		for _, lo := range edges {
			for spare := range 22 {
				for _, form := range []int{0, 1, 2, 3, 4, 5} {
					fixed = append(fixed, sample(rng, hi, lo, 2, spare, form))
				}
			}
		}
	}
	fixed = append(fixed, fixture{}, fixture{Backing: []byte{}})
	run(fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				rows[i] = sample(rng, rng.Uint64(), rng.Uint64(), rng.Intn(15), rng.Intn(31), rng.Intn(6))
			}
			run(rows, enc)
		}
	}
	must(buf.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_PCG_CASES", file.Name()))
	const data = "difftest/math/rand/v2/pcg/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanPCG", data+"bridge_test.go", data+"pcg_test.go"))
	fmt.Printf("%d boundary cases and 20,000 per seed (17, 2718, 65537) passed.\n", len(fixed))
}

func sample(r *rand.Rand, hi, lo uint64, length, spare, form int) fixture {
	seedHi, seedLo := r.Uint64(), r.Uint64()
	backing := make([]byte, length+spare)
	r.Read(backing)
	valid := binary.BigEndian.AppendUint64([]byte("pcg:"), seedHi)
	valid = binary.BigEndian.AppendUint64(valid, seedLo)
	input := slices.Clone(valid)
	switch form {
	case 1:
		input = input[:r.Intn(20)]
	case 2:
		input = append(input, byte(r.Uint32()))
	case 3:
		input[r.Intn(4)] ^= 128
	case 4:
		input = nil
	case 5:
		input = []byte{}
	}
	return fixture{hi, lo, seedHi, seedLo, length, backing, input, ""}
}

func run(cases []fixture, out *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d %d", c.Hi, c.Lo, c.SeedHi, c.SeedLo, c.Length, len(c.Backing))
		for _, section := range [][]byte{c.Backing, c.Input} {
			for _, b := range section {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-rand-v2-pcg", len(cases), &input) {
		c := cases[i]
		c.Want = row
		must(out.Encode(c))
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
