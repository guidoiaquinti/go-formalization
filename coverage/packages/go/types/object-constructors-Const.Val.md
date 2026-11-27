# Independent review of `go/types.Const.Val`

**Verdict: FAITHFUL**

Evidence SHA-256: `5b9b7c9aa4d6c563c3e56f3342d1f55b50a909f47f68d770c6f4ced4b1846e7a`

The accessor returns the stored interface exactly and faults on a nil receiver. A typed-nil value and a nonnil constant separate retrieval from a nil default.

I checked the pinned Go 1.27.0 source in `src/go/types/object.go` against the model and replayed `go run ./difftest/go/types/object-constructors` on darwin/arm64. The direct comparison passed 864 boundary rows and 60,000 sampled rows with seeds 17, 2718 and 65537. The scoped Lean kernel audit passed.

`go run ./difftest/go/types/object-constructors/review` checks six fixed Go and Lean cases with literal expected outcomes. The fixed cases reject the output mutant for this declaration in 4 rows. Lean represents reference identity by tokens, and the native overlay checks actual pointers, allocation freshness and constructor defaults. Allocator internals and exact panic diagnostics are outside scope.
