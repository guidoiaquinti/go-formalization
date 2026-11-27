# Independent review of `go/types.Package.MarkComplete`

**Verdict: FAITHFUL**

Reviewed author commit `d2dee3a` against Go 1.27.0 `src/go/types/package.go`. Sets the completion flag to true and preserves the other modeled fields. Repeating on a complete package preserves its state. Nil receiver panics.

Pinned Go and Lean matched 150 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures agreed in both implementations. The control detected skipping the assignment or accepting a nil receiver in 2 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/package-fields/review` after building `godiff-go-types-package-fields`.

Lean uses scope identity tokens. Native compares actual scope identity and observes the constructor's scope and imports. Covered receivers include nil and live packages, including a synthetic live package with no stored scope.

Evidence SHA-256: `fbe243fe55cb2c48f557a6b132ef800bf2e34f361bca6e410b7ab6cdd1791377`
