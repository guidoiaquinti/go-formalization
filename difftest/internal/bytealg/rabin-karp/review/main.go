package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func hashes(raw []byte) [4]uint32 {
	p := big.NewInt(16777619)
	f, r := new(big.Int), new(big.Int)
	for i, b := range raw {
		weight := new(big.Int).Exp(p, big.NewInt(int64(i)), nil)
		f.Add(f, new(big.Int).Mul(weight, big.NewInt(int64(raw[len(raw)-1-i]))))
		r.Add(r, weight.Mul(weight, big.NewInt(int64(b))))
	}
	modulus := new(big.Int).Lsh(big.NewInt(1), 32)
	power := uint32(new(big.Int).Exp(p, big.NewInt(int64(len(raw))), modulus).Uint64())
	return [4]uint32{uint32(f.Mod(f, modulus).Uint64()), power, uint32(r.Mod(r, modulus).Uint64()), power}
}
func main() {
	a := []byte{0x28, 0x9e, 0x1a, 0x30, 0xb8, 0xe4, 0x2a}
	b := []byte{0x20, 0x0c, 0xf2, 0x59, 0x3d, 0x9d, 0x05}
	ra, rb := bytes.Clone(a), bytes.Clone(b)
	slices.Reverse(ra)
	slices.Reverse(rb)
	if hashes(a)[0] != hashes(b)[0] || hashes(ra)[2] != hashes(rb)[2] || bytes.Equal(a, b) {
		log.Fatal("collision witness does not collide")
	}
	cases := []struct {
		Input, Separator []byte
		First, Last      string
	}{
		{nil, nil, "0", "0"},
		{[]byte{0xff, 0x80}, nil, "0", "2"},
		{nil, []byte{1}, "panic", "panic"},
		{[]byte("aaa"), []byte("aa"), "0", "1"},
		{[]byte{0xff, 0, 0x80, 0, 0x80}, []byte{0, 0x80}, "1", "3"},
		{a, b, "-1", "-1"},
		{append(append([]byte{1}, a...), b...), b, "8", "8"},
		{ra, rb, "-1", "-1"},
		{append(append(bytes.Clone(rb), ra...), 0), rb, "0", "0"},
		{bytes.Repeat([]byte{0, 255}, 64), bytes.Repeat([]byte{0, 255}, 31), "0", "66"},
	}
	f, err := os.CreateTemp("", "golean-review-rabin-karp-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	var wants []string
	for _, c := range cases {
		h := hashes(c.Separator)
		must(enc.Encode(struct {
			Input, Separator []byte
			Hashes           [4]uint32
			First, Last      string
		}{c.Input, c.Separator, h, c.First, c.Last}))
		fmt.Fprint(&model, len(c.Separator))
		for _, raw := range [][]byte{c.Separator, c.Input} {
			for _, v := range raw {
				fmt.Fprintf(&model, " %d", v)
			}
		}
		fmt.Fprintln(&model)
		wants = append(wants, fmt.Sprintf("%d,%d,%d,%d,%s,%s", h[0], h[1], h[2], h[3], c.First, c.Last))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_RABIN_KARP_CASES", f.Name()))
	const dir = "difftest/internal/bytealg/rabin-karp/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanRabinKarp", dir+"bridge_test.go", dir+"rabin_karp_test.go"))
	for i, got := range check.LeanRows("godiff-internal-bytealg-rabin-karp", len(cases), &model) {
		if got != wants[i] {
			log.Fatalf("row %d got %q want %q", i, got, wants[i])
		}
	}
	fmt.Printf("Rabin-Karp passed %d independent native and model cases for all four declarations\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
