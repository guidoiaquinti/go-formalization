# Independent review of `math/rand/v2.Rand.Float64`

**Verdict: FAITHFUL**

Evidence SHA-256: `297dafb58a2767027128bceb8ad5b8609ba742c9fe0341708ddcae355bbec692`

One draw interprets the low 53 bits as an exact numerator over 2^53. Replacing the selected interval with the top 53 bits failed 64 basis cases. The native observer compares the exact binary64 bits. The native controls verified the identical source panic object after exactly one pull and nil dereference without a pull for nil source and receiver. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
