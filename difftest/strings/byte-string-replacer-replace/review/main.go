package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type replacement struct {
	old   byte
	value string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	entries := []replacement{{'a', "b\xff"}, {'b', ""}}
	cases := []struct {
		present     bool
		entries     []replacement
		input, want string
	}{
		{false, nil, "", ""},
		{false, nil, "ab", ""},
		{true, nil, "\x00\xff", "\x00\xff"},
		{true, entries, "", ""},
		{true, entries, "abb\x00\xff", "b\xff\x00\xff"},
		{true, entries, strings.Repeat("a", 14) + "b", strings.Repeat("b\xff", 14)},
		{true, entries, strings.Repeat("a", 15) + "b", strings.Repeat("b\xff", 15)},
		{true, entries, strings.Repeat("b", 17), ""},
		{true, entries, strings.Repeat("\xff", 18), strings.Repeat("\xff", 18)},
		{true, []replacement{{0xa9, "!"}}, "é", "\xc3!"},
		{true, []replacement{{0, "x\x00"}}, "\x00", "x\x00"},
		{true, []replacement{{'a', "a"}}, "aa", "aa"},
	}
	var leanInput, expected bytes.Buffer
	wants := make([]string, 0, len(cases))
	deleteAsAbsent, cascade := 0, 0
	for _, c := range cases {
		marker := 0
		if c.present {
			marker = 1
		}
		row := fmt.Appendf(nil, "%d %d", marker, len(c.entries))
		for _, e := range c.entries {
			row = fmt.Appendf(row, " %d %d", e.old, len(e.value))
			for _, b := range []byte(e.value) {
				row = fmt.Appendf(row, " %d", b)
			}
		}
		for _, b := range []byte(c.input) {
			row = fmt.Appendf(row, " %d", b)
		}
		fmt.Fprintln(&leanInput, string(row))
		if !c.present {
			fmt.Fprintf(&expected, "%s\tpanic\n", row)
			wants = append(wants, "panic")
			continue
		}
		fmt.Fprintf(&expected, "%s\t%x\n", row, c.want)
		values := make([]string, len(c.want))
		for i := range len(c.want) {
			values[i] = strconv.Itoa(int(c.want[i]))
		}
		wants = append(wants, "bytes:"+strings.Join(values, ","))
		if replaceIndividually(c.input, c.entries, true) != c.want {
			deleteAsAbsent++
		}
		if replaceIndividually(c.input, c.entries, false) != c.want {
			cascade++
		}
	}
	for i, row := range check.LeanRows("godiff-strings-byte-string-replacer-replace", len(cases), &leanInput) {
		if row != wants[i] {
			return fmt.Errorf("Lean fixture %d returned %q, want %q", i, row, wants[i])
		}
	}
	if deleteAsAbsent != 4 || cascade != 3 {
		return fmt.Errorf("mutation counts delete-as-absent=%d cascade=%d, want 4 and 3", deleteAsAbsent, cascade)
	}
	dir, err := os.MkdirTemp("", "golean-byte-string-replacer-replace-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, expected.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_STRINGS_BYTE_STRING_REPLACER_REPLACE_CASES", file); err != nil {
		return err
	}
	const data = "difftest/strings/byte-string-replacer-replace/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanByteStringReplacerReplace", data+"bridge_test.go", data+"byte_string_replacer_replace_test.go"); err != nil {
		return err
	}
	fmt.Printf("private byteStringReplacer.Replace passed %d fixed Go and Lean fixtures; delete-as-absent and cascading mutants failed %d and %d.\n", len(cases), deleteAsAbsent, cascade)
	return nil
}

func replaceIndividually(input string, entries []replacement, skipEmpty bool) string {
	if skipEmpty {
		var result strings.Builder
		for _, b := range []byte(input) {
			value := string([]byte{b})
			for _, e := range entries {
				if e.old == b && e.value != "" {
					value = e.value
				}
			}
			result.WriteString(value)
		}
		return result.String()
	}
	for _, e := range entries {
		input = strings.ReplaceAll(input, string([]byte{e.old}), e.value)
	}
	return input
}
