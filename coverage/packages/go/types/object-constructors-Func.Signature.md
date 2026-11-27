# Independent review of `go/types.Func.Signature`

**Verdict: FAITHFUL**

Evidence SHA-256: `2c4ea2d181691bfaa1a74af53436eb2fbf784bb0e3d8378f3d50d385a1d04928`

An absent type gives a fresh zero signature on each call; typed nil gives nil, wrong type and nil receiver panic, and a live signature keeps its pointer identity.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 5 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
