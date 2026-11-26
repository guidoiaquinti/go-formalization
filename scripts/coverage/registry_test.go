package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewedEvidence(t *testing.T) {
	t.Chdir(t.TempDir())
	const version = "1.27.0"
	key := symbol{"pkg", "Example"}
	r := record{GoVersion: version, Model: "lean/GoLean/Example.lean", Theorems: []string{"GoLean.Ω"},
		ModelSymbol: "GoLean.example", Claim: "Return the expected value", Limits: "Values only",
		Configurations: "Native", Check: "difftest/example/check.go", Control: "example/control.go", Sources: []string{"example/🧪&.go"}, Review: "review.md"}
	write := func(path, text string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	paths := []string{
		".go-version", "go.mod", "lean/lean-toolchain", "lean/lakefile.toml",
		"lean/lake-manifest.json", "lean/Audit.lean", "Makefile", ".github/workflows/verify.yml",
		r.Model, "lean/GoLean/Helper.lean", "lean/GoDiff/Example.lean",
		r.Check, "difftest/internal/check/helper.go", r.Sources[0], r.Control, "scripts/coverage/main.go",
	}
	for _, p := range paths {
		write(p, p)
	}
	sum, err := fingerprint(key, r)
	if err != nil {
		t.Fatal(err)
	}
	// Golden digest covers canonical encoding of non-ASCII and HTML characters.
	const wantSum = "6b81d90d30688f30669f20221ce0af0896035c2ddbaf77efb134511a1c305ca7"
	if sum != wantSum {
		t.Fatalf("fingerprint %s, want %s", sum, wantSum)
	}
	const verdict = "# Independent review of `pkg.Example`\n\n**Verdict: FAITHFUL**\n"
	write(r.Review, verdict+"\nEvidence SHA-256: `"+sum+"`\n")
	expect := func(r record, want string) {
		t.Helper()
		state, err := reviewState(version, key, r)
		if err != nil || state != want {
			t.Fatalf("state %q, error %v; want %s", state, err, want)
		}
	}
	expect(r, "FAITHFUL")
	for _, p := range paths {
		write(p, p+"\n")
		expect(r, "AWAITING_REVIEW")
		write(p, p)
	}
	changed := r
	changed.Theorems = []string{"GoLean.other_spec"}
	expect(changed, "AWAITING_REVIEW")
	changed.Theorems = nil
	expect(changed, "AWAITING_REVIEW")
	changed = r
	changed.Claim = "A changed claim"
	expect(changed, "AWAITING_REVIEW")
	changed = r
	changed.SourceFallback = true
	expect(changed, "AWAITING_REVIEW")
	changed = r
	changed.Check = r.Sources[0]
	expect(changed, "AWAITING_REVIEW")
	expect(r, "FAITHFUL")
	write(r.Review, verdict)
	expect(r, "AWAITING_REVIEW")
}

func TestInvalidRecords(t *testing.T) {
	for _, data := range []string{
		`{"model":"first","model":"second"}`, `{"unknown":true}`,
		`{"Model":"wrong case"}`, `{"theorems":null}`, `{"review":""}`, `{} {}`,
	} {
		if _, err := readRecord([]byte(data)); err == nil {
			t.Errorf("accepted %s", data)
		}
	}
}

func TestTargetOrderAndOwnership(t *testing.T) {
	fields, err := object([]byte(`{"z":{},"a":{}}`))
	if err != nil || len(fields) != 2 || fields[0].name != "z" {
		t.Fatalf("object order %v, error %v", fields, err)
	}
	targets := []target{{Name: "cmd/compile", Packages: []string{"cmd/compile/..."}}, {Name: "sort"}}
	for _, p := range []string{"cmd/compile", "cmd/compile/internal/test", "sort"} {
		got, err := owner(p, targets)
		if err != nil || !strings.HasPrefix(p, got) {
			t.Errorf("%s owner %s, error %v", p, got, err)
		}
	}
	for _, p := range []string{"cmd/compiler", "sort/internal"} {
		if _, err := owner(p, targets); err == nil {
			t.Errorf("unexpected owner for %s", p)
		}
	}
}
