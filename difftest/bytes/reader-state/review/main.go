package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	cases := []struct{ row, want string }{
		{"0 0 9223372036854775807 4", "nil:0:-1::|nil:0:-1::|panic|0|panic"},
		{"1 0 0 -1", "slice:0:-1::|slice:0:-1::|panic|0|panic"},
		{"1 0 4 0 255 0", "slice:0:-1::255,0|slice:0:-1::255,0|panic|0|panic"},
		{"1 3 9223372036854775807 1 195 169 255 9 8", "slice:0:-1:195,169,255:9,8|slice:0:-1:195,169,255:9,8|panic|3|panic"},
		{"1 1 1 0 0 255", "slice:0:-1:0:255|slice:0:-1:0:255|panic|1|panic"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.row)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-reader-setup")
	cmd.Stdin = &input
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(lines) != len(cases) {
		log.Fatalf("got %d Lean outcomes", len(lines))
	}
	for i, got := range lines {
		if got != cases[i].want {
			log.Fatalf("Lean fixture %d: got %q want %q", i, got, cases[i].want)
		}
	}
	checkPublic()
	fmt.Printf("Reader state reviewer control: %d Lean fixtures and native constructor/reset/size interactions passed.\n", len(cases))
}

func checkPublic() {
	old := []byte{9, 195, 169, 255, 7}
	r := bytes.NewReader(old[1:4])
	old[2] = 168
	ch, width, err := r.ReadRune()
	if err != nil || ch != 'è' || width != 2 || r.Size() != 3 || r.Len() != 1 {
		log.Fatalf("constructor did not adopt live input or report size: %q %d %v", ch, width, err)
	}
	if err := r.UnreadRune(); err != nil || r.Len() != 3 {
		log.Fatalf("Size changed unread history: %v", err)
	}
	if _, _, err := r.ReadRune(); err != nil {
		log.Fatal(err)
	}
	fresh := []byte{23, 15, 244, 9}
	r.Reset(fresh[1:3:4])
	if r.Size() != 2 || r.Len() != 2 || r.UnreadRune() == nil {
		log.Fatal("Reset retained previous offset or rune state")
	}
	fresh[1] = 90
	old[1] = 0
	if value, err := r.ReadByte(); err != nil || value != 90 {
		log.Fatalf("Reset did not adopt replacement: %d %v", value, err)
	}
	if _, err := r.Seek(math.MaxInt64, io.SeekStart); err != nil || r.Size() != 2 || r.Len() != 0 {
		log.Fatalf("size at past-end cursor: %v", err)
	}
	r.Reset(nil)
	if value, err := r.ReadByte(); value != 0 || err != io.EOF || r.Size() != 0 {
		log.Fatalf("reset to nil gave %d %v", value, err)
	}
	var nilReader *bytes.Reader
	mustNilPanic(func() { nilReader.Reset(nil) })
	mustNilPanic(func() { nilReader.Size() })
}

func mustNilPanic(f func()) {
	defer func() {
		got, ok := recover().(runtime.Error)
		if !ok || got.Error() != "runtime error: invalid memory address or nil pointer dereference" {
			log.Fatal("missing nil reader panic")
		}
	}()
	f()
}
