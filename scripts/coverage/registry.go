// Registry and review evidence for the pinned Go source tree.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
)

type member struct {
	name string
	data json.RawMessage
}

// JSON objects in these files are at most two levels deep. Read each level in
// file order (the dashboard uses target order) and reject duplicate keys.
func object(data []byte) ([]member, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	if t, err := d.Token(); err != nil || t != json.Delim('{') {
		return nil, fmt.Errorf("expected JSON object")
	}
	seen := map[string]bool{}
	var fields []member
	for d.More() {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		name := t.(string)
		if seen[name] {
			return nil, fmt.Errorf("duplicate JSON key: %s", name)
		}
		seen[name] = true
		var value json.RawMessage
		if err := d.Decode(&value); err != nil {
			return nil, err
		}
		fields = append(fields, member{name, value})
	}
	if _, err := d.Token(); err != nil {
		return nil, err
	}
	if _, err := d.Token(); err != io.EOF {
		return nil, fmt.Errorf("unexpected trailing JSON")
	}
	return fields, nil
}

func readMembers(file string) ([]member, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fields, err := object(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return fields, nil
}

type target struct {
	Name     string
	Focus    string
	Packages []string
}

func readTargets() ([]target, error) {
	fields, err := readMembers("coverage/targets.json")
	if err != nil {
		return nil, err
	}
	var targets []target
	for _, field := range fields {
		if _, err := object(field.data); err != nil {
			return nil, err
		}
		var t target
		if err := json.Unmarshal(field.data, &t); err != nil {
			return nil, err
		}
		t.Name = field.name
		targets = append(targets, t)
	}
	return targets, nil
}

func (t target) patterns() []string {
	if t.Packages == nil {
		return []string{t.Name}
	}
	return t.Packages
}

func owner(pkg string, targets []target) (string, error) {
	var matches []string
	for _, t := range targets {
		for _, pattern := range t.patterns() {
			base, tree := strings.CutSuffix(pattern, "/...")
			if pkg == base || (tree && strings.HasPrefix(pkg, base+"/")) {
				matches = append(matches, t.Name)
			}
		}
	}
	if len(matches) != 1 {
		return "", fmt.Errorf("%s must belong to exactly one target, got %v", pkg, matches)
	}
	return matches[0], nil
}

type symbol struct{ Package, Name string }
type record struct {
	GoVersion      string   `json:"go_version"`
	ModelSymbol    string   `json:"model_symbol"`
	Claim          string   `json:"claim"`
	Limits         string   `json:"limits"`
	Configurations string   `json:"configurations"`
	Control        string   `json:"control"`
	SourceFallback bool     `json:"source_fallback"`
	Model          string   `json:"model"`
	Theorems       []string `json:"theorems"`
	Check          string   `json:"check"`
	Review         string   `json:"review"`
	Sources        []string `json:"review_sources"`
}
type evidence struct {
	record
	target, state string
}

func readRecord(data []byte) (record, error) {
	var r record
	fields, err := object(data)
	if err != nil {
		return r, err
	}
	for _, f := range fields {
		switch f.name {
		case "go_version", "model", "model_symbol", "check", "control", "review", "claim", "limits", "configurations":
			if string(f.data) == `""` {
				return r, fmt.Errorf("empty record field %q", f.name)
			}
		case "review_sources", "theorems", "source_fallback":
		default:
			return r, fmt.Errorf("unknown record field %q", f.name)
		}
		if bytes.Equal(f.data, []byte("null")) {
			return r, fmt.Errorf("null record field %q", f.name)
		}
	}
	err = json.Unmarshal(data, &r)
	return r, err
}

func isFile(path string) error {
	if !fs.ValidPath(path) || !filepath.IsLocal(path) {
		return fmt.Errorf("invalid repository path %q", path)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("expected a regular file: %s", path)
	}
	return nil
}

type evidenceFingerprint struct {
	shared func() (string, error)
}

func newEvidenceFingerprint() evidenceFingerprint {
	return evidenceFingerprint{shared: sync.OnceValues(sharedFingerprint)}
}

func sharedFingerprint() (string, error) {
	paths := []string{".go-version", "go.mod", "lean/lean-toolchain",
		"lean/lakefile.toml", "lean/lake-manifest.json", "Makefile", ".github/workflows/verify.yml"}
	// Hash a conservative source closure. Exclude the generated registry,
	// whose contents are already derived from the record and checked in CI.
	for _, root := range []string{"lean", "difftest", "scripts"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			switch filepath.ToSlash(path) {
			case "lean/.lake":
				return filepath.SkipDir
			case "lean/GoLean.lean":
				return nil
			}
			if !entry.IsDir() {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	if _, err := os.Stat("go.sum"); err == nil {
		paths = append(paths, "go.sum")
	} else if !os.IsNotExist(err) {
		return "", err
	}
	files, err := evidenceFileHashes(paths)
	if err != nil {
		return "", err
	}
	return canonicalFingerprint(files)
}

func evidenceFileHashes(paths []string) (map[string]string, error) {
	files := map[string]string{}
	for _, p := range paths {
		if _, ok := files[p]; ok {
			continue
		}
		if err := isFile(p); err != nil {
			return nil, err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		files[p] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	return files, nil
}

func (f evidenceFingerprint) fingerprint(key symbol, r record) (string, error) {
	shared, err := f.shared()
	if err != nil {
		return "", err
	}
	paths := append([]string{r.Model, r.Check}, r.Sources...)
	if r.Control != "" {
		paths = append(paths, r.Control)
	}
	files, err := evidenceFileHashes(paths)
	if err != nil {
		return "", err
	}
	theorems := append([]string{}, r.Theorems...)
	slices.Sort(theorems)
	manifest := map[string]any{"package": key.Package, "symbol": key.Name, "go_version": r.GoVersion,
		"model": r.Model, "model_symbol": r.ModelSymbol, "check": r.Check, "theorems": theorems,
		"control": r.Control, "source_fallback": r.SourceFallback, "claim": r.Claim,
		"limits": r.Limits, "configurations": r.Configurations, "files": files, "shared_files_sha256": shared}
	return canonicalFingerprint(manifest)
}

func canonicalFingerprint(manifest any) (string, error) {
	// Keep the existing canonical JSON format: sorted keys, compact ASCII
	// escapes (including surrogate pairs), and no HTML escaping.
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(manifest); err != nil {
		return "", err
	}
	var canonical strings.Builder
	for _, c := range strings.TrimSuffix(buf.String(), "\n") {
		switch {
		case c < 0x7f:
			canonical.WriteRune(c)
		case c <= 0xffff:
			fmt.Fprintf(&canonical, "\\u%04x", c)
		default:
			c -= 0x10000
			fmt.Fprintf(&canonical, "\\u%04x\\u%04x", 0xd800+c/1024, 0xdc00+c%1024)
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(canonical.String()))), nil
}

var verdictLine = regexp.MustCompile(`(?m)^\*\*Verdict: (FAITHFUL|HOLLOW)\*\*`)
var digestLine = regexp.MustCompile("(?m)^Evidence SHA-256: `([0-9a-f]{64})`$")

func (f evidenceFingerprint) reviewState(version string, key symbol, r record) (string, error) {
	if r.GoVersion != version {
		return "", fmt.Errorf("record must be checked for go%s", version)
	}
	if !strings.HasPrefix(r.Model, "lean/GoLean/") || !strings.HasSuffix(r.Model, ".lean") {
		return "", fmt.Errorf("invalid model path %q", r.Model)
	}
	if len(r.Theorems) > 0 && !strings.HasPrefix(r.ModelSymbol, "GoLean.") {
		return "", fmt.Errorf("theorems require a GoLean model_symbol")
	}
	if r.Control != "" {
		if err := isFile(r.Control); err != nil {
			return "", err
		}
	}
	for _, name := range r.Theorems {
		if !strings.HasPrefix(name, "GoLean.") {
			return "", fmt.Errorf("invalid theorem name %q", name)
		}
	}
	if err := isFile(r.Model); err != nil {
		return "", err
	}
	if r.Check == "" {
		if r.Review != "" {
			return "", fmt.Errorf("reviewed record needs a differential check")
		}
		return "MODELLED", nil
	}
	if err := isFile(r.Check); err != nil {
		return "", err
	}
	if r.Claim == "" || r.Limits == "" || r.Configurations == "" {
		return "", fmt.Errorf("checked record requires claim, limits and configurations")
	}
	if r.Review == "" {
		return "AWAITING_REVIEW", nil
	}
	if err := isFile(r.Review); err != nil {
		return "", err
	}
	data, err := os.ReadFile(r.Review)
	if err != nil {
		return "", err
	}
	review := strings.ReplaceAll(string(data), "\r\n", "\n")
	header, _, _ := strings.Cut(review, "\n")
	if header != fmt.Sprintf("# Independent review of `%s.%s`", key.Package, key.Name) {
		return "", fmt.Errorf("wrong review heading in %s", r.Review)
	}
	verdicts := verdictLine.FindAllStringSubmatch(review, -1)
	if len(verdicts) != 1 {
		return "", fmt.Errorf("%s: expected one FAITHFUL or HOLLOW verdict", r.Review)
	}
	digests := digestLine.FindAllStringSubmatch(review, -1)
	if len(digests) > 1 {
		return "", fmt.Errorf("%s: expected a single evidence digest", r.Review)
	}
	sum, err := f.fingerprint(key, r)
	if err != nil {
		return "", err
	}
	if len(digests) == 0 || digests[0][1] != sum {
		return "AWAITING_REVIEW", nil
	}
	verdict := verdicts[0][1]
	if verdict == "FAITHFUL" && (len(r.Theorems) == 0 || r.Control == "") {
		return "", fmt.Errorf("FAITHFUL record needs a theorem and a reviewer control")
	}
	return verdict, nil
}

func loadRecords(version string, targets []target) (map[symbol]evidence, error) {
	records := map[symbol]evidence{}
	fingerprints := newEvidenceFingerprint()
	const root = "coverage/packages"
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "records.json" {
			return nil
		}
		rel, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		pkg := filepath.ToSlash(rel)
		target, err := owner(pkg, targets)
		if err != nil {
			return err
		}
		fields, err := readMembers(path)
		if err != nil {
			return err
		}
		for _, f := range fields {
			key := symbol{pkg, f.name}
			r, err := readRecord(f.data)
			if err != nil {
				return fmt.Errorf("%s.%s: %w", pkg, f.name, err)
			}
			state, err := fingerprints.reviewState(version, key, r)
			if err != nil {
				return fmt.Errorf("%s.%s: %w", pkg, f.name, err)
			}
			records[key] = evidence{r, target, state}
		}
		return nil
	})
	return records, err
}

func sortedKeys[V any](m map[string]V) []string { return slices.Sorted(maps.Keys(m)) }
