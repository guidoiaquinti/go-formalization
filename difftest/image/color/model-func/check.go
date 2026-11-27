package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type token uint32

func (t token) RGBA() (r, g, b, a uint32) { return 0, 0, 0, 0xffff }

type callbackPanic uint32

type input struct {
	kind    int
	arg     token
	capture callbackPanic
}

func main() {
	var edges []input
	for kind := range 3 {
		for _, arg := range []token{0, 1, 0x7fff, 0x8000, 0xffff} {
			for _, capture := range []callbackPanic{0, 1, 0x7fff, 0x8000, 0xffff} {
				edges = append(edges, input{kind, arg, capture})
			}
		}
	}
	check.CompareSamples("godiff-image-color-model-func", edges, func(r *rand.Rand, i int) input {
		return input{i % 3, token(r.Uint32() >> 16), callbackPanic(r.Uint32() >> 16)}
	}, func(c input) string {
		return fmt.Sprintf("%d %d %d", c.kind, c.arg, c.capture)
	}, actual, expected)
}

func actual(c input) string {
	var f func(color.Color) color.Color
	var calls int
	var last token
	if c.kind != 0 {
		f = func(arg color.Color) color.Color {
			calls++
			last = arg.(token)
			if c.kind == 2 {
				panic(c.capture)
			}
			return last + token(c.capture)
		}
	}
	first, second := color.ModelFunc(f), color.ModelFunc(f)
	if first == nil || second == nil {
		return "nil interface"
	}
	for _, m := range []color.Model{first, second} {
		v := reflect.ValueOf(m)
		if v.Kind() != reflect.Pointer || v.Type().Elem().Name() != "modelFunc" ||
			v.Type().Elem().PkgPath() != "image/color" {
			return fmt.Sprintf("unexpected concrete type %T", m)
		}
		if v.IsNil() {
			return "typed nil"
		}
	}
	identity := "distinct"
	if first == second {
		identity = "same"
	}
	initialCalls := calls
	observe := func(m color.Model) string {
		last = 0
		priorCalls := calls
		tag, value := convert(m, c.arg)
		return fmt.Sprintf("%s %d %d %d", tag, value, calls-priorCalls, last)
	}
	return fmt.Sprintf("pointer nonnil %s %d %s %s", identity, initialCalls, observe(first), observe(second))
}

func convert(m color.Model, arg token) (tag string, value uint32) {
	defer func() {
		switch p := recover().(type) {
		case nil:
		case callbackPanic:
			tag, value = "callback-panic", uint32(p)
		case runtime.Error:
			tag = "nil-panic"
		default:
			tag = fmt.Sprintf("unexpected panic %T", p)
		}
	}()
	return "ok", uint32(m.Convert(arg).(token))
}

func expected(c input) string {
	result := "nil-panic 0 0 0"
	switch c.kind {
	case 1:
		result = fmt.Sprintf("ok %d 1 %d", c.arg+token(c.capture), c.arg)
	case 2:
		result = fmt.Sprintf("callback-panic %d 1 %d", c.capture, c.arg)
	}
	return fmt.Sprintf("pointer nonnil distinct 0 %s %s", result, result)
}
