# Independent review of `path/filepath.WalkDir`

**Verdict: FAITHFUL**

Evidence SHA-256: `4a73c12fd764f26cca8d1e665c40bf05590d19d37525036371e3232d23abe38f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `071d886551ced29b97d84c43c42b14ea5b61a5a1849eb9837074207767cd8d21`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `527675fe2649a8bc86ca5a1dd694013d6005f523423fb7f80454abbad9d8faa3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a1228aad763dfcf1cab1f8dff1145a2b5654c70e12de08ebbc0affbb04b86200`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3add418795a302842a63832a96553b1f94aef9f661861460aac564658745a940`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3241ffc3da3e150326dc7908e12fda2be314ae5208c72aa55a2e61280d6b7df`.

Reviewed c6756b3 against Go 1.27.0. The public wrapper uses Lstat and reports an initial lookup error through a single callback with a nil entry and the original path. Successful lookup delegates ordered traversal to the private walker. The wrapper normalizes exact SkipDir or SkipAll sentinels and returns other errors unchanged, including an error that wraps a sentinel.

Six fixed model rows and six native sequences passed independently. They distinguish SkipDir at a file from SkipAll through the remainder of a nested walk, preserve partial-directory read error and continuation as separate callback outcomes, and check a missing root's nil entry. Treating wrapped SkipAll as the sentinel would suppress the independent control's error. Treating SkipDir as a global stop would omit the later root-level sibling.

The complete replay passed all 4,096 callback maps on the fixed tree and 60,000 seeded repeats under 17, 2718 and 65537, plus four native missing-root observations and eight fixed model rows. Filesystem outcomes and a valid depth bound are model inputs. Partial-directory read errors were supplied to the model. Native Windows and Plan 9 remain outside the local test matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/walk-dir/review/main.go
```
