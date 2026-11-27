# Independent review of `archive/zip.fileEntryCompare`

**Verdict: FAITHFUL**

Evidence SHA-256: `1485892b1dc0efa3d9f583f762a47b6a51795b51fad4b4a982a8da0d1610c299`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. Pinned source compares the split directories first and compares the elements only when directories match. The controls check precedence, unsigned raw bytes and that a trailing slash gives no sort distinction.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
