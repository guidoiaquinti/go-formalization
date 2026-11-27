# Independent review of `sort.Float64Slice.Sort`

**Verdict: FAITHFUL**

Evidence SHA-256: `f01a10407bc01b78fd06448c9234636a3d0a1f978f9d7fae190f023dbbb778cf`

I checked evidence commit `a6a654e` against Go 1.27.0. The generic slices entry and Interface method use the same order with NaNs first. The model ranks all NaNs alike and both zero encodings alike, then replays swaps on original words. The independent short fixture checks exact NaN payload placement and both zero signs, with signed subnormals and infinities. It detects plain floating comparison and a rank splitting the two zero encodings.

I ran `go run ./difftest/sort/typed-sorts` and replayed 102 boundaries and 20,000 cases each at seeds 17, 2718, and 65537 against pinned native Go and compiled Lean on darwin/arm64. Nine additional literal cases passed both implementations with `go run ./difftest/sort/typed-sorts/review` across the six declarations. Exact float output observations preserve payload bits and reflect the pinned permutation. They make no stability guarantee for other inputs or versions. [Independent control](../../../difftest/sort/typed-sorts/review/main.go)
