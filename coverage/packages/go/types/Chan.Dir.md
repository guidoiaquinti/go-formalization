# Independent review of `go/types.Chan.Dir`

**Verdict: FAITHFUL**

Evidence SHA-256: `a231ea9b878b6176ea830edfb59d1fb4dd2b90d01fedec9efb8986b8ac2e5108`

Returns the raw stored direction, even outside the documented enum values. A zero-value channel has direction zero and a nil receiver panics. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
