package filepath_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestGoLeanUnixAbsDeepCwd(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("deep-cwd construction verified on Darwin")
	}
	if os.Getenv("GOLEAN_UNIX_ABS_DEEP_CWD") == "1" {
		t.Setenv("PWD", "")
		const component = "golean-deep-step"
		for range 360 {
			if err := os.Mkdir(component, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(component); err != nil {
				t.Fatal(err)
			}
		}
		if got, err := filepath.GoLeanUnixAbs("/x/../bad\xff"); got != "/bad\xff" || err != nil {
			t.Fatalf("absolute path with inaccessible cwd name: %q, %v", got, err)
		}
		_, want := os.Getwd()
		if want == nil {
			t.Fatal("Getwd unexpectedly resolved the deep cwd")
		}
		for _, path := range []string{"", "../bad\xff"} {
			if got, err := filepath.GoLeanUnixAbs(path); got != "" || !reflect.DeepEqual(err, want) {
				t.Fatalf("unixAbs(%q) returned %q, %#v, want empty and %#v", path, got, err, want)
			}
		}
		return
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "-test.run=^TestGoLeanUnixAbsDeepCwd$")
	cmd.Dir = t.TempDir()
	cmd.Env = append(os.Environ(), "GOLEAN_UNIX_ABS_DEEP_CWD=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("deep cwd subprocess: %v\n%s", err, out)
	}
	t.Log("native Getwd error was propagated verbatim for two relative paths and bypassed for an absolute path")
}
