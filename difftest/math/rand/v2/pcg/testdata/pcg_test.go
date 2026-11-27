package rand_test

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type pcgFixture struct {
	Hi, Lo, SeedHi, SeedLo uint64
	Length                 int
	Backing, Input         []byte
	Want                   string
}

func pcgNative(t *testing.T, c pcgFixture) string {
	gen := rand.NewPCG(c.Hi, c.Lo)
	hi, lo := rand.GoLeanPCGState(gen)
	initial := stateWire(hi, lo)
	gen.Seed(c.SeedHi, c.SeedLo)
	hi, lo = rand.GoLeanPCGState(gen)
	seeded := stateWire(hi, lo)

	gen = rand.NewPCG(c.Hi, c.Lo)
	h, l := rand.GoLeanPCGNext(gen)
	hi, lo = rand.GoLeanPCGState(gen)
	if h != hi || l != lo {
		t.Fatal("next returned state differs from receiver")
	}
	next := stateWire(h, l)
	gen = rand.NewPCG(c.Hi, c.Lo)
	first, second := gen.Uint64(), gen.Uint64()
	hi, lo = rand.GoLeanPCGState(gen)
	output := fmt.Sprintf("%d,%d:%s", first, second, stateWire(hi, lo))

	gen = rand.NewPCG(c.Hi, c.Lo)
	raw, err := gen.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	backing := slices.Clone(c.Backing)
	out, err := gen.AppendBinary(backing[:c.Length:len(backing)])
	if err != nil {
		t.Fatal(err)
	}
	reused := len(backing) > 0 && &out[0] == &backing[0]
	appended := fmt.Sprintf("%t:%s:%s", reused, byteWire(out), byteWire(backing))
	if hi, lo := rand.GoLeanPCGState(gen); hi != c.Hi || lo != c.Lo {
		t.Fatal("serialization modified receiver state")
	}
	input := slices.Clone(c.Input)
	err = gen.UnmarshalBinary(input)
	if err != nil && err != rand.GoLeanPCGUnmarshalError() {
		t.Fatalf("wrong invalid-encoding singleton: %v", err)
	}
	if !slices.Equal(input, c.Input) {
		t.Fatal("UnmarshalBinary mutated input")
	}
	hi, lo = rand.GoLeanPCGState(gen)
	restored := fmt.Sprintf("%t:%s", err == nil, stateWire(hi, lo))
	return strings.Join([]string{initial, seeded, next, output, byteWire(raw), appended, restored}, "|")
}

func stateWire(hi, lo uint64) string { return fmt.Sprintf("%d,%d", hi, lo) }

func byteWire(b []byte) string {
	fields := make([]string, len(b))
	for i, value := range b {
		fields[i] = strconv.Itoa(int(value))
	}
	return strings.Join(fields, ",")
}

func oracleNext(hi, lo uint64) (uint64, uint64) {
	join := func(hi, lo uint64) *big.Int {
		h := new(big.Int).SetUint64(hi)
		h.Lsh(h, 64)
		return h.Or(h, new(big.Int).SetUint64(lo))
	}
	state := join(hi, lo)
	state.Mul(state, join(2549297995355413924, 4865540595714422341))
	state.Add(state, join(6364136223846793005, 1442695040888963407))
	low := state.Uint64()
	state.Rsh(state, 64)
	return state.Uint64(), low
}

func oracleOutput(hi, lo uint64) uint64 {
	step := (hi ^ (hi >> 32)) * 0xda942042e4dd58b5
	return (step ^ (step >> 48)) * (lo | 1)
}

func pcgOracle(c pcgFixture) string {
	hi, lo := oracleNext(c.Hi, c.Lo)
	h, l := oracleNext(hi, lo)
	output := fmt.Sprintf("%d,%d:%s", oracleOutput(hi, lo), oracleOutput(h, l), stateWire(h, l))
	payload := make([]byte, 20)
	copy(payload, "pcg:")
	binary.BigEndian.PutUint64(payload[4:], c.Hi)
	binary.BigEndian.PutUint64(payload[12:], c.Lo)
	backing := slices.Clone(c.Backing)
	for _, span := range [][2]int{{0, 4}, {4, 12}, {12, 20}} {
		if c.Length+span[1] > len(backing) {
			break
		}
		copy(backing[c.Length+span[0]:], payload[span[0]:span[1]])
	}
	out := append(slices.Clone(c.Backing[:c.Length]), payload...)
	appended := fmt.Sprintf("%t:%s:%s", c.Length+len(payload) <= len(backing), byteWire(out), byteWire(backing))
	valid := len(c.Input) == 20 && string(c.Input[:4]) == "pcg:"
	u, v := c.Hi, c.Lo
	if valid {
		u, v = binary.BigEndian.Uint64(c.Input[4:]), binary.BigEndian.Uint64(c.Input[12:])
	}
	return strings.Join([]string{stateWire(c.Hi, c.Lo), stateWire(c.SeedHi, c.SeedLo), stateWire(hi, lo),
		output, byteWire(payload), appended, fmt.Sprintf("%t:%s", valid, stateWire(u, v))}, "|")
}

func TestGoLeanPCG(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_PCG_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		var c pcgFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		got, want := pcgNative(t, c), pcgOracle(c)
		if got != c.Want || got != want {
			t.Fatalf("fixture %+v: Go %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d native rows matched seven PCG declarations and independent oracles", count)
}
