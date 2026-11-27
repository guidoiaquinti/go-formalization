# Independent review of `go/types.Func.Pkg`

**Verdict: FAITHFUL**

Evidence SHA-256: `2a04eddda91b9ee48cc60a812b357921a0c0d50ce9534404faf82a9fd7904f35`

The accessor follows the object owner pointer, including nil, and faults on nil receivers. Distinct owner and imported pointers prevent conflating them.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 4 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
