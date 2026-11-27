# Independent review of `archive/zip.headerFileInfo.Info`

**Verdict: FAITHFUL**

Evidence SHA-256: `c2fb8b333aacedd6e0c358cd407c5e1c094a5d2355b8fc87e8e6c1d13161d39a`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. The method returns an interface containing the same adapter value and a nil error, including an adapter holding a nil header. The native overlay checks both interface equality and the error.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,521 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
