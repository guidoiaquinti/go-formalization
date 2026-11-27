# Independent review of `go/types.Chan.Elem`

**Verdict: FAITHFUL**

Evidence SHA-256: `36e3afd5f2d673dceeb8ccd2efe40a28ba5174a74a4d6cd30f8d73956050923d`

Returns the identical stored interface, preserving typed nil values. A nil channel receiver panics. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
