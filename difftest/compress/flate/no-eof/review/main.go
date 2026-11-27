package main

import (
	"log"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := check.TestInPackage("compress/flate", "TestReviewGoLeanNoEOF",
		"difftest/compress/flate/no-eof/review/testdata/control_test.go"); err != nil {
		log.Fatal(err)
	}
	want := []string{"2", "8", "9", "10", "11"}
	rows := check.LeanRows("godiff-flate-no-eof", len(want), strings.NewReader("1\n8\n9\n10\n11\n"))
	if !slices.Equal(rows, want) {
		log.Fatalf("Lean returned %v, want %v", rows, want)
	}
	for _, mutation := range []struct {
		row  int
		want string
	}{
		{0, "1"}, // leave exact EOF alone
		{1, "2"}, // normalize joined EOF
		{2, "2"}, // normalize an error matching via Is, Unwrap or text
		{3, "0"}, // drop the incomparable error
		{4, "0"}, // collapse a typed nil to a nil interface
	} {
		if rows[mutation.row] == mutation.want {
			log.Fatalf("row %d accepted mutated result %q", mutation.row, mutation.want)
		}
	}
	log.Print("five independent Go and Lean cases passed and rejected five wrong observations")
}
