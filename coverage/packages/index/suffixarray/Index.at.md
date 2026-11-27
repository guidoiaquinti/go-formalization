# Independent review of `index/suffixarray.Index.at`

**Verdict: FAITHFUL**

Evidence SHA-256: `a86c0b89458bfb26518707a43cd388e53395490f7d1bd723d05fa6ed78e5aba0`

I checked evidence commit `f00c370` against the Go 1.27.0 declarations. The chosen signed suffix index becomes the start of a one-sided byte slice. An offset at the visible end yields a zero-length view with remaining spare capacity. An offset above visible length panics even when within capacity.

I ran `go run ./difftest/index/suffixarray/accessors` and reproduced 56,448 boundary inputs and 20,000 cases each from seeds 17, 2718, and 65537 on pinned native darwin/arm64 and compiled Lean. Six independent literal fixtures passed both with `go run ./difftest/index/suffixarray/accessors/review`. The checks distinguish signed conversion, selected width, exact view capacity and nilness, index and nil-receiver panic classes, and backing contents. The scope assumes legal slice headers and sequential access. [Independent control](../../../../difftest/index/suffixarray/accessors/review/main.go)
