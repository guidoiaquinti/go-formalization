package color

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestGoLeanSqDiff(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SQ_DIFF_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows := bufio.NewReader(f)
	for {
		var x, y, want uint32
		n, err := fmt.Fscan(rows, &x, &y, &want)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if got := sqDiff(x, y); got != want {
			t.Fatalf("sqDiff(%d, %d) = %d, want %d", x, y, got, want)
		}
	}
}
