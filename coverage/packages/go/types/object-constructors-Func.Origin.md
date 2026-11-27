# Independent review of `go/types.Func.Origin`

**Verdict: FAITHFUL**

Evidence SHA-256: `9886022de7351c0f3713f8e0d60f296c1994e1a2b71ca0688e39a30650220269`

A set origin takes precedence over the receiver; absent origin returns self and nil receiver panics.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 3 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
