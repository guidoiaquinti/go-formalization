package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const goldenHash = "bfec3d418b829afe5df2d8887d1508348409c293b73758d7efd841dd995fe021"
const firstGolden = 0xb773b6063d4616a5
const secondGolden = 0x1160af22a66abc3c

type fixture struct {
	seed, replacement                 [32]byte
	draw, first, second, front, spare int
}

func main() {
	golden := [32]byte([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456"))
	inverse := golden
	for i := range inverse {
		inverse[i] ^= 0xff
	}
	cases := []fixture{
		{seed: golden, replacement: golden, draw: 372},
		{seed: golden, replacement: golden, first: 2976},
	}
	for _, spare := range []int{0, 7, 8, 9, 12, 13, 60, 61} {
		cases = append(cases, fixture{seed: golden, replacement: golden, first: 1, second: 3, front: 2, spare: spare})
	}
	for _, draws := range []int{31, 32, 95, 96, 123, 124, 247, 248} {
		cases = append(cases, fixture{seed: inverse, replacement: golden, draw: draws, first: 7, second: 9, front: 3, spare: 64})
	}
	var input bytes.Buffer
	var want []string
	for i, c := range cases {
		fmt.Fprint(&input, c.draw, " ", c.first, " ", c.second, " ", c.front)
		for _, key := range [][32]byte{c.seed, c.replacement} {
			for off := 0; off < 32; off += 8 {
				fmt.Fprint(&input, " ", binary.LittleEndian.Uint64(key[off:off+8]))
			}
		}
		for range c.front + c.spare {
			fmt.Fprint(&input, " 165")
		}
		fmt.Fprintln(&input)
		want = append(want, native(i, c, golden))
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(input.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-rand-v2-chacha-stream")
	cmd.Stdin = &input
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(got) != len(want) {
		log.Fatalf("got %d rows, want %d", len(got), len(want))
	}
	for i, row := range got {
		if row != want[i] {
			log.Fatalf("Lean row %d differs", i)
		}
	}
	checkZeroValue(golden)
	fmt.Printf("ChaCha8 review: %d pinned Go and Lean rows, golden transcript, and native seed/reset passed\n", len(cases))
}

func native(row int, c fixture, golden [32]byte) string {
	r := rand.NewChaCha8(c.seed)
	marshal := func() []byte {
		b, err := r.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}
		return b
	}
	initial := marshal()
	if !bytes.Equal(initial, plainState(c.seed, 0)) {
		log.Fatal("fresh state serialization differed")
	}
	drawn := make([]string, c.draw)
	var transcript []byte
	for i := range drawn {
		x := r.Uint64()
		drawn[i] = strconv.FormatUint(x, 10)
		if row == 0 {
			transcript = binary.LittleEndian.AppendUint64(transcript, x)
		}
	}
	read := func(length int) []byte {
		data := make([]byte, length)
		n, err := r.Read(data)
		if n != length || err != nil {
			log.Fatalf("Read(%d) = %d, %v", length, n, err)
		}
		return data
	}
	a := read(c.first)
	middle := r.Uint64()
	b := read(c.second)
	marshaled := marshal()
	if row == 0 && fmt.Sprintf("%x", sha256.Sum256(transcript)) != goldenHash {
		log.Fatal("golden Uint64 transcript differed")
	}
	if row == 1 && fmt.Sprintf("%x", sha256.Sum256(a)) != goldenHash {
		log.Fatal("golden Read transcript differed")
	}
	if row >= 2 && row < 10 {
		word := binary.LittleEndian.AppendUint64(nil, firstGolden)
		want := append([]byte("readbuf:\x04"), word[4:]...)
		want = append(want, plainState(golden, 2)...)
		if !bytes.Equal(a, word[:1]) || !bytes.Equal(b, word[1:4]) || middle != secondGolden ||
			!bytes.Equal(marshaled, want) {
			log.Fatal("mixed stream buffer layout differed")
		}
	}
	backing := bytes.Repeat([]byte{165}, c.front+c.spare)
	appended, err := r.AppendBinary(backing[:c.front])
	if err != nil {
		log.Fatal(err)
	}
	reused := len(backing) > 0 && &appended[0] == &backing[0]
	if !bytes.Equal(appended[:c.front], bytes.Repeat([]byte{165}, c.front)) ||
		!bytes.Equal(appended[c.front:], marshaled) {
		log.Fatal("append payload or front differed")
	}
	if row >= 2 && row < 10 {
		written := 0
		for _, n := range []int{8, 1, 4, 48} {
			if written+n > c.spare {
				break
			}
			written += n
		}
		want := bytes.Repeat([]byte{165}, len(backing))
		copy(want[c.front:], marshaled[:written])
		if reused != (written == len(marshaled)) || !bytes.Equal(backing, want) {
			log.Fatal("append staging differed")
		}
	}
	r.Seed(c.replacement)
	reset := marshal()
	resetWord := r.Uint64()
	if !bytes.Equal(reset, plainState(golden, 0)) || resetWord != firstGolden {
		log.Fatal("Seed failed to reset counter, buffered bytes or output")
	}
	return strings.Join([]string{wire(initial), strings.Join(drawn, ","), wire(a), strconv.FormatUint(middle, 10),
		wire(b), wire(marshaled), fmt.Sprintf("%t:%s:%s", reused, wire(appended), wire(backing)), wire(reset),
		strconv.FormatUint(resetWord, 10)}, "|")
}

func plainState(seed [32]byte, position uint64) []byte {
	out := binary.BigEndian.AppendUint64([]byte("chacha8:"), position)
	return append(out, seed[:]...)
}

func wire(data []byte) string {
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func checkZeroValue(seed [32]byte) {
	var zero rand.ChaCha8
	initial := zero.Uint64()
	if initial == rand.NewChaCha8([32]byte{}).Uint64() {
		log.Fatal("pinned zero receiver unexpectedly behaved as initialized with a zero seed")
	}
	var reset rand.ChaCha8
	reset.Seed(seed)
	if reset.Uint64() != firstGolden {
		log.Fatal("Seed failed from a previously unseeded receiver")
	}
}
