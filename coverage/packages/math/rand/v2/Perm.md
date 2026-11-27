# Independent review of `math/rand/v2.Perm`

**Verdict: FAITHFUL**

Evidence SHA-256: `55ef6978c1e85f1f0c8cb87a48dadebab05a7f4a7f42a4fb10339760b858024e`

Returns exact descending-swap permutation, including a nonnil empty result. An exhausting source or a negative length prevents a returned slice. The exact package-level Go entry ran against a temporary serial test source. All 60,540 author cases and eight independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalShuffle.lean) · [Comparison](../../../../difftest/math/rand/v2/global-shuffle/check.go)

The fixed controls reject reversing swap order, continuing after callback panic, miscounting the exhausting draw, checking a nil callback before a draw, and returning a nil empty permutation. The test overlay restores the original source on exit. This review covers the delegated behavior on native 64-bit Go with deterministic callbacks. Runtime generator output is outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-shuffle/review/main.go)
