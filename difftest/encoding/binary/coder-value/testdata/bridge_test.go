package binary

import (
	"fmt"
	"reflect"
	"strings"
)

func GoLeanCoderValue(encoding bool, order ByteOrder, buf []byte, offset int, v reflect.Value) (ok bool, next int) {
	c := coder{order: order, buf: buf, offset: offset}
	if encoding {
		e := encoder(c)
		defer func() {
			next = e.offset
			if failure := recover(); failure != nil {
				goLeanCoderBounds(failure)
			}
		}()
		e.value(v)
	} else {
		d := decoder(c)
		defer func() {
			next = d.offset
			if failure := recover(); failure != nil {
				goLeanCoderBounds(failure)
			}
		}()
		d.value(v)
	}
	return true, offset
}

func goLeanCoderBounds(failure any) {
	msg := fmt.Sprint(failure)
	if !strings.HasPrefix(msg, "runtime error: slice bounds out of range") &&
		!strings.HasPrefix(msg, "runtime error: index out of range") {
		panic(failure)
	}
}
