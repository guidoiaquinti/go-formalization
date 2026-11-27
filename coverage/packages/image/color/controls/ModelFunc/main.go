package main

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"runtime"
)

type marker int

func (*marker) RGBA() (uint32, uint32, uint32, uint32) { panic("unexpected RGBA call") }

type payload struct{ id int }
type constructor func(func(color.Color) color.Color) color.Model

const (
	nilCallback = iota
	returningCallback
	panickingCallback
)

func main() {
	candidates := []struct {
		name  string
		fresh func() constructor
		want  int
	}{
		{"Go", func() constructor { return color.ModelFunc }, 0},
		{"reuse wrapper", reusedConstructor, 3},
		{"call during construction", func() constructor { return eagerConstructor }, 2},
		{"return nil for absent callback", func() constructor { return nilConstructor }, 1},
		{"replace stored callback", func() constructor { return ignoreConstructor }, 2},
	}
	for _, candidate := range candidates {
		count := 0
		for mode := nilCallback; mode <= panickingCallback; mode++ {
			if !matches(candidate.fresh(), mode) {
				count++
			}
		}
		if count != candidate.want {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, count, candidate.want)
		}
		fmt.Printf("%s disagreed on %d of 3 fixtures\n", candidate.name, count)
	}
}

func matches(makeModel constructor, mode int) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	arg, output := marker(7), marker(7)
	throw := &payload{10}
	var seen []color.Color
	var callback func(color.Color) color.Color
	if mode != nilCallback {
		callback = func(c color.Color) color.Color {
			seen = append(seen, c)
			if mode == panickingCallback {
				panic(throw)
			}
			return &output
		}
	}
	first, second := makeModel(callback), makeModel(callback)
	for _, got := range []color.Model{first, second} {
		v := reflect.ValueOf(got)
		if v.Kind() != reflect.Pointer || v.IsNil() ||
			v.Type().Elem().Name() != "modelFunc" || v.Type().Elem().PkgPath() != "image/color" {
			return false
		}
	}
	if first == second || len(seen) != 0 {
		return false
	}
	for _, model := range []color.Model{first, second} {
		before := len(seen)
		value, caught := observe(model, &arg)
		if mode == nilCallback {
			if _, ok := caught.(runtime.Error); !ok || len(seen) != before {
				return false
			}
			continue
		}
		if len(seen) != before+1 || seen[before] != &arg {
			return false
		}
		if mode == panickingCallback {
			if caught != throw {
				return false
			}
		} else if value != &output || caught != nil {
			return false
		}
	}
	return true
}

func observe(model color.Model, c color.Color) (result color.Color, caught any) {
	defer func() { caught = recover() }()
	return model.Convert(c), nil
}

func reusedConstructor() constructor {
	var saved color.Model
	return func(f func(color.Color) color.Color) color.Model {
		if saved == nil {
			saved = color.ModelFunc(f)
		}
		return saved
	}
}

func eagerConstructor(f func(color.Color) color.Color) color.Model {
	if f != nil {
		f(nil)
	}
	return color.ModelFunc(f)
}

func nilConstructor(f func(color.Color) color.Color) color.Model {
	if f == nil {
		return nil
	}
	return color.ModelFunc(f)
}

func ignoreConstructor(f func(color.Color) color.Color) color.Model {
	if f != nil {
		f = func(c color.Color) color.Color { return c }
	}
	return color.ModelFunc(f)
}
