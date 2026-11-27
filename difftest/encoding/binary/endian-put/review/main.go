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

func observed(out []byte, write func()) (status string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || !strings.HasPrefix(err.Error(), "runtime error: index out of range") {
				panic(value)
			}
			status = wire("panic", out)
		}
	}()
	write()
	return wire("ok", out)
}
func wire(status string, backing []byte) string {
	fields := make([]string, len(backing))
	for i, b := range backing {
		fields[i] = strconv.Itoa(int(b))
	}
	return status + ":" + strings.Join(fields, ",")
}
func main() {
	var model bytes.Buffer
	var expectations []string
	value := uint64(0x8070605040302010)
	encoded := [6][]byte{
		{0x20, 0x10}, {0x40, 0x30, 0x20, 0x10}, {0x80, 0x70, 0x60, 0x50, 0x40, 0x30, 0x20, 0x10},
		{0x10, 0x20}, {0x10, 0x20, 0x30, 0x40}, {0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
	}
	for _, size := range []int{0, 1, 2, 3, 4, 7, 8, 10} {
		for _, offset := range []int{0, 3, 15} {
			initial := bytes.Repeat([]byte{0xee}, offset+size+8)
			fmt.Fprintf(&model, "%d %d %d", value, offset, size)
			for _, b := range initial {
				fmt.Fprintf(&model, " %d", b)
			}
			fmt.Fprintln(&model)
			cells := make([]string, 6)
			for i, word := range encoded {
				actual := bytes.Clone(initial)
				dst := actual[offset : offset+size]
				write := []func(){
					func() { binary.BigEndian.PutUint16(dst, uint16(value)) },
					func() { binary.BigEndian.PutUint32(dst, uint32(value)) },
					func() { binary.BigEndian.PutUint64(dst, value) },
					func() { binary.LittleEndian.PutUint16(dst, uint16(value)) },
					func() { binary.LittleEndian.PutUint32(dst, uint32(value)) },
					func() { binary.LittleEndian.PutUint64(dst, value) },
				}[i]
				expected := bytes.Clone(initial)
				status := "panic"
				if size >= len(word) {
					copy(expected[offset:], word)
					status = "ok"
				}
				cells[i] = wire(status, expected)
				if got := observed(actual, write); got != cells[i] {
					panic(fmt.Sprintf("width %d offset %d writer %d got %q want %q", size, offset, i, got, cells[i]))
				}
			}
			expectations = append(expectations, strings.Join(cells, "|"))
		}
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-endian-put")
	cmd.Stdin = &model
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(expectations) {
		panic("wrong Lean row count")
	}
	for i, want := range expectations {
		if rows[i] != want {
			panic(fmt.Sprintf("row %d got %q want %q", i, rows[i], want))
		}
	}
	fmt.Printf("PutUint methods passed %d independent native and model geometry cases each\n", len(rows))
}
