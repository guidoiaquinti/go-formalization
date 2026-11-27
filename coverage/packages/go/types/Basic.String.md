# Independent review of `go/types.Basic.String`

**Verdict: FAITHFUL**

Reviewed author commit `e0067b2` against Go 1.27.0 `src/go/types/basic.go` and the Basic arm of `typestring.go`. The Basic TypeString branch qualifies an exported name only when the pinned unsafe scope has a TypeName for it. Pointer is the sole applicable name. Sizeof is a builtin and remains unchanged. Other raw names are preserved, and a nil Basic receiver faults.

Pinned native and Lean matched 576 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal fixtures matched both. The control detected omitting the unsafe.Pointer qualification in 1 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/basic-fields/review` after building `godiff-go-types-basic-fields`.

Lean abstracts the Basic pointer as an identity token. Native comparison tests the concrete pointer in the initialized unsafe scope and signed native-width fields.

Evidence SHA-256: `b3e2360c7dc8dffcb2eb03875300987c2f929a9a27efb130f5625dc06fcdb695`
