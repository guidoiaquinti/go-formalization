# Independent review of `math/rand/v2.Rand.uint32n`

**Verdict: FAITHFUL**

Evidence SHA-256: `01880dca42887a86f49716b3656aafe210409c949d369ecd6a093002f1dcb3f1`

The private helper uses the low 32-bit mask at power-of-two bounds, including a wrapped mask at zero. Combining the two 32-bit products is equivalent to the full 64-bit numerator. The comparison uses an independent full-product oracle. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
