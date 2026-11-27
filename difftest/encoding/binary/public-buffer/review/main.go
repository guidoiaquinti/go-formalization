package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"runtime"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type octet uint8

func must(ok bool, label string) {
	if !ok {
		log.Fatal(label)
	}
}

func bounds(call func()) (hit bool) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); !ok || !strings.Contains(e.Error(), "index out of range") {
				panic(p)
			}
			hit = true
		}
	}()
	call()
	return
}

func main() {
	order := binary.LittleEndian
	bools := []bool{true, false}
	signed := []int8{10, 11}
	raw := []byte{0, 255, 3}
	must(bounds(func() { binary.Decode(raw, order, bools) }) && slices.Equal(bools, []bool{false, true}), "oversized bool decode lost its writes or did not panic")
	must(bounds(func() { binary.Decode(raw, order, signed) }) && slices.Equal(signed, []int8{0, -1}), "oversized signed byte decode lost its writes or did not panic")
	ordinary := []byte{10, 11}
	n, e := binary.Decode(raw, order, ordinary)
	must(n == 2 && e == nil && slices.Equal(ordinary, raw[:2]), "oversized ordinary byte slice")
	var defined = []octet{10, 11}
	n, e = binary.Decode(raw, order, defined)
	must(n == 2 && e == nil && slices.Equal(defined, []octet{0, 255}), "oversized defined byte slice")

	_, short := binary.Decode(nil, order, uint8(9))
	_, bad := binary.Decode(nil, order, octet(9))
	must(short != nil && bad != nil && strings.Contains(short.Error(), "buffer too small") && strings.Contains(bad.Error(), "invalid type"), "scalar preflight priority")

	backing := []byte{71, 72, 73, 74, 75}
	n, e = binary.Encode(backing[:1], order, uint16(0xabcd))
	must(n == 0 && e != nil && slices.Equal(backing, []byte{71, 72, 73, 74, 75}), "short encode touched spare backing")
	n, e = binary.Encode(backing[:3], order, octet(7))
	must(n == 1 && e == nil && slices.Equal(backing, []byte{7, 72, 73, 74, 75}), "defined scalar encode touched suffix")

	out, e := binary.Append(backing[:1:4], order, uint16(0xabcd))
	must(e == nil && slices.Equal(out, []byte{7, 205, 171}) && &out[0] == &backing[0] && backing[3] == 74 && backing[4] == 75, "append reuse or spare backing mismatch")
	before := slices.Clone(backing)
	out, e = binary.Append(backing[:1:2], order, uint16(0xabcd))
	must(e == nil && slices.Equal(out, []byte{7, 205, 171}) && &out[0] != &backing[0] && slices.Equal(backing, before), "append growth changed original")
	out, e = binary.Append(backing[:1], order, "unsupported")
	must(e != nil && out == nil && slices.Equal(backing, before), "unsupported append changed original or returned input")
	nilOut, nilErr := binary.Append(nil, order, []byte(nil))
	emptyOut, emptyErr := binary.Append([]byte{}, order, []byte{})
	must(nilErr == nil && emptyErr == nil && nilOut == nil && emptyOut != nil && len(emptyOut) == 0, "nil and nonnil empty append")

	cases := []struct{ wire, want string }{
		{"0 2 0 2 3 0 1 0 0 255 3", "bounds:0:0,1|ok:2:1,0,3|ok:true:false:0,255,3,1,0:0,255,3"},
		{"1 2 0 2 3 0 10 11 0 255 3", "bounds:0:0,255|ok:2:10,11,3|ok:true:false:0,255,3,10,11:0,255,3"},
		{"2 2 0 2 3 0 10 11 0 255 3", "ok:2:0,255|ok:2:10,11,3|ok:true:false:0,255,3,10,11:0,255,3"},
		{"11 2 0 2 3 0 10 11 0 255 3", "ok:2:0,255|ok:2:10,11,3|ok:true:false:0,255,3,10,11:0,255,3"},
		{"2 0 0 1 0 1 9", "short:0:9|short:0:|ok:true:false:9:"},
		{"11 0 0 1 0 1 9", "type:0:9|short:0:|ok:true:false:9:"},
		{"4 0 0 1 1 0 43981 7 72 73 74", "short:0:43981|short:0:7,72,73,74|ok:false:false:7,205,171:7,205,171,74"},
		{"4 0 0 1 1 0 43981 7 72", "short:0:43981|short:0:7,72|ok:true:false:7,205,171:7,72"},
		{"2 2 0 0 0 1", "ok:0:|ok:0:|ok:false:true::"},
		{"2 2 0 0 0 0", "ok:0:|ok:0:|ok:false:false::"},
	}
	var wire strings.Builder
	for _, c := range cases {
		fmt.Fprintln(&wire, c.wire)
	}
	for i, got := range check.LeanRows("godiff-binary-public-buffer", len(cases), strings.NewReader(wire.String())) {
		if got != cases[i].want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].want)
		}
	}
	fmt.Println("Public buffer review: 10 independent Lean rows plus pinned native controls for prefix mutation, preflight, backing and nilness.")
}
