# Independent review of `archive/zip.fileListEntry.Name`

**Verdict: FAITHFUL**

Evidence SHA-256: `4aa789addd3a138adc2232b8717f1f4e3c21fa5c3a10479d7e7d2a90a07af025`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. Pinned source returns the split element verbatim. The controls distinguish empty elements, malformed UTF-8 and an embedded NUL from a cleaned basename.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
