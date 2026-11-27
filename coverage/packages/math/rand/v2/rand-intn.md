# Independent review of `math/rand/v2.Rand.IntN`

**Verdict: FAITHFUL**

Evidence SHA-256: `ba83159fa6fc56c6c904650280d8fcb772d191a2ae3031173255e081fe20ee75`

Native signed interpretation rejects zero or negative bounds without drawing. Positive bounds preserve exact bounded values and pull counts. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
