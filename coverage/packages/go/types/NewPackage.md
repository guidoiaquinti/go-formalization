# Independent review of `go/types.NewPackage`

**Verdict: FAITHFUL**

Reviewed author commit `d2dee3a` against Go 1.27.0 `src/go/types/package.go`. New package and child scope identities are fresh in native checks. The scope starts empty with Universe as parent, no range and a comment based on the quoted path. Path and name are retained. Completion is false, imports are nil and Go version is empty.

Pinned Go and Lean matched 150 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Four independent literal fixtures agreed in both implementations. The control detected starting complete in 4 cases. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/go/types/package-fields/review` after building `godiff-go-types-package-fields`.

Lean uses scope identity tokens. Native compares actual scope identity and observes the constructor's scope and imports. Covered receivers include nil and live packages, including a synthetic live package with no stored scope.

Evidence SHA-256: `ce653b9a20b698afb68f607f561ec7cefb806480fb3407289bc82bb329e8adc6`
