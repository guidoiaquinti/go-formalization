package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Hi, Lo, SeedHi, SeedLo uint64
	Length                 int
	Backing, Input         []byte
	Want                   string
}

func join(hi, lo uint64) *big.Int {
	data := binary.BigEndian.AppendUint64(nil, hi)
	data = binary.BigEndian.AppendUint64(data, lo)
	return new(big.Int).SetBytes(data)
}

func steps(hi, lo uint64) (once, twice, output string) {
	modulus := new(big.Int).Lsh(big.NewInt(1), 128)
	m := join(2549297995355413924, 4865540595714422341)
	inc := join(6364136223846793005, 1442695040888963407)
	s := join(hi, lo)
	var states, outputs []string
	for range 2 {
		s.Mod(s.Add(s.Mul(s, m), inc), modulus)
		cells := s.FillBytes(make([]byte, 16))
		h, l := binary.BigEndian.Uint64(cells), binary.BigEndian.Uint64(cells[8:])
		states = append(states, pair(h, l))
		mixed := (h ^ (h / 4294967296)) * 0xda942042e4dd58b5
		result := (mixed ^ (mixed / 281474976710656)) * (l - l%2 + 1)
		outputs = append(outputs, strconv.FormatUint(result, 10))
	}
	return states[0], states[1], strings.Join(outputs, ",")
}
func pair(a, b uint64) string { return fmt.Sprintf("%d,%d", a, b) }
func csv(data []byte) string {
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func main() {
	// Choose the low word so the low half of the product is all ones.
	mod := new(big.Int).Lsh(big.NewInt(1), 64)
	inverse := new(big.Int).ModInverse(new(big.Int).SetUint64(4865540595714422341), mod)
	lo := new(big.Int).Sub(mod, inverse).Uint64()
	if lo*4865540595714422341 != math.MaxUint64 {
		log.Fatal("carry witness malformed")
	}
	hi := uint64(0x91a2b3c4d5e6f708)
	once, twice, output := steps(hi, lo)
	highText, lowText, _ := strings.Cut(once, ",")
	nextHi, err := strconv.ParseUint(highText, 10, 64)
	must(err)
	if lowText != "1442695040888963406" {
		log.Fatal("carry witness did not wrap at low half")
	}
	withoutCarry := strconv.FormatUint(nextHi-1, 10) + "," + lowText
	payload := binary.BigEndian.AppendUint64([]byte("pcg:"), hi)
	payload = binary.BigEndian.AppendUint64(payload, lo)
	valid := binary.BigEndian.AppendUint64([]byte("pcg:"), 7)
	valid = binary.BigEndian.AppendUint64(valid, 9)
	cases := []struct {
		spare, written int
		input          []byte
		restored       string
	}{
		{3, 0, nil, "false:" + pair(hi, lo)},
		{4, 4, valid[:19], "false:" + pair(hi, lo)},
		{11, 4, append(bytes.Clone(valid), 1), "false:" + pair(hi, lo)},
		{12, 12, append([]byte("PCG:"), valid[4:]...), "false:" + pair(hi, lo)},
		{19, 12, valid, "true:7,9"},
		{20, 20, valid, "true:7,9"},
		{21, 20, valid, "true:7,9"},
	}
	var fixtures []fixture
	var wire bytes.Buffer
	var rejects [3]int
	for _, c := range cases {
		backing := bytes.Repeat([]byte{91}, 2+c.spare)
		original := bytes.Clone(backing)
		copy(original[2:], payload[:c.written])
		out := append([]byte{91, 91}, payload...)
		appended := fmt.Sprintf("%t:%s:%s", c.spare >= 20, csv(out), csv(original))
		want := strings.Join([]string{pair(hi, lo), "11,13", once, output + ":" + twice, csv(payload), appended, c.restored}, "|")
		fixtures = append(fixtures, fixture{hi, lo, 11, 13, 2, backing, c.input, want})
		fmt.Fprintf(&wire, "%d %d 11 13 2 %d", hi, lo, len(backing))
		for _, part := range [][]byte{backing, c.input} {
			for _, v := range part {
				fmt.Fprintf(&wire, " %d", v)
			}
		}
		fmt.Fprintln(&wire)
		if c.spare < 20 && !bytes.Equal(original, backing) {
			rejects[0]++
		}
		if strings.HasPrefix(c.restored, "false:") && c.restored != "false:0,0" {
			rejects[1]++
		}
		if once != withoutCarry {
			rejects[2]++
		}
	}
	if rejects != [3]int{4, 4, 7} {
		log.Fatalf("unexpected control counts: %v", rejects)
	}
	for i, got := range check.LeanRows("godiff-rand-v2-pcg", len(cases), &wire) {
		if got != fixtures[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, fixtures[i].Want)
		}
	}
	file, err := os.CreateTemp("", "golean-pcg-review-")
	must(err)
	defer os.Remove(file.Name())
	rows := json.NewEncoder(file)
	for _, f := range fixtures {
		must(rows.Encode(f))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_PCG_CASES", file.Name()))
	const dir = "difftest/math/rand/v2/pcg/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanPCG", dir+"bridge_test.go", dir+"pcg_test.go"))
	fmt.Printf("PCG: %d reviewer fixtures and control counts %v\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
