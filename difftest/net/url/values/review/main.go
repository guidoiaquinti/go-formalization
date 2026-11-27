package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type entry struct {
	key   string
	words []string
}

type fixture struct {
	present       bool
	entries       []entry
	key, value    string
	get           string
	has           bool
	set, add, del url.Values
	encoded       string
}

func main() {
	cases := []fixture{
		{false, nil, "x", "v", "", false, nil, nil, nil, ""},
		{true, nil, "x", "v", "", false, url.Values{"x": {"v"}}, url.Values{"x": {"v"}}, url.Values{}, ""},
		{true, []entry{{"x", nil}, {"X", []string{}}}, "x", "", "", true,
			url.Values{"x": {""}, "X": {}}, url.Values{"x": {""}, "X": {}}, url.Values{"X": {}}, ""},
		{true, []entry{{"x", []string{}}, {"X", nil}}, "x", "v", "", true,
			url.Values{"x": {"v"}, "X": nil}, url.Values{"x": {"v"}, "X": nil}, url.Values{"X": nil}, ""},
		{true, []entry{{"x", []string{"first", "last"}}, {"X", []string{"other"}}}, "x", "new", "first", true,
			url.Values{"x": {"new"}, "X": {"other"}}, url.Values{"x": {"first", "last", "new"}, "X": {"other"}},
			url.Values{"X": {"other"}}, "X=other&x=first&x=last"},
		{true, []entry{{"\xff", []string{"\x00", " "}}, {"~", []string{"?"}}, {"A", []string{"", "+"}}, {"", []string{""}}, {"0", nil}},
			"a", "&", "", false,
			url.Values{"\xff": {"\x00", " "}, "~": {"?"}, "A": {"", "+"}, "": {""}, "0": nil, "a": {"&"}},
			url.Values{"\xff": {"\x00", " "}, "~": {"?"}, "A": {"", "+"}, "": {""}, "0": nil, "a": {"&"}},
			url.Values{"\xff": {"\x00", " "}, "~": {"?"}, "A": {"", "+"}, "": {""}, "0": nil},
			"=&A=&A=%2B&~=%3F&%FF=%00&%FF=+"},
	}
	var protocol bytes.Buffer
	var expected []string
	for i, c := range cases {
		tag := 0
		if c.present {
			tag = 1
		}
		fmt.Fprint(&protocol, tag, " ", len(c.entries))
		writeWord(&protocol, c.key)
		writeWord(&protocol, c.value)
		for _, e := range c.entries {
			writeWord(&protocol, e.key)
			if e.words == nil {
				fmt.Fprint(&protocol, " -1")
				continue
			}
			fmt.Fprint(&protocol, " ", len(e.words))
			for _, w := range e.words {
				writeWord(&protocol, w)
			}
		}
		fmt.Fprintln(&protocol)
		initial := func() url.Values {
			if !c.present {
				return nil
			}
			m := make(url.Values)
			for _, e := range c.entries {
				m[e.key] = slices.Clone(e.words)
			}
			return m
		}
		set, add := state(c.set), state(c.add)
		if !c.present {
			set, add = "panic", "panic"
		}
		want := strings.Join([]string{wire(c.get), strconv.FormatBool(c.has), set, add,
			state(c.del), state(initial()), wire(c.encoded)}, "|")
		m := initial()
		gotSet := mutation(initial(), func(values url.Values) { values.Set(c.key, c.value) })
		gotAdd := mutation(initial(), func(values url.Values) { values.Add(c.key, c.value) })
		deleted := initial()
		deleted.Del(c.key)
		cloned := m.Clone()
		got := strings.Join([]string{wire(m.Get(c.key)), strconv.FormatBool(m.Has(c.key)), gotSet, gotAdd,
			state(deleted), state(cloned), wire(m.Encode())}, "|")
		if got != want {
			log.Fatalf("native row %d got %q want %q", i, got, want)
		}
		expected = append(expected, want)
	}
	checkCloneIsolation()
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(protocol.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-net-url-values")
	cmd.Stdin = &protocol
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(got) != len(expected) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(expected))
	}
	for i, row := range got {
		if row != expected[i] {
			log.Fatalf("Lean row %d got %q want %q", i, row, expected[i])
		}
	}
	fmt.Printf("net/url.Values review: %d fixed native/Lean rows and clone isolation passed\n", len(cases))
}

func writeWord(w *bytes.Buffer, s string) {
	fmt.Fprint(w, " ", len(s))
	for i := range len(s) {
		fmt.Fprint(w, " ", s[i])
	}
}

func wire(s string) string {
	parts := make([]string, len(s))
	for i := range parts {
		parts[i] = strconv.Itoa(int(s[i]))
	}
	return fmt.Sprintf("%d:%s", len(s), strings.Join(parts, ","))
}

func state(values url.Values) string {
	if values == nil {
		return "nil"
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	var parts []string
	for _, key := range keys {
		entry := "nil"
		if words := values[key]; words != nil {
			items := make([]string, len(words))
			for i, w := range words {
				items[i] = wire(w)
			}
			entry = "list:" + strings.Join(items, ";")
		}
		parts = append(parts, wire(key)+"="+entry)
	}
	return "map:" + strings.Join(parts, "/")
}

func mutation(values url.Values, op func(url.Values)) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); !ok || err.Error() != "assignment to entry in nil map" {
				panic(p)
			}
			result = "panic"
		}
	}()
	op(values)
	return state(values)
}

func checkCloneIsolation() {
	shared := []string{"value", "spare"}
	original := url.Values{"x": shared[:1], "y": shared[:1], "empty": {}, "nil": nil}
	cloned := original.Clone()
	cloned["x"][0] = "different"
	if original["x"][0] != "value" || cloned["y"][0] != "value" {
		log.Fatal("Clone aliases an original or cloned value list")
	}
	delete(cloned, "x")
	original["y"][0] = "later"
	if !original.Has("x") || cloned["y"][0] != "value" || cloned["empty"] == nil || cloned["nil"] != nil {
		log.Fatal("Clone failed map isolation or value nilness")
	}
}
