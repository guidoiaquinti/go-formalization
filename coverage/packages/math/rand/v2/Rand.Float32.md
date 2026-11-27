# Independent review of `math/rand/v2.Rand.Float32`

**Verdict: FAITHFUL**

Evidence SHA-256: `8c356317c40b9e1df80e836a4234a8465d38ecc25bf3f6a6c340742ce762486c`

One draw interprets bits 32 through 55 as an exact numerator over 2^24. Using the top 24 bits failed 32 basis cases. The native observer compares the exact binary32 bits. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
