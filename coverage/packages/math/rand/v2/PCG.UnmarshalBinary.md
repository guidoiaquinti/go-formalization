# Independent review of `math/rand/v2.PCG.UnmarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `860fa86de2ae9812e764a6cc1da1862517ffc336bcc8624e10a5249dc32a6b10`

Requires exact length and header before changing either word. Malformed short, long and header inputs returned the exact sentinel and retained the original state. Valid payloads replaced both words. Clearing state on error fails four controls. The pinned replay passed all 68,450 author fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../../lean/GoLean/Math/RandV2/PCG.lean) · [Comparison](../../../../difftest/math/rand/v2/pcg/check.go)

Seven independent native and Lean cases passed, including a deliberate low-limb carry, exact capacity transitions, and valid and malformed wire encodings. Exact allocator capacity is outside the append claim. [Reviewer control](../../../../difftest/math/rand/v2/pcg/review/main.go)
