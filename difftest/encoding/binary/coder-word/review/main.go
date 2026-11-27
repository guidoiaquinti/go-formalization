package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Width, Order, Length, Offset int
	Word                         uint64
	Backing                      []byte
	Want                         string
}

func snapshot(status, value string, offset int, capArg string, raw []byte) string {
	text := make([]string, len(raw))
	for i, b := range raw {
		text[i] = strconv.Itoa(int(b))
	}
	return fmt.Sprintf("%s:%s:%d:%s:%s", status, value, offset, capArg, strings.Join(text, ","))
}

func main() {
	var cases []fixture
	var mutations [6]int
	for _, width := range []int{2, 4, 8} {
		// Position the signed high bit entirely beyond the logical one-byte buffer.
		raw := make([]byte, width+3)
		raw[0], raw[width], raw[width+1], raw[width+2] = 37, 128, 57, 79
		sign := uint64(1) << (8*width - 1)
		unsigned := strconv.FormatUint(sign, 10)
		signed := strconv.FormatInt(-int64(sign), 10)
		makeCase := func(order, length, offset int, word uint64, rows ...string) {
			cases = append(cases, fixture{width, order, length, offset, word, raw, strings.Join(rows, "|")})
		}

		enc := slices.Clone(raw)
		enc[1] = 1
		rd := snapshot("ok", unsigned, width+1, "-", raw)
		rs := snapshot("ok", signed, width+1, "-", raw)
		wr := snapshot("ok", "0", width+1, "-", enc)
		makeCase(0, 1, 1, sign|1, rd, rs, wr, wr)

		bump := slices.Clone(raw)
		bump[1], bump[len(bump)-1] = 165, 90
		capArg := strconv.Itoa(len(raw) - 1)
		failed := snapshot("callback", "0", 1, capArg, bump)
		makeCase(3, 0, 1, sign|1, failed, failed, failed, failed)
		success := snapshot("ok", "0", width+1, capArg, bump)
		makeCase(2, 0, 1, sign|1,
			snapshot("ok", strconv.FormatUint((sign-1)*2+1, 10), width+1, capArg, bump),
			snapshot("ok", "-1", width+1, capArg, bump), success, success)

		missing := snapshot("nil", "0", 1, "-", raw)
		makeCase(4, 0, 1, sign, missing, missing, missing, missing)
		bounds := snapshot("bounds", "0", 4, "-", raw)
		makeCase(4, 0, 4, sign, bounds, bounds, bounds, bounds)

		if rd != snapshot("bounds", "0", 1, "-", raw) {
			mutations[0] += 4
		}
		if failed != snapshot("callback", "0", width+1, capArg, bump) {
			mutations[1] += 4
		}
		if failed != snapshot("callback", "0", 1, capArg, raw) {
			mutations[2] += 4
		}
		if bounds != snapshot("nil", "0", 4, "-", raw) {
			mutations[3] += 4
		}
		if rs != rd {
			mutations[4]++
		}
		if failed != snapshot("callback", "0", 1, strconv.Itoa(width), bump) {
			mutations[5] += 4
		}
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Width, c.Order, c.Length, c.Offset, c.Word)
		for _, b := range c.Backing {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-binary-coder-word", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("case %d: Lean %q, want %q", i, got, cases[i].Want)
		}
	}
	if mutations != [6]int{12, 12, 12, 12, 3, 12} {
		log.Fatalf("incorrect mutant rejections %v", mutations)
	}
	f, err := os.CreateTemp("", "golean-coder-word-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_WORD_CASES", f.Name()))
	const dir = "difftest/encoding/binary/coder-word/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderWord", dir+"bridge_test.go", dir+"coder_word_test.go"))
	fmt.Printf("binary coder words: %d independent native and Lean rows, mutant rejections %v\n", len(cases), mutations)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
