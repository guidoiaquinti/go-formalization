# Independent review of `math/big.Int.Div`

**Verdict: FAITHFUL**

Reviewed author commit `2eb40ea` against Go 1.27.0 `src/math/big/int.go`. When the truncated remainder is negative, decrements for a positive divisor and increments for a negative divisor. Go saves the divisor sign before an aliased overwrite.

Pinned Go, Lean and the author bitwise long-division oracle matched 210 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. Seven additional literal cases matched in the standalone control, which detected always decrementing the quotient for a negative remainder in 1 cases. All six methods panicked with the exact diagnostic on division by zero, leaving operands and destinations unchanged. The zero-success mutation failed one fixture. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/big/int-division/review` after building `godiff-math-big-int-division`.

Scope covers normalized nonnil integers, the primary destination separate from or aliasing either operand, and distinct remainder destinations for paired calls. Pinned native comparison checks mutation and pointer identity, while Lean models values and the zero-divisor failure outcome.

Evidence SHA-256: `0c58988b741cfe8a42c96da07369e13bbbd51f6125637f5e6314c7d52d30cc04`
