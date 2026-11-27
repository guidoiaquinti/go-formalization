# Independent review of `math/rand/v2.NewPCG`

**Verdict: FAITHFUL**

Evidence SHA-256: `2a3a039a01a81888c895371c228947fdf4930f5e633ecfe772d7fc6fb72d306f`

Construction stores both seed words in order, including high bits. The model and pinned private state read agree. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
