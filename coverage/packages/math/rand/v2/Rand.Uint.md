# Independent review of `math/rand/v2.Rand.Uint`

**Verdict: FAITHFUL**

Evidence SHA-256: `908f3a0eac85d866d372e627ee6cec79234625b226159ec3ccd1fd7d9ee4bc9e`

One draw takes the low native uint width. The native replay used a 64-bit uint. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
