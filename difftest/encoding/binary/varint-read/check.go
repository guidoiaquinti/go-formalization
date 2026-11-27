package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	data []byte
	tag  int
}

var endErrors = []error{io.EOF, io.ErrUnexpectedEOF, errors.New("reader failure"), fmt.Errorf("wrapped: %w", io.EOF)}

type scriptedReader struct {
	scenario
	calls int
}

func (r *scriptedReader) ReadByte() (byte, error) {
	index := r.calls
	r.calls++
	if index >= len(r.data) {
		// A byte accompanying a non-nil error has no effect on the result.
		return 0xff, endErrors[r.tag]
	}
	return r.data[index], nil
}

func main() {
	var boundaries []scenario
	for tag := range endErrors {
		boundaries = append(boundaries, scenario{nil, tag})
		for b := range 256 {
			for _, n := range []int{0, 1, 8, 9, 10} {
				prefix := slices.Repeat([]byte{0xff}, n)
				boundaries = append(boundaries, scenario{append(prefix, byte(b)), tag})
			}
		}
	}
	check.CompareSamples("godiff-binary-varint-read", boundaries, sample, encode, native, oracle)
}

func sample(r *rand.Rand, i int) scenario {
	tag := r.Intn(len(endErrors))
	data := make([]byte, r.Intn(16))
	r.Read(data)
	if i%3 == 0 {
		for j := range data {
			data[j] |= 128
		}
	}
	if i%3 == 1 {
		data = binary.AppendVarint(data[:0], int64(r.Uint64())>>uint(i%64))
	}
	return scenario{data, tag}
}

func encode(c scenario) string {
	row := strconv.AppendInt(nil, int64(c.tag), 10)
	for _, b := range c.data {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return string(row)
}

func native(c scenario) string {
	first, second := &scriptedReader{scenario: c}, &scriptedReader{scenario: c}
	u, err := binary.ReadUvarint(first)
	s, signedErr := binary.ReadVarint(second)
	return fmt.Sprintf("%d:%s:%d|%d:%s:%d", u, classify(err, c.tag), first.calls, s, classify(signedErr, c.tag), second.calls)
}

func classify(err error, tag int) string {
	switch err {
	case nil:
		return "ok"
	case io.EOF:
		return "eof"
	case io.ErrUnexpectedEOF:
		return "unexpected"
	case endErrors[tag]:
		return "other"
	default:
		if err.Error() == "binary: varint overflows a 64-bit integer" {
			return "overflow"
		}
		return fmt.Sprintf("unexpected error %T %q", err, err)
	}
}

func oracle(c scenario) string {
	clipped := c.data[:min(10, len(c.data))]
	stop := slices.IndexFunc(clipped, func(b byte) bool { return b < 128 })
	errorLabel, calls, groups := "ok", stop+1, clipped
	switch {
	case stop == 9 && clipped[stop] > 1:
		errorLabel, groups = "overflow", clipped[:stop]
	case stop >= 0:
		groups = clipped[:stop+1]
	case len(clipped) == 10:
		errorLabel, calls = "overflow", 10
	default:
		errorLabel, calls = "other", len(clipped)+1
		if c.tag == 0 && len(clipped) == 0 {
			errorLabel = "eof"
		} else if c.tag == 0 || c.tag == 1 {
			errorLabel = "unexpected"
		}
	}
	var value uint64
	for _, b := range slices.Backward(groups) {
		value = 128*value + uint64(b&127)
	}
	signed := int64(value / 2)
	if value%2 == 1 {
		signed = -signed - 1
	}
	return fmt.Sprintf("%d:%s:%d|%d:%s:%d", value, errorLabel, calls, signed, errorLabel, calls)
}
