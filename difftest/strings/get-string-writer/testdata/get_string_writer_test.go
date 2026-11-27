package strings_test

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

var selectionFailures = [...]error{nil, errors.New("first selection error"), errors.New("second selection error")}

type selectionReply struct{ limit, failure int }

type byteSelectionWriter struct {
	reply selectionReply
	calls []string
}

type dualSelectionWriter struct {
	*byteSelectionWriter
	reply selectionReply
}

func (w *byteSelectionWriter) record(payload string, method byte, reply selectionReply) (int, error) {
	event := []byte{method}
	for i := range len(payload) {
		event = strconv.AppendUint(append(event, ','), uint64(payload[i]), 10)
	}
	w.calls = append(w.calls, string(event))
	return min(len(payload), reply.limit), selectionFailures[reply.failure]
}

func (w *byteSelectionWriter) Write(raw []byte) (int, error) {
	return w.record(string(raw), 'B', w.reply)
}

func (w *dualSelectionWriter) WriteString(raw string) (int, error) {
	return w.record(raw, 'S', w.reply)
}

func TestGoLeanGetStringWriter(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_GET_STRING_WRITER_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 7 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		var numbers [5]int
		for i := range numbers {
			numbers[i], err = strconv.Atoi(fields[i])
			if err != nil {
				t.Fatal(err)
			}
		}
		raw, err := hex.DecodeString(fields[5])
		if err != nil {
			t.Fatal(err)
		}
		mode := numbers[0]
		byteWriter := &byteSelectionWriter{reply: selectionReply{numbers[1], numbers[2]}}
		var writer io.Writer
		switch mode {
		case 0:
			writer = byteWriter
		case 1:
			writer = &dualSelectionWriter{byteWriter, selectionReply{numbers[3], numbers[4]}}
		case 2:
			writer = (*dualSelectionWriter)(nil)
		case 3:
			writer = (*byteSelectionWriter)(nil)
		case 4:
			writer = nil
		default:
			t.Fatalf("unknown writer mode %d", mode)
		}
		selected := strings.GoLeanGetStringWriter(writer)
		if selected == nil {
			t.Fatalf("getStringWriter(%d) returned a nil interface", mode)
		}
		original, hasStringMethod := writer.(io.StringWriter)
		route := 0
		if hasStringMethod && selected == original {
			route = 1
		}
		got := strconv.Itoa(route)
		if mode <= 1 {
			n, fault := selected.WriteString(string(raw))
			actual := -1
			for i, sentinel := range selectionFailures {
				if fault == sentinel {
					actual = i
					break
				}
			}
			got = fmt.Sprintf("%d:%d:%d:%s", route, n, actual, strings.Join(byteWriter.calls, ";"))
		}
		if got != fields[6] {
			t.Fatalf("getStringWriter(%d) input %x: Go=%s want=%s", mode, raw, got, fields[6])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.getStringWriter matched %d oracle rows\n", count)
}
