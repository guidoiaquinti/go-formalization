# Independent review of `index/suffixarray.ints.slice`

**Verdict: FAITHFUL**

Evidence SHA-256: `19a7c72013b4703b605ad84dc5ca5f8a23fe09810f044975966c17ed599700d5`

I checked evidence commit `f00c370` against the Go 1.27.0 declarations. The selected two-index slice aliases its original backing and may extend into spare storage. The independent window length three over a visible length one detects clamping to the visible length.

I ran `go run ./difftest/index/suffixarray/accessors` and reproduced 56,448 boundary inputs and 20,000 cases each from seeds 17, 2718, and 65537 on pinned native darwin/arm64 and compiled Lean. Six independent literal fixtures passed both with `go run ./difftest/index/suffixarray/accessors/review`. The checks distinguish signed conversion, selected width, exact view capacity and nilness, index and nil-receiver panic classes, and backing contents. The scope assumes legal slice headers and sequential access. [Independent control](../../../../difftest/index/suffixarray/accessors/review/main.go)
