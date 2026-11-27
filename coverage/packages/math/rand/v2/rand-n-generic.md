# Independent review of `math/rand/v2.Rand.N`

**Verdict: FAITHFUL**

Evidence SHA-256: `392f89aceb14264a356b9433f020a5b5d2a94344dc2baf8c388a64aab16bc2e1`

The control exercises a defined signed int8 type and uint16. Conversion to each width precedes checking zero or a signed negative bound. Rejected input makes no source call. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
