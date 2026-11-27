# Independent review of `math/rand/v2.Rand.UintN`

**Verdict: FAITHFUL**

Evidence SHA-256: `636ae3bbfa11dda45cf3d637acca1e24185dc8dcbccd3e5c6cc7a58ce342b7cf`

Native unsigned zero is rejected before drawing. A maximum native bound and small power-of-two bounds keep their exact draw sequence. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
