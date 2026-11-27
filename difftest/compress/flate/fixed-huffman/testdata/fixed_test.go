package flate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGoLeanFixedHuffman(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_FLATE_FIXED_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	input := bufio.NewScanner(file)
	for _, generate := range []func() *huffmanEncoder{generateFixedLiteralEncoding, generateFixedOffsetEncoding} {
		if !input.Scan() {
			t.Fatal("missing expected table")
		}
		first, second := generate(), generate()
		if first == second || &first.codes[0] == &second.codes[0] {
			t.Fatal("separate calls shared a table or backing")
		}
		for _, table := range []*huffmanEncoder{first, second} {
			zero := table.bitCount == [17]int32{} && table.freqcache == [literalCount + 1]literalNode{}
			var row strings.Builder
			fmt.Fprint(&row, len(table.codes), " ", cap(table.codes), " ", zero)
			for _, word := range table.codes[:cap(table.codes)] {
				fmt.Fprint(&row, " ", uint32(word))
			}
			if row.String() != input.Text() {
				t.Fatalf("table with %d symbols differs from canonical words", len(table.codes))
			}
		}
	}
	if input.Scan() || input.Err() != nil {
		t.Fatal("invalid trailing table data")
	}
}
