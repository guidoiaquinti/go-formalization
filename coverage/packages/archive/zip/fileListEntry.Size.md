# Independent review of `archive/zip.fileListEntry.Size`

**Verdict: FAITHFUL**

Evidence SHA-256: `80e3c3a63bc9e06dff6501fc482c44bb1fd7d77ad0cab5f4bd53d0cd467b2adf`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. The method always returns zero. The control supplies an attached file with a nonzero size.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
