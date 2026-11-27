package binary

import (
	"fmt"
	"reflect"
)

func GoLeanCoderSkip(v reflect.Value, backing []byte, length, offset int) (read, write int, status string) {
	buf := backing[:length:len(backing)]
	d := decoder{buf: buf, offset: offset}
	d.skip(v)
	e := encoder{buf: buf, offset: offset}
	defer func() {
		read, write = d.offset, e.offset
		if failure := recover(); failure != nil {
			status = fmt.Sprint(failure)
		}
	}()
	e.skip(v)
	status = "ok"
	return
}
