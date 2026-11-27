# Independent review of `net/url.URL.Redacted`

**Verdict: FAITHFUL**

Reviewed author commit `e3f3dbc` against Go 1.27.0 `src/net/url/url.go`. Replaces only a set password with xxxxx, preserves username and all other rendered portions and leaves the receiver unchanged. A nil receiver returns empty. A control detects exposing the original password.

Pinned Go and Lean matched 603 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Five independent literal cases passed in both implementations. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/serialize/review` after building `godiff-net-url-serialize`.

Scope covers direct URL values and valid slice geometry. Lean abstracts allocator-selected capacity and zero-capacity backing identity. Native checks nil receiver outcomes.

Evidence SHA-256: `36d2d984aefd553b99d7030b7fc4b957bcc3646553cb08c2b717c251b5217ae7`
