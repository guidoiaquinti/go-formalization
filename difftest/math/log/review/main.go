package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	bits := math.Float64bits
	type fixture struct {
		word uint64
		want [2]uint64
	}
	same := func(word, want uint64) fixture {
		return fixture{word, [2]uint64{want, want}}
	}
	cases := []fixture{
		same(0, sign|inf), same(sign, sign|inf),
		same(inf, inf), same(sign|inf, inf+1),
		same(inf+1, inf+1), same(sign|inf|2, inf+1),
		same(sign|1, inf+1), same(bits(-1), inf+1),
		same(bits(1), 0), same(bits(0.5), 0xbfe62e42fefa39ef),
		same(bits(2), 0x3fe62e42fefa39ef),
		same(0x0010000000000000, 0xc086232bdd7abcd2),
		{1, [2]uint64{0xc08628b76e3a7b61, 0xc0874385446d71c3}},
		{2, [2]uint64{0xc08628b76e3a7b61, 0xc0873df9b3adb335}},
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	profile := 0
	if runtime.GOARCH == "arm64" {
		profile = 1
	}

	var inputs bytes.Buffer
	var flushFailures, exchangedFailures int
	for _, c := range cases {
		got := observe(bits(math.Log(math.Float64frombits(c.word))))
		if want := observe(c.want[profile]); got != want {
			log.Fatalf("Log(%#x)=%s want %s", c.word, got, want)
		}
		if 0 < c.word && c.word < 1<<52 && got != observe(sign|inf) {
			flushFailures++
		}
		if observe(c.want[1-profile]) != got {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-log", 2*len(cases), &inputs) {
		if want := observe(cases[i/2].want[i%2]); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/2, i%2, row, want)
		}
	}
	if flushFailures != 2 || exchangedFailures != 2 {
		log.Fatalf("flush failed %d, exchanged profile failed %d", flushFailures, exchangedFailures)
	}
	fmt.Printf("Log passed %d fixtures; flushing subnormal inputs failed %d; exchanging profiles failed %d\n",
		len(cases), flushFailures, exchangedFailures)
}
