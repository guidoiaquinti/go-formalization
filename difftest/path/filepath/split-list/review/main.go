package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input                string
		unix, windows, plan9 []string
	}{
		{"", []string{}, []string{}, []string{}},
		{":", []string{"", ""}, []string{":"}, []string{":"}},
		{";;", []string{";;"}, []string{"", "", ""}, []string{";;"}},
		{"\x00\x00", []string{"\x00\x00"}, []string{"\x00\x00"}, []string{"", "", ""}},
		{"\x00a:\x00;\x00", []string{"\x00a", "\x00;\x00"}, []string{"\x00a:\x00", "\x00"}, []string{"", "a:", ";", ""}},
		{"\"a;b\";c", []string{"\"a;b\";c"}, []string{"a;b", "c"}, []string{"\"a;b\";c"}},
		{"a\";b;c", []string{"a\";b;c"}, []string{"a;b;c"}, []string{"a\";b;c"}},
		{"a\"\";b", []string{"a\"\";b"}, []string{"a", "b"}, []string{"a\"\";b"}},
		{"\";\x00\xff\";\"", []string{"\";\x00\xff\";\""}, []string{";\x00\xff", ""}, []string{"\";", "\xff\";\""}},
		{"a\\\";b\";c", []string{"a\\\";b\";c"}, []string{"a\\;b", "c"}, []string{"a\\\";b\";c"}},
		{"\"a:b\";c", []string{"\"a", "b\";c"}, []string{"a:b", "c"}, []string{"\"a:b\";c"}},
	}
	var input bytes.Buffer
	want := make([]string, 0, len(cases))
	windows, plan9 := 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, 0)
		for _, b := range []byte(c.input) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		var native []string
		switch filepath.ListSeparator {
		case ':':
			native = c.unix
		case ';':
			native = c.windows
		case 0:
			native = c.plan9
		default:
			log.Fatalf("unknown profile %q", filepath.ListSeparator)
		}
		if got := filepath.SplitList(c.input); !reflect.DeepEqual(got, native) {
			log.Fatalf("Go fixture %d got %q, want %q", i, got, native)
		}
		if c.input != "" && !reflect.DeepEqual(strings.Split(strings.ReplaceAll(c.input, "\"", ""), ";"), c.windows) {
			windows++
		}
		if !reflect.DeepEqual(c.unix, c.plan9) {
			plan9++
		}
		want = append(want, wire(c.unix)+"|"+wire(c.windows)+"|"+wire(c.plan9))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-split-list", len(cases), &input) {
		if row != want[i] {
			log.Fatalf("Lean fixture %d got %q, want %q", i, row, want[i])
		}
	}
	if windows != 4 || plan9 != 5 {
		log.Fatalf("mutation counts windows=%d plan9=%d, want 4 and 5", windows, plan9)
	}
	fmt.Printf("filepath.SplitList passed %d pinned native and three-profile Lean fixtures; ignoring quote parity and using colon on Plan 9 failed %d and %d.\n", len(cases), windows, plan9)
}

func wire(parts []string) string {
	encoded := make([]string, len(parts))
	for i, part := range parts {
		values := make([]string, len(part))
		for j := range len(part) {
			values[j] = strconv.Itoa(int(part[j]))
		}
		encoded[i] = strings.Join(values, ",")
	}
	return fmt.Sprintf("slice:%d:%s", len(parts), strings.Join(encoded, ";"))
}
