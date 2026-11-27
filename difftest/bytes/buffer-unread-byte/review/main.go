package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const errorText = "bytes.Buffer: UnreadByte: previous operation was not a successful read"

type input struct {
	Tag, Offset, Flag int
	Stored, Spare     []byte
}

type fixture struct {
	Input input
	Want  string
}

func decimal(raw []byte) string {
	values := make([]string, len(raw))
	for i, octet := range raw {
		values[i] = strconv.Itoa(int(octet))
	}
	return strings.Join(values, ",")
}

func expected(c input, status string, flag, offset int) fixture {
	marker := "slice"
	if c.Tag == 1 {
		marker = "nil"
	}
	unread := c.Stored[offset:]
	result := fmt.Sprintf("%s:%d:%d:%d:%s:%s:%s:%s", status, flag, offset, len(unread)+len(c.Spare),
		marker, decimal(unread), decimal(c.Spare), decimal(c.Stored))
	return fixture{c, result}
}

func main() {
	original := []byte("abcdef")
	spare := []byte{0xff, 0}
	state := func(offset, flag int) input {
		return input{Tag: 2, Offset: offset, Flag: flag, Stored: original, Spare: spare}
	}
	tests := []fixture{
		{Input: input{}, Want: "panic"},
		expected(input{Tag: 1}, "sentinel-error", 0, 0),
		expected(input{Tag: 1, Flag: 127}, "ok", 0, 0),
		expected(state(3, -128), "ok", 0, 2),
		expected(state(2, -1), "ok", 0, 1),
		expected(state(2, 0), "sentinel-error", 0, 2),
		expected(state(6, 4), "ok", 0, 5),
		expected(state(3, 3), "ok", 0, 2),
		expected(state(1, 1), "ok", 0, 0),
		expected(state(0, 3), "ok", 0, 0),
		expected(state(1, 127), "ok", 0, 0),
	}
	var requests bytes.Buffer
	var rejectNegative, rewindWidth, rejectZeroOffset int
	for _, test := range tests {
		c := test.Input
		if c.Tag == 0 {
			fmt.Fprintln(&requests, 0)
			continue
		}
		fmt.Fprintf(&requests, "%d %d %d %d", c.Tag, len(c.Stored), c.Offset, c.Flag)
		for _, raw := range [][]byte{c.Stored, c.Spare} {
			for _, value := range raw {
				fmt.Fprintf(&requests, " %d", value)
			}
		}
		fmt.Fprintln(&requests)
		if c.Flag < 0 && expected(c, "sentinel-error", c.Flag, c.Offset).Want != test.Want {
			rejectNegative++
		}
		if c.Flag > 0 && expected(c, "ok", 0, max(0, c.Offset-c.Flag)).Want != test.Want {
			rewindWidth++
		}
		if c.Offset == 0 && c.Flag != 0 && expected(c, "sentinel-error", c.Flag, 0).Want != test.Want {
			rejectZeroOffset++
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-unread-byte", len(tests), &requests) {
		if want := tests[i].Want; got != want {
			log.Fatalf("Lean fixture %d: UnreadByte=%s want=%s", i, got, want)
		}
	}

	b := bytes.NewBufferString("éZ")
	first := b.UnreadByte()
	if first == nil || first.Error() != errorText {
		log.Fatalf("initial UnreadByte returned %v", first)
	}
	value, width, err := b.ReadRune()
	if err != nil || value != 'é' || width != 2 {
		log.Fatalf("ReadRune returned %q, %d, %v", value, width, err)
	}
	if err := b.UnreadByte(); err != nil || string(b.Bytes()) != "\xa9Z" {
		log.Fatal("UnreadByte after ReadRune did not rewind exactly one byte")
	}
	for range 3 {
		if err := b.UnreadByte(); err != first || string(b.Bytes()) != "\xa9Z" {
			log.Fatal("rejected UnreadByte changed the cursor or sentinel")
		}
	}

	dir, err := os.MkdirTemp("", "golean-unread-byte-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	payload, err := json.Marshal(tests)
	must(err)
	must(os.WriteFile(path, payload, 0600))
	must(os.Setenv("GOLEAN_BUFFER_UNREAD_BYTE_CASES", path))
	const overlay = "difftest/bytes/buffer-unread-byte/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferUnreadByte", overlay+"bridge_test.go", overlay+"unread_byte_test.go"))
	if rejectNegative != 2 || rewindWidth != 2 || rejectZeroOffset != 2 {
		log.Fatalf("unexpected controls: reject negative flag=%d, rewind rune width=%d, reject zero offset=%d",
			rejectNegative, rewindWidth, rejectZeroOffset)
	}
	fmt.Printf("Buffer.UnreadByte: %d fixed native and Lean fixtures and repeated sentinel identity; three incorrect variants fail %d/%d/%d cases\n",
		len(tests), rejectNegative, rewindWidth, rejectZeroOffset)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
