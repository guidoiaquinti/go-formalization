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

	words := []uint32{0, 1 << 31, 1, 0x3f800000, 0x7f800000, 0x7fc00000, 0xffc00002, 0x7f800001}
	var bitsFailures, fromFailures, canonicalized int
	for _, word := range words {
		var f float32
		binary.NativeEndian.PutUint32((*[4]byte)(unsafe.Pointer(&f))[:], word)
		if math.Float32bits(f) != word {
			bitsFailures++
		}
		f = math.Float32frombits(word)
		if binary.NativeEndian.Uint32((*[4]byte)(unsafe.Pointer(&f))[:]) != word {
			fromFailures++
		}
		mutant := word
		if word>>23&0xff == 0xff && word&((1<<23)-1) != 0 {
			mutant = 0x7fc00000
		}
		if mutant != word {
			canonicalized++
		}
	}
	if bitsFailures != 0 || fromFailures != 0 || canonicalized != 2 {
		log.Fatalf("Float32 conversion failures: bits=%d frombits=%d canonicalization=%d",
			bitsFailures, fromFailures, canonicalized)
	}
	fmt.Printf("%s %s/%s: Float32bits=%d/8 Float32frombits=%d/8 canonicalization=%d/8\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, bitsFailures, fromFailures, canonicalized)
}
