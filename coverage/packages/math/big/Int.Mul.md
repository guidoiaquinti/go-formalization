# Independent review of `math/big.Int.Mul`

**Verdict: FAITHFUL**

Evidence SHA-256: `2ead1f3992eb2a2880c9cbb440257e3960c4233facc4d703102ffb99952a8304`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. Pinned Go uses the square path for identical operand pointers and otherwise multiplies magnitudes with sign selection. Both paths return the same mathematical product for equal values. Independent fixtures exercise self-squaring and a signed product of magnitude 2^64. The Lean theorem proves multiplication by zero.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
