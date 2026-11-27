package bytes

import (
	"fmt"
	"strconv"
	"strings"
)

func GoLeanGrowSlice(raw []byte, capacity, count int, nilInput bool) string {
	var input []byte
	if !nilInput {
		backing := make([]byte, capacity)
		for i := range backing {
			backing[i] = 0x99
		}
		copy(backing, raw)
		input = backing[:len(raw)]
	}
	got := growSlice(input, count)
	if got == nil {
		return "nil"
	}
	if cap(input) > 0 && &input[:cap(input)][0] == &got[:cap(got)][0] {
		panic("growSlice returned original backing")
	}
	wire := func(raw []byte) string {
		fields := make([]string, len(raw))
		for i, b := range raw {
			fields[i] = strconv.Itoa(int(b))
		}
		return strings.Join(fields, ",")
	}
	return fmt.Sprintf("%d:%s:%s", cap(got), wire(got), wire(got[len(got):cap(got)]))
}
