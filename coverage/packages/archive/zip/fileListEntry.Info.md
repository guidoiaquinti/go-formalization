# Independent review of `archive/zip.fileListEntry.Info`

**Verdict: FAITHFUL**

Evidence SHA-256: `7918133041f93a57997fa94d32384def48e79403f6cdde73049f7fdb29a92ef5`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. The method returns the exact entry pointer as fs.FileInfo and a nil error. The native comparison checks pointer identity under the interface.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
