// Command coverage reports and checks Go verification evidence. Run from the repository root.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
)

func run(args []string) error {
	actions := []string{"report", "update", "check", "inventory", "diff", "diff-build", "lake-update", "fingerprint"}
	if len(args) == 0 || !slices.Contains(actions, args[0]) {
		return fmt.Errorf("usage: coverage <%s> [flags]", strings.Join(actions, "|"))
	}
	action := args[0]
	flags := flag.NewFlagSet(action, flag.ContinueOnError)
	filter := flags.String("target", "", "target to inventory or compare")
	externalBlockers := flags.Bool("external-blockers", false, "list unstarted declarations inaccessible to external Go drivers")
	pkg := flags.String("package", "", "package containing the review")
	name := flags.String("symbol", "", "exact inventory symbol")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}
	if *externalBlockers && action != "inventory" {
		return fmt.Errorf("--external-blockers requires inventory")
	}
	if action == "lake-update" || action == "update" || action == "check" {
		if err := syncLakeConfig(action != "check"); err != nil || action == "lake-update" {
			return err
		}
	}
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		return err
	}
	version := strings.TrimSpace(string(pin))
	targets, err := readTargets()
	if err != nil {
		return err
	}
	if *filter != "" && !slices.ContainsFunc(targets, func(t target) bool { return t.Name == *filter }) {
		return fmt.Errorf("unknown target %q", *filter)
	}
	if action == "fingerprint" {
		if *pkg == "" || *name == "" {
			return fmt.Errorf("fingerprint requires --package and --symbol")
		}
		if _, err := owner(*pkg, targets); err != nil {
			return err
		}
		if !fs.ValidPath(*pkg) {
			return fmt.Errorf("invalid package path %q", *pkg)
		}
		fields, err := readMembers(filepath.Join("coverage/packages", *pkg, "records.json"))
		if err != nil {
			return err
		}
		for _, f := range fields {
			if f.name == *name {
				r, err := readRecord(f.data)
				if err != nil {
					return err
				}
				digest, err := newEvidenceFingerprint().fingerprint(symbol{*pkg, *name}, r)
				if err != nil {
					return err
				}
				fmt.Printf("Evidence SHA-256: `%s`\n", digest)
				return nil
			}
		}
		return fmt.Errorf("unknown record %s.%s", *pkg, *name)
	}
	records, err := loadRecords(version, targets)
	if err != nil {
		return err
	}
	if action == "diff" || action == "diff-build" {
		checks := map[string]map[string]bool{}
		controls := map[string]bool{}
		for key, r := range records {
			if r.Check != "" && (*filter == "" || r.target == *filter) {
				if checks[r.Check] == nil {
					checks[r.Check] = map[string]bool{}
				}
				if r.SourceFallback {
					checks[r.Check][key.Package] = true
				}
				if r.Control != "" {
					controls[r.Control] = true
				}
			}
		}
		if action == "diff-build" {
			if err := buildDifferentials(checks, *filter == ""); err != nil {
				return err
			}
		}
		if len(checks) != 0 {
			if err := runGo(version, "env", "-json", "GOVERSION", "GOOS", "GOARCH", "GO386", "GOAMD64", "GOARM", "GOARM64", "GOMIPS", "GOMIPS64",
				"GOPPC64", "GORISCV64", "GOWASM", "CGO_ENABLED", "GOEXPERIMENT", "GOFLAGS"); err != nil {
				return err
			}
		}
		var jobs [][]string
		for _, path := range sortedKeys(checks) {
			jobs = append(jobs, []string{"run", path})
			if len(checks[path]) > 0 {
				args := []string{"run", "-gcflags=command-line-arguments=-d=ssa/intrinsics/off"}
				for _, pkg := range sortedKeys(checks[path]) {
					args = append(args, "-gcflags="+pkg+"=-d=ssa/intrinsics/off")
				}
				jobs = append(jobs, append(args, path))
			}
		}
		for _, path := range sortedKeys(controls) {
			jobs = append(jobs, []string{"run", path})
		}
		if err := runDifferentials(version, jobs); err != nil {
			return err
		}
		fmt.Printf("Ran %d registered differential checks and %d reviewer controls.\n", len(checks), len(controls))
		return nil
	}
	rows, err := inventory(version, targets)
	if err != nil {
		return err
	}
	seen := map[symbol]bool{}
	totals := map[string]counts{}
	if action == "inventory" {
		header := "Target\tPackage\tSymbol\tState\tSource"
		if *externalBlockers {
			header += "\tExternal call restriction"
		}
		fmt.Println(header)
	}
	for _, row := range rows {
		key := symbol{row.Package, row.Symbol}
		if seen[key] {
			return fmt.Errorf("duplicate function in inventory: %s.%s", row.Package, row.Symbol)
		}
		seen[key] = true
		targetName, err := owner(row.Package, targets)
		if err != nil {
			return err
		}
		c := totals[targetName]
		c.total++
		state := "UNSTARTED"
		if r, ok := records[key]; ok {
			state = r.state
			c.started++
			if state == "FAITHFUL" {
				c.faithful++
			}
		}
		totals[targetName] = c
		if action == "inventory" && (*filter == "" || *filter == targetName) {
			line := fmt.Sprintf("%s\t%s\t%s\t%s\thttps://github.com/golang/go/blob/go%s/src/%s/%s#L%d",
				targetName, row.Package, row.Symbol, state, version, row.Package, row.File, row.Line)
			if *externalBlockers {
				reasons := externalCallRestriction(row)
				if state != "UNSTARTED" || reasons == "" {
					continue
				}
				line += "\t" + reasons
			}
			fmt.Println(line)
		}
	}
	for key := range records {
		if !seen[key] {
			return fmt.Errorf("%s.%s absent from go%s inventory", key.Package, key.Name, version)
		}
	}
	if action == "inventory" {
		return nil
	}
	md, err := report(version, targets, totals, records)
	if err != nil {
		return err
	}
	if action == "report" {
		fmt.Print(md)
		return nil
	}
	for _, file := range []struct{ path, text string }{{"coverage/README.md", md}, {"lean/GoLean.lean", leanRegistry(records)}} {
		if action == "update" {
			if err := os.WriteFile(file.path, []byte(file.text), 0644); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(file.path)
			if err != nil || string(data) != file.text {
				return fmt.Errorf("stale %s. Run make coverage-update", file.path)
			}
		}
	}
	if action == "check" {
		fmt.Println("Coverage report and Lean registry match the pinned source and verification records.")
	}
	return nil
}

// Keep the number of simultaneous Go builds bounded and print each
// command's output together so failures and seeds are easy to attribute.
func runDifferentials(version string, jobs [][]string) error {
	type result struct {
		output bytes.Buffer
		err    error
	}
	results := make([]result, len(jobs))
	queue := make(chan int)
	var workers sync.WaitGroup
	for range min(4, runtime.GOMAXPROCS(0), len(jobs)) {
		workers.Go(func() {
			for i := range queue {
				res := &results[i]
				fmt.Fprintf(&res.output, "go %s\n", strings.Join(jobs[i], " "))
				res.err = runGoWithOutput(version, &res.output, &res.output, jobs[i]...)
			}
		})
	}
	for i := range jobs {
		queue <- i
	}
	close(queue)
	workers.Wait()
	var failures []error
	for i := range results {
		_, err := results[i].output.WriteTo(os.Stdout)
		failures = append(failures, err, results[i].err)
	}
	return errors.Join(failures...)
}

func runGo(version string, args ...string) error {
	return runGoWithOutput(version, os.Stdout, os.Stderr, args...)
}

func runGoWithOutput(version string, stdout, stderr io.Writer, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=go"+version, "GOWORK=off", "GOENV=off",
		"GOOS=", "GOARCH=", "GOFLAGS=", "GOEXPERIMENT=", "CGO_ENABLED=0",
		"GO386=", "GOAMD64=", "GOARM=", "GOARM64=", "GOMIPS=",
		"GOMIPS64=", "GOPPC64=", "GORISCV64=", "GOWASM=")
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
