package main

import (
	"fmt"
	"strings"
	"unsafe"
)

func main() {
	const initial = "é\xff"
	var b strings.Builder
	b.Grow(64)
	if got := b.String(); got != "" {
		panic("reserved space leaked into an empty result")
	}
	b.WriteString(strings.Clone(initial))
	old := b.String()
	again := b.String()
	if old != initial || again != initial || unsafe.StringData(old) != unsafe.StringData(again) {
		panic("Builder.String did not expose the same accumulated bytes")
	}
	copyMutant := strings.Clone(old)
	if unsafe.StringData(copyMutant) == unsafe.StringData(old) {
		panic("copy mutant unexpectedly reused the buffer")
	}
	b.WriteString(strings.Repeat("x", 128))
	if b.String() != initial+strings.Repeat("x", 128) || old != initial {
		panic("old result changed after growth")
	}
	b.Reset()
	b.WriteString("new")
	if old != initial || b.String() != "new" {
		panic("reset changed a previously returned string")
	}
	fmt.Println("Builder.String: Go passed reserve, shared view, growth and reset checks; copy mutant rejected")
}
