package binary

import "reflect"

func GoLeanReflectTypeSize(t reflect.Type) int { return sizeof(t) }

func GoLeanReflectValueSize(v reflect.Value) (cold, warm int) {
	structSize.Clear()
	cold = dataSize(v)
	warm = dataSize(v)
	return
}
