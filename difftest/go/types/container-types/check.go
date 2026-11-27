package main

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	key, elem, mode int
	sample          int64
}

func token(tag int) string {
	if tag == 0 {
		return "nil"
	}
	if tag < 6 {
		return "typed:" + strconv.Itoa(tag)
	}
	return strconv.Itoa(tag)
}

func sampleType(tag int, sample int64) types.Type {
	switch tag {
	case 0:
		return nil
	case 1:
		return (*types.Array)(nil)
	case 2:
		return (*types.Chan)(nil)
	case 3:
		return (*types.Slice)(nil)
	case 4:
		return (*types.Pointer)(nil)
	case 5:
		return (*types.Map)(nil)
	case 6:
		return types.Typ[types.Int]
	case 7:
		return types.NewArray(types.Typ[types.String], sample)
	case 8:
		return types.NewSlice(types.NewPointer(types.Typ[types.Byte]))
	default:
		return types.NewNamed(types.NewTypeName(0, nil, "T", nil), types.Typ[types.Int], nil)
	}
}

func render(got, expected types.Type, expectedID string) string {
	switch got {
	case nil:
		return "nil"
	case (*types.Array)(nil):
		return "typed:1"
	case (*types.Chan)(nil):
		return "typed:2"
	case (*types.Slice)(nil):
		return "typed:3"
	case (*types.Pointer)(nil):
		return "typed:4"
	case (*types.Map)(nil):
		return "typed:5"
	}
	if got == expected {
		return expectedID
	}
	return "wrong"
}

func observe(read func() types.Type, expected types.Type, expectedID string) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); !ok ||
				err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(p)
			}
			result = "panic"
		}
	}()
	return render(read(), expected, expectedID)
}

func native(c fixture) string {
	key, elem := sampleType(c.key, c.sample), sampleType(c.elem, c.sample)
	k, e := token(c.key), token(c.elem)
	s, p, m := types.NewSlice(elem), types.NewPointer(elem), types.NewMap(key, elem)
	fields := []string{
		strconv.FormatBool(s != types.NewSlice(elem)), render(s.Elem(), elem, e),
		strconv.FormatBool(p != types.NewPointer(elem)), render(p.Elem(), elem, e),
		strconv.FormatBool(m != types.NewMap(key, elem)), render(m.Key(), key, k), render(m.Elem(), elem, e),
	}
	switch c.mode {
	case 0:
		s, p, m = nil, nil, nil
	case 2:
		s, p, m = new(types.Slice), new(types.Pointer), new(types.Map)
	}
	return strings.Join(append(fields,
		observe(s.Elem, elem, e), render(s.Underlying(), s, "31"),
		observe(p.Elem, elem, e), render(p.Underlying(), p, "32"),
		observe(m.Key, key, k), observe(m.Elem, elem, e), render(m.Underlying(), m, "33")), "|")
}

func independent(c fixture) string {
	k, e := token(c.key), token(c.elem)
	fields := []string{"true", e, "true", e, "true", k, e}
	if c.mode == 0 {
		return strings.Join(append(fields, "panic", "typed:3", "panic", "typed:4", "panic", "panic", "typed:5"), "|")
	}
	if c.mode == 2 {
		k, e = "nil", "nil"
	}
	return strings.Join(append(fields, e, "31", e, "32", k, e, "33"), "|")
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.key, c.elem, c.mode)
	}
	for i, row := range check.LeanRows("godiff-go-types-container-types", len(cases), &input) {
		c := cases[i]
		goResult, want := native(c), independent(c)
		if row != want || goResult != want {
			log.Fatalf("%s %+v: Go %q Lean %q want %q", label, c, goResult, row, want)
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean rows agree\n", label, len(cases))
}

func main() {
	var fixed []fixture
	for k := range 10 {
		for e := range 10 {
			for mode := range 3 {
				fixed = append(fixed, fixture{k, e, mode, -1})
			}
		}
	}
	compare("all representative receiver and element profiles", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{r.Intn(10), r.Intn(10), r.Intn(3), int64(r.Uint64())}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
}
