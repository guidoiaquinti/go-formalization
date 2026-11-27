# Independent review of `archive/zip.fileListEntry.Mode`

**Verdict: FAITHFUL**

Evidence SHA-256: `f1acc8bd8d0b455f5a42a2e38c913151c17d5df16883fecf07921b16d49e234c`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. The method always returns the directory mode with 0555 permission bits. The native fixtures vary both entry flags and attachment.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
