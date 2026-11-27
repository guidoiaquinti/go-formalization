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

const errorText = "bytes.Buffer: UnreadRune: previous operation was not a successful ReadRune"

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
		expected(input{Tag: 1}, "fresh-error", 0, 0),
		expected(input{Tag: 1, Flag: 127}, "ok", 0, 0),
		expected(state(3, -128), "fresh-error", -128, 3),
		expected(state(2, -1), "fresh-error", -1, 2),
		expected(state(2, 0), "fresh-error", 0, 2),
		expected(state(6, 4), "ok", 0, 2),
		expected(state(3, 3), "ok", 0, 0),
		expected(state(1, 1), "ok", 0, 0),
		expected(state(2, 3), "ok", 0, 2),
		expected(state(1, 127), "ok", 0, 1),
	}
	var requests bytes.Buffer
	var clearOnError, unconditionalRewind int
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
		if c.Flag < 0 {
			clearOnError++
		}
		if c.Offset > 0 && c.Offset < c.Flag {
			unconditionalRewind++
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-unread-rune", len(tests), &requests) {
		if want := tests[i].Want; got != want {
			log.Fatalf("Lean fixture %d: UnreadRune=%s want=%s", i, got, want)
		}
	}

	// Three simultaneous errors on one receiver must have distinct
	// identities. The rejected rune rewind must not consume the ordinary
	// read flag, which still permits a one-byte rewind afterwards.
	b := bytes.NewBufferString("abcd")
	b.Next(2)
	var failures []error
	for range 3 {
		err := b.UnreadRune()
		if err == nil || err.Error() != errorText || string(b.Bytes()) != "cd" {
			log.Fatal("rejected UnreadRune changed the cursor or error result")
		}
		for _, earlier := range failures {
			if earlier == err {
				log.Fatal("UnreadRune reused a live error")
			}
		}
		failures = append(failures, err)
	}
	if err := b.UnreadByte(); err != nil || string(b.Bytes()) != "bcd" {
		log.Fatal("rejected UnreadRune cleared the prior ordinary read flag")
	}
	if err := b.UnreadByte(); err == nil || err == failures[0] {
		log.Fatal("UnreadRune error identity matched UnreadByte's error")
	}

	dir, err := os.MkdirTemp("", "golean-unread-rune-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	payload, err := json.Marshal(tests)
	must(err)
	must(os.WriteFile(path, payload, 0600))
	must(os.Setenv("GOLEAN_BUFFER_UNREAD_RUNE_CASES", path))
	const overlay = "difftest/bytes/buffer-unread-rune/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferUnreadRune", overlay+"bridge_test.go", overlay+"unread_rune_test.go"))
	if clearOnError != 2 || unconditionalRewind != 2 {
		log.Fatalf("unexpected controls: clear on error=%d, unguarded rewind=%d", clearOnError, unconditionalRewind)
	}
	fmt.Printf("Buffer.UnreadRune: %d fixed native and Lean fixtures and %d distinct repeated errors; clear-on-error=%d, unguarded rewind=%d failures\n",
		len(tests), len(failures), clearOnError, unconditionalRewind)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
