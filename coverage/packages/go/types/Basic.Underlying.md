# Independent review of `go/types.Basic.Underlying`

**Verdict: FAITHFUL**

Reviewed author commit `e0067b2` against Go 1.27.0 `src/go/types/basic.go` and the Basic arm of `typestring.go`. Returns the receiver pointer as a Type interface. A nil receiver yields a typed nil Basic pointer.

Pinned native and Lean matched 576 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal fixtures matched both. The control detected returning a nil interface for a nil receiver in 1 fixtures. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/basic-fields/review` after building `godiff-go-types-basic-fields`.

Lean abstracts the Basic pointer as an identity token. Native comparison tests the concrete pointer in the initialized unsafe scope and signed native-width fields.

Evidence SHA-256: `c0fbea5baf475ec5144e76a2ad33b6f788e8ded05a57aeb1e84f0789ae504448`
