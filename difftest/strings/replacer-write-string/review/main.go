package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

var failure = errors.New("independent sentinel")

type writer struct {
	calls  []string
	script [][2]int
}
type stringWriter struct{ *writer }

func (w *writer) record(mode byte, raw string) (int, error) {
	event := []byte{mode}
	for i := range len(raw) {
		event = fmt.Appendf(event, ",%d", raw[i])
	}
	i := len(w.calls)
	w.calls = append(w.calls, string(event))
	n := len(raw)
	if i < len(w.script) {
		reply := w.script[i]
		if reply[0] >= 0 {
			n = min(n, reply[0])
		}
		if reply[1] != 0 {
			return n, failure
		}
	}
	return n, nil
}
func (w *writer) Write(p []byte) (int, error) {
	if p == nil {
		panic("unexpected nil payload")
	}
	return w.record('B', string(p))
}
func (w stringWriter) WriteString(s string) (int, error) { return w.record('S', s) }

type fixture struct {
	mode          int
	nativeString  bool
	pairs         []string
	first, second string
	script        [][2]int
	want          string
}

func native(c fixture) (out string) {
	var r *strings.Replacer
	if c.mode == 1 {
		r = new(strings.Replacer)
	}
	if c.mode == 2 {
		r = strings.NewReplacer(c.pairs...)
	}
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.mode != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			out = "panic"
		}
	}()
	run := func(input string, special bool, script [][2]int) string {
		w := &writer{script: script}
		var target io.Writer = w
		if special {
			target = stringWriter{w}
		}
		n, err := r.WriteString(target, input)
		identity := 0
		if err == failure {
			identity = 1
		} else if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%d:%d:%s", n, identity, strings.Join(w.calls, ";"))
	}
	return run(c.first, c.nativeString, c.script) + "|" + run(c.second, !c.nativeString, nil)
}
func main() {
	cases := []fixture{
		{0, false, nil, "a", "", nil, "panic"},
		{1, false, nil, "\xff", "", nil, "1:0:B,255|0:0:"},
		{2, true, []string{"a", "X", "a", "Y"}, "zaq", "a", nil, "3:0:S,122;B,88;S,113|1:0:B,88"},
		{2, true, []string{"a", ""}, "aq", "az", [][2]int{{-1, 1}}, "0:1:B|1:0:B;B,122"},
		{2, true, []string{"ab", "XY"}, "zabq", "ab", [][2]int{{0, 0}, {1, 1}}, "1:1:S,122;S,88,89|2:0:B;B,88,89;B"},
		{2, true, []string{"", "R", "a", "K"}, "a", "", nil, "3:0:S;S,82;S;S,75;S;S,82|1:0:B;B,82"},
		{2, true, []string{"a", "a", "b", ""}, "a", "a", nil, "1:0:B,97|1:0:B,97"},
		{2, false, []string{"a", "B"}, "zaq", "z", [][2]int{{0, 0}, {0, 0}}, "1:0:B,122;B,66;B,113|1:0:S,122"},
	}
	var model bytes.Buffer
	add := func(s string) {
		for i := range len(s) {
			fmt.Fprintf(&model, " %d", s[i])
		}
	}
	for i, c := range cases {
		if got := native(c); got != c.want {
			panic(fmt.Sprintf("native row %d got %q want %q", i, got, c.want))
		}
		mode := 0
		if c.nativeString {
			mode = 1
		}
		fmt.Fprintf(&model, "%d %d %d %d", c.mode, mode, len(c.script), len(c.pairs)/2)
		for _, reply := range c.script {
			fmt.Fprintf(&model, " %d %d", reply[0], reply[1])
		}
		for j := 0; j < len(c.pairs); j += 2 {
			key, val := c.pairs[j], c.pairs[j+1]
			fmt.Fprintf(&model, " %d %d", len(key), len(val))
			add(key)
			add(val)
		}
		fmt.Fprintf(&model, " %d", len(c.first))
		add(c.first)
		add(c.second)
		fmt.Fprintln(&model)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-strings-replacer-write-string")
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic(fmt.Sprintf("Lean output length %d", len(rows)))
	}
	for i, c := range cases {
		if rows[i] != c.want {
			panic(fmt.Sprintf("Lean row %d got %q want %q", i, rows[i], c.want))
		}
	}
	if strings.ReplaceAll(cases[2].want, "B,88", "S,88") == cases[2].want ||
		strings.ReplaceAll(cases[4].want, ";S,88,89", ";S,88,89;S,113") == cases[4].want {
		panic("writer trace mutant survived")
	}
	fmt.Printf("Replacer.WriteString passed %d independent pinned Go and Lean sequences\n", len(cases))
}
