package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	var parseFailures, formatFailures, appendFailures int
	for _, input := range []string{
		"1", "0", "t", "f", "T", "F", "True", "False", "TRUE", "FALSE", "true", "false",
		"tRuE", "FaLSE", " true", "false\x00", "\xff", "",
	} {
		got, err := strconv.ParseBool(input)
		switch input {
		case "1", "t", "T", "True", "TRUE", "true":
			if !got || err != nil {
				log.Fatalf("ParseBool(%q)=(%t,%v), want true", input, got, err)
			}
		case "0", "f", "F", "False", "FALSE", "false":
			if got || err != nil {
				log.Fatalf("ParseBool(%q)=(%t,%v), want false", input, got, err)
			}
		default:
			e, ok := err.(*strconv.NumError)
			if !ok || e.Func != "ParseBool" || e.Num != input || e.Err != strconv.ErrSyntax || got {
				log.Fatalf("ParseBool(%q)=(%t,%#v), want untouched input and ErrSyntax", input, got, err)
			}
			if strings.EqualFold(input, "true") || strings.EqualFold(input, "false") {
				parseFailures++
			}
		}
	}
	for _, c := range []struct {
		flag bool
		want string
	}{{false, "false"}, {true, "true"}} {
		got := strconv.FormatBool(c.flag)
		if got != c.want {
			log.Fatalf("FormatBool(%t)=%q", c.flag, got)
		}
		if strings.ToUpper(got) != c.want {
			formatFailures++
		}
		for _, dst := range [][]byte{nil, {}, {0x00, 0xff}, make([]byte, 2, 8)} {
			want := string(dst) + c.want
			if got := strconv.AppendBool(dst, c.flag); string(got) != want {
				log.Fatalf("AppendBool(%q,%t)=%q, want %q", dst, c.flag, got, want)
			}
			if c.want != want {
				appendFailures++
			}
		}
	}
	if parseFailures != 2 || formatFailures != 2 || appendFailures != 4 {
		log.Fatalf("mutation failures parse=%d format=%d append=%d", parseFailures, formatFailures, appendFailures)
	}
	fmt.Printf("boolean conversions: 18 parse, 2 format, 8 append fixtures pass; mutant failures=%d/%d/%d\n",
		parseFailures, formatFailures, appendFailures)
}
