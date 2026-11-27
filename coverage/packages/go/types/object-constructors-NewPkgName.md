# Independent review of `go/types.NewPkgName`

**Verdict: FAITHFUL**

Evidence SHA-256: `a17123d362bb731a64f717b64b369368189ec5d4bf76c651487fdf9047e0e630`

The constructor installs Typ[Invalid] and keeps the imported package pointer distinct from the owner. The nonnil import rows reject dropping that pointer.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 5 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
