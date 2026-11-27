# Independent review of `math/rand/v2.Int64N`

**Verdict: FAITHFUL**

Evidence SHA-256: `b72a1e0e0c89f160af680bfecc6116e5a515843bb8c2eea42fe81c3c5243279b`

Rejects signed zero and negative bounds before drawing. The exact package-level Go entry ran against a temporary serial test source. All 61,536 author cases and eight independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalBounded.lean) · [Comparison](../../../../difftest/math/rand/v2/global-bounded/check.go)

The fixed controls reject skipping source rejection, rejecting a word exactly at threshold, failing to count the exhausting draw, and allowing invalid zero to reach the source. The test overlay restores the original source on exit. This review covers the delegated behavior on native 64-bit Go with deterministic callbacks. Runtime generator output is outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-bounded/review/main.go)
