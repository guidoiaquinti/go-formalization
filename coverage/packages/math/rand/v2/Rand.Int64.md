# Independent review of `math/rand/v2.Rand.Int64`

**Verdict: FAITHFUL**

Evidence SHA-256: `cbcfe90ca7ee9a14926e06e87c9ce268cf0e2eed1c33eced3dfe3c74efec9501`

One draw retains bits 0 through 62 and discards the sign bit. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
