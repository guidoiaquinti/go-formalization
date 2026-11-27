package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	N, Hist uint16
	Filled  int
	Lit     byte
	Data    []byte
	Full    bool
}

type expectedState struct {
	N             uint16
	Filled        int
	Literal       [256]uint16
	Extra, Offset [32]uint16
	Writes        map[uint16]uint32
}

type comparison struct {
	Case    fixture
	Indices []uint16
	States  []expectedState
}

func oracle(c fixture) comparison {
	indices := make([]uint16, 0, 3+len(c.Data))
	indices = append(indices, 0, c.N, 65535)
	for i := range c.Data {
		indices = append(indices, uint16(int(c.N)+i))
	}
	makeState := func() expectedState {
		s := expectedState{N: c.N, Filled: c.Filled, Writes: make(map[uint16]uint32)}
		for i := range s.Literal {
			s.Literal[i] = c.Hist
		}
		for i := range s.Extra {
			s.Extra[i], s.Offset[i] = c.Hist, c.Hist
		}
		return s
	}
	reset := makeState()
	if c.N != 0 {
		reset = expectedState{}
	}
	literal := makeState()
	literal.N = uint16(uint32(c.N) + 1)
	literal.Literal[c.Lit] = uint16(uint32(c.Hist) + 1)
	literal.Writes[c.N] = uint32(c.Lit)
	end := makeState()
	end.N = literal.N
	end.Extra[0] = uint16(uint32(c.Hist) + 1)
	end.Writes[c.N] = 256
	emitted := makeState()
	emitted.N = uint16(int(c.N) + len(c.Data))
	var frequencies [256]uint64
	for i, b := range c.Data {
		frequencies[b]++
		emitted.Writes[uint16(int(c.N)+i)] = uint32(b)
	}
	for i, count := range frequencies {
		emitted.Literal[i] = uint16(uint64(c.Hist) + count)
	}
	return comparison{c, indices, []expectedState{reset, literal, end, emitted}}
}

func encode(c comparison) string {
	parts := make([]string, 0, len(c.States)*(6+len(c.Case.Data)+len(c.Indices))+6)
	add := func(n int64) { parts = append(parts, strconv.FormatInt(n, 10)) }
	for _, s := range c.States {
		add(int64(s.N))
		add(int64(s.Filled))
		add(int64(s.Literal[c.Case.Lit]))
		add(int64(s.Extra[0]))
		add(int64(s.Offset[0]))
		add(int64(s.Literal[c.Case.Lit]))
		for _, b := range c.Case.Data {
			add(int64(s.Literal[b]))
		}
		for _, i := range c.Indices {
			v, changed := s.Writes[i]
			if !changed {
				v = 65536 + uint32(i)
			}
			add(int64(v))
		}
	}
	for _, value := range []int64{int64(c.Case.N), 65536, 131071, 73, 73, 29} {
		add(value)
	}
	return strings.Join(parts, " ")
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-token-state")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.N, " ", c.Filled, " ", c.Hist, " ", c.Lit)
			for _, b := range c.Data {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
		}
		for i, got := range check.LeanRows("godiff-flate-token-state", len(cases), &input) {
			want := oracle(cases[i])
			if row := encode(want); got != row {
				log.Fatalf("%s row %d: Lean %q, oracle %q", label, i, got, row)
			}
			must(enc.Encode(want))
		}
		fmt.Printf("%s: %d token-state rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, n := range []uint16{0, 1, 2, 32768, 65531, 65534, 65535} {
		for _, h := range []uint16{0, 1, 65534, 65535} {
			for _, data := range [][]byte{nil, {}, {0}, {255}, {0, 255, 0, 255, 0, 255}} {
				for _, literal := range []byte{0, 255} {
					fixed = append(fixed, fixture{n, h, -27, literal, data, true})
				}
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for batch := range 10 {
			cases := make([]fixture, 2000)
			for i := range cases {
				n := uint16(r.Uint32())
				if i%4 == 0 {
					n = ^uint16(r.Intn(16))
				}
				data := make([]byte, r.Intn(17))
				_, err := r.Read(data)
				must(err)
				cases[i] = fixture{n, uint16(r.Uint32()), int(int32(r.Uint32())), byte(r.Uint32()), data, false}
			}
			compare(fmt.Sprintf("seed %d batch %d", seed, batch), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_TOKEN_STATE_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanTokenState", "difftest/compress/flate/token-state/testdata/token_state_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
