# Independent review of `math/rand/v2.absInt32`

**Verdict: FAITHFUL**

Evidence SHA-256: `c1a012bfd2a976defadf03616de9a2712ffd2478cea39b351339224268b4ef54`

I reviewed evidence commit `60ea45d` against Go 1.27.0 `src/math/rand/v2/normal.go`. The sign test is the most significant bit of the input. Go's signed negation and conversion to `uint32` preserve the word representation of the wrapped value, which yields 2147483648 for the minimum signed input. Lean proves that applying the original sign recovers all 32 input bits and that the minimum input yields 2147483648. The widened signed oracle checks exact magnitude separately.

I ran `go run ./difftest/math/rand/v2/abs-int32` on darwin/arm64 and reproduced nine boundaries and 20,000 sampled words per seed for 17, 2718 and 65537. The comparison calls the pinned private declaration using the package's test overlay and replays the same words against the compiled Lean model. The model accepts all `int32` inputs.

I ran `go run ./difftest/math/rand/v2/abs-int32/review` on six independent literal cases, including the minimum signed input and equal-magnitude inputs of opposite signs. The checks found two wrong results when returning the input word unchanged, one with saturation at the maximum signed input and two when clearing the sign bit. [Independent control](../../../../../difftest/math/rand/v2/abs-int32/review/main.go)
