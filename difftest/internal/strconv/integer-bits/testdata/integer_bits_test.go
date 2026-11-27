package strconv_test

import (
	"bufio"
	"fmt"
	strconv "internal/strconv"
	"math/bits"
	"os"
	"strings"
	"testing"
)

func TestGoLeanIntegerBits(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_STRCONV_BITS_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := bufio.NewScanner(f)
	count := 0
	for rows.Scan() {
		fields := strings.Fields(rows.Text())
		if len(fields) != 3 {
			t.Fatalf("bad input %q", rows.Text())
		}
		var width int
		var word uint64
		if _, err := fmt.Sscan(fields[0]+" "+fields[1], &width, &word); err != nil {
			t.Fatal(err)
		}
		var got string
		switch width {
		case 8:
			got = fmt.Sprint(strconv.GoLeanLower(byte(word)))
		case bits.UintSize:
			got = fmt.Sprint(strconv.GoLeanIsPowerOfTwo(int(word)))
		default:
			t.Fatalf("unsupported width %d", width)
		}
		if got != fields[2] {
			t.Fatalf("width %d word %#x: got %s want %s", width, word, got, fields[2])
		}
		count++
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("empty case file")
	}
}
