# Independent review of `go/constant.stringVal.len`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Counts raw bytes without flattening or altering any reachable tree node. Both appearances of a shared subtree contribute to the count. Nil receiver faults.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 1 disagreement(s) when counting a shared subtree only once. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `3d794241c39ba2c37b60a45e518a884d6f385ff1362a1e16b293757b44777d85`
