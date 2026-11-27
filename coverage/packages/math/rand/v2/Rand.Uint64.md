# Independent review of `math/rand/v2.Rand.Uint64`

**Verdict: FAITHFUL**

Evidence SHA-256: `a5dd8bae0fd35b2f1844e11a37079e61d35b9977a7b3acd9721a1ae8a2a18819`

One draw returns the source word unchanged, including all 64 positions. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
