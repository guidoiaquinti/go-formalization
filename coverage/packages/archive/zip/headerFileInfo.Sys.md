# Independent review of `archive/zip.headerFileInfo.Sys`

**Verdict: FAITHFUL**

Evidence SHA-256: `397ca206aa1f7caf9231abfc4d3ef52adafe3d06d7b3efba07a1ee316688500e`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. The method returns the stored header as an interface. The native nil-receiver fixture distinguishes a nonnil interface containing a typed nil pointer from a nil interface. The Lean identity model represents the nil pointer abstractly.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,521 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
