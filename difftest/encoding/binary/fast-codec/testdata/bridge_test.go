package binary

import "fmt"

func GoLeanFastDecode(bs []byte, order ByteOrder, v any) (status string) {
	defer func() {
		if p := recover(); p != nil {
			status = fmt.Sprint(p)
		}
	}()
	if decodeFast(bs, order, v) {
		return "ok"
	}
	return "unsupported"
}

func GoLeanFastEncode(bs []byte, order ByteOrder, v any) (status string) {
	defer func() {
		if p := recover(); p != nil {
			status = fmt.Sprint(p)
		}
	}()
	encodeFast(bs, order, v)
	return "ok"
}
