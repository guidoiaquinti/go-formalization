package main

import (
	"fmt"
	"go/types"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	length           int64
	dir, token, mode int
	want             string
}

func read[T any](f func() T, render func(T) string) (out string) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); !ok || !strings.Contains(e.Error(), "nil pointer dereference") {
				panic(p)
			}
			out = "panic"
		}
	}()
	return render(f())
}

func native(c fixture, elem types.Type) string {
	a, ch := types.NewArray(elem, c.length), types.NewChan(types.ChanDir(c.dir), elem)
	show := func(got types.Type) string {
		switch got {
		case nil:
			return "nil"
		case (*types.Array)(nil):
			return "nilArray"
		case (*types.Chan)(nil):
			return "nilChan"
		case elem:
			return strconv.Itoa(c.token)
		}
		return "wrong identity"
	}
	fields := []string{strconv.FormatBool(a != types.NewArray(elem, c.length)), fmt.Sprint(a.Len()), show(a.Elem()),
		strconv.FormatBool(ch != types.NewChan(types.ChanDir(c.dir), elem)), fmt.Sprint(ch.Dir()), show(ch.Elem())}
	if c.mode == 0 {
		a, ch = nil, nil
	} else if c.mode == 2 {
		a, ch = new(types.Array), new(types.Chan)
	}
	ar, cr := show(a.Underlying()), show(ch.Underlying())
	if a != nil && a.Underlying() == a {
		ar = "31"
	}
	if ch != nil && ch.Underlying() == ch {
		cr = "32"
	}
	fields = append(fields, read(a.Len, func(n int64) string { return fmt.Sprint(n) }), read(a.Elem, show), ar,
		read(ch.Dir, func(n types.ChanDir) string { return fmt.Sprint(n) }), read(ch.Elem, show), cr)
	return strings.Join(fields, "|")
}

func main() {
	elems := []types.Type{nil, (*types.Array)(nil), (*types.Chan)(nil), types.Typ[types.String],
		types.NewNamed(types.NewTypeName(0, nil, "Named", nil), types.Typ[types.Bool], nil)}
	cases := []fixture{
		{-7, -123, 4, 1, "true|-7|4|true|-123|4|-7|4|31|-123|4|32"},
		{-1, 3, 1, 1, "true|-1|nilArray|true|3|nilArray|-1|nilArray|31|3|nilArray|32"},
		{8, 2, 2, 0, "true|8|nilChan|true|2|nilChan|panic|panic|nilArray|panic|panic|nilChan"},
		{math.MinInt64, math.MaxInt, 3, 2, fmt.Sprintf("true|%d|3|true|%d|3|0|nil|31|0|nil|32", int64(math.MinInt64), math.MaxInt)},
		{0, 0, 0, 1, "true|0|nil|true|0|nil|0|nil|31|0|nil|32"},
	}
	mutations := []struct {
		row, col int
		alt      string
	}{
		{0, 0, "false"}, {0, 1, "0"}, {0, 3, "false"}, {0, 4, "0"},
		{1, 7, "nil"}, {1, 8, "nilArray"}, {2, 6, "0"}, {2, 8, "nil"}, {2, 9, "0"},
		{2, 11, "nil"}, {3, 10, "3"}, {4, 5, "nilChan"},
	}
	var wire strings.Builder
	for _, c := range cases {
		if got := native(c, elems[c.token]); got != c.want {
			log.Fatalf("native %q want %q", got, c.want)
		}
		fmt.Fprintf(&wire, "%d %d %d %d %d\n", strconv.IntSize, c.length, c.dir, c.token, c.mode)
	}
	n := 0
	for i, got := range check.LeanRows("godiff-go-types-array-chan", len(cases), strings.NewReader(wire.String())) {
		if got != cases[i].want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].want)
		}
		fields := strings.Split(got, "|")
		for _, m := range mutations {
			if i == m.row && fields[m.col] != m.alt {
				n++
			}
		}
	}
	if n != len(mutations) {
		log.Fatalf("rejected %d of %d mutants", n, len(mutations))
	}
	fmt.Printf("go/types array and channel: %d fixed native and Lean fixtures, %d altered results rejected\n", len(cases), n)
}
