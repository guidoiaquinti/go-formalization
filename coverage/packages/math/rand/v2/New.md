# Independent review of `math/rand/v2.New`

**Verdict: FAITHFUL**

Evidence SHA-256: `36426fce5d47a6b6acb0d3880ce26a16979ca48c4ca3aec99d9af25c72984087`

The native constructor creates a distinct Rand and keeps the exact supplied source interface without pulling a word. A nil interface is stored unchanged. All 61,536 author cases and 69 independent native and Lean cases passed. The author seeds were 17, 2718 and 65537. The reviewer control includes zero, all ones, and all 64 one-hot words. [Model](../../../../lean/GoLean/Math/RandV2/Word.lean) · [Comparison](../../../../difftest/math/rand/v2/word/check.go)

Native 32-bit execution and arbitrary callback effects are outside the review scope. [Reviewer control](../../../../difftest/math/rand/v2/word/review/main.go)
