package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type fixture struct {
	left, right           []byte
	leftRunes, rightRunes []rune
	fold                  bool
}

func main() {
	cases := []fixture{
		{nil, []byte{}, nil, nil, true},
		{[]byte("a\x00ZK"), []byte("A\x00zk"), []rune{'a', 0, 'Z', 'K'}, []rune{'A', 0, 'z', 'k'}, true},
		{[]byte{0xff}, []byte("�"), []rune{0xfffd}, []rune{0xfffd}, true},
		{[]byte{0xe2, 0x82, 0xef, 0xbf, 0xbd}, []byte("��"), []rune{0xfffd, 0xfffd, 0xfffd}, []rune{0xfffd, 0xfffd}, false},
		{[]byte{0xe2, 'A', 0x80}, []byte("�a�"), []rune{0xfffd, 'A', 0xfffd}, []rune{0xfffd, 'a', 0xfffd}, true},
		{[]byte("ß"), []byte("SS"), []rune{'ß'}, []rune{'S', 'S'}, false},
		{[]byte("İ"), []byte("i"), []rune{'İ'}, []rune{'i'}, false},
	}
	orbit := []rune{0x345, 0x399, 0x3b9, 0x1fbe}
	for _, left := range orbit {
		for _, right := range orbit {
			if left != right {
				cases = append(cases, fixture{[]byte(string(left)), []byte(string(right)), []rune{left}, []rune{right}, true})
			}
		}
	}
	var input bytes.Buffer
	var want []string
	for i, r := range orbit {
		next := orbit[(i+1)%len(orbit)]
		if got := unicode.SimpleFold(r); got != next {
			log.Fatalf("SimpleFold(%U) = %U, want %U", r, got, next)
		}
		fmt.Fprintln(&input, 0, r)
		want = append(want, strconv.Itoa(int(next)))
	}
	for _, c := range cases {
		fmt.Fprint(&input, 1, " ", len(c.left))
		for _, cells := range [][]byte{c.left, c.right} {
			for _, b := range cells {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		if got := bytes.EqualFold(c.left, c.right); got != c.fold {
			log.Fatalf("EqualFold(%x,%x) = %t, want %t", c.left, c.right, got, c.fold)
		}
		left := observeRunes(c.left, c.leftRunes)
		right := observeRunes(c.right, c.rightRunes)
		want = append(want, fmt.Sprintf("%t|%s|%s", c.fold, left, right))
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(input.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-fold-runes")
	cmd.Stdin = &input
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(got) != len(want) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(want))
	}
	for i, row := range got {
		if row != want[i] {
			log.Fatalf("Lean row %d = %q, want %q", i, row, want[i])
		}
	}
	fmt.Printf("bytes fold/runes reviewer control: %d pairs and %d fold successors passed\n", len(cases), len(orbit))
}

func observeRunes(input []byte, want []rune) string {
	storage := append(bytes.Clone(input), 0x55)
	var visible []byte
	if input != nil {
		visible = storage[:len(input)]
	}
	before := bytes.Clone(storage)
	runes := bytes.Runes(visible)
	if runes == nil || cap(runes) != len(want) || !slices.Equal(runes, want) {
		log.Fatalf("Runes(%x) = %U cap %d, want %U cap %d", input, runes, cap(runes), want, len(want))
	}
	values := make([]string, len(runes))
	for i, r := range runes {
		values[i] = strconv.Itoa(int(r))
		runes[i] = -1
	}
	if !bytes.Equal(storage, before) {
		log.Fatal("mutating decoded runes changed their source or spare byte")
	}
	return fmt.Sprintf("slice:%d:%s", cap(runes), strings.Join(values, ","))
}
