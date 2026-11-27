package zip

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanEntryStat(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ZIP_ENTRY_STAT_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	for row := 0; ; row++ {
		var c struct {
			Flags uint8
			Name  string
			Want  string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d rows agree for all three stat methods", row)
				return
			}
			t.Fatal(err)
		}
		e := &fileListEntry{name: c.Name, isDup: c.Flags&1 != 0, isDir: c.Flags&2 != 0}
		var fh *FileHeader
		if c.Flags&4 != 0 {
			e.file = &File{FileHeader: FileHeader{Name: "initial"}}
			fh = &e.file.FileHeader
		}
		opened := &openDir{offset: 8}
		if c.Flags&8 != 0 {
			opened.e = e
		}
		reader := &checksumReader{f: e.file, nread: 27, err: io.ErrUnexpectedEOF}
		observe := func(run func() (fs.FileInfo, error)) (result string) {
			defer func() {
				if p := recover(); p != nil {
					if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "nil pointer") {
						result = "nil"
					} else {
						panic(p)
					}
				}
			}()
			resultInfo, err := run()
			if err != nil {
				again, secondErr := run()
				if resultInfo != nil || again != nil || secondErr == nil || secondErr == err || secondErr.Error() != err.Error() {
					t.Fatalf("row %d: error did not return nil info and a fresh message per call", row)
				}
				text := err.Error()
				cells := make([]string, len(text))
				for i := range len(text) {
					cells[i] = strconv.Itoa(int(text[i]))
				}
				return "duplicate:" + strings.Join(cells, ",")
			}
			switch view := resultInfo.(type) {
			case *fileListEntry:
				if view != e {
					t.Fatalf("row %d: directory identity changed", row)
				}
				return "directory:17"
			case headerFileInfo:
				if view.fh == nil || view.fh != fh || view.Sys() != fh {
					t.Fatalf("row %d: header identity changed", row)
				}
				fh.Name = "updated"
				if view.Name() != "updated" {
					t.Fatalf("row %d: header view did not see mutation", row)
				}
				return "header:29"
			default:
				t.Fatalf("row %d: unexpected result type %T", row, resultInfo)
				return ""
			}
		}
		first := observe(func() (fs.FileInfo, error) { return e.stat() })
		got := fmt.Sprintf("%s|%s|%s", first, observe(opened.Stat), observe(reader.Stat))
		if got != c.Want || opened.offset != 8 || reader.nread != 27 || reader.err != io.ErrUnexpectedEOF {
			t.Fatalf("row %d: Go %q want %q or wrapper state changed", row, got, c.Want)
		}
	}
}
