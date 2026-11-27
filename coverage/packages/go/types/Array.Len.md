# Independent review of `go/types.Array.Len`

**Verdict: FAITHFUL**

Evidence SHA-256: `65ed1708a287a402c5047a205335c705409a139c7a25d0d54e31c4c87756b330`

Returns the stored signed length even if negative. A zero-value array has length zero and a nil receiver panics. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
