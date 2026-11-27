# Independent review of `math/rand/v2.PCG.MarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `04fadff2e7624b67048db30463f7263fed127b659cf95a020fc6927cd582e3f1`

The result contains the header and both big-endian words in order. The native serialization leaves both receiver words intact. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
