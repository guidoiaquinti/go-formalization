# Independent review of `math/rand/v2.Float64`

**Verdict: FAITHFUL**

Evidence SHA-256: `5c08f7e3e98a975b5a12df5e3671a8cc7e2cf76ac487c2049da7ea397881fb5d`

Uses the low 53 bits as a numerator scaled by 2^-53. The native observation checks the precise binary64 encoding. Direct native calls with a temporary serial test source verified one draw, exact source panic identity, and a nil source failure before drawing. The overlay restores the global source after the observation. All 61,152 author cases and 68 independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalWord.lean) · [Comparison](../../../../difftest/math/rand/v2/global-word/check.go)

The review checks the delegated result and source effects on native 64-bit Go. Runtime generator output, concurrent mutation of the global source, and native 32-bit verification are outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-word/review/main.go)
