package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	present bool
	offset  int64
	mode    int
	message string
}

type callback struct {
	message string
	calls   *int
	failure *int
}

func (c *callback) Error() string {
	*c.calls++
	if c.failure != nil {
		panic(c.failure)
	}
	return c.message
}

type typedNilCallback struct{}

var typedNilCalls int

func (*typedNilCallback) Error() string {
	typedNilCalls++
	return "\x00\xff"
}

func wire(value string) string {
	octets := []byte(value)
	fields := make([]string, len(octets))
	for i, b := range octets {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func goIO(c fixture, read bool) (result string) {
	var calls, token int
	typedNilBefore := typedNilCalls
	defer func() {
		if p := recover(); p != nil {
			calls += typedNilCalls - typedNilBefore
			if p == &token {
				result = fmt.Sprintf("panic:%d", calls)
			} else if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "nil pointer") {
				result = fmt.Sprintf("nil:%d", calls)
			} else {
				panic(p)
			}
		}
	}()
	var e error
	switch c.mode {
	case 1:
		e = &callback{message: c.message, calls: &calls}
	case 2:
		e = &callback{calls: &calls, failure: &token}
	case 3:
		// A nonnil interface can invoke a nil receiver successfully.
		e = (*typedNilCallback)(nil)
	}
	var text string
	if read {
		var err *flate.ReadError
		if c.present {
			err = &flate.ReadError{Offset: c.offset, Err: e}
		}
		text = err.Error()
	} else {
		var err *flate.WriteError
		if c.present {
			err = &flate.WriteError{Offset: c.offset, Err: e}
		}
		text = err.Error()
	}
	calls += typedNilCalls - typedNilBefore
	return fmt.Sprintf("ok:%d:%s", calls, wire(text))
}

func expected(c fixture) string {
	result := func(kind string) string {
		switch {
		case !c.present || c.mode == 0:
			return "nil:0"
		case c.mode == 2:
			return "panic:1"
		default:
			return "ok:1:" + wire(fmt.Sprintf("flate: %s error at offset %d: %s", kind, c.offset, c.message))
		}
	}
	return strings.Join([]string{wire(fmt.Sprintf("flate: corrupt input before offset %d", c.offset)),
		wire("flate: internal error: " + c.message), result("read"), result("write")}, "|")
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", c.offset, " ", c.mode)
		for _, b := range []byte(c.message) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, row := range check.LeanRows("godiff-flate-error-text", len(cases), &input) {
		c := cases[i]
		got := strings.Join([]string{wire(flate.CorruptInputError(c.offset).Error()),
			wire(flate.InternalError(c.message).Error()), goIO(c, true), goIO(c, false)}, "|")
		want := expected(c)
		if got != want || row != want {
			log.Fatalf("%s row %d: Go=%q Lean=%q want=%q", label, i, got, row, want)
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean rows agree\n", label, len(cases))
}

func main() {
	var fixed []fixture
	for _, offset := range []int64{math.MinInt64, -100, -1, 0, 1, 99, math.MaxInt64} {
		for _, present := range []bool{false, true} {
			for mode := range 4 {
				for _, message := range []string{"", "a", "\x00\xff", "ä", "EOF"} {
					if mode == 3 {
						message = "\x00\xff"
					}
					fixed = append(fixed, fixture{present, offset, mode, message})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			bytes := make([]byte, r.Intn(18))
			r.Read(bytes)
			mode := r.Intn(4)
			if mode == 3 {
				bytes = []byte{0, 255}
			}
			cases[i] = fixture{r.Intn(7) != 0, int64(r.Uint64()), mode, string(bytes)}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
}
