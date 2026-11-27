package heap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type upProbe struct {
	decisions uint64
	trace     strings.Builder
}

func (*upProbe) Len() int { panic("unexpected Len") }
func (*upProbe) Push(any) { panic("unexpected Push") }
func (*upProbe) Pop() any { panic("unexpected Pop") }

func (p *upProbe) Less(child, parent int) bool {
	fmt.Fprintf(&p.trace, " L:%d:%d", child, parent)
	proceed := p.decisions&1 != 0
	p.decisions >>= 1
	return proceed
}

func (p *upProbe) Swap(parent, child int) {
	fmt.Fprintf(&p.trace, " S:%d:%d", parent, child)
}

func TestGoLeanUp(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_HEAP_UP_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var rows int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index, row, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			t.Fatalf("invalid fixture line %d", rows+1)
		}
		decisions, want, ok := strings.Cut(row, " ")
		if !ok {
			t.Fatalf("invalid fixture line %d", rows+1)
		}
		start, err := strconv.Atoi(index)
		if err != nil {
			t.Fatal(err)
		}
		mask, err := strconv.ParseUint(decisions, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		probe := upProbe{decisions: mask}
		probe.trace.WriteString("trace")
		up(&probe, start)
		if got := probe.trace.String(); got != want {
			t.Fatalf("up(%d) mask %#x: %q, want %q", start, mask, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no up fixtures")
	}
}
