# Independent review of `go/types.Chan.Underlying`

**Verdict: FAITHFUL**

Evidence SHA-256: `c7db00eb86297d08544bd4882de9ad603601135a7c80e1c8ff8cc2f15364540d`

Returns the same channel pointer through a Type interface, including the typed nil channel on a nil receiver. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
