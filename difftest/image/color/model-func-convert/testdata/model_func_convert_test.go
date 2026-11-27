package color

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type goLeanConvertColor struct{ token uint64 }

func (*goLeanConvertColor) RGBA() (uint32, uint32, uint32, uint32) {
	panic("unexpected RGBA call")
}

type goLeanConvertPanic uint64

func TestGoLeanModelFuncConvert(t *testing.T) {
	dir := os.Getenv("GOLEAN_MODEL_FUNC_CONVERT_DIR")
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
	r, w := bufio.NewReader(in), bufio.NewWriter(out)
	for {
		var mode, argument, state, delta uint64
		n, err := fmt.Fscan(r, &mode, &argument, &state, &delta)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintln(w, goLeanObserveConvert(mode, argument, state, delta)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func goLeanObserveConvert(mode, argument, state, delta uint64) (row string) {
	var input Color
	if argument != 0 {
		input = &goLeanConvertColor{argument}
	}
	distinct := &goLeanConvertColor{argument}
	var calls, seen uint64
	var m *modelFunc
	if mode != 0 {
		m = &modelFunc{}
	}
	if mode > 1 {
		m.f = func(c Color) Color {
			calls++
			if c != input {
				panic("callback received a different argument")
			}
			if c != nil {
				seen = c.(*goLeanConvertColor).token
			}
			state += delta + seen
			switch mode {
			case 2:
				return c
			case 3:
				return nil
			case 4:
				return distinct
			default:
				panic(goLeanConvertPanic(delta))
			}
		}
	}
	format := func(tag string, payload uint64) string {
		return fmt.Sprintf("%s %d %d %d %d", tag, payload, calls, seen, state)
	}
	defer func() {
		switch p := recover().(type) {
		case nil:
		case runtime.Error:
			row = format("nil-panic", 0)
		case goLeanConvertPanic:
			row = format("callback-panic", uint64(p))
		default:
			row = fmt.Sprintf("unexpected panic: %v", p)
		}
	}()
	switch result := m.Convert(input); result {
	case nil:
		return format("ok", 0)
	case input:
		return format("ok", argument)
	case distinct:
		return format("ok", 4294967297)
	default:
		return "unexpected result identity"
	}
}
