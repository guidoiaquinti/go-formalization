# Independent review of `go/types.Basic.Kind`

**Verdict: FAITHFUL**

Reviewed author commit `e0067b2` against Go 1.27.0 `src/go/types/basic.go` and the Basic arm of `typestring.go`. Reads the stored signed BasicKind word and faults on a nil receiver.

Pinned native and Lean matched 576 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal fixtures matched both. The control detected returning zero unconditionally in 4 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/basic-fields/review` after building `godiff-go-types-basic-fields`.

Lean abstracts the Basic pointer as an identity token. Native comparison tests the concrete pointer in the initialized unsafe scope and signed native-width fields.

Evidence SHA-256: `360c7902402155954c1bde2bc2e95c8a7f15334459115f04e829f5888a21a663`
