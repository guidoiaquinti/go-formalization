# Independent review of `slices.isNaN`

**Verdict: FAITHFUL**

Evidence SHA-256: `6f58e8011f88e5ca69b2c64e06440be31832ea1c78d532b8d0a0064de438992c`

I checked evidence commit `f15e6f3` against pinned Go 1.27.0. The generic helper uses self-inequality. The Lean exponent and payload rule gives that result for binary32 and binary64, including signaling NaNs and both signs. A zero payload with all exponent bits set describes infinity and yields false. Ordered integer and string values are reflexive and receive the false branch in the model.

I ran `go run ./difftest/slices/sort-primitives` on darwin/arm64 and replayed 35,526 boundaries and 20,000 per seed 17, 2718, and 65537 against pinned private Go and compiled Lean. Eight literal fixtures passed both through `go run ./difftest/slices/sort-primitives/review`. Treating an all-ones exponent as sufficient misclassifies positive binary64 infinity. The fixed control also includes negative binary32 infinity, positive binary32 signaling NaN, and negative binary64 signaling NaN. The native overlay instantiates the helper at named float32, float64, int64, and string types. [Independent control](../../../difftest/slices/sort-primitives/review/main.go)
