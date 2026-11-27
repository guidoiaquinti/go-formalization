package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	words := []uint64{0, 1 << 63, 1, 2, 1 << 51, (1 << 52) - 1, 1 << 52, (1 << 52) + 1,
		0x7ff0000000000000, 0xfff0000000000000, 0x7ff8000000000061, 0xbff0000000000000,
		0x3fefffffffffffff, 0x3ff0000000000000, 0x3ff0000000000001}
	for _, x := range []float64{0.125, 0.5, 0.75, 1.25, 2, 3, 100, math.MaxFloat64} {
		words = append(words, math.Float64bits(x))
	}
	var input bytes.Buffer
	for p := range 2 {
		for _, word := range words {
			fmt.Fprintln(&input, p, word)
		}
	}
	rows := check.LeanRows("godiff-math-arch-log", 2*len(words), &input)
	var native bytes.Buffer
	for i, word := range words {
		x := math.Float64frombits(word)
		got := rows[i]
		switch {
		case math.IsNaN(x) || x < 0:
			if got != "NaN" {
				log.Fatalf("input %#x: got %s, want NaN", word, got)
			}
		case x == 0:
			assertWord(word, got, math.Inf(-1))
		case math.IsInf(x, 1):
			assertWord(word, got, x)
		default:
			scale := 0
			if word < 1<<52 {
				x = math.Float64frombits(word | 1<<52)
				scale = -1
			}
			fraction, exponent := math.Frexp(x)
			want := math.Log1p(2*fraction-1) + float64(exponent+scale-1)*math.Ln2
			if x > 0.75 && x < 1.25 {
				want = math.Log1p(x - 1)
			}
			bits, err := strconv.ParseUint(got, 10, 64)
			must(err)
			value := math.Float64frombits(bits)
			if !(math.Abs(want-value) <= 4e-15*math.Abs(want)) {
				log.Fatalf("input %#x: model %.17g, reference %.17g", word, value, want)
			}
		}
		fallback := rows[i+len(words)]
		if fallback != "panic:not implemented" {
			log.Fatalf("fallback input %#x: %s", word, fallback)
		}
		if runtime.GOARCH != "amd64" {
			got = fallback
		}
		fmt.Fprintf(&native, "%d %q\n", word, got)
	}
	if runtime.GOARCH == "s390x" {
		return
	}
	out, err := os.CreateTemp("", "golean-arch-log-review-*")
	must(err)
	defer os.Remove(out.Name())
	_, err = out.Write(native.Bytes())
	must(err)
	must(out.Close())
	must(os.Setenv("GOLEAN_MATH_ARCH_LOG_CASES", out.Name()))
	const data = "difftest/math/arch-log/testdata/"
	must(check.TestInPackage("math", "TestGoLeanArchLog", data+"bridge_test.go", data+"arch_log_test.go"))
	fmt.Printf("archLog: %d native, %d model observations passed.\n", len(words), len(rows))
}

func assertWord(input uint64, got string, want float64) {
	if encoded := strconv.FormatUint(math.Float64bits(want), 10); got != encoded {
		log.Fatalf("input %#x: model %s, want %s", input, got, encoded)
	}
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
