package zip

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanHeaderInfo(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_HEADER_INFO"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	rows := 0
	for reader.Scan() {
		fixture, want, found := strings.Cut(reader.Text(), "\t")
		if !found {
			t.Fatal("missing Lean results")
		}
		var creator uint16
		var attrs, legacy uint32
		var wide uint64
		var present bool
		var name string
		if _, err := fmt.Sscanf(fixture, "%d %d %d %d %t %q", &creator, &attrs, &legacy, &wide, &present, &name); err != nil {
			t.Fatal(err)
		}
		var h *FileHeader
		if present {
			h = &FileHeader{Name: "before/creation", UncompressedSize64: 37, CreatorVersion: 0xff00}
		}
		original := h.FileInfo()
		view, typed := original.(headerFileInfo)
		if !typed {
			t.Fatalf("row %d: FileInfo returned %T", rows, original)
		}
		next, err := view.Info()
		got := fmt.Sprintf("%t %t %t", view.fh == h, view.Sys() == h, err == nil && next == original)
		reference := "true true true"
		if present {
			*h = FileHeader{Name: name, CreatorVersion: creator, ExternalAttrs: attrs,
				UncompressedSize: legacy, UncompressedSize64: wide}
			mode := view.Mode()
			got += fmt.Sprintf(" %s %d %d %d %t", encodeName(view.Name()), view.Size(), mode,
				view.Type(), view.IsDir())
			expectedMode := referenceMode(creator, attrs, name)
			kind := expectedMode & (fs.ModeDir | fs.ModeSymlink | fs.ModeDevice | fs.ModeNamedPipe |
				fs.ModeSocket | fs.ModeCharDevice | fs.ModeIrregular)
			size := int64(wide)
			if wide == 0 {
				size = int64(legacy)
			}
			reference += fmt.Sprintf(" %s %d %d %d %t", encodeName(baseName(name)), size,
				expectedMode, kind, expectedMode&fs.ModeDir != 0)
		}
		if got != want || got != reference {
			t.Fatalf("row %d: Go %q, Lean %q, oracle %q", rows, got, want, reference)
		}
		rows++
	}
	if err := reader.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no rows")
	}
	t.Logf("%d pinned metadata rows", rows)
}

func encodeName(raw string) string {
	cells := make([]string, len(raw))
	for i, b := range []byte(raw) {
		cells[i] = strconv.Itoa(int(b))
	}
	return strings.Join(cells, ",")
}

func baseName(raw string) string {
	if raw == "" {
		return "."
	}
	parts := strings.Split(raw, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "/"
}

func referenceMode(creator uint16, attrs uint32, name string) fs.FileMode {
	var mode fs.FileMode
	switch creator / 256 {
	case 3, 19:
		encoded := attrs / 65536
		types := [16]fs.FileMode{1: fs.ModeNamedPipe, 2: fs.ModeDevice | fs.ModeCharDevice,
			4: fs.ModeDir, 6: fs.ModeDevice, 10: fs.ModeSymlink, 12: fs.ModeSocket}
		mode = types[encoded/4096] | fs.FileMode(encoded%512)
		for i, flag := range []fs.FileMode{fs.ModeSticky, fs.ModeSetgid, fs.ModeSetuid} {
			if encoded/(512<<i)%2 != 0 {
				mode |= flag
			}
		}
	case 0, 11, 14:
		dir, readOnly := (attrs/16)%2 == 1, attrs%2 != 0
		switch {
		case dir && readOnly:
			mode = fs.ModeDir | 0o555
		case dir:
			mode = fs.ModeDir | 0o777
		case readOnly:
			mode = 0o444
		default:
			mode = 0o666
		}
	}
	if strings.HasSuffix(name, "/") {
		mode |= fs.ModeDir
	}
	return mode
}
