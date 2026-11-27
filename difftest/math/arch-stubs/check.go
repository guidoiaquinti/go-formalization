package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type operandShape int

const (
	unary operandShape = iota
	floatPair
	floatInt
)

type stub struct {
	name  string
	shape operandShape
}

type operands struct{ a, b uint64 }

var stubs = []stub{
	{"archAcos", unary},
	{"archAcosh", unary},
	{"archAsin", unary},
	{"archAsinh", unary},
	{"archAtan", unary},
	{"archAtan2", floatPair},
	{"archAtanh", unary},
	{"archCbrt", unary},
	{"archCos", unary},
	{"archCosh", unary},
	{"archErf", unary},
	{"archErfc", unary},
	{"archExpm1", unary},
	{"archFrexp", unary},
	{"archLdexp", floatInt},
	{"archLog10", unary},
	{"archLog2", unary},
	{"archLog1p", unary},
	{"archMod", floatPair},
	{"archPow", floatPair},
	{"archRemainder", floatPair},
	{"archSin", unary},
	{"archSinh", unary},
	{"archTan", unary},
	{"archTanh", unary},
}

const (
	executable = "godiff-math-arch-stubs"
	expected   = "false:panic:not implemented"
	sampleSize = 20_000
)

var seeds = []int64{17, 2718, 65537}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var profiles bytes.Buffer
	for i := range stubs {
		fmt.Fprintln(&profiles, 1, i, 0, 0)
	}
	for i, output := range check.LeanRows(executable, len(stubs), &profiles) {
		if output != "outside-scope" {
			return fmt.Errorf("%s s390x profile gave %q", stubs[i].name, output)
		}
	}
	file, err := os.CreateTemp("", "golean-arch-stubs-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	fixtures := bufio.NewWriter(file)
	for i, s := range stubs {
		edges := boundaries(s.shape)
		if err := writeHeader(fixtures, s.name, len(edges)+len(seeds)*sampleSize); err != nil {
			return err
		}
		if err := compare(i, s, edges, fixtures); err != nil {
			return err
		}
		for _, seed := range seeds {
			rng := rand.New(rand.NewSource(seed))
			cases := make([]operands, sampleSize)
			for j := range cases {
				cases[j] = sample(rng, s.shape)
			}
			if err := compare(i, s, cases, fixtures); err != nil {
				return fmt.Errorf("seed %d: %w", seed, err)
			}
		}
		fmt.Printf("%s: %d Lean boundaries and %d samples agree.\n", s.name, len(edges), len(seeds)*sampleSize)
	}
	if err := fixtures.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if runtime.GOARCH == "s390x" {
		fmt.Println("The private non-s390x bodies are outside this native build.")
		return nil
	}
	if err := os.Setenv("GOLEAN_MATH_ARCH_STUBS_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/math/arch-stubs/testdata/"
	return check.TestInPackage("math", "TestGoLeanArchStubs", data+"bridge_test.go", data+"stubs_test.go")
}

func writeHeader(f *bufio.Writer, name string, count int) error {
	if err := f.WriteByte(byte(len(name))); err != nil {
		return err
	}
	if _, err := f.WriteString(name); err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, uint32(count))
}

func compare(index int, s stub, cases []operands, f *bufio.Writer) error {
	var input bytes.Buffer
	for _, c := range cases {
		if s.shape == floatInt {
			fmt.Fprintln(&input, 0, index, c.a, int(c.b))
		} else {
			fmt.Fprintln(&input, 0, index, c.a, c.b)
		}
		var words [16]byte
		binary.LittleEndian.PutUint64(words[:8], c.a)
		binary.LittleEndian.PutUint64(words[8:], c.b)
		if _, err := f.Write(words[:]); err != nil {
			return err
		}
	}
	for i, output := range check.LeanRows(executable, len(cases), &input) {
		if output != expected {
			return fmt.Errorf("%s case %#v gave Lean result %q, want %q", s.name, cases[i], output, expected)
		}
	}
	return nil
}

func boundaries(shape operandShape) []operands {
	words := binary64.Boundaries()
	var cases []operands
	switch shape {
	case unary:
		for _, a := range words {
			cases = append(cases, operands{a, 0})
		}
	case floatPair:
		for _, a := range words {
			for _, b := range words {
				cases = append(cases, operands{a, b})
			}
		}
	case floatInt:
		minimum := -1 << (strconv.IntSize - 1)
		maximum := -1 - minimum
		for _, a := range words {
			for _, b := range []int{minimum, minimum + 1, -1075, -1074, -1, 0, 1, 1023, 1024, maximum - 1, maximum} {
				cases = append(cases, operands{a, uint64(b)})
			}
		}
	}
	return cases
}

func sample(rng *rand.Rand, shape operandShape) operands {
	a := rng.Uint64()
	if shape == unary {
		return operands{a, 0}
	}
	return operands{a, rng.Uint64()}
}
