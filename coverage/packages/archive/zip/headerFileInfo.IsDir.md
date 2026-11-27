# Independent review of `archive/zip.headerFileInfo.IsDir`

**Verdict: FAITHFUL**

Evidence SHA-256: `cdd6407e3d614df03860f9ed73cb07142817275a8e88b650400fbb446fd37de1`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. The method tests the directory bit in the live resolved mode. The source's trailing-slash rule applies even to an unknown creator and can coexist with another file type bit.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,520 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
