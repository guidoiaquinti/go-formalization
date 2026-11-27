package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func classification(x float64) string {
	switch {
	case math.IsNaN(x):
		return "NaN"
	case math.IsInf(x, 1):
		return "+Inf"
	case math.IsInf(x, -1):
		return "-Inf"
	case math.Float64bits(x) == 0:
		return "zero"
	case x > 0:
		return "positive"
	default:
		return "negative"
	}
}

func profile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("missing build info")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unknown Erfcinv profile %s", runtime.GOARCH)
		return -1
	}
}

func main() {
	fixtures := []struct {
		word, stage uint64
		want        string
	}{
		{0, 0x3ff0000000000000, "+Inf"}, {0x8000000000000000, 0x3ff0000000000000, "+Inf"},
		{0x8000000000000001, 0x3ff0000000000000, "+Inf"}, {0xbca0000000000000, 0x3ff0000000000000, "+Inf"},
		{0xbca0000000000001, 0x3ff0000000000001, "NaN"},
		{1, 0x3ff0000000000000, "+Inf"}, {0x3c90000000000000, 0x3ff0000000000000, "+Inf"},
		{0x3c90000000000001, 0x3fefffffffffffff, "positive"},
		{0x3ff0000000000000, 0, "zero"}, {0x4000000000000000, 0xbff0000000000000, "-Inf"},
		{0x4000000000000001, 0xbff0000000000002, "NaN"},
		{0x7ff0000000000000, 0, "NaN"}, {0xfff0000000000000, 0, "NaN"},
	}
	var rows bytes.Buffer
	wants := make([]string, len(fixtures))
	strictDomain, narrowStage := 0, 0
	selected := profile()
	for i, f := range fixtures {
		x := math.Float64frombits(f.word)
		got := math.Erfcinv(x)
		if classification(got) != f.want {
			log.Fatalf("Erfcinv(%016x)=%016x, want %s", f.word, math.Float64bits(got), f.want)
		}
		if !math.IsInf(x, 0) {
			exact := new(big.Rat).Sub(big.NewRat(1, 1), new(big.Rat).SetFloat64(x))
			stage, _ := exact.Float64()
			if math.Float64bits(stage) != f.stage || math.Float64bits(1-x) != f.stage {
				log.Fatalf("stage at %016x: exact rational rounded to %016x and Go to %016x, want %016x",
					f.word, math.Float64bits(stage), math.Float64bits(1-x), f.stage)
			}
		}
		if x < 0 && f.want != "NaN" {
			strictDomain++
		}
		if classification(math.Erfinv(float64(float32(1-x)))) != f.want {
			narrowStage++
		}
		wants[i] = strconv.FormatUint(math.Float64bits(got), 10)
		if math.IsNaN(got) {
			wants[i] = "NaN"
		}
		fmt.Fprintln(&rows, selected, f.word)
	}
	for i, got := range check.LeanRows("godiff-math-erfcinv", len(fixtures), &rows) {
		if got != wants[i] {
			log.Fatalf("Lean Erfcinv(%016x)=%s, want %s", fixtures[i].word, got, wants[i])
		}
	}
	if strictDomain != 2 || narrowStage != 3 {
		log.Fatalf("mutation counts strict domain=%d, binary32 stage=%d, want 2, 3", strictDomain, narrowStage)
	}
	fmt.Printf("Erfcinv passed %d pinned Go and Lean fixtures; strict negative rejection and binary32 stage failed %d / %d\n",
		len(fixtures), strictDomain, narrowStage)
}
