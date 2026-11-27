# Independent review of `math/rand/v2.Rand.uint64n`

**Verdict: FAITHFUL**

Evidence SHA-256: `0dcc123a01b8f50b845b947d4bdaa38057f51d799d5ae240bf4196a7bcb7a598`

The private helper treats zero as a wrapped all-ones mask. Non-power bounds reject low products below the threshold and accept at equality. On the pinned 64-bit profile, private zero returns the full source word. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
