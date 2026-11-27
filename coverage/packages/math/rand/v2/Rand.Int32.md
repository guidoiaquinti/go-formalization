# Independent review of `math/rand/v2.Rand.Int32`

**Verdict: FAITHFUL**

Evidence SHA-256: `aa68cdc702f625de4911baa0ecdda3f5632ba17e2f49a0d6e2560e5fad55688f`

One draw uses bits 33 through 63. Substituting the low 31-bit interval failed 62 basis cases. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
