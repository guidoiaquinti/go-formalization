# Independent review of `archive/zip.fileListEntry.IsDir`

**Verdict: FAITHFUL**

Evidence SHA-256: `2f7a2537f20f4ad02d2d9e87fc45b5a986ff47416d4aa90f2d80ae4c6ed677af`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. The method always returns true. The fixed fixture with a false internal directory flag confirms that it is ignored.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
