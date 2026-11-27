# Independent review of `archive/zip.headerFileInfo.Name`

**Verdict: FAITHFUL**

Evidence SHA-256: `dec321881440d0c4ca935b7b5c85af2a72c26d3aa3ab0c842c892d05d30273e8`

I reviewed evidence commit `f8984ef` against Go 1.27.0 `src/archive/zip/struct.go`. Pinned Go delegates to path.Base. The model trims trailing slashes and retains the last slash-delimited raw byte component. The fixtures distinguish malformed bytes, a literal backslash, empty input and all slashes.

The pinned native and Lean comparison at `go run ./difftest/archive/zip/header-info` passed 2,520 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Six independent rows passed with `go run ./difftest/archive/zip/header-info/review` and rejected nine incorrect outcomes across the group. Native checks update header fields after view creation. The Lean identity model abstracts native object layout. Nil header inputs cover FileInfo, Sys and Info. Field accessors are scoped to nonnil headers. [Independent control](../../../../difftest/archive/zip/header-info/review/main.go)
