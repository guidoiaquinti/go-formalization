package main

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"runtime"
)

type marker int

func (*marker) RGBA() (uint32, uint32, uint32, uint32) {
	panic("unexpected RGBA call")
}

type callbackPanic struct{ id int }

type trace struct {
	calls      int
	seen       color.Color
	state      int
	output     color.Color
	panicValue any
}

type converter func(color.Model, color.Color) color.Color

func main() {
	implementations := []struct {
		name          string
		convert       converter
		disagreements int
	}{
		{"Go", color.Model.Convert, 0},
		{"invoke twice", invokeTwice, 4},
		{"replace argument with nil", replaceArgument, 5},
		{"use input as returned color", returnInput, 2},
		{"swallow panic", swallowPanic, 3},
	}
	for _, candidate := range implementations {
		mismatches := 0
		for mode := range 7 {
			if !check(candidate.convert, mode) {
				mismatches++
			}
		}
		if mismatches != candidate.disagreements {
			log.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, mismatches, candidate.disagreements)
		}
		fmt.Printf("%s disagreed on %d of 7 fixtures\n", candidate.name, mismatches)
	}
}

func check(call converter, mode int) bool {
	inputValue, otherValue := marker(3), marker(3)
	var input color.Color = &inputValue
	if mode == 6 {
		input = (*marker)(nil)
	}
	const initialState = 11
	payload := &callbackPanic{30}
	actual := trace{state: initialState}
	cb := func(c color.Color) color.Color {
		actual.calls++
		actual.seen = c
		actual.state += 2
		switch mode {
		case 3:
			return &otherValue
		case 4:
			return nil
		case 5:
			panic(payload)
		default:
			return c
		}
	}
	var model color.Model
	switch mode {
	case 0:
		// Construct a typed nil of the private pointer type returned by ModelFunc.
		model = reflect.Zero(reflect.TypeOf(color.ModelFunc(nil))).Interface().(color.Model)
	case 1:
		model = color.ModelFunc(nil)
	default:
		model = color.ModelFunc(cb)
	}
	func() {
		defer func() { actual.panicValue = recover() }()
		actual.output = call(model, input)
	}()
	switch mode {
	case 0, 1:
		_, panicked := actual.panicValue.(runtime.Error)
		return panicked && actual.calls == 0 && actual.state == initialState && actual.output == nil
	case 2, 3, 4, 5, 6:
		if actual.calls != 1 || actual.seen != input || actual.state != initialState+2 {
			return false
		}
		switch mode {
		case 3:
			return actual.output == &otherValue && actual.panicValue == nil
		case 4:
			return actual.output == nil && actual.panicValue == nil
		case 5:
			return actual.panicValue == payload
		default:
			return actual.output == input && actual.panicValue == nil
		}
	}
	return false
}

func invokeTwice(model color.Model, input color.Color) color.Color {
	model.Convert(input)
	return model.Convert(input)
}

func replaceArgument(model color.Model, _ color.Color) color.Color {
	return model.Convert(nil)
}

func returnInput(model color.Model, input color.Color) color.Color {
	model.Convert(input)
	return input
}

func swallowPanic(model color.Model, input color.Color) color.Color {
	defer func() { recover() }()
	return model.Convert(input)
}
