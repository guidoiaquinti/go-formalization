package constant

import (
	"encoding/json"
	"fmt"
	"go/token"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanStringValues(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_CONSTANT_STRING_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	wire := func(s string) string {
		out := make([]string, len(s))
		for i, b := range []byte(s) {
			out[i] = strconv.Itoa(int(b))
		}
		return strings.Join(out, ",")
	}
	for n := 0; ; n++ {
		var c struct {
			Shape, Variant int
			Bytes          [4][]byte
			Want           string
		}
		err := dec.Decode(&c)
		if err == io.EOF && n > 0 {
			t.Logf("%d pinned string value rows", n)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		makeLeaf := func(b []byte) *stringVal {
			x := MakeString(string(b)).(*stringVal)
			if x.l != nil || x.r != nil || x.s != string(b) {
				t.Fatal("MakeString returned incorrect leaf")
			}
			return x
		}
		initial := makeLeaf(c.Bytes[0])
		if (initial == makeLeaf(c.Bytes[0])) != (initial == &emptyString) {
			t.Fatal("wrong constructor storage identity")
		}
		tree := func() *stringVal {
			if c.Variant == 4 {
				return nil
			}
			a, b, c1, d := makeLeaf(c.Bytes[0]), makeLeaf(c.Bytes[1]), makeLeaf(c.Bytes[2]), makeLeaf(c.Bytes[3])
			add := func(l, r *stringVal) *stringVal { return BinaryOp(l, token.ADD, r).(*stringVal) }
			switch c.Shape {
			case 0:
				return a
			case 1:
				return add(add(add(a, b), c1), d)
			case 2:
				return add(a, add(b, add(c1, d)))
			default:
				sub := add(a, b)
				return add(sub, sub)
			}
		}
		observe := func(f func() string) (result string) {
			defer func() {
				if p := recover(); p != nil {
					switch e := p.(type) {
					case runtime.Error:
						if e.Error() != "runtime error: invalid memory address or nil pointer dereference" {
							panic(p)
						}
						result = "nilpanic"
					case string:
						result = "panic:" + e
					default:
						panic(p)
					}
				}
			}()
			return f()
		}
		state := func(x *stringVal, raw string) string {
			if x.l != nil || x.r != nil {
				return "uncached"
			}
			return wire(raw) + ":leaf:" + wire(x.s)
		}
		type nodeSnapshot struct {
			left, right *stringVal
			bytes       string
		}
		snapshot := func(root *stringVal) map[*stringVal]nodeSnapshot {
			saved := make(map[*stringVal]nodeSnapshot)
			var walk func(*stringVal)
			walk = func(node *stringVal) {
				if node == nil {
					return
				}
				if _, seen := saved[node]; seen {
					return
				}
				saved[node] = nodeSnapshot{node.l, node.r, node.s}
				walk(node.l)
				walk(node.r)
			}
			walk(root)
			return saved
		}
		assertUnchanged := func(saved map[*stringVal]nodeSnapshot) {
			for node, prior := range saved {
				if node.l != prior.left || node.r != prior.right || node.s != prior.bytes {
					t.Fatal("unexpected change to a string tree node")
				}
			}
		}
		lenRead := func(x *stringVal, f func() int64) string {
			return observe(func() string {
				before := snapshot(x)
				out := f()
				assertUnchanged(before)
				return strconv.FormatInt(out, 10)
			})
		}
		x := tree()
		length := lenRead(x, x.len)
		raw := observe(func() string {
			before := snapshot(x)
			delete(before, x)
			value := x.string()
			assertUnchanged(before)
			return state(x, value)
		})
		y := tree()
		var value Value = y
		switch c.Variant {
		case 1:
			value = MakeUnknown()
		case 2:
			value = MakeBool(true)
		case 3:
			value = nil
		}
		publicLength := lenRead(y, func() int64 { return StringLen(value) })
		publicValue := observe(func() string {
			before := snapshot(y)
			delete(before, y)
			result := StringVal(value)
			assertUnchanged(before)
			if c.Variant == 1 {
				return wire(result) + ":unknown"
			}
			return state(y, result)
		})
		got := fmt.Sprintf("%t|%d|%s|%s|%s|%s", initial == &emptyString, x.Kind(), length, raw, publicLength, publicValue)
		if got != c.Want {
			t.Fatalf("row %d: pinned Go %q want %q", n, got, c.Want)
		}
	}
}
