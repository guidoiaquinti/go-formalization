package filepath_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestGoLeanWalkSymlinksBoundary(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		t.Skip("native Unix fixture")
	}
	t.Chdir(t.TempDir())
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.MkdirAll("real/child", 0700))
	must(os.WriteFile("real/file", nil, 0600))
	must(os.WriteFile("file", nil, 0600))
	must(os.Symlink("real/child", "alias"))
	must(os.Symlink("real/file", "filelink"))
	for i := 255; i >= 0; i-- {
		target := "file"
		if i < 255 {
			target = fmt.Sprintf("link%03d", i+1)
		}
		must(os.Symlink(target, fmt.Sprintf("link%03d", i)))
	}
	for _, c := range []struct{ path, want string }{{"alias/../file", "real/file"}, {"link001", "file"}} {
		got, err := filepath.GoLeanNativeWalkSymlinks(c.path)
		if got != c.want || err != nil {
			t.Fatalf("%q got %q %v want %q", c.path, got, err, c.want)
		}
	}
	if got, err := filepath.GoLeanNativeWalkSymlinks("link000"); got != "" || err == nil || err.Error() != "EvalSymlinks: too many links" {
		t.Fatalf("256 links got %q %v", got, err)
	}
	if got, err := filepath.GoLeanNativeWalkSymlinks("filelink/../file"); got != "" || !errors.Is(err, syscall.ENOTDIR) {
		t.Fatalf("file before parent got %q %v", got, err)
	}
}
