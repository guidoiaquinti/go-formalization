package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	randv2 "math/rand/v2"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	draws, first, second, front int
	seed, replacement           [4]uint64
	backing                     []byte
}

func main() {
	rng := rand.New(rand.NewSource(17))
	var fixed []fixture
	for _, count := range []int{0, 1, 30, 31, 32, 62, 63, 64, 94, 95, 96, 122, 123, 124, 125, 126, 127, 248} {
		for a := range 10 {
			for _, spare := range []int{0, 7, 8, 9, 10, 15, 47, 48, 59, 64} {
				c := sample(rng, 0)
				c.draws, c.first, c.second, c.front = count, a, a+1, 1
				c.backing = make([]byte, 1+spare)
				rng.Read(c.backing)
				fixed = append(fixed, c)
			}
		}
	}
	check.CompareSamples("godiff-rand-v2-chacha-stream", fixed, sample, encode, native, reference)
}

func sample(r *rand.Rand, _ int) fixture {
	c := fixture{draws: r.Intn(4), first: r.Intn(15), second: r.Intn(15), front: r.Intn(10)}
	c.backing = make([]byte, c.front+r.Intn(75))
	r.Read(c.backing)
	for i := range c.seed {
		c.seed[i] = r.Uint64()
		c.replacement[i] = r.Uint64()
	}
	return c
}

func seedBytes(words [4]uint64) (seed [32]byte) {
	for i, word := range words {
		binary.LittleEndian.PutUint64(seed[8*i:], word)
	}
	return
}

func freshState(words [4]uint64) []byte {
	state := make([]byte, 48)
	copy(state, "chacha8:")
	seed := seedBytes(words)
	copy(state[16:], seed[:])
	return state
}

func marshal(c *randv2.ChaCha8) []byte {
	state, err := c.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	return state
}

func read(c *randv2.ChaCha8, count int) []byte {
	out := make([]byte, count)
	if n, err := c.Read(out); n != count || err != nil {
		log.Fatalf("Read requested %d and returned %d, %v", count, n, err)
	}
	return out
}

func native(f fixture) string {
	c := randv2.NewChaCha8(seedBytes(f.seed))
	initial := marshal(c)
	words := make([]uint64, f.draws)
	for i := range words {
		words[i] = c.Uint64()
	}
	first := read(c, f.first)
	middle := c.Uint64()
	second := read(c, f.second)
	serialized := marshal(c)
	backing := slices.Clone(f.backing)
	out, err := c.AppendBinary(backing[:f.front:len(backing)])
	if err != nil || !slices.Equal(serialized, marshal(c)) {
		log.Fatalf("AppendBinary returned %v or changed the receiver", err)
	}
	reused := len(backing) > 0 && &out[0] == &backing[0]
	c.Seed(seedBytes(f.replacement))
	afterSeed := marshal(c)
	return wire(initial, words, first, middle, second, serialized, reused, out, backing, afterSeed, c.Uint64())
}

func reference(f fixture) string {
	c := randv2.NewChaCha8(seedBytes(f.seed))
	draws := make([]uint64, f.draws)
	for i := range draws {
		draws[i] = c.Uint64()
	}
	var pending []byte
	collect := func(n int) []byte {
		for len(pending) < n {
			pending = binary.LittleEndian.AppendUint64(pending, c.Uint64())
		}
		out := slices.Clone(pending[:n])
		pending = pending[n:]
		return out
	}
	a := collect(f.first)
	middle := c.Uint64()
	b := collect(f.second)
	rawState := marshal(c)
	if len(rawState) != 48 || string(rawState[:8]) != "chacha8:" {
		log.Fatal("unexpected state format")
	}
	var chunks [][]byte
	if len(pending) > 0 {
		chunks = append(chunks, []byte("readbuf:"), []byte{byte(len(pending))}, pending)
	}
	chunks = append(chunks, rawState)
	payload := slices.Concat(chunks...)
	old := slices.Clone(f.backing)
	position := f.front
	for _, chunk := range chunks {
		if position+len(chunk) > len(old) {
			break
		}
		copy(old[position:], chunk)
		position += len(chunk)
	}
	out := append(slices.Clone(f.backing[:f.front]), payload...)
	reset := randv2.NewChaCha8(seedBytes(f.replacement))
	return wire(freshState(f.seed), draws, a, middle, b, payload, len(out) <= len(old), out, old, freshState(f.replacement), reset.Uint64())
}

func encode(f fixture) string {
	var row strings.Builder
	fmt.Fprintf(&row, "%d %d %d %d", f.draws, f.first, f.second, f.front)
	for _, block := range [][4]uint64{f.seed, f.replacement} {
		for _, w := range block {
			fmt.Fprintf(&row, " %d", w)
		}
	}
	for _, b := range f.backing {
		fmt.Fprintf(&row, " %d", b)
	}
	return row.String()
}

func wire(initial []byte, words []uint64, first []byte, middle uint64, second, state []byte,
	reused bool, out, old, reset []byte, final uint64) string {
	bytesWire := func(data []byte) string {
		parts := make([]string, len(data))
		for i, cell := range data {
			parts[i] = strconv.Itoa(int(cell))
		}
		return strings.Join(parts, ",")
	}
	parts := make([]string, len(words))
	for i, word := range words {
		parts[i] = strconv.FormatUint(word, 10)
	}
	return strings.Join([]string{
		bytesWire(initial), strings.Join(parts, ","), bytesWire(first),
		strconv.FormatUint(middle, 10), bytesWire(second), bytesWire(state),
		fmt.Sprintf("%t:%s:%s", reused, bytesWire(out), bytesWire(old)),
		bytesWire(reset), strconv.FormatUint(final, 10),
	}, "|")
}
