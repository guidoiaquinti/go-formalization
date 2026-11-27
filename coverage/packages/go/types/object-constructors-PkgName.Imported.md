# Independent review of `go/types.PkgName.Imported`

**Verdict: FAITHFUL**

Evidence SHA-256: `e82a1ed587a099ce66d56e781331b5c9fdb4605920a8c28d6f88076851b336a0`

The accessor returns the stored imported package pointer and faults on a nil receiver. Fixtures with distinct owner and import pointers verify identity.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 6 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
