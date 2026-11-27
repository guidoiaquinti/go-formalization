# Independent review of `compress/flate.dictDecoder.init`

**Verdict: FAITHFUL**

The pinned source resets stale cursor and full state, retains the old slice storage whenever its capacity suffices, and sizes the visible slice before copying. An oversized preset contributes its final fitting bytes. A short preset leaves the unfilled part of reused backing intact, and fresh backing is zero-filled outside the copied bytes. The reader and writer positions are equal. Both reset to zero with full set when the copied size fills the visible window, including when that size is zero.

The exact author scope is a nonnil decoder with valid prior geometry, feasible nonnegative sizes, and a disjoint preset. The native overlay inspects the complete old and new backing and compares allocation identity for nonempty storage. At zero capacity it distinguishes nil from nonnil views and abstracts identity between distinct zero-capacity allocations. The Lean theorems check zero-size state and cursor equality. The scoped audit found only permitted axiom dependencies.

Pinned Go 1.27.0 and Lean passed 900 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Nine independent literal controls passed on both implementations. Taking the preset prefix instead of its suffix disagrees on 2 cases. Basing reuse on visible length disagrees on 2, and suppressing the full flag at size zero disagrees on 3.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/dictionary-init/review
```

Evidence SHA-256: `940b8cce7bf39e5318cd90f267e1080e1ce7eb0118b07c15083fb105c6621c3f`
