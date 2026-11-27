package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	kind, tag, reported int
	values              []int
	want                string
}

type trace struct {
	c      fixture
	reads  int
	events []string
}

func (t *trace) Len() int {
	t.reads++
	return t.c.reported
}

func (t *trace) Less(a, b int) bool {
	t.events = append(t.events, fmt.Sprintf("L%d,%d", a, b))
	return t.c.values[a] < t.c.values[b]
}

func (t *trace) Swap(a, b int) {
	t.events = append(t.events, fmt.Sprintf("S%d,%d", a, b))
	t.c.values[a], t.c.values[b] = t.c.values[b], t.c.values[a]
}

func observe(t *trace) (status string) {
	status = "ok"
	defer func() {
		switch p := recover(); {
		case p == nil:
		case t.c.kind == 0 && t.c.tag == 1:
			if e, ok := p.(runtime.Error); !ok || !strings.Contains(e.Error(), "nil pointer dereference") {
				panic(p)
			}
			status = "nil-interface"
		case t.c.kind == 2 && t.c.tag == 2:
			if fmt.Sprint(p) != "reflect: call of Swapper on zero Value" {
				panic(p)
			}
			status = "invalid"
		case t.c.kind == 2 && t.c.tag == 3:
			if fmt.Sprint(p) != "reflect: call of Swapper on array Value" {
				panic(p)
			}
			status = "non-slice"
		case t.c.kind == 2 && t.c.tag == 1:
			if e, ok := p.(runtime.Error); !ok || !strings.Contains(e.Error(), "nil pointer dereference") {
				panic(p)
			}
			status = "nil-less"
		default:
			panic(p)
		}
	}()
	switch t.c.kind {
	case 0:
		if t.c.tag == 1 {
			sort.Sort(nil)
		} else {
			sort.Sort(t)
		}
	case 1:
		sort.IntSlice(t.c.values).Sort()
	case 2:
		switch t.c.tag {
		case 0:
			type named []int
			sort.Slice(named(t.c.values), t.Less)
		case 1:
			sort.Slice(t.c.values, nil)
		case 2:
			sort.Slice(nil, t.Less)
		case 3:
			sort.Slice([2]int{}, t.Less)
		}
	}
	return
}

func main() {
	cases := []fixture{
		{0, 0, 2, []int{5, 1, 90}, "ok:1:1,5,90:L1,0;S1,0"},
		{0, 0, -1, []int{5, 1}, "ok:1:5,1:"},
		{0, 1, 2, []int{5, 1}, "nil-interface:0:5,1:"},
		{1, 0, 0, []int{5, 1, 3}, "ok:0:1,3,5:"},
		{1, 1, 0, nil, "ok:0::"},
		{2, 0, 0, []int{5, 1}, "ok:0:1,5:L1,0"},
		{2, 1, 0, []int{5}, "ok:0:5:"},
		{2, 1, 0, []int{5, 1}, "nil-less:0:5,1:"},
		{2, 2, 0, nil, "invalid:0::"},
		{2, 3, 0, nil, "non-slice:0::"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d", c.kind, strconv.IntSize, c.tag, c.reported)
		for _, v := range c.values {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-sort-entrypoints", len(cases), &input)
	for i, c := range cases {
		t := &trace{c: c}
		status := observe(t)
		var values []string
		for _, v := range t.c.values {
			values = append(values, strconv.Itoa(v))
		}
		goRow := fmt.Sprintf("%s:%d:%s:%s", status, t.reads, strings.Join(values, ","), strings.Join(t.events, ";"))
		lean := rows[i]
		if c.kind == 2 {
			fields := strings.SplitN(lean, ":", 4)
			var visible []string
			for _, call := range strings.Split(fields[3], ";") {
				if strings.HasPrefix(call, "L") {
					visible = append(visible, call)
				}
			}
			fields[3] = strings.Join(visible, ";")
			lean = strings.Join(fields, ":")
		}
		if goRow != c.want || lean != c.want {
			log.Fatalf("fixture %d: Go %q, Lean %q, expected %q", i, goRow, lean, c.want)
		}
	}
	mutants := []struct {
		row int
		bad string
	}{
		{0, "ok:2:1,5,90:L1,0;S1,0"},
		{3, "ok:0:5,1,3:"},
		{6, "nil-less:0:5:"},
	}
	for _, m := range mutants {
		if rows[m.row] == m.bad {
			log.Fatal("wrong wrapper behavior accepted")
		}
	}
	fmt.Println("Ten independent public-entry fixtures passed; three mutant observations rejected.")
}
