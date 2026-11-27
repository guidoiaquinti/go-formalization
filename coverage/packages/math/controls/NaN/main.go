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
	bits := math.Float64bits(math.NaN())
	withoutPayload := bits &^ 1
	withoutQuietBit := bits &^ (1 << 51)
	if !matches(bits) || matches(withoutPayload) || matches(withoutQuietBit) {
		log.Fatalf("NaN bit checks: Go=%#016x payload-cleared=%#016x quiet-bit-cleared=%#016x", bits, withoutPayload, withoutQuietBit)
	}
	fmt.Printf("%s %s/%s: NaN bits=%#016x; both mutants rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, bits)
}

func matches(bits uint64) bool {
	const payloadMask = (1 << 51) - 1
	return bits>>63 == 0 && (bits>>52)&0x7ff == 0x7ff &&
		(bits>>51)&1 == 1 && bits&payloadMask == 1
}
