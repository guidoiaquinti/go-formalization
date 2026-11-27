# Independent review of `encoding/binary.encoder.skip`

**Verdict: FAITHFUL**

Evidence SHA-256: `3e53a240e4643c7ed7b70ab45ae208060d06961acfdaad421fc517ba59e04f50`

Go 1.27.0 evaluates the reflected width, bounds-checks a two-index view against capacity, clears that view, then advances the offset. A bounds failure preserves the complete backing and original offset, including unsupported widths and arithmetic wrap. An empty clear is permitted exactly at capacity and still checks the bounds. All 63,840 author cases passed, including 60,000 sampled inputs at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/CoderSkip.lean) · [Comparison](../../../difftest/encoding/binary/coder-skip/check.go)

Ten independent fixtures passed on native Go and Lean. Checking logical length at the start instead of capacity failed two controls, advancing the offset before the clear failed five, and bypassing bounds for an empty clear failed one. The controls also checked partial and full spare-byte clears and preserved trailing bytes. Native verification exercised 64-bit int arithmetic on constructed nonnil receivers. [Reviewer control](../../../difftest/encoding/binary/coder-skip/review/main.go)
