package math_test

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanArchStubs(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_ARCH_STUBS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := bufio.NewReader(file)
	seen := make(map[string]bool)
	total := 0
	for {
		size, err := rows.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		nameBytes := make([]byte, size)
		if _, err := io.ReadFull(rows, nameBytes); err != nil {
			t.Fatal(err)
		}
		name := string(nameBytes)
		if seen[name] {
			t.Fatalf("repeated stub group %q", name)
		}
		seen[name] = true
		var count uint32
		if err := binary.Read(rows, binary.LittleEndian, &count); err != nil {
			t.Fatal(err)
		}
		if count < 60_024 {
			t.Fatalf("%s has %d cases, want boundary and seeded inputs", name, count)
		}
		for i := range count {
			var words [16]byte
			if _, err := io.ReadFull(rows, words[:]); err != nil {
				t.Fatal(err)
			}
			a := binary.LittleEndian.Uint64(words[:8])
			b := binary.LittleEndian.Uint64(words[8:])
			available, status := math.GoLeanCallArchStub(name, a, b)
			if available || status != "panic:not implemented" {
				t.Fatalf("%s case %d (%#016x, %#016x): haveArch=%t, %q", name, i, a, b, available, status)
			}
		}
		total += int(count)
	}
	if len(seen) != 25 {
		t.Fatalf("checked %d distinct stubs, want 25", len(seen))
	}
	t.Logf("all %d cases across %d private stubs observed availability false and exact string panic", total, len(seen))
}
