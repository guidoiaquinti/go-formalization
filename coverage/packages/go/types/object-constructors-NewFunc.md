# Independent review of `go/types.NewFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `402b6db50e5ccb5dfbda02f3b0c998a30be94122500106406df84bbf58894ba0`

A nil signature produces an absent Type interface and a live signature retains its identity. Native fixtures also check new object identity, zero-valued common fields, empty origin and false function flags.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 4 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
