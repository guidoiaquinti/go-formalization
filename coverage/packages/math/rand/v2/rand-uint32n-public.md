# Independent review of `math/rand/v2.Rand.Uint32N`

**Verdict: FAITHFUL**

Evidence SHA-256: `4c5f36a9ab49c60969913e7fc6b29bdc2a31cdc5039e2ac203ab3cf5a98e2267`

Zero panics before a pull. A non-power-of-two bound checks the source’s complete word and retries below the threshold. Pinned Go, Lean and an independent big-integer oracle agreed on 61,536 author cases per native instantiation at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/Bounded.lean) · [Comparison](../../../../difftest/math/rand/v2/bounded/check.go)

Eight independent native and Lean fixtures passed. The shared control rejects skipping rejection, rejecting at equality, undercounting the exhausting draw, and treating checked zero as a valid one-draw bound. The review used native 64-bit Go. [Reviewer control](../../../../difftest/math/rand/v2/bounded/review/main.go)
