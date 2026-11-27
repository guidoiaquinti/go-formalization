# Independent review of `go/types.NewConst`

**Verdict: FAITHFUL**

Evidence SHA-256: `50d2c06a2b2d83129d5f74686c53b1f05c8a4855ee6c03c7cab1f750fa6316e9`

The constructor stores the exact supplied object fields and constant value; the typed-nil value distinguishes erased payloads from preserved interfaces.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 5 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
