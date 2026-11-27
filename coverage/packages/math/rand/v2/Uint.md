# Independent review of `math/rand/v2.Uint`

**Verdict: FAITHFUL**

Evidence SHA-256: `b1202769347ca3a676e68ddf1a89d9b1cb5d8dd07fe0c16c81c3ed4be02a4413`

Returns the low native unsigned word. Direct native calls with a temporary serial test source verified one draw, exact source panic identity, and a nil source failure before drawing. The overlay restores the global source after the observation. All 61,152 author cases and 68 independent native and Lean cases passed. Author seeds were 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/GlobalWord.lean) · [Comparison](../../../../difftest/math/rand/v2/global-word/check.go)

The review checks the delegated result and source effects on native 64-bit Go. Runtime generator output, concurrent mutation of the global source, and native 32-bit verification are outside scope. [Reviewer control](../../../../difftest/math/rand/v2/global-word/review/main.go)
