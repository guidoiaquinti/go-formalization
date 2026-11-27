package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedOrder interface {
	fmt.Stringer
	fmt.GoStringer
}

func main() {
	fixtures := []struct {
		order namedOrder
		name  string
	}{
		{binary.BigEndian, "BigEndian"},
		{binary.LittleEndian, "LittleEndian"},
		{binary.NativeEndian, "NativeEndian"},
	}
	rows := check.LeanRows("godiff-binary-endian-names", 2*len(fixtures), nil)
	for i, f := range fixtures {
		plain, quoted := f.order.String(), f.order.GoString()
		wantGo := "binary." + f.name
		if plain != f.name || quoted != wantGo || rows[2*i] != f.name || rows[2*i+1] != wantGo ||
			fmt.Sprint(f.order) != f.name || fmt.Sprintf("%#v", f.order) != wantGo {
			log.Fatalf("%T: native %q/%q, Lean %q/%q, want %q/%q", f.order, plain, quoted, rows[2*i], rows[2*i+1], f.name, wantGo)
		}
	}
	fmt.Println("All six encoding/binary display methods match pinned Go, Lean, and formatting.")
}
