package heap

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type downProbe struct {
	low, high uint64
	trace     strings.Builder
}

func (*downProbe) Len() int { panic("unexpected Len") }
func (*downProbe) Push(any) { panic("unexpected Push") }
func (*downProbe) Pop() any { panic("unexpected Pop") }

func (p *downProbe) Less(a, b int) bool {
	fmt.Fprintf(&p.trace, " L:%d:%d", a, b)
	proceed := p.low&1 != 0
	p.low = p.low>>1 | p.high<<63
	p.high >>= 1
	return proceed
}

func (p *downProbe) Swap(a, b int) {
	fmt.Fprintf(&p.trace, " S:%d:%d", a, b)
}

func TestGoLeanDown(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_HEAP_DOWN_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		line := strings.NewReader(scanner.Text())
		var index, length int
		var p downProbe
		if _, err := fmt.Fscan(line, &index, &length, &p.low, &p.high); err != nil {
			t.Fatal(err)
		}
		expected, err := io.ReadAll(line)
		if err != nil {
			t.Fatal(err)
		}
		moved := down(&p, index, length)
		got := fmt.Sprint(moved) + p.trace.String()
		if want := strings.TrimSpace(string(expected)); got != want {
			t.Fatalf("down(%d,%d): %q, want %q", index, length, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no down fixtures")
	}
}
