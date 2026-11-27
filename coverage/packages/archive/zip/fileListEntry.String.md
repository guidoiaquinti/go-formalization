# Independent review of `archive/zip.fileListEntry.String`

**Verdict: FAITHFUL**

Evidence SHA-256: `58dfa88bd8863fa9f457f56cdfb600d7e52649c26dfbf4b6b099655fc73e9190`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. Pinned source delegates to fs.FormatDirEntry. For these directory entries the formatting concatenates d and a space, the raw name and a slash. The controls include an empty name and preserve malformed bytes and a NUL.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
