# Independent review of `math/big.Int.Set`

**Verdict: FAITHFUL**

Evidence SHA-256: `47c4092dd2ccb9d8d8a7f81e342657382f97fc74d202520e88a5019b145b478a`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. The source copies normalized magnitude and sign, skipping the copy when receiver and operand are the same pointer. Lean preserves the signed value. Native checks verify that the original retains its value when the destination is distinct.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
