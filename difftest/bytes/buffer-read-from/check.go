package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type step struct{ Mode, Count, Writes, Value int }

type input struct {
	Tag, Len, Off, Cap, Flag int
	NilReader                bool
	Steps                    []step
}

type fixture struct {
	Input input
	Want  string
}

func wire(values []byte) string {
	var out []string
	for len(values) > 0 {
		b := values[0]
		end := 1
		for end < len(values) && values[end] == b {
			end++
		}
		out = append(out, fmt.Sprintf("%d.%d", b, end))
		values = values[end:]
	}
	return strings.Join(out, ",")
}

func expected(c input) string {
	if c.Tag == 0 {
		return "nilReceiver"
	}
	length, offset, capacity := c.Len, c.Off, c.Cap
	data := append(bytes.Repeat([]byte{11}, offset), bytes.Repeat([]byte{22}, length-offset)...)
	data = append(data, bytes.Repeat([]byte{33}, capacity-length)...)
	var trace []string
	publish := func(status string, count int) string {
		return fmt.Sprintf("%s|%d|slice|%d|0|%s|%s|%s", status, count, offset,
			wire(data[:length]), wire(data[length:]), strings.Join(trace, ";"))
	}
	total := 0
	for i := 0; ; i++ {
		unread := length - offset
		if unread == 0 {
			length, offset = 0, 0
		}
		if capacity-length < 512 {
			if unread+512 <= capacity/2 {
				copy(data, data[offset:length])
			} else {
				request := max(length+512, 2*(capacity-offset))
				// Measure the pinned runtime's capacity class independently of
				// the Buffer method and the Lean table.
				next := append([]byte(nil), make([]byte, request)...)
				next = next[:cap(next)]
				copy(next, data[offset:length])
				data, capacity = next, len(next)
			}
			length, offset = unread, 0
		}
		if c.NilReader {
			return publish("nilReader", 0)
		}
		p := data[length:]
		trace = append(trace, fmt.Sprintf("%d:%s", len(p), wire(p)))
		reply := c.Steps[i]
		for j := range min(len(p), reply.Writes) {
			p[j] = byte(reply.Value)
		}
		switch {
		case reply.Mode == 3:
			return publish("readerPanic", 0)
		case reply.Count < 0:
			return publish("negativeRead", 0)
		case reply.Count > len(p):
			return publish("bounds", 0)
		}
		length += reply.Count
		total += reply.Count
		switch reply.Mode {
		case 1:
			return publish("ok", total)
		case 2:
			return publish("readError", total)
		}
	}
}

func batch(label string, cases []input, encoder *json.Encoder) {
	const shard = 512
	for first := 0; first < len(cases); first += shard {
		current := cases[first:min(first+shard, len(cases))]
		var protocol bytes.Buffer
		rows := make([]fixture, len(current))
		for i, c := range current {
			nilReader := 0
			if c.NilReader {
				nilReader = 1
			}
			fmt.Fprintf(&protocol, "%d %d %d %d %d %d %d", c.Tag, c.Len, c.Off, c.Cap, c.Flag, nilReader, math.MaxInt)
			for _, s := range c.Steps {
				fmt.Fprintf(&protocol, " %d %d %d %d", s.Mode, s.Count, s.Writes, s.Value)
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, got := range check.LeanRows("godiff-bytes-buffer-read-from", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d (%+v): got %q want %q", label, first+i, current[i], got, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("Buffer.ReadFrom %s: %d Lean rows agree.\n", label, len(cases))
}

func (c input) withSteps(steps ...step) input {
	c.Steps = steps
	return c
}

func main() {
	geometries := []input{{Tag: 1, Flag: 3}, {Tag: 2, Flag: 1}}
	for _, capacity := range []int{511, 512, 513, 575, 576, 639, 640, 1023, 1024, 1025, 1200, 1536, 2048} {
		for _, length := range []int{0, 1, max(0, capacity-512), capacity - 1, capacity} {
			for _, offset := range []int{0, length / 2, max(0, length-1), length} {
				geometries = append(geometries, input{Tag: 2, Cap: capacity, Len: length, Off: offset, Flag: 4})
			}
		}
	}
	edges := []input{{}, {NilReader: true}}
	for _, g := range geometries {
		missing := g
		missing.NilReader = true
		edges = append(edges, missing)
		for _, s := range []step{{1, 0, 0, 0}, {1, 0, 1, 77}, {1, 1, 4, 77}, {2, 1, 4, 98}, {3, 0, 7, 41},
			{1, -1, 4, 97}, {2, math.MinInt, 4, 97}, {1, math.MaxInt, 4, 97}, {1, 512, 503, 91}} {
			edges = append(edges, g.withSteps(s), g.withSteps(step{0, 0, 3, 93}, step{0, 17, 25, 32}, s))
		}
		edges = append(edges, g.withSteps(step{0, 512, 512, 15}, step{0, 512, 527, 16}, step{1, 9, 2, 17}))
	}
	file, err := os.CreateTemp("", "golean-buffer-read-from-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			capacity := rng.Intn(1601)
			length := rng.Intn(capacity + 1)
			c := input{Tag: 2, Cap: capacity, Len: length, Off: rng.Intn(length + 1), Flag: rng.Intn(6) - 1, NilReader: i%23 == 0}
			switch {
			case i%97 == 0:
				c = input{}
			case i%19 == 0:
				c = input{Tag: 1, Flag: 4}
			}
			for range rng.Intn(3) {
				c.Steps = append(c.Steps, step{0, rng.Intn(513), rng.Intn(650), rng.Intn(256)})
			}
			count := rng.Intn(601)
			if i%13 == 0 {
				count = -1
			}
			c.Steps = append(c.Steps, step{1 + rng.Intn(3), count, rng.Intn(650), rng.Intn(256)})
			cases[i] = c
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_BUFFER_READ_FROM_CASES", file.Name()))
	const data = "difftest/bytes/buffer-read-from/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReadFrom", data+"bridge_test.go", data+"buffer_read_from_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
