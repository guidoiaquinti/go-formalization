# Independent review of `go/types.Package.Path`

**Verdict: FAITHFUL**

Reviewed author commit `d2dee3a` against Go 1.27.0 `src/go/types/package.go`. Returns the raw path exactly. Nil receiver panics.

Pinned Go and Lean matched 150 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures agreed in both implementations. The control detected always returning an empty path in 3 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/package-fields/review` after building `godiff-go-types-package-fields`.

Lean uses scope identity tokens. Native compares actual scope identity and observes the constructor's scope and imports. Covered receivers include nil and live packages, including a synthetic live package with no stored scope.

Evidence SHA-256: `408d6ad933f126b6b4aefe423a24cb41e7d436756fd6b2eb7e4a16db5b0dbcb5`
