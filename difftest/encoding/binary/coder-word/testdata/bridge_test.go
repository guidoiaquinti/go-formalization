package binary

import (
	"fmt"
	"strconv"
)

func GoLeanCoderWord(order ByteOrder, width, method, length, offset int, backing []byte, word uint64) (value, status string, after int) {
	buf := backing[:length:len(backing)]
	d := decoder{order: order, buf: buf, offset: offset}
	e := encoder{order: order, buf: buf, offset: offset}
	value, status = "0", "ok"
	defer func() {
		if p := recover(); p != nil {
			status = fmt.Sprint(p)
		}
		after = d.offset
		if method >= 2 {
			after = e.offset
		}
	}()
	switch width {
	case 2:
		switch method {
		case 0:
			value = strconv.FormatUint(uint64(d.uint16()), 10)
		case 1:
			value = strconv.FormatInt(int64(d.int16()), 10)
		case 2:
			e.uint16(uint16(word))
		case 3:
			e.int16(int16(word))
		}
	case 4:
		switch method {
		case 0:
			value = strconv.FormatUint(uint64(d.uint32()), 10)
		case 1:
			value = strconv.FormatInt(int64(d.int32()), 10)
		case 2:
			e.uint32(uint32(word))
		case 3:
			e.int32(int32(word))
		}
	case 8:
		switch method {
		case 0:
			value = strconv.FormatUint(d.uint64(), 10)
		case 1:
			value = strconv.FormatInt(d.int64(), 10)
		case 2:
			e.uint64(word)
		case 3:
			e.int64(int64(word))
		}
	}
	return
}
