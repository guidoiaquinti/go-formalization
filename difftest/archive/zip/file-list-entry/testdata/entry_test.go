package zip

import (
	"bufio"
	"cmp"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGoLeanZipFileList(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_FILE_LIST_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	locationByOffset := make(map[int]*time.Location)
	for row := 0; ; row++ {
		var c struct {
			Name, Other                   string
			Present, Directory, Duplicate bool
			Seconds                       int64
			Nanos, Offset                 int
			Want                          string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d rows agree for all eleven pinned Go declarations", row)
				return
			}
			t.Fatal(err)
		}
		e := &fileListEntry{name: c.Name, isDir: c.Directory, isDup: c.Duplicate}
		wantSeconds, wantNanos := int64(-62135596800), 0
		if c.Present {
			location := locationByOffset[c.Offset]
			if location == nil {
				location = time.FixedZone("other", c.Offset)
				locationByOffset[c.Offset] = location
			}
			e.file = &File{FileHeader: FileHeader{Name: "unrelated", Modified: time.Unix(c.Seconds, int64(c.Nanos)).In(location),
				ModifiedTime: 0xffff, ModifiedDate: 0xffff, UncompressedSize64: 9900}}
			wantSeconds, wantNanos = c.Seconds, c.Nanos
		}
		modified := e.ModTime()
		if modified.Location() != time.UTC || (!c.Present && !modified.IsZero()) {
			t.Fatalf("row %d: modified time has wrong location or zero state", row)
		}
		info, failure := e.Info()
		dir, elem, tail := split(c.Name)
		order := fileEntryCompare(c.Name, c.Other)
		got := fmt.Sprintf("%s|%s|%t|%d|%s|%d|%d|%d|%t|%t|%d|%d|%t|%s",
			listEntryWire(dir), listEntryWire(elem), tail, order, listEntryWire(e.Name()), e.Size(), e.Mode(), e.Type(),
			e.IsDir(), e.Sys() == nil, modified.Unix(), modified.Nanosecond(), info == e && failure == nil, listEntryWire(e.String()))
		wantDir, wantElem, wantTail := listEntrySplitOracle(c.Name)
		otherDir, otherElem, _ := listEntrySplitOracle(c.Other)
		wantOrder := cmp.Compare(wantDir, otherDir)
		if wantOrder == 0 {
			wantOrder = cmp.Compare(wantElem, otherElem)
		}
		ref := fmt.Sprintf("%s|%s|%t|%d|%s|0|%d|%d|true|true|%d|%d|true|%s",
			listEntryWire(wantDir), listEntryWire(wantElem), wantTail, wantOrder, listEntryWire(wantElem),
			fs.ModeDir|0o555, fs.ModeDir, wantSeconds, wantNanos, listEntryWire("d "+wantElem+"/"))
		if got != c.Want || got != ref {
			t.Fatalf("row %d name=%q other=%q: Go %q Lean %q oracle %q", row, c.Name, c.Other, got, c.Want, ref)
		}
	}
}

func listEntrySplitOracle(name string) (string, string, bool) {
	parts := strings.Split(name, "/")
	trailing := len(parts) > 1 && parts[len(parts)-1] == ""
	if trailing {
		parts = parts[:len(parts)-1]
	}
	elem := parts[len(parts)-1]
	if len(parts) == 1 {
		return ".", elem, trailing
	}
	return strings.Join(parts[:len(parts)-1], "/"), elem, trailing
}

func listEntryWire(raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(fields, ",")
}
