# Independent review of `go/types.Array.Elem`

**Verdict: FAITHFUL**

Evidence SHA-256: `f80f13e45365c12e10ac487087261e0bb7bd042df40613a33543638637f514db`

Returns the identical stored interface, including typed nil values, and panics when the array receiver is nil. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
