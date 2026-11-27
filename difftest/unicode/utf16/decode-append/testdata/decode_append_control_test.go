package utf16

import (
	"slices"
	"testing"
)

func TestGoLeanDecodeAppendControl(t *testing.T) {
	prefix := make([]rune, 2, 8)
	prefix[0], prefix[1] = -1, 0xd800
	words := []uint16{0xd800, 0xd801, 0xdc02, 0xdc00}
	want := []rune{-1, 0xd800, 0xfffd, 0x10402, 0xfffd}
	got := decode(words, prefix)
	if !slices.Equal(got, want) {
		t.Fatalf("decode = %x, want %x", got, want)
	}
	if slices.Equal(decode(words, nil), want) {
		t.Fatal("discarding the destination prefix unexpectedly matched")
	}
	t.Log("Go matched; discard-prefix mutant rejected")
}
