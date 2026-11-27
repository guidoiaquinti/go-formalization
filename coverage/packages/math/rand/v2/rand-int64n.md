# Independent review of `math/rand/v2.Rand.Int64N`

**Verdict: FAITHFUL**

Evidence SHA-256: `a6cb68a3976ac0fe414d18cfe55830164f115f294040cda3d1d1dcbfbc30bcf2`

Signed zero and negative bounds panic with the method-specific text before drawing. A positive bound uses the modeled rejection and scaling path. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
