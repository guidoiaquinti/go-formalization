# Independent review of `math/rand/v2.IntN`

**Verdict: FAITHFUL**

Evidence SHA-256: `e385d537d00a70259f23a5b1e87c08fd6efde39aacee3b56c9b461f790a01621`

Uses the native signed width for preflight and makes no draw for zero or negative bounds. The exact package-level Go entry ran against a temporary serial test source. All 61,536 author cases and eight independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalBounded.lean) · [Comparison](../../../../difftest/math/rand/v2/global-bounded/check.go)

The fixed controls reject skipping source rejection, rejecting a word exactly at threshold, failing to count the exhausting draw, and allowing invalid zero to reach the source. The test overlay restores the original source on exit. This review covers the delegated behavior on native 64-bit Go with deterministic callbacks. Runtime generator output is outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-bounded/review/main.go)
