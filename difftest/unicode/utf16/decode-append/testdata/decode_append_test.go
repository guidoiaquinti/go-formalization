package utf16

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGoLeanDecodeAppend(t *testing.T) {
	dir := os.Getenv("GOLEAN_UTF16_DECODE_DIR")
	if dir == "" {
		t.Fatal("missing GOLEAN_UTF16_DECODE_DIR")
	}
	in, err := os.Open(filepath.Join(dir, "inputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(filepath.Join(dir, "outputs"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	input, output := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var spare, length, count int
		if _, err := fmt.Fscan(input, &spare, &length, &count); err == io.EOF {
			break
		} else if err != nil || spare < 0 || length < -1 || count < -1 {
			t.Fatalf("invalid sequence header: %d %d %d: %v", spare, length, count, err)
		}
		var dst []rune
		if length >= 0 || spare != 0 {
			length = max(length, 0)
			dst = make([]rune, length, length+spare)
		}
		var words []uint16
		if count >= 0 {
			words = make([]uint16, count)
		}
		for i := range dst {
			if _, err := fmt.Fscan(input, &dst[i]); err != nil {
				t.Fatal(err)
			}
		}
		for i := range words {
			if _, err := fmt.Fscan(input, &words[i]); err != nil {
				t.Fatal(err)
			}
		}
		result := decode(words, dst)
		fmt.Fprint(output, "runes")
		for _, r := range result {
			fmt.Fprint(output, " ", r)
		}
		fmt.Fprintln(output)
	}
	if err := output.Flush(); err != nil {
		t.Fatal(err)
	}
}
