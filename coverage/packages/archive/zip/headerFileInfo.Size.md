# Independent review of `archive/zip.headerFileInfo.Size`

**Verdict: FAITHFUL**

Evidence SHA-256: `cca7c61bc2868ee72303a309b15a2514597269ba2a5aabaad53c38261a6f0167`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. Pinned Go uses any nonzero wide field and casts to signed int64, then falls back to the zero-extended legacy field only when the wide field is zero. The fixtures cover both unsigned maxima and verify that changes after adapter creation are visible.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,520 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
