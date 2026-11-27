# Independent review of `sort.isNaN`

**Verdict: FAITHFUL**

Evidence SHA-256: `dcfbb8eafab93c4ffbbe51fd6dd0970a5117c062ae958fa1cc5305ca280cdf87`

I checked evidence commit `a6a654e` against Go 1.27.0. The private function compares the float with itself. A binary64 value differs from itself exactly when the exponent is all ones and the fraction is nonzero. The fixed negative infinity witness rejects the exponent-only mutation, while positive and negative signaling NaNs and negative zero separate sign and payload effects.

I ran `go run ./difftest/sort/typed-sorts` and replayed 16,384 boundaries and 20,000 cases each at seeds 17, 2718, and 65537 against pinned native Go and compiled Lean on darwin/arm64. Nine additional literal cases passed both implementations with `go run ./difftest/sort/typed-sorts/review` across the six declarations. Exact float output observations preserve payload bits and reflect the pinned permutation. They make no stability guarantee for other inputs or versions. [Independent control](../../../difftest/sort/typed-sorts/review/main.go)
