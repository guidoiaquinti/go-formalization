package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"unsafe"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}

	words := []uint64{0, 1 << 63, 1, 0x3ff0000000000000, 0x7ff0000000000000,
		0x7ff8000000000001, 0xfff8000000000002, 0x7ff0000000000001}
	var bitsFailures, fromFailures, canonicalized int
	for _, word := range words {
		var f float64
		binary.NativeEndian.PutUint64((*[8]byte)(unsafe.Pointer(&f))[:], word)
		if math.Float64bits(f) != word {
			bitsFailures++
		}
		f = math.Float64frombits(word)
		if binary.NativeEndian.Uint64((*[8]byte)(unsafe.Pointer(&f))[:]) != word {
			fromFailures++
		}
		mutant := word
		if word>>52&0x7ff == 0x7ff && word&((1<<52)-1) != 0 {
			mutant = 0x7ff8000000000001
		}
		if mutant != word {
			canonicalized++
		}
	}
	if bitsFailures != 0 || fromFailures != 0 || canonicalized != 2 {
		log.Fatalf("Float64 conversion failures: bits=%d frombits=%d canonicalization=%d",
			bitsFailures, fromFailures, canonicalized)
	}
	fmt.Printf("%s %s/%s: Float64bits=%d/8 Float64frombits=%d/8 canonicalization=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, bitsFailures, fromFailures, canonicalized)
}
