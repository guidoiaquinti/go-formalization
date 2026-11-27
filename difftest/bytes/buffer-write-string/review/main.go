package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Flag int
	Stored, Spare     []byte
	Data              []byte
}

type fixture struct {
	Input input
	Want  string
}

func result(stored, spare []byte, width, offset int, reused bool) string {
	return fmt.Sprintf("ok:%d:0:%d:%d:slice:%t:%s:%s", width, offset, len(stored)+len(spare), reused, decimal(stored), decimal(spare))
}

func decimal(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func main() {
	full := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	raw := []byte{0, 255}
	large := bytes.Repeat([]byte{0x91}, 65)
	cases := []fixture{
		{input{}, "panic"},
		{input{Tag: 1, Flag: -128}, "ok:0:0:0:0:nil:false::"},
		{input{Tag: 2, Flag: 127}, result(nil, nil, 0, 0, false)},
		{input{Tag: 2, Offset: 8, Flag: -1, Stored: full}, result(full, nil, 0, 8, true)},
		{input{Tag: 2, Offset: 8, Flag: -1, Stored: full, Spare: []byte{17, 18, 19}, Data: raw}, result(append(bytes.Clone(full), raw...), []byte{19}, 2, 8, true)},
		{input{Tag: 2, Offset: 8, Flag: -1, Stored: full, Spare: []byte{17}, Data: raw}, result(raw, []byte{3, 4, 5, 6, 7, 8, 17}, 2, 0, true)},
		{input{Tag: 2, Offset: 6, Flag: 4, Stored: full, Data: raw}, result([]byte{7, 8, 0, 255}, full[4:], 2, 0, true)},
		{input{Tag: 2, Offset: 5, Flag: 3, Stored: full, Data: raw}, result([]byte{6, 7, 8, 0, 255}, make([]byte, 11), 2, 0, false)},
		{input{Tag: 1, Data: large}, result(large, make([]byte, 15), 65, 0, false)},
		{input{Tag: 1, Data: raw}, result(raw, make([]byte, 62), 2, 0, false)},
		{input{Tag: 2, Data: raw}, result(raw, make([]byte, 6), 2, 0, false)},
	}
	var rows bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprint(&rows, c.Tag)
		if c.Tag != 0 {
			fmt.Fprintf(&rows, " %d %d %d %d", len(c.Stored), len(c.Spare), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare, c.Data} {
				for _, b := range block {
					fmt.Fprintf(&rows, " %d", b)
				}
			}
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-string", len(cases), &rows) {
		if want := cases[i].Want; got != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, got, want)
		}
	}

	// A stable string copied before the write retains its bytes even when
	// compaction on the Buffer changes the original byte slice.
	original := bytes.Clone(full)
	b := bytes.NewBuffer(original[:len(original):len(original)])
	saved := string(original[:2])
	b.Next(6)
	n, err := b.WriteString(saved)
	if n != 2 || err != nil || !bytes.Equal(original, []byte{7, 8, 1, 2, 5, 6, 7, 8}) || b.UnreadByte() == nil {
		log.Fatal("compaction failed to append detached string bytes or clear unread state")
	}
	b = bytes.NewBufferString("abc")
	b.Next(3)
	n, err = b.WriteString("")
	if n != 0 || err != nil || b.UnreadByte() == nil || b.Len() != 0 {
		log.Fatal("empty write failed to clear unread state")
	}

	if cases[3].Want == result(nil, full, 0, 0, true) ||
		cases[6].Want == result([]byte{7, 8, 0, 255}, make([]byte, 12), 2, 0, false) ||
		cases[9].Want == cases[10].Want {
		log.Fatal("one incorrect variant escaped")
	}
	trace, err := os.CreateTemp("", "golean-write-string-review-")
	must(err)
	defer os.Remove(trace.Name())
	must(json.NewEncoder(trace).Encode(cases))
	must(trace.Close())
	must(os.Setenv("GOLEAN_BUFFER_WRITE_STRING_CASES", trace.Name()))
	const data = "difftest/bytes/buffer-write-string/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferWriteString", data+"bridge_test.go", data+"write_string_test.go"))
	fmt.Printf("Buffer.WriteString: %d native and Lean fixtures plus two public state checks; empty reset, strict compaction threshold, and nil/empty storage equivalence each fail a fixture\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
