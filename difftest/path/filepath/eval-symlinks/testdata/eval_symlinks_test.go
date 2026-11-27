package filepath_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoLeanPublicEvalSymlinks(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_PUBLIC_EVAL_SYMLINK_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	t.Chdir(os.Getenv("GOLEAN_PUBLIC_EVAL_SYMLINK_TREE"))
	sc := bufio.NewScanner(file)
	count := 0
	for sc.Scan() {
		raw, err := hex.DecodeString(sc.Text())
		if err != nil {
			t.Fatal(err)
		}
		path := string(raw)
		got, gotErr := filepath.EvalSymlinks(path)
		want, wantErr := filepath.GoLeanWalkSymlinks(path)
		if got != want || !reflect.DeepEqual(gotErr, wantErr) {
			t.Fatalf("evalSymlinks(%q)=(%q,%#v), walk=(%q,%#v)", path, got, gotErr, want, wantErr)
		}
		count++
	}
	if err := sc.Err(); err != nil || count == 0 {
		t.Fatalf("read %d cases: %v", count, err)
	}
	t.Logf("public EvalSymlinks matched native walk delegation on %d paths", count)
}
