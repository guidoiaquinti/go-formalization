# Independent review of `go/types.Array.Underlying`

**Verdict: FAITHFUL**

Evidence SHA-256: `044685c4f02228368a790f74107818df05a21b9bf75682cf7cfe8b954c2ed7c3`

Returns the same array pointer through a Type interface. A nil receiver is boxed as a typed nil array. Pinned Go and Lean agreed on 61,008 observations, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Go/Types/ArrayChan.lean) · [Comparison](../../../difftest/go/types/array-chan/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/array-chan/review` passed five native and Lean fixtures and rejected 12 altered results across the group. The native fixtures check object identity. Lean represents identities by symbols and models the signed native width. The host replay used darwin/arm64. [Reviewer control](../../../difftest/go/types/array-chan/review/main.go)
