# Independent review of `math/rand/v2.Rand.Int32N`

**Verdict: FAITHFUL**

Evidence SHA-256: `7e3636b797ec55465189d8657fcce0f146f9307112f019adc16795bd7adf252d`

Signed zero and negative bounds introduced by native 32-bit conversion panic before any pull. Positive inputs retain the exact retry and scale behavior. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
