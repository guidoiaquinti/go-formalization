# Independent review of `go/constant.StringVal`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Flattens and caches a string root, returns empty for Unknown and preserves wrong-kind panic text. A typed nil string faults separately.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 2 disagreement(s) when omitting lazy string caching. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `2a26576a3de80be2a88008163efef0b7e1c207d1fcd5caf77ddb787bac7ccc36`
