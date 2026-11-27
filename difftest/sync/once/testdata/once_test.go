package sync_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type oncePayload struct{ code int64 }

func TestGoLeanSyncOnce(t *testing.T) {
	t.Setenv("GODEBUG", "panicnil=0")
	file, err := os.Open(os.Getenv("GOLEAN_SYNC_ONCE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Tag, Mode, Second int
			Payload, X, Y     int64
			Want              string
		}
		if err := dec.Decode(&c); err == io.EOF {
			if row == 0 {
				t.Fatal("missing Once fixtures")
			}
			t.Logf("five pinned sync Once declarations matched %d rows", row)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		marker := &oncePayload{c.Payload}
		once := func(slow bool) string {
			var o *sync.Once
			if c.Tag != 0 {
				o = new(sync.Once)
			}
			if c.Tag == 2 {
				o.Do(func() {})
			}
			var results []string
			for _, mode := range []int{c.Mode, c.Second} {
				entered := false
				var callback func()
				if mode != 0 {
					callback = func() {
						entered = true
						if sync.GoLeanOnceDone(o) || !sync.GoLeanOnceLocked(o) {
							t.Fatal("Once callback did not run before done under the mutex")
						}
						fireOnce(mode, marker)
					}
				}
				status := catchOnce(o == nil, marker, func() string {
					if slow {
						sync.GoLeanOnceDoSlow(o, callback)
					} else {
						o.Do(callback)
					}
					return "ok"
				})
				done := "nil"
				if o != nil {
					done = fmt.Sprint(sync.GoLeanOnceDone(o))
					if sync.GoLeanOnceLocked(o) {
						t.Fatal("Once returned with locked mutex")
					}
				}
				results = append(results, fmt.Sprintf("%s:%t:%s", status, entered, done))
			}
			return strings.Join(results, ",")
		}
		wrapper := func(kind int) string {
			calls := 0
			var callback func() (int64, int64)
			if c.Mode != 0 {
				callback = func() (int64, int64) {
					calls++
					fireOnce(c.Mode, marker)
					return c.X, c.Y
				}
			}
			var invoke func() string
			switch kind {
			case 0:
				var f func()
				if callback != nil {
					f = func() { callback() }
				}
				g := sync.OnceFunc(f)
				invoke = func() string { g(); return "ok" }
			case 1:
				var f func() int64
				if callback != nil {
					f = func() int64 { x, _ := callback(); return x }
				}
				g := sync.OnceValue(f)
				invoke = func() string { return fmt.Sprint(g()) }
			default:
				g := sync.OnceValues(callback)
				invoke = func() string { x, y := g(); return fmt.Sprintf("%d/%d", x, y) }
			}
			if calls != 0 {
				t.Fatal("Once closure eagerly evaluated its callback")
			}
			var results []string
			for range 2 {
				before := calls
				observed := catchOnce(false, marker, invoke)
				results = append(results, fmt.Sprintf("%s:%t", observed, before != calls))
			}
			return strings.Join(results, ",")
		}
		got := strings.Join([]string{once(false), once(true), wrapper(0), wrapper(1), wrapper(2)}, "|")
		if got != c.Want {
			t.Fatalf("row %d %+v: got %q", row, c, got)
		}
	}
}

func fireOnce(mode int, marker *oncePayload) {
	switch mode {
	case 2:
		panic(marker)
	case 3:
		panic(nil)
	}
}

func catchOnce(nilReceiver bool, marker *oncePayload, f func() string) (report string) {
	defer func() {
		switch p := recover().(type) {
		case nil:
		case *oncePayload:
			if p != marker {
				panic("Once did not retain panic identity")
			}
			report = fmt.Sprintf("panic:%d", p.code)
		case *runtime.PanicNilError:
			report = "nilPanic"
		case runtime.Error:
			if p.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(p)
			}
			report = "nilCallback"
			if nilReceiver {
				report = "nilReceiver"
			}
		default:
			panic(p)
		}
	}()
	return f()
}
