package strings_test

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

var adapterFailures = [...]error{nil, errors.New("first adapter error"), errors.New("second adapter error")}

type adapterWriter struct {
	count   int
	bounded bool
	failure error
	calls   []string
}

func (w *adapterWriter) Write(raw []byte) (int, error) {
	row := []byte("B")
	if raw == nil {
		row = []byte("nil")
	}
	for _, value := range raw {
		row = strconv.AppendUint(append(row, ','), uint64(value), 10)
	}
	w.calls = append(w.calls, string(row))
	count := w.count
	if w.bounded {
		count = min(count, len(raw))
	}
	return count, w.failure
}

func (w *adapterWriter) WriteString(raw string) (int, error) {
	w.calls = append(w.calls, "unexpected StringWriter dispatch")
	return 0, nil
}

func TestGoLeanStringWriterWriteString(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_STRING_WRITER_WRITE_STRING_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 5 {
			t.Fatalf("invalid fixture row %q", scanner.Text())
		}
		value, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		code, err := strconv.Atoi(fields[2])
		if err != nil || code < 0 || code >= len(adapterFailures) {
			t.Fatalf("invalid failure field %q", fields[2])
		}
		raw, err := hex.DecodeString(fields[3])
		if err != nil {
			t.Fatal(err)
		}
		writer := &adapterWriter{count: value, bounded: fields[0] == "1", failure: adapterFailures[code]}
		n, fault := strings.GoLeanStringWriterWriteString(writer, string(raw))
		if fault != writer.failure {
			t.Fatalf("WriteString(%x): returned error %v, expected identity %v", raw, fault, writer.failure)
		}
		got := fmt.Sprintf("%d:%d:%s", n, code, strings.Join(writer.calls, ";"))
		if got != fields[4] {
			t.Fatalf("WriteString(%x): Go=%s want=%s", raw, got, fields[4])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.stringWriter.WriteString matched %d oracle traces\n", count)
}
