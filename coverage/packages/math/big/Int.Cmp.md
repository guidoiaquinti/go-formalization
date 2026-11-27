# Independent review of `math/big.Int.Cmp`

**Verdict: FAITHFUL**

Evidence SHA-256: `bb8a00dc997b638a20b31e8b624f58d5875ab76d52550c313a1b27f4f85f18fd`

I checked evidence commit `7f5f439` against Go 1.27.0 `src/math/big/int.go`. Go handles identical pointers and distinguishes signs before comparing same-sign magnitudes. The Lean theorem relates the two nonzero results to signed order. The controls test opposite signs with reversed magnitude order and verify neither operand changes.

The pinned Go and Lean comparison at `go run ./difftest/math/big/int-arithmetic` passed 240 fixed pairs and 20,000 sampled pairs per seed for 17, 2718 and 65537 on darwin/arm64. The native replay checks destination pointer identity, unchanged distinct operands and legal alias profiles. Five independent pairs passed with `go run ./difftest/math/big/int-arithmetic/review` and rejected nine incorrect outcomes across the group. Lean covers exact signed values. The native scope is nonnil normalized operands and destinations. Raw backing, allocation and arithmetic implementation details such as limb algorithms are outside scope. [Independent control](../../../../difftest/math/big/int-arithmetic/review/main.go)
