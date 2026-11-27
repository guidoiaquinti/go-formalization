# Independent review of `archive/zip.FileHeader.FileInfo`

**Verdict: FAITHFUL**

Evidence SHA-256: `a666be5030b00ac60905b6da8979fbe8d2fae85a9fafa57cc6399a995967e974`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. The adapter stores the original header pointer in a concrete value. Calling the method on a nil pointer still returns a concrete adapter. The native overlay independently checks the dynamic adapter type and its stored pointer.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,521 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
