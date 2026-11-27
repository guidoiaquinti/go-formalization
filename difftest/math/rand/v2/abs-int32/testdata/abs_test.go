package rand

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanAbsInt32(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ABS_INT32_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for i := 0; ; i++ {
		var word, want uint32
		n, err := fmt.Fscan(input, &word, &want)
		if err == io.EOF && n == 0 && i > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := absInt32(int32(word)); got != want {
			t.Fatalf("word %#08x: Go = %d, expected %d", word, got, want)
		}
	}
}
