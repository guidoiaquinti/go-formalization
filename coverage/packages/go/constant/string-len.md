# Independent review of `go/constant.StringLen`

**Verdict: FAITHFUL**

Reviewed author commit `cf7c39d` against Go 1.27.0 `src/go/constant/value.go`. Uses nonmutating string length, returns zero for Unknown and preserves the exact wrong-kind panic text. A typed nil string faults separately.

Pinned Go and Lean matched 500 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal fixtures agreed in both implementations. The control found 1 disagreement(s) when panicking for Unknown. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/constant/string-values/review` after building `godiff-go-constant-string-values`.

The scope covers sequential acyclic trees and shared subtrees, with representable byte counts. Native snapshots check immutable descendants and exact root rewrites. The model abstracts synchronization and takes the formatted argument as an input for wrong-kind errors.

Evidence SHA-256: `4cc5597c389ad9a42a97fe65d478bfd52a23fd8eb9abb61f6e9eb0d7861bc6bc`
