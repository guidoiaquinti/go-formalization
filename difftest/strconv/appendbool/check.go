package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	dst  []byte
	flag bool
}

func decimalBytes(b []byte) string {
	fields := make([]string, len(b))
	for i, v := range b {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, " ")
}

func main() {
	var edges []input
	for _, flag := range []bool{false, true} {
		edges = append(edges, input{nil, flag}, input{make([]byte, 0), flag})
		for spare := range 7 {
			dst := make([]byte, 3, 3+spare)
			copy(dst, []byte{0, 10, 255})
			edges = append(edges, input{dst, flag})
		}
	}
	check.CompareSamples("godiff-appendbool", edges, func(r *rand.Rand, i int) input {
		n := r.Intn(65)
		dst := make([]byte, n, n+r.Intn(9))
		r.Read(dst)
		return input{dst, i%2 == 0}
	}, func(c input) string {
		flag := "0"
		if c.flag {
			flag = "1"
		}
		if len(c.dst) == 0 {
			return flag
		}
		return flag + " " + decimalBytes(c.dst)
	}, func(c input) string {
		return decimalBytes(strconv.AppendBool(c.dst, c.flag))
	}, func(c input) string {
		tail := "false"
		if c.flag {
			tail = "true"
		}
		return decimalBytes(append(append([]byte(nil), c.dst...), tail...))
	})
}
