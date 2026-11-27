# Independent review of `math/rand/v2.PCG.AppendBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `4828912116f31730268a8d2b71ffe9e5330e2f8b0bcb4f848bccfb8886f77306`

Appends the header and two complete word writes. A complete chunk written before reallocating remains in the original backing. Fixed comparisons bracket spare capacities 4, 12 and 20. Omitting the writes made before growth fails four. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
