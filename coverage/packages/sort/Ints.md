# Independent review of `sort.Ints`

**Verdict: FAITHFUL**

Evidence SHA-256: `6491f3291308a94bd2cd60034f88bc20e962345b2a5e592cb055fd099ceb01d5`

I checked evidence commit `a6a654e` against Go 1.27.0. The public wrapper calls slices.Sort on signed native ints. The signed extrema fixture detects sorting by unsigned word order.

I ran `go run ./difftest/sort/typed-sorts` and replayed 102 boundaries and 20,000 cases each at seeds 17, 2718, and 65537 against pinned native Go and compiled Lean on darwin/arm64. Nine additional literal cases passed both implementations with `go run ./difftest/sort/typed-sorts/review` across the six declarations. Exact float output observations preserve payload bits and reflect the pinned permutation. They make no stability guarantee for other inputs or versions. [Independent control](../../../difftest/sort/typed-sorts/review/main.go)
