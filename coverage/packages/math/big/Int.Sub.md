# Independent review of `math/big.Int.Sub`

**Verdict: FAITHFUL**

Evidence SHA-256: `9a8ebc7f7e46f8d562fc299f0657f230c28e3a2b0f9c50dd7aa883fa373ee81c`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. Pinned Go adds magnitudes for different signs and subtracts with the correct orientation for equal signs. The Lean result is mathematical subtraction, with self-subtraction proved. The fixtures distinguish reversed operands.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
