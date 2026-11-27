# Independent review of `go/types.NewTypeName`

**Verdict: FAITHFUL**

Evidence SHA-256: `017d91f82e3927648abeabbc3af90240b92538263ab7f256f8b361331379b08d`

The constructor retains the supplied type interface, including both absent and typed-nil values, and preserves the supplied metadata.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 5 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
