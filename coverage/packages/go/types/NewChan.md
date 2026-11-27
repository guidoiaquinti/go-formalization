# Independent review of `go/types.NewChan`

**Verdict: FAITHFUL**

Evidence SHA-256: `c816ea0b39cef27bae6b2ae7b0b90072b85a1e02aefacede71f3bd16b3ac0d5b`

Preserves the supplied direction and element interface, including raw negative direction and typed nil elements. Two native constructions return distinct objects. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
