# Independent review of `index/suffixarray.ints.len`

**Verdict: FAITHFUL**

Evidence SHA-256: `0535b3d42998f64dd548b2fb9ca7dac7b64210e5dc292247e9c5714dc7734baa`

I checked evidence commit `f00c370` against the Go 1.27.0 declarations. The method adds both visible lengths, even for synthetic input with both storage widths populated. The fixture with visible lengths one and two detects selecting a single width or using capacity.

I ran `go run ./difftest/index/suffixarray/accessors` and reproduced 56,448 boundary inputs and 20,000 cases each from seeds 17, 2718, and 65537 on pinned native darwin/arm64 and compiled Lean. Six independent literal fixtures passed both with `go run ./difftest/index/suffixarray/accessors/review`. The checks distinguish signed conversion, selected width, exact view capacity and nilness, index and nil-receiver panic classes, and backing contents. The scope assumes legal slice headers and sequential access. [Independent control](../../../../difftest/index/suffixarray/accessors/review/main.go)
