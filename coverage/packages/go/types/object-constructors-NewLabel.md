# Independent review of `go/types.NewLabel`

**Verdict: FAITHFUL**

Evidence SHA-256: `1404ca790120b57b26da501971a00aaac2284243acd8af8811803f981a9589c0`

The constructor fixes the type to Typ[Invalid] and starts unused. All six rows reject an initially used label.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 4 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
