# Independent review of `encoding/binary.decoder.skip`

**Verdict: FAITHFUL**

Evidence SHA-256: `f41e61c3c97591e2e55cc494f0823f4e4f14396c2a3cd329b83237f39e9f1513`

Go 1.27.0 adds the reflected data size directly to its native signed offset. It does not inspect the buffer, and an unsupported size subtracts one. Ten independent fixtures checked the unsupported decrement, valid zero and positive widths, and wrap in both directions at native int extrema. Retaining the offset for unsupported input failed three fixtures. Saturating at the signed extrema failed two. All 63,840 author cases passed, including 60,000 sampled inputs at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/CoderSkip.lean) · [Comparison](../../../difftest/encoding/binary/coder-skip/check.go)

The model and control use constructed nonnil receivers and representable reflected widths. Native verification exercised 64-bit int arithmetic. [Reviewer control](../../../difftest/encoding/binary/coder-skip/review/main.go)
