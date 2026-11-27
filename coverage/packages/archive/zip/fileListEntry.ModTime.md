# Independent review of `archive/zip.fileListEntry.ModTime`

**Verdict: FAITHFUL**

Evidence SHA-256: `41e7b9699010b516877dfcae328fbd34606c38b3736345110c2dcf5343aad872`

I checked evidence commit `b3bf2c6` against Go 1.27.0 `src/archive/zip/reader.go`. An absent file yields Go's zero time. Otherwise the method reads the backing Modified instant and converts the location to UTC. Native replay checks the UTC location, while Lean tracks seconds and nanoseconds.

Pinned same-package Go and Lean agreed on 2,888 fixed inputs and 20,000 sampled inputs per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/archive/zip/file-list-entry`. Five independent rows passed with `go run ./difftest/archive/zip/file-list-entry/review` and rejected fourteen incorrect outcomes across the group. The scoped kernel audit passed. Scope covers byte-level names and nonnil entry metadata. ZIP parsing and monotonic clock state are outside scope. [Independent control](../../../../difftest/archive/zip/file-list-entry/review/main.go)
