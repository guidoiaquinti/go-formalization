package check

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// TestInPackage runs test fixtures alongside a pinned Go standard library
// package. The installed toolchain files are only read.
func TestInPackage(pkg, testName string, fixtures ...string) error {
	if len(fixtures) == 0 {
		return fmt.Errorf("no overlay fixtures for %s.%s", pkg, testName)
	}
	tmp, err := os.MkdirTemp("", "golean-overlay-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	env := append(os.Environ(), "GOENV=off", "GOWORK=off", "GOOS=", "GOARCH=",
		"GOFLAGS=", "GOEXPERIMENT=", "CGO_ENABLED=0")
	locate := exec.Command("go", "env", "GOROOT")
	locate.Env = append(env, "GOTOOLCHAIN="+runtime.Version())
	locate.Stderr = os.Stderr
	result, err := locate.Output()
	if err != nil {
		return fmt.Errorf("locate pinned GOROOT: %w", err)
	}
	goroot := strings.TrimSpace(string(result))
	if !filepath.IsAbs(goroot) {
		return fmt.Errorf("expected an absolute GOROOT, got %q", goroot)
	}

	// Use unique virtual names so installed tests are never replaced.
	stem := strings.ReplaceAll(filepath.Base(tmp), "-", "_")
	replacements := make(map[string]string, len(fixtures))
	for i, fixture := range fixtures {
		path, err := filepath.Abs(fixture)
		if err != nil {
			return err
		}
		virtualName := fmt.Sprintf("%s_%d_test.go", stem, i)
		virtual := filepath.Join(goroot, "src", filepath.FromSlash(pkg), virtualName)
		replacements[virtual] = path
	}
	overlay := struct {
		Replace map[string]string
	}{Replace: replacements}
	data, err := json.Marshal(overlay)
	if err != nil {
		return err
	}
	overlayPath := filepath.Join(tmp, "overlay.json")
	if err := os.WriteFile(overlayPath, data, 0600); err != nil {
		return err
	}

	// The source index does not observe new overlay files. An independent
	// module cache also keeps the virtual path outside the test's GOMODCACHE.
	goBinary := filepath.Join(goroot, "bin", "go")
	if runtime.GOOS == "windows" {
		goBinary += ".exe"
	}
	test := exec.Command(goBinary, "test", "-count=1", "-json", "-overlay", overlayPath,
		"-run", "^"+regexp.QuoteMeta(testName)+"$", pkg)
	test.Env = append(env, "GOROOT="+goroot, "GOTOOLCHAIN=local",
		"GOMODCACHE="+filepath.Join(tmp, "modcache"), "GODEBUG=goindex=0")
	test.Stderr = os.Stderr
	output, err := test.Output()
	scanner := bufio.NewScanner(bytes.NewReader(output))
	passed := false
	for scanner.Scan() {
		var event struct {
			Action, Package, Test, Output string
		}
		if parseErr := json.Unmarshal(scanner.Bytes(), &event); parseErr != nil {
			return fmt.Errorf("decode in-package test output for %s: %w", pkg, parseErr)
		}
		fmt.Print(event.Output)
		if event.Action == "pass" && event.Package == pkg && event.Test == testName {
			passed = true
		}
	}
	if err != nil {
		return fmt.Errorf("in-package test for %s: %w", pkg, err)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read in-package test output for %s: %w", pkg, err)
	}
	if !passed {
		return fmt.Errorf("in-package test %s.%s did not pass", pkg, testName)
	}
	return nil
}
