package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

func observeAppendSliceWriteString(buffer *[]byte, input string) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = "panic"
		}
	}()
	originalCapacity := 0
	if buffer != nil {
		originalCapacity = cap(*buffer)
	}
	written, err := strings.GoLeanAppendSliceWriterWriteString(buffer, input)
	if buffer == nil {
		return "missing panic"
	}
	data := *buffer
	capacity := strconv.Itoa(cap(data))
	if len(data) > originalCapacity && cap(data) >= len(data) {
		capacity = "growth"
	}
	nilness := "slice"
	if data == nil {
		nilness = "nil"
	}
	row := fmt.Appendf(nil, "%d:%t:%s:%s:", written, err != nil, nilness, capacity)
	for i, value := range data {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(value), 10)
	}
	return string(row)
}

func TestGoLeanAppendSliceWriterWriteString(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_APPEND_SLICE_WRITER_WRITE_STRING_CASES"))
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
		width, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		source, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		backing, err := hex.DecodeString(fields[2])
		if err != nil {
			t.Fatal(err)
		}
		input, err := hex.DecodeString(fields[3])
		if err != nil {
			t.Fatal(err)
		}
		// Allocate exactly the supplied backing capacity, including a nonnil empty slice.
		storage := make([]byte, len(backing))
		copy(storage, backing)
		var data []byte
		buffer := &data
		switch width {
		case -1:
			buffer = nil
		case -2:
			data = nil
		default:
			data = storage[:width]
		}
		text := string(input)
		if source >= 0 && len(input) != 0 {
			end := width + len(input)
			if width < 0 || source+len(input) > len(storage) ||
				(end <= len(storage) && source < end && width < source+len(input)) {
				t.Fatalf("fixture would modify the string backing: %q", scanner.Text())
			}
			text = unsafe.String(&storage[source], len(input))
			if text != string(input) {
				t.Fatalf("fixture source differs from the requested input: %q", scanner.Text())
			}
		}
		if got := observeAppendSliceWriteString(buffer, text); got != fields[4] {
			t.Fatalf("appendSliceWriter.WriteString(%s): Go=%s want=%s", scanner.Text(), got, fields[4])
		}
		if text != string(input) {
			t.Fatalf("append modified the string source: %q", scanner.Text())
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.appendSliceWriter.WriteString matched %d oracle rows\n", count)
}
