# Independent review of `archive/zip.split`

**Verdict: FAITHFUL**

Evidence SHA-256: `6d27d81ae768b75acb1cc6f2e6e2259ecad20c4f64e8789c44943a53b187d5d3`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. Pinned source removes exactly one final slash and then splits at the last remaining slash. The two-slash fixture leaves an empty element and the plain slash fixture defaults to dot directory.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
