package slices_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type goLeanRotateWire []int

func goLeanRotateRight(r int, cells goLeanRotateWire) (outcome string) {
	defer func() {
		if p := recover(); p != nil {
			if _, bounds := p.(runtime.Error); !bounds || (r >= 0 && r <= len(cells)) {
				panic(p)
			}
			outcome = "panic"
		}
	}()
	slices.GoLeanRotateRight(cells, r)
	return "ok"
}

func TestGoLeanRotateRight(t *testing.T) {
	dir := os.Getenv("GOLEAN_ROTATE_RIGHT_DIR")
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
	lines, writer := bufio.NewScanner(in), bufio.NewWriter(out)
	for lines.Scan() {
		fields := strings.Fields(lines.Text())
		r, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatal(err)
		}
		const before, after = 9001, 9002
		length, err := strconv.Atoi(fields[1])
		if err != nil {
			t.Fatal(err)
		}
		backing := make(goLeanRotateWire, len(fields))
		backing[0], backing[len(backing)-1] = before, after
		capacity := len(backing) - 2
		view := backing[1 : 1+length : 1+capacity]
		for i, f := range fields[2:] {
			v, err := strconv.Atoi(f)
			if err != nil {
				t.Fatal(err)
			}
			view[:capacity][i] = v
		}
		outcome := goLeanRotateRight(r, view)
		if backing[0] != before || backing[len(backing)-1] != after {
			t.Fatal("rotateRight modified a guard cell outside the view")
		}
		fmt.Fprint(writer, outcome)
		for _, v := range view[:capacity] {
			fmt.Fprintf(writer, " %d", v)
		}
		fmt.Fprintln(writer)
	}
	if err := lines.Err(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
}
