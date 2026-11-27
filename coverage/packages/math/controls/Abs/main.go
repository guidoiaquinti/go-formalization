package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	const sign uint64 = 1 << 63
	magnitudes := []uint64{0, 0x3ff0000000000000, 0x7ff0000000000000,
		0x7ff0000000000001, 0x7ff8000000000001}
	var real, orderedComparison int
	for _, word := range magnitudes {
		for _, bits := range []uint64{word, word | sign} {
			x := math.Float64frombits(bits)
			if math.Float64bits(math.Abs(x)) != word {
				real++
			}
			if x < 0 {
				x = -x
			}
			if math.Float64bits(x) != word {
				orderedComparison++
			}
		}
	}
	if real != 0 || orderedComparison != 3 {
		log.Fatalf("Abs failures: Go=%d, ordered comparison=%d, want 0 and 3", real, orderedComparison)
	}
	fmt.Printf("%s %s/%s: Abs Go=%d/10, ordered comparison=%d/10\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, orderedComparison)
}
