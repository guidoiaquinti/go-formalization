package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	type named []string
	empty := func(func(string) bool) {}
	if slices.Collect(empty) != nil || slices.AppendSeq(named(nil), empty) != nil || slices.AppendSeq(named{}, empty) == nil {
		panic("empty sequence nilness differed")
	}
	backing := named{"start", "old1", "old2"}
	result := slices.AppendSeq(backing[:1], func(yield func(string) bool) {
		if !yield("one") || !yield("two") || !yield("three") {
			panic("append rejected a yield")
		}
	})
	if !slices.Equal(result, []string{"start", "one", "two", "three"}) ||
		!slices.Equal(backing, []string{"start", "one", "two"}) || &result[0] == &backing[0] {
		panic("growth or retained backing differed")
	}
	marker := new(int)
	for _, collect := range []bool{false, true} {
		old := named{"root", "spare"}
		calls := 0
		func() {
			defer func() {
				if recover() != marker {
					panic("missing original panic")
				}
			}()
			seq := func(yield func(string) bool) {
				if !yield("first") || !yield("second") {
					panic("yield returned false")
				}
				calls = 2
				panic(marker)
			}
			if collect {
				_ = slices.Collect(seq)
			} else {
				_ = slices.AppendSeq(old[:1], seq)
			}
		}()
		want := "first"
		if collect {
			want = "spare"
		}
		if calls != 2 || old[1] != want {
			panic("panic changed yields or original backing")
		}
	}
	cases := []struct{ input, want string }{
		{"0 0 0 0 0", "nil/0/false/|nil/0/false/"},
		{"1 0 0 0 0", "slice:/0/false/|nil/0/false/"},
		{"1 1 3 0 0 11 90 91 4 5", "slice:11,4,5/2/true/11,4,5|slice:4,5/2/false/"},
		{"1 1 3 1 72 11 90 91 4 5 6", "producer:72/3/false/11,4,5|producer:72/3/false/"},
		{"1 1 2 2 0 11 90", "nil-seq/0/false/11,90|nil-seq/0/false/"},
	}
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&model, c.input)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-slices-append-collect")
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic("wrong model row count")
	}
	for i, c := range cases {
		if rows[i] != c.want {
			panic(fmt.Sprintf("row %d got %q want %q", i, rows[i], c.want))
		}
	}
	fmt.Println("AppendSeq and Collect passed five model rows and native nilness, growth, acknowledgment and panic checks")
}
