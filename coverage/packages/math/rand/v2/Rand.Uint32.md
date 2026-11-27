# Independent review of `math/rand/v2.Rand.Uint32`

**Verdict: FAITHFUL**

Evidence SHA-256: `4ee88aa09fdc4cd832d31b934973da93d01ada7c51fec90c6c5e5f167ada532f`

One draw returns the upper 32 bits. Substituting the low 32 bits failed 64 basis cases. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
