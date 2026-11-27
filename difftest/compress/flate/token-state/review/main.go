package main

import (
	"log"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	err := check.TestInPackage("compress/flate", "TestReviewGoLeanTokenState",
		"difftest/compress/flate/token-state/review/testdata/control_test.go")
	if err != nil {
		log.Fatal(err)
	}
	input := strings.NewReader("0 9 7 5\n65535 -4 65535 5 5 7\n1 12 3 0\n")
	wants := []string{
		"0 9 7 7 7 7 65536 65536 131071 | 1 9 8 7 7 8 5 5 131071 | 1 9 7 8 7 7 256 256 131071 | 0 9 7 7 7 7 65536 65536 131071 | 0 65536 131071 73 73 29",
		"0 0 0 0 0 0 0 0 65536 131071 131071 131071 65536 | 0 -4 0 65535 65535 0 0 65535 65536 5 5 5 65536 | 0 -4 65535 0 65535 65535 65535 65535 65536 256 256 256 65536 | 1 -4 0 65535 65535 0 0 0 7 5 5 5 7 | 65535 65536 131071 73 73 29",
		"0 0 0 0 0 0 65536 65537 131071 | 2 12 4 3 3 4 65536 0 131071 | 2 12 3 4 3 3 65536 256 131071 | 1 12 3 3 3 3 65536 65537 131071 | 1 65536 131071 73 73 29",
	}
	rows := check.LeanRows("godiff-flate-token-state", len(wants), input)
	for i, row := range rows {
		if row != strings.ReplaceAll(wants[i], " | ", " ") {
			log.Fatalf("Lean row %d: %q", i, row)
		}
	}
	for _, m := range []struct {
		row, group, field int
		value             string
	}{
		{0, 0, 1, "0"},      // clear filled count when n is already zero
		{2, 0, 8, "0"},      // clear token backing on nonzero reset
		{1, 1, 0, "65535"},  // saturate the token count
		{1, 1, 2, "65535"},  // saturate the literal histogram
		{1, 2, 3, "65535"},  // saturate the EOB histogram
		{1, 2, 9, "255"},    // choose a literal instead of EOB
		{1, 3, 8, "65536"},  // stop emitting at the old backing boundary
		{0, 4, 1, "0"},      // cap returned view to visible count
		{0, 4, 3, "131071"}, // slice write detached from receiver
		{0, 4, 5, "73"},     // receiver write hidden from slice
	} {
		groups := strings.Split(wants[m.row], " | ")
		start := 0
		for _, group := range groups[:m.group] {
			start += len(strings.Fields(group))
		}
		if strings.Fields(rows[m.row])[start+m.field] == m.value {
			log.Fatalf("row %d group %d accepted wrong field %d", m.row, m.group, m.field)
		}
	}
	log.Print("independent native state and three Lean rows passed with ten wrong observations rejected")
}
