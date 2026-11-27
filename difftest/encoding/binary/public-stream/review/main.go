package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type reviewByte byte

type reader struct {
	rest     []byte
	requests []int
	stop     error
}

func (r *reader) Read(p []byte) (int, error) {
	r.requests = append(r.requests, len(p))
	n := copy(p[:min(1, len(p))], r.rest)
	r.rest = r.rest[n:]
	if len(r.rest) == 0 {
		return n, r.stop
	}
	return n, nil
}

type writer struct {
	original *byte
	payload  []byte
	calls    int
	alias    bool
	stop     error
}

func (w *writer) Write(p []byte) (int, error) {
	w.calls++
	w.payload = slices.Clone(p)
	w.alias = len(p) > 0 && &p[0] == w.original
	return 0, w.stop
}

func must(ok bool, label string) {
	if !ok {
		log.Fatal(label)
	}
}

func main() {
	order := binary.LittleEndian
	custom := errors.New("EOF")
	wrapped := fmt.Errorf("wrapped: %w", io.EOF)
	retries := []struct {
		data           []byte
		terminal, want error
		requests       []int
		value          uint32
	}{
		{nil, io.EOF, io.EOF, []int{4}, 99},
		{[]byte{19, 20}, io.EOF, io.ErrUnexpectedEOF, []int{4, 3}, 99},
		{[]byte{19, 20}, custom, custom, []int{4, 3}, 99},
		{[]byte{19, 20}, wrapped, wrapped, []int{4, 3}, 99},
		{[]byte{1, 2, 3, 4}, custom, nil, []int{4, 3, 2, 1}, 67305985},
	}
	for i, c := range retries {
		value := uint32(99)
		r := &reader{rest: c.data, stop: c.terminal}
		err := binary.Read(r, order, &value)
		must(err == c.want && value == c.value && slices.Equal(r.requests, c.requests), fmt.Sprintf("reader control %d", i))
	}
	r := &reader{stop: custom}
	must(binary.Read(r, order, uint16(99)) == custom && slices.Equal(r.requests, []int{2}), "ordinary scalar I/O happens before invalid type")
	r = &reader{stop: custom}
	err := binary.Read(r, order, reviewByte(99))
	must(err != nil && err != custom && len(r.requests) == 0, "defined scalar fails before I/O")
	r = &reader{stop: custom}
	must(binary.Read(r, order, []uint16(nil)) == nil && len(r.requests) == 0, "zero width invokes reader")

	original := []byte{9, 11}
	w := &writer{original: &original[0]}
	must(binary.Write(w, order, original) == nil && w.alias && w.calls == 1 && slices.Equal(w.payload, original), "byte alias and ignored successful short count")
	defined := []reviewByte{9, 11}
	w = &writer{original: (*byte)(&defined[0]), stop: custom}
	must(binary.Write(w, order, defined) == custom && !w.alias && w.calls == 1 && slices.Equal(w.payload, original), "defined elements and writer identity")
	w = &writer{stop: custom}
	must(binary.Write(w, order, []uint16{}) == custom && w.calls == 1 && len(w.payload) == 0, "empty payload must reach writer")
	w = &writer{stop: custom}
	err = binary.Write(w, order, "unsupported")
	must(err != nil && err != custom && w.calls == 0, "bad type reaches writer")

	cases := []struct{ wire, want string }{
		{"6 1 0 1 2 1 0 0 99 19 20", "unexpected:99:4,3|ok:1:false:99,0,0,0"},
		{"6 1 0 1 2 1 2 0 99 19 20", "wrapped:99:4,3|ok:1:false:99,0,0,0"},
		{"6 1 0 1 4 1 1 0 99 1 2 3 4", "ok:67305985:4,3,2,1|ok:1:false:99,0,0,0"},
		{"4 0 0 1 0 1 1 0 99", "other:99:2|ok:1:false:99,0"},
		{"11 0 0 1 0 1 1 0 99", "type:99:|ok:1:false:99"},
		{"2 2 0 2 2 1 0 0 9 11 19 20", "ok:19,20:2,1|ok:1:true:9,11"},
		{"11 2 0 2 2 1 0 1 9 11 19 20", "ok:19,20:2,1|other:1:false:9,11"},
		{"4 2 0 0 0 1 1 1", "ok::|other:1:false:"},
		{"2 3 0 0 0 1 1 1", "type::|type:0:false:"},
	}
	var wire strings.Builder
	for _, c := range cases {
		fmt.Fprintln(&wire, c.wire)
	}
	for i, got := range check.LeanRows("godiff-binary-public-stream", len(cases), strings.NewReader(wire.String())) {
		must(got == cases[i].want, fmt.Sprintf("case %d: Lean %q want %q", i, got, cases[i].want))
	}
	fmt.Printf("Public stream: %d independent Lean rows and pinned direct I/O controls\n", len(cases))
}
