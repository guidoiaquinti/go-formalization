package main

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	key, elem, mode int
	want            string
}

func main() {
	cases := []fixture{
		{6, 7, 1, "true|7|true|7|true|6|7|7|31|7|32|6|7|33"},
		{4, 5, 1, "true|typed:5|true|typed:5|true|typed:4|typed:5|typed:5|31|typed:5|32|typed:4|typed:5|33"},
		{7, 6, 0, "true|6|true|6|true|7|6|panic|typed:3|panic|typed:4|panic|panic|typed:5"},
		{5, 4, 2, "true|typed:4|true|typed:4|true|typed:5|typed:4|nil|31|nil|32|nil|nil|33"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.key, c.elem, c.mode)
	}
	var rejected [10]int
	for i, lean := range check.LeanRows("godiff-go-types-container-types", len(cases), &input) {
		c := cases[i]
		if goResult := native(c); lean != c.want || goResult != c.want {
			log.Fatalf("row %d: Go %q Lean %q want %q", i, goResult, lean, c.want)
		}
		f := strings.Split(lean, "|")
		for j, mutation := range []struct {
			field int
			value string
		}{
			{0, "false"}, {2, "false"}, {5, f[6]}, {7, "nil"}, {8, "nil"},
			{9, "nil"}, {10, "nil"}, {11, f[12]}, {12, "nil"}, {13, "nil"},
		} {
			if f[mutation.field] != mutation.value {
				rejected[j]++
			}
		}
	}
	if want := [10]int{4, 4, 4, 3, 4, 3, 4, 2, 3, 4}; rejected != want {
		log.Fatalf("mutant counts %v want %v", rejected, want)
	}
	fmt.Printf("go/types containers: %d independent native and Lean fixtures, mutation counts %v\n", len(cases), rejected)
}

func native(c fixture) string {
	integer := types.Typ[types.Int]
	array := types.NewArray(types.Typ[types.String], -17)
	refs := map[int]types.Type{4: (*types.Pointer)(nil), 5: (*types.Map)(nil), 6: integer, 7: array}
	key, elem := refs[c.key], refs[c.elem]
	s, p, m := types.NewSlice(elem), types.NewPointer(elem), types.NewMap(key, elem)
	fields := []string{
		fmt.Sprint(s != types.NewSlice(elem)), observe(s.Elem, array, nil),
		fmt.Sprint(p != types.NewPointer(elem)), observe(p.Elem, array, nil),
		fmt.Sprint(m != types.NewMap(key, elem)), observe(m.Key, array, nil), observe(m.Elem, array, nil),
	}
	switch c.mode {
	case 0:
		s, p, m = nil, nil, nil
	case 2:
		s, p, m = new(types.Slice), new(types.Pointer), new(types.Map)
	}
	return strings.Join(append(fields,
		observe(s.Elem, array, nil), observe(s.Underlying, array, s),
		observe(p.Elem, array, nil), observe(p.Underlying, array, p),
		observe(m.Key, array, nil), observe(m.Elem, array, nil), observe(m.Underlying, array, m)), "|")

}

func observe(read func() types.Type, array, self types.Type) (result string) {
	defer func() {
		if v := recover(); v != nil {
			if _, ok := v.(runtime.Error); !ok {
				panic(v)
			}
			result = "panic"
		}
	}()
	got := read()
	switch got {
	case nil:
		return "nil"
	case (*types.Slice)(nil):
		return "typed:3"
	case (*types.Pointer)(nil):
		return "typed:4"
	case (*types.Map)(nil):
		return "typed:5"
	case types.Typ[types.Int]:
		return "6"
	case array:
		return "7"
	}
	if got == self {
		switch got.(type) {
		case *types.Slice:
			return "31"
		case *types.Pointer:
			return "32"
		case *types.Map:
			return "33"
		}
	}
	return "wrong"
}
