package main

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	length int64
	dir    int
	elem   int
	mode   int
}

func renderElem(got, expected types.Type, token string) string {
	switch got {
	case nil:
		return "nil"
	case (*types.Array)(nil):
		return "nilArray"
	case (*types.Chan)(nil):
		return "nilChan"
	default:
		if got == expected {
			return token
		}
		return "other"
	}
}

func read[T any](readValue func() T, render func(T) string) (result string) {
	defer func() {
		if fault := recover(); fault != nil {
			if err, ok := fault.(runtime.Error); !ok ||
				err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(fault)
			}
			result = "panic"
		}
	}()
	return render(readValue())
}

func native(c fixture, elem types.Type, token string) string {
	a, a2 := types.NewArray(elem, c.length), types.NewArray(elem, c.length)
	ch, ch2 := types.NewChan(types.ChanDir(c.dir), elem), types.NewChan(types.ChanDir(c.dir), elem)
	elemWire := func(t types.Type) string { return renderElem(t, elem, token) }
	fields := []string{strconv.FormatBool(a != a2), strconv.FormatInt(a.Len(), 10), elemWire(a.Elem()),
		strconv.FormatBool(ch != ch2), strconv.Itoa(int(ch.Dir())), elemWire(ch.Elem())}
	switch c.mode {
	case 0:
		a, ch = nil, nil
	case 2:
		a, ch = new(types.Array), new(types.Chan)
	}
	return strings.Join(append(fields,
		read(a.Len, func(x int64) string { return strconv.FormatInt(x, 10) }), read(a.Elem, elemWire),
		renderElem(a.Underlying(), a, "31"),
		read(ch.Dir, func(x types.ChanDir) string { return strconv.Itoa(int(x)) }), read(ch.Elem, elemWire),
		renderElem(ch.Underlying(), ch, "32")), "|")
}

func oracle(c fixture, elem string) string {
	length, dir := strconv.FormatInt(c.length, 10), strconv.Itoa(c.dir)
	prefix := []string{"true", length, elem, "true", dir, elem}
	var access []string
	switch c.mode {
	case 0:
		access = []string{"panic", "panic", "nilArray", "panic", "panic", "nilChan"}
	case 1:
		access = []string{length, elem, "31", dir, elem, "32"}
	case 2:
		access = []string{"0", "nil", "31", "0", "nil", "32"}
	}
	return strings.Join(append(prefix, access...), "|")
}

func main() {
	elements := []types.Type{nil, (*types.Array)(nil), (*types.Chan)(nil), types.Typ[types.Int],
		types.NewArray(types.Typ[types.String], 2), types.NewChan(types.SendOnly, types.Typ[types.Byte]),
		types.NewNamed(types.NewTypeName(0, nil, "T", nil), types.Typ[types.Int], nil)}
	tokens := []string{"nil", "nilArray", "nilChan", "3", "4", "5", "6"}
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, strconv.IntSize, c.length, c.dir, c.elem, c.mode)
		}
		for i, actual := range check.LeanRows("godiff-go-types-array-chan", len(cases), &input) {
			c := cases[i]
			token := tokens[c.elem]
			want := oracle(c, token)
			goResult := native(c, elements[c.elem], token)
			if actual != want || goResult != want {
				log.Fatalf("%s row %d %+v: Go %q Lean %q want %q", label, i, c, goResult, actual, want)
			}
		}
		fmt.Printf("%s: %d pinned Go and Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, length := range []int64{math.MinInt64, -3, -1, 0, 1, math.MaxInt64} {
		for _, dir := range []int{math.MinInt, -3, -1, 0, 1, 2, 3, math.MaxInt} {
			for elem := range elements {
				for mode := range 3 {
					fixed = append(fixed, fixture{length, dir, elem, mode})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{int64(r.Uint64()), int(r.Uint64()), r.Intn(len(elements)), r.Intn(3)}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}
