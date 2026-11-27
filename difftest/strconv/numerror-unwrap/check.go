package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pointerError struct{}

func (*pointerError) Error() string { panic("unexpected Error call") }

type input struct {
	tag     int
	fn, num byte
}

func actual(c input, values []error) (result string) {
	var receiver *strconv.NumError
	if c.tag != 0 {
		receiver = &strconv.NumError{Func: string([]byte{c.fn}), Num: string([]byte{c.num}), Err: values[c.tag]}
	}
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(runtime.Error); ok && e.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				result = "nil receiver panic"
			} else {
				result = fmt.Sprintf("unexpected panic %T", p)
			}
		}
	}()
	got := receiver.Unwrap()
	for tag := 1; tag < len(values); tag++ {
		if got == values[tag] {
			return strconv.Itoa(tag)
		}
	}
	return fmt.Sprintf("unexpected error type %T", got)
}

func main() {
	var typedNil *pointerError
	values := []error{nil, nil, strconv.ErrSyntax, strconv.ErrRange, typedNil, errors.New("same"), errors.New("same")}
	var edges []input
	for tag := range values {
		edges = append(edges, input{tag, 0, 0}, input{tag, 255, 10})
	}
	check.CompareSamples("godiff-numerror-unwrap", edges, func(r *rand.Rand, _ int) input {
		return input{r.Intn(len(values)), byte(r.Uint32()), byte(r.Uint32())}
	}, func(c input) string {
		return fmt.Sprintf("%d %d %d", c.tag, c.fn, c.num)
	}, func(c input) string {
		return actual(c, values)
	}, func(c input) string {
		if c.tag == 0 {
			return "nil receiver panic"
		}
		return strconv.Itoa(c.tag)
	})
}
