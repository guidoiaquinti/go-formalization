# Independent review of `math/big.Int.Rem`

**Verdict: FAITHFUL**

Reviewed author commit `2eb40ea` against Go 1.27.0 `src/math/big/int.go`. The nonzero truncated remainder follows the dividend sign, independent of divisor sign.

Pinned Go, Lean and the author bitwise long-division oracle matched 210 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. Seven additional literal cases matched in the standalone control, which detected forcing a nonnegative remainder in 3 cases. All six methods panicked with the exact diagnostic on division by zero, leaving operands and destinations unchanged. The zero-success mutation failed one fixture. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/big/int-division/review` after building `godiff-math-big-int-division`.

Scope covers normalized nonnil integers, the primary destination separate from or aliasing either operand, and distinct remainder destinations for paired calls. Pinned native comparison checks mutation and pointer identity, while Lean models values and the zero-divisor failure outcome.

Evidence SHA-256: `662048bdcb28dc47e1e5b6525eb5c03d7da1c358157f3c6ec392fdf094c909fc`
