# Independent review of `index/suffixarray.ints.get`

**Verdict: FAITHFUL**

Evidence SHA-256: `a06603b3d4504b9e802401992db76b6940f0812e4c0421812852274bc52675fc`

I checked evidence commit `f00c370` against the Go 1.27.0 declarations. The int32 branch is selected by nonnilness and sign-extends the stored word. A nonnil empty int32 slice still wins over populated int64 storage and gives a bounds panic. The negative fixture detects zero extension.

I ran `go run ./difftest/index/suffixarray/accessors` and reproduced 56,448 boundary inputs and 20,000 cases each from seeds 17, 2718, and 65537 on pinned native darwin/arm64 and compiled Lean. Six independent literal fixtures passed both with `go run ./difftest/index/suffixarray/accessors/review`. The checks distinguish signed conversion, selected width, exact view capacity and nilness, index and nil-receiver panic classes, and backing contents. The scope assumes legal slice headers and sequential access. [Independent control](../../../../difftest/index/suffixarray/accessors/review/main.go)
