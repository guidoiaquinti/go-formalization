package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type fixture struct {
	input, old, replacement []byte
	limit                   int
	clone, replace, all     string
}

func main() {
	cases := []fixture{
		{nil, nil, nil, 0, "nil:0:", "nil:0:", "slice:0:"},
		{[]byte{}, []byte("x"), nil, -1, "slice:0:", "nil:0:", "nil:0:"},
		{[]byte("aa"), []byte("a"), nil, -1, "slice:append:97,97", "slice:0:", "slice:0:"},
		{[]byte("aaaaa"), []byte("aa"), []byte("X"), 1,
			"slice:append:97,97,97,97,97", "slice:4:88,97,97,97", "slice:3:88,88,97"},
		{[]byte{226, 130, 239, 191, 189}, nil, []byte{255}, 2,
			"slice:append:226,130,239,191,189", "slice:7:255,226,255,130,239,191,189", "slice:9:255,226,255,130,255,239,191,189,255"},
		{[]byte{255, 239, 191, 189}, []byte{239, 191, 189}, []byte{0}, -2,
			"slice:append:255,239,191,189", "slice:2:255,0", "slice:2:255,0"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		tag := 0
		if c.input != nil {
			tag = 1
		}
		fmt.Fprint(&input, tag, " ", len(c.input), " ", len(c.old), " ", len(c.replacement), " ", c.limit)
		for _, data := range [][]byte{c.input, c.old, c.replacement} {
			for _, b := range data {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(input.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-copy-replace")
	cmd.Stdin = &input
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		log.Fatalf("got %d rows, want %d", len(rows), len(cases))
	}
	for i, c := range cases {
		actual := []struct {
			output []byte
			want   string
		}{
			{bytes.Clone(c.input), c.clone},
			{bytes.Replace(c.input, c.old, c.replacement, c.limit), c.replace},
			{bytes.ReplaceAll(c.input, c.old, c.replacement), c.all},
		}
		for j, result := range actual {
			if got := wire(result.output, strings.Contains(result.want, ":append:")); got != result.want {
				log.Fatalf("fixture %d method %d got %q want %q", i, j, got, result.want)
			}
		}
		if want := strings.Join([]string{c.clone, c.replace, c.all}, "|"); rows[i] != want {
			log.Fatalf("fixture %d Lean=%q want %q", i, rows[i], want)
		}
	}
	checkCopyIsolation()
	fmt.Printf("bytes copy reviewer control: %d native and Lean fixtures and three native alias probes passed\n", len(cases))
}

func wire(data []byte, appendAllocation bool) string {
	tag := "slice"
	if data == nil {
		tag = "nil"
	}
	capacity := fmt.Sprint(cap(data))
	if appendAllocation && cap(data) >= len(data) {
		capacity = "append"
	}
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = fmt.Sprint(b)
	}
	return tag + ":" + capacity + ":" + strings.Join(fields, ",")
}

func checkCopyIsolation() {
	for index, copy := range []func([]byte) []byte{
		bytes.Clone,
		func(source []byte) []byte { return bytes.Replace(source, source[0:2], source[2:4], 0) },
		func(source []byte) []byte { return bytes.ReplaceAll(source, source[0:2], source[2:4]) },
	} {
		storage := []byte("abCDabzz")
		source := storage[:6]
		before := slices.Clone(storage)
		result := copy(source)
		if !bytes.Equal(storage, before) {
			log.Fatalf("method %d mutated aliased source during copy", index)
		}
		want := "abCDab"
		if index == 2 {
			want = "CDCDCD"
		}
		if string(result) != want {
			log.Fatalf("method %d result %q want %q", index, result, want)
		}
		result[0] ^= 0xff
		if !bytes.Equal(storage, before) {
			log.Fatalf("method %d output aliases source or spare", index)
		}
		beforeResult := slices.Clone(result)
		source[2] ^= 0xff
		if !bytes.Equal(result, beforeResult) {
			log.Fatalf("method %d source aliases result", index)
		}
	}
}
