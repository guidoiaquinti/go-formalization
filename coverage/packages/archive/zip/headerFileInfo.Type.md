# Independent review of `archive/zip.headerFileInfo.Type`

**Verdict: FAITHFUL**

Evidence SHA-256: `e71735261ba73cfbe6bfd252222b0c037ace1394068b9407032b16510c970f3c`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. FileMode.Type masks the live mode bits. The independent Unix fixture retains both symlink and directory type bits while discarding permissions.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,520 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
