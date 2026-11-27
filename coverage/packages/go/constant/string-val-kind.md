# Independent review of `go/constant.stringVal.Kind`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Returns code two even on a typed nil string receiver because the method never dereferences it.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 1 disagreement(s) when faulting on a nil receiver. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `51e4abf3e3996857e37753c386679935de9392798d5ef9f8e7b388e390f2a12c`
