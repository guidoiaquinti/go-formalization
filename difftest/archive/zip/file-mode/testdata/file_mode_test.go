package zip

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
	"testing"
)

func TestGoLeanFileMode(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_MODE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for row := 0; ; row++ {
		var version, flags, newVersion uint16
		var attrs, mode, headerMode, newAttrs, dos, unix, toUnix uint32
		var slash, descriptor uint8
		n, err := fmt.Fscan(input, &version, &attrs, &mode, &flags, &slash,
			&headerMode, &newVersion, &newAttrs, &descriptor, &dos, &unix, &toUnix)
		if n == 0 && err == io.EOF && row > 0 {
			t.Logf("%d rows agree with all six Go entry points", row)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		names := [...]string{"", "x", "dir/child", "\\", "\xff\x00", "/", "\xff\x00/", "//"}
		name := names[row%5]
		if slash != 0 {
			name = names[5+row%3]
		}
		h := FileHeader{
			Name: name, Comment: "kept", NonUTF8: true, ReaderVersion: 42, CreatorVersion: version,
			Flags: flags, Method: 8, ModifiedTime: 1, ModifiedDate: 2, CRC32: 3,
			CompressedSize: 4, UncompressedSize: 5, CompressedSize64: 6,
			UncompressedSize64: 7, Extra: []byte{8, 9}, ExternalAttrs: attrs,
		}
		want := h
		want.Extra = []byte{8, 9}
		if got := h.Mode(); got != fs.FileMode(headerMode) {
			t.Fatalf("row %d Mode: Go %#x, Lean %#x", row, uint32(got), headerMode)
		}
		if got := h.hasDataDescriptor(); got != (descriptor != 0) {
			t.Fatalf("row %d hasDataDescriptor: Go %t, Lean %d", row, got, descriptor)
		}
		if got := msdosModeToFileMode(attrs); got != fs.FileMode(dos) {
			t.Fatalf("row %d msdosModeToFileMode: Go %#x, Lean %#x", row, uint32(got), dos)
		}
		if got := unixModeToFileMode(attrs); got != fs.FileMode(unix) {
			t.Fatalf("row %d unixModeToFileMode: Go %#x, Lean %#x", row, uint32(got), unix)
		}
		if got := fileModeToUnixMode(fs.FileMode(mode)); got != toUnix {
			t.Fatalf("row %d fileModeToUnixMode: Go %#x, Lean %#x", row, got, toUnix)
		}
		want.CreatorVersion, want.ExternalAttrs = newVersion, newAttrs
		h.SetMode(fs.FileMode(mode))
		if !reflect.DeepEqual(h, want) {
			t.Fatalf("row %d SetMode: Go %+v, Lean want %+v", row, h, want)
		}
	}
}
