# Independent review of `math/big.Int.Abs`

**Verdict: FAITHFUL**

Evidence SHA-256: `0452ff0d37b000f08225a4aa7167bd36cfd1641d96a165156c5557e7c41b2a6d`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. The source copies the value into the receiver and clears its negative flag. The Lean model gives the mathematical magnitude and the fixed negative input distinguishes sign removal.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
