# Independent review of `go/types.Basic.Name`

**Verdict: FAITHFUL**

Reviewed author commit `e0067b2` against Go 1.27.0 `src/go/types/basic.go` and the Basic arm of `typestring.go`. Returns the exact stored bytes and faults on a nil receiver. The fixtures include a zero-value name and malformed UTF-8.

Pinned native and Lean matched 576 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal fixtures matched both. The control detected discarding the name and returning an empty string in 4 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/basic-fields/review` after building `godiff-go-types-basic-fields`.

Lean abstracts the Basic pointer as an identity token. Native comparison tests the concrete pointer in the initialized unsafe scope and signed native-width fields.

Evidence SHA-256: `ea90d96dd4e981d22aeedb2f5f1f558739bca572fee51681c5db417fd1203a0c`
