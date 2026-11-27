# Independent review of `archive/zip.fileListEntry.Sys`

**Verdict: FAITHFUL**

Evidence SHA-256: `4b483b20ff99dd2ec53f986b76d09021bd8c950ffa25969ee8094a3471b60566`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. The method returns a nil interface even with a backing file. The native replay checks interface nilness directly.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
