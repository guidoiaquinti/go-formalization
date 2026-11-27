# Independent review of `archive/zip.headerFileInfo.Mode`

**Verdict: FAITHFUL**

Evidence SHA-256: `e47bf5e5ec434fb2730da24c2910159fca046f3af9198ab8fed61d7a8d5f2fa9`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. Pinned Go reads FileHeader.Mode using the current header pointer. The control includes a DOS read-only mode, an unknown creator with slash forcing directory status and a Unix symlink with a trailing slash.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,520 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
