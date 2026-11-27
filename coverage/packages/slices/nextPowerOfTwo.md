# Independent review of `slices.nextPowerOfTwo`

**Verdict: FAITHFUL**

Evidence SHA-256: `dcd722b8228c52a56af55d63476afa36ceeb1383df73e3edd4544df9e36c8d7f`

I checked evidence commit `f15e6f3` against pinned Go 1.27.0. Go converts the signed input to native `uint` before computing its significant-bit count and uses a native-width unsigned result. Lean applies the same fixed-width reinterpretation and shift. A shift equal to the result width produces zero. Zero has bit length zero, and an exact positive power of two doubles if it fits.

I ran `go run ./difftest/slices/sort-primitives` on darwin/arm64 and replayed 35,526 boundaries and 20,000 per seed 17, 2718, and 65537 against pinned private Go and compiled Lean. Eight literal fixtures passed both through `go run ./difftest/slices/sort-primitives/review`. An inclusive round-up mutation gives 32 on input 32, where native Go and Lean both return 64. The control also checks zero, minus one and native MaxInt. This replay exercised native 64-bit Go. The Lean theorem includes a 32-bit boundary, while native 32-bit execution is unverified here. [Independent control](../../../difftest/slices/sort-primitives/review/main.go)
