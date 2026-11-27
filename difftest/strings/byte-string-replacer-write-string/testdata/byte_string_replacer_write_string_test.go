package strings_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

var writeFailures = [...]error{nil, errors.New("first scripted error"), errors.New("second scripted error")}

type replacerReply struct{ limit, failure int }

type replacerWriter struct {
	replies []replacerReply
	calls   []string
}

type replacerStringWriter struct{ *replacerWriter }

func (w *replacerWriter) reply(payload string, mode byte) (int, error) {
	event := []byte{mode}
	for i := range len(payload) {
		event = strconv.AppendUint(append(event, ','), uint64(payload[i]), 10)
	}
	n := len(payload)
	var failure error
	if i := len(w.calls); i < len(w.replies) {
		reply := w.replies[i]
		if reply.limit >= 0 {
			n = min(n, reply.limit)
		}
		failure = writeFailures[reply.failure]
	}
	w.calls = append(w.calls, string(event))
	return n, failure
}

func (w *replacerWriter) Write(payload []byte) (int, error) {
	if payload == nil {
		return w.reply("", 'N')
	}
	return w.reply(string(payload), 'B')
}

func (w replacerStringWriter) WriteString(payload string) (int, error) {
	return w.reply(payload, 'S')
}

func TestGoLeanByteStringReplacerWriteString(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_BYTE_STRING_REPLACER_WRITE_STRING_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		input, want, ok := strings.Cut(scanner.Text(), "\t")
		if !ok {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		values := strings.Fields(input)
		numbers := make([]int, len(values))
		for i, value := range values {
			numbers[i], err = strconv.Atoi(value)
			if err != nil {
				t.Fatal(err)
			}
		}
		if len(numbers) < 4 {
			t.Fatalf("invalid writer fixture %q", input)
		}
		stringWriter, hasTable, scriptCount := numbers[0], numbers[1], numbers[2]
		if scriptCount < 0 || 4+2*scriptCount > len(numbers) {
			t.Fatalf("invalid writer fixture lengths %q", input)
		}
		writer := &replacerWriter{replies: make([]replacerReply, scriptCount)}
		fields := numbers[4:]
		for i := range writer.replies {
			writer.replies[i] = replacerReply{fields[0], fields[1]}
			fields = fields[2:]
		}
		var table *[256][]byte
		if hasTable == 1 {
			table = new([256][]byte)
		}
		for range numbers[3] {
			if len(fields) < 2 || fields[0] < 0 || fields[0] > 255 ||
				fields[1] < 0 || fields[1] > len(fields)-2 || table == nil {
				t.Fatalf("invalid replacement entry %q", input)
			}
			key, width := fields[0], fields[1]
			entry := make([]byte, width)
			for i := range entry {
				entry[i] = byte(fields[i+2])
			}
			table[key] = entry
			fields = fields[width+2:]
		}
		raw := make([]byte, len(fields))
		for i, value := range fields {
			raw[i] = byte(value)
		}
		var native io.Writer = writer
		if stringWriter == 1 {
			native = replacerStringWriter{writer}
		}
		got := observeByteStringWriter(table, native, writer, string(raw))
		if got != want {
			t.Fatalf("WriteString input %q: Go=%s want=%s", input, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.byteStringReplacer.WriteString matched %d oracle traces\n", count)
}

func observeByteStringWriter(table *[256][]byte, native io.Writer, writer *replacerWriter, input string) (result string) {
	defer func() {
		if value := recover(); value != nil {
			fault, ok := value.(runtime.Error)
			if table != nil || len(writer.calls) != 0 || !ok ||
				fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = "panic"
		}
	}()
	written, failure := strings.GoLeanByteStringReplacerWriteString(table, native, input)
	code := -1
	for i, sentinel := range writeFailures {
		if failure == sentinel {
			code = i
			break
		}
	}
	return fmt.Sprintf("%d:%d:%s", written, code, strings.Join(writer.calls, ";"))
}
