# Independent review of `go/types.Package.Scope`

**Verdict: FAITHFUL**

Reviewed author commit `d2dee3a` against Go 1.27.0 `src/go/types/package.go`. Returns the stored scope exactly for any live receiver, including a nil stored pointer. A nil receiver returns Universe.

Pinned Go and Lean matched 150 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures agreed in both implementations. The control detected confusing live and nil receiver scope behavior in 4 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/package-fields/review` after building `godiff-go-types-package-fields`.

Lean uses scope identity tokens. Native compares actual scope identity and observes the constructor's scope and imports. Covered receivers include nil and live packages, including a synthetic live package with no stored scope.

Evidence SHA-256: `5e06aa806be2151e65be4ff785efc4aad8ec1dc5b57c5e91d25040d6de73f7e1`
