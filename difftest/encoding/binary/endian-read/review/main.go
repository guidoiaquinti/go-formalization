package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func observe(f func() uint64) (out string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || !strings.HasPrefix(err.Error(), "runtime error: index out of range") {
				panic(value)
			}
			out = "panic"
		}
	}()
	return strconv.FormatUint(f(), 10)
}
func main() {
	cases := []struct {
		raw  []byte
		want [6]string
	}{
		{nil, [6]string{"panic", "panic", "panic", "panic", "panic", "panic"}},
		{[]byte{}, [6]string{"panic", "panic", "panic", "panic", "panic", "panic"}},
		{[]byte{0xff}, [6]string{"panic", "panic", "panic", "panic", "panic", "panic"}},
		{[]byte{0x80, 1}, [6]string{"32769", "panic", "panic", "384", "panic", "panic"}},
		{[]byte{0x80, 1, 0xff}, [6]string{"32769", "panic", "panic", "384", "panic", "panic"}},
		{[]byte{1, 2, 3, 4, 255}, [6]string{"258", "16909060", "panic", "513", "67305985", "panic"}},
		{[]byte{1, 2, 3, 4, 5, 6, 7}, [6]string{"258", "16909060", "panic", "513", "67305985", "panic"}},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8}, [6]string{"258", "16909060", "72623859790382856", "513", "67305985", "578437695752307201"}},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 0xff}, [6]string{"258", "16909060", "72623859790382856", "513", "67305985", "578437695752307201"}},
	}
	var model bytes.Buffer
	for i, c := range cases {
		storage := append([]byte{0xff, 0xfe, 0xfd}, c.raw...)
		storage = append(storage, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22)
		view := storage[3 : 3+len(c.raw)]
		if c.raw == nil {
			view = nil
		}
		before := bytes.Clone(storage)
		observed := [6]string{
			observe(func() uint64 { return uint64(binary.BigEndian.Uint16(view)) }),
			observe(func() uint64 { return uint64(binary.BigEndian.Uint32(view)) }),
			observe(func() uint64 { return binary.BigEndian.Uint64(view) }),
			observe(func() uint64 { return uint64(binary.LittleEndian.Uint16(view)) }),
			observe(func() uint64 { return uint64(binary.LittleEndian.Uint32(view)) }),
			observe(func() uint64 { return binary.LittleEndian.Uint64(view) }),
		}
		if observed != c.want || !bytes.Equal(storage, before) {
			panic(fmt.Sprintf("native row %d: got %v want %v", i, observed, c.want))
		}
		fmt.Fprint(&model, "0")
		for _, b := range view {
			fmt.Fprintf(&model, " %d", b)
		}
		fmt.Fprintln(&model)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-endian-read")
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		panic("unexpected output count")
	}
	for i, c := range cases {
		if rows[i] != strings.Join(c.want[:], "|") {
			panic(fmt.Sprintf("Lean row %d got %q want %v", i, rows[i], c.want))
		}
	}
	fmt.Printf("encoding/binary endian getters passed %d independent native and model cases\n", len(cases))
}
