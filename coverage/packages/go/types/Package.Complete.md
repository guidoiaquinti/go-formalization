# Independent review of `go/types.Package.Complete`

**Verdict: FAITHFUL**

Reviewed author commit `d2dee3a` against Go 1.27.0 `src/go/types/package.go`. Returns the stored Boolean flag. Nil receiver panics.

Pinned Go and Lean matched 150 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures agreed in both implementations. The control detected always returning false in 3 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/package-fields/review` after building `godiff-go-types-package-fields`.

Lean uses scope identity tokens. Native compares actual scope identity and observes the constructor's scope and imports. Covered receivers include nil and live packages, including a synthetic live package with no stored scope.

Evidence SHA-256: `7bde1afea12f1407803722de35b7408f7e248fe803141b9111cb37dd3216ab44`
