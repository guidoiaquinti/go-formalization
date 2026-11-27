# Independent review of `math/rand/v2.Shuffle`

**Verdict: FAITHFUL**

Evidence SHA-256: `ead1254b3348bd408c80fda270bb652826993f435b5647086e1a4f86f23a48ba`

Invokes descending swaps at the delegated indices. Rejection precedes a swap, and source exhaustion or callback panic stops the trace at the exact point with the corresponding panic payload. The exact package-level Go entry ran against a temporary serial test source. All 60,540 author cases and eight independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalShuffle.lean) · [Comparison](../../../../difftest/math/rand/v2/global-shuffle/check.go)

The fixed controls reject reversing swap order, continuing after callback panic, miscounting the exhausting draw, checking a nil callback before a draw, and returning a nil empty permutation. The test overlay restores the original source on exit. This review covers the delegated behavior on native 64-bit Go with deterministic callbacks. Runtime generator output is outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-shuffle/review/main.go)
