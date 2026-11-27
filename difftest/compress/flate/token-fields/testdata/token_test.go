package flate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanTokenFields(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_TOKEN_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for row := 0; ; row++ {
		var word, typ, offset uint32
		var literal, length uint8
		parsed, err := fmt.Fscan(input, &word, &typ, &literal, &offset, &length)
		if parsed == 0 && err == io.EOF && row > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		tok := token(word)
		if tok.typ() != typ || tok.literal() != literal || tok.offset() != offset || tok.length() != length {
			t.Fatalf("row %d word %#08x: got %d %d %d %d; want %d %d %d %d",
				row, word, tok.typ(), tok.literal(), tok.offset(), tok.length(), typ, literal, offset, length)
		}
	}
}
