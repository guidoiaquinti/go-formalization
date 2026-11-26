package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input bool
		want  string
	}{{false, "false"}, {true, "true"}}
	rows := check.LeanRows("godiff-formatbool", len(cases), nil)
	for i, tc := range cases {
		got := strconv.FormatBool(tc.input)
		if got != tc.want || rows[i] != tc.want {
			log.Fatalf("input %t: Go=%q Lean=%q want=%q", tc.input, got, rows[i], tc.want)
		}
	}
	fmt.Println("strconv.FormatBool: Go and Lean agree on both boolean inputs.")
}
