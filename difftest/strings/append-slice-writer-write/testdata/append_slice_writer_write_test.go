package strings_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func observeAppendSliceWrite(buffer *[]byte, input []byte) (result string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = "panic"
		}
	}()
	originalCapacity := 0
	if buffer != nil {
		originalCapacity = cap(*buffer)
	}
	written, err := strings.GoLeanAppendSliceWriterWrite(buffer, input)
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

func TestGoLeanAppendSliceWriterWrite(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_STRINGS_APPEND_SLICE_WRITER_WRITE_CASES"))
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
		switch {
		case source == -1:
			input = nil
		case source >= 0:
			input = storage[source : source+len(input)]
		}
		if got := observeAppendSliceWrite(buffer, input); got != fields[4] {
			t.Fatalf("appendSliceWriter.Write(%s): Go=%s want=%s", scanner.Text(), got, fields[4])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty comparison fixture")
	}
	fmt.Printf("private strings.appendSliceWriter.Write matched %d oracle rows\n", count)
}
