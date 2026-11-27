# Independent review of `math/rand/v2.PCG.Seed`

**Verdict: FAITHFUL**

Evidence SHA-256: `8f0c171e2b13d1a44173d3dd40faa34758ff90dc47d44c8977159e03a37926b0`

Reseeding replaces both original words. The control uses asymmetric replacement seeds, with native receiver state inspected directly. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
