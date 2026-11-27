# Independent review of `go/constant.stringVal.string`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Flattens raw bytes in leaf order, writes the flattened root and drops its child links. Nonroot nodes retain their exact state. Nil receiver faults.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 3 disagreement(s) when returning a lazy uncached root. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `8d46a6564ac462920b2875129fa02660789ac63509a3472df691b30b7b5548ed`
