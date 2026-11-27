# Independent review of `math/rand/v2.PCG.next`

**Verdict: FAITHFUL**

Evidence SHA-256: `fa1fb887c327f48da0cb779abb0135834afc4a4346284c9e3fdb50f1cb739c68`

The native return matches the updated receiver. A modular inverse constructs a low product of all ones, so the low increment wraps and contributes one to the high word. The control distinguishes losing this carry from the full modular integer result. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
