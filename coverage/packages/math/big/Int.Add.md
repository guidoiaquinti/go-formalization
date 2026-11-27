# Independent review of `math/big.Int.Add`

**Verdict: FAITHFUL**

Evidence SHA-256: `7cae77f17148a7e122e0065d17da2584665f1de4b2c30ab13fb2abb6e46a4b62`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. Pinned Go adds magnitudes for equal signs and subtracts with sign selection for different signs. The Lean result is exact addition and its theorem proves cancellation. The controls exercise mixed signs and aliased destinations.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
