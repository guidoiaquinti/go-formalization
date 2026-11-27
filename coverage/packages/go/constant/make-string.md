# Independent review of `go/constant.MakeString`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Reuses the empty singleton and otherwise returns a new leaf holding the exact input bytes.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 2 disagreement(s) when allocating without the empty singleton. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `fd0c509553f466e634d5705ca2fe913a3b9352e890fa8bd84f5ed71bce0f96f8`
