package filepath_test

import (
	"bufio"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

func TestGoLeanWalkSymlinks(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_WALK_SYMLINK_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	t.Chdir(os.Getenv("GOLEAN_WALK_SYMLINK_TREE"))
	sc := bufio.NewScanner(file)
	count := 0
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 || !strings.HasPrefix(fields[0], "x") {
			t.Fatal("invalid row")
		}
		raw, err := hex.DecodeString(fields[0][1:])
		if err != nil {
			t.Fatal(err)
		}
		got, err := filepath.GoLeanNativeWalkSymlinks(string(raw))
		wire := ""
		switch {
		case err == nil:
			parts := make([]string, len(got))
			for i := range len(got) {
				parts[i] = strconv.Itoa(int(got[i]))
			}
			wire = "ok:" + strings.Join(parts, ",")
		case errors.Is(err, syscall.ENOTDIR):
			wire = "notDirectory"
		case errors.Is(err, os.ErrNotExist):
			wire = "missing"
		case errors.Is(err, syscall.EINVAL):
			wire = "invalid"
		case err.Error() == "EvalSymlinks: too many links":
			wire = "tooManyLinks"
		default:
			t.Fatalf("%q: unexpected native error %v", raw, err)
		}
		if wire != fields[1] {
			t.Fatalf("walkSymlinks(%q): Go %s Lean %s", raw, wire, fields[1])
		}
		count++
	}
	if err := sc.Err(); err != nil || count == 0 {
		t.Fatalf("read %d cases: %v", count, err)
	}
	t.Logf("walkSymlinks matched %d native outputs", count)
}
