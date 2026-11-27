# Independent review of `go/types.Func.Scope`

**Verdict: FAITHFUL**

Evidence SHA-256: `c6f095451f75bc6d98bd44e6064af5a639a21b015d23bccf81d9404c8d45875b`

A live signature yields its exact scope, including nil. The modeled error denotes a generic panic, with pinned checks for absent type, typed nil, wrong dynamic type and nil receiver.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 5 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
