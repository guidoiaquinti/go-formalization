# Independent review of `math/rand/v2.PCG.Uint64`

**Verdict: FAITHFUL**

Evidence SHA-256: `88bb65e9036a1c724ccfe3b43851a7ce1412c95e4ae63a0eac242f29d1e3b193`

Two successive outputs use each newly advanced state. The independent full integer state and arithmetic form of the DXSM output mixer match both outputs and the exact final receiver state. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
