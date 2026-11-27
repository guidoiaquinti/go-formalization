package binary

import "fmt"

func GoLeanCoderByte(method, length, position int, backing []byte, input byte, isNil bool) (result string, offset int) {
	d := &decoder{buf: backing[:length:len(backing)], offset: position}
	e := &encoder{buf: backing[:length:len(backing)], offset: position}
	if isNil {
		d, e = nil, nil
	}
	defer func() {
		if failure := recover(); failure != nil {
			result = fmt.Sprint(failure)
		}
		if d != nil && method < 3 {
			offset = d.offset
		} else if e != nil {
			offset = e.offset
		}
	}()
	switch method {
	case 0:
		result = fmt.Sprint(d.bool())
	case 1:
		result = fmt.Sprint(d.uint8())
	case 2:
		result = fmt.Sprint(d.int8())
	case 3:
		e.bool(input != 0)
		result = "ok"
	case 4:
		e.uint8(input)
		result = "ok"
	case 5:
		e.int8(int8(input))
		result = "ok"
	default:
		panic("unknown byte coder")
	}
	return
}
