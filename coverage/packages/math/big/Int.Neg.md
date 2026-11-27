# Independent review of `math/big.Int.Neg`

**Verdict: FAITHFUL**

Evidence SHA-256: `5007826609e5c059220353e0b4113ba1cd09bfe7d6a060cb27a70edeca8bfc29`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. The source copies the value and flips sign only for nonzero magnitude. The involution theorem and zero fixture check exact inverse value and zero normalization.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
