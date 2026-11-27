package main

import (
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const fixture = "difftest/unicode/utf16/decode-append/testdata/decode_append_control_test.go"
	if err := check.TestInPackage("unicode/utf16", "TestGoLeanDecodeAppendControl", fixture); err != nil {
		log.Fatal(err)
	}
}
